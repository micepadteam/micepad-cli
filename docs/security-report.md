# Security Report: Micepad CLI (Go Rewrite)

**Date:** 2026-03-26
**Scope:** Terminalwire protocol client — file access, environment variables, browser, stdin/stdout

## Architecture

The CLI is a thin client. The server drives all interactions over WebSocket using the Terminalwire protocol. The client executes resource requests (file read/write, directory operations, stdin/stdout, browser open, env vars) and returns results.

```
Server ──[WebSocket]──▶ CLI Client ──▶ Local filesystem, env vars, browser
```

The client declares **entitlement paths** in the init message, but historically did not enforce them. The server was implicitly trusted.

## Threat Model

| Direction | Risk Level | Description |
|-----------|------------|-------------|
| Server → Client | **Critical** | A compromised server can instruct the client to read/write arbitrary files, leak env vars, or open malicious URLs on the user's machine |
| Client → Server | Low | Standard input validation concerns (CSV injection, oversized payloads) — handled by Rails |

## Vulnerabilities Found & Fixed

### 1. Arbitrary File Read/Write (Critical)

**Before:** `handleFile` and `handleDirectory` accepted any path from the server with no validation. A compromised server could read `~/.ssh/id_ed25519`, `~/.aws/credentials`, `/etc/passwd`, or write to `~/.zshrc`.

**Attack vectors:**
- Absolute paths: `{command:"read", path:"/etc/passwd"}`
- Tilde expansion: `{command:"read", path:"~/.ssh/id_ed25519"}`
- Path traversal: `{command:"read", path:"authorities/../../.ssh/id_ed25519"}`

**Verified:** All three vectors successfully read files in testing.

**Fix:** Added `isAllowedPath()` check — all file and directory operations must resolve within `~/.micepad/authorities/<authority>/storage/`. Uses `filepath.Clean()` to neutralize `..` traversal before prefix comparison.

**Files:** `internal/terminalwire/resources.go` — `isAllowedPath()`, enforced in `handleFile()` and `handleDirectory()`

### 2. Environment Variable Leak (High)

**Before:** `handleEnvVar` read any environment variable the server requested — `AWS_SECRET_ACCESS_KEY`, `GITHUB_TOKEN`, `DATABASE_URL`, etc.

**Fix:** Allowlist restricted to `TERMINALWIRE_HOME` and `MICEPAD_HOME`. All other requests return "access denied."

**File:** `internal/terminalwire/resources.go` — `allowedEnvVars` map

### 3. Browser Scheme Abuse (Medium)

**Before:** `handleBrowser` opened any URL — including `file:///`, `javascript:`, or phishing URLs.

**Fix:** Only `http://` and `https://` schemes are allowed.

**File:** `internal/terminalwire/resources.go` — `handleBrowser()`

### 4. Relative Path Resolution (Medium)

**Before:** `expandPath()` only handled `~/` prefix. The server sends relative paths (e.g., `authorities/.../storage/file.csv`), which resolved against CWD instead of `~/.micepad/`. This broke file operations and could expose files in the working directory.

**Fix:** Non-absolute paths now resolve against `~/.micepad/` (the Terminalwire home directory).

**File:** `internal/terminalwire/resources.go` — `expandPath()`

### 5. File Size Limit (Denial of Service)

**Before:** No size limit on file reads or uploads. A malicious file or server request could cause memory exhaustion.

**Fix:** 10 MB limit enforced at three layers:
- **CLI `handleFile` read:** Checks `os.Stat` size before `os.ReadFile`
- **CLI `prepareFileArgs`:** Skips files over 10 MB when copying to storage
- **Server `import_terminal.rb`:** Checks `file_content.bytesize` after reading, rejects with user-friendly message
- **Server `application_terminal.rb`:** `resolve_file_content` checks `File.size` before reading

**Files:** `resources.go`, `client.go`, `import_terminal.rb`, `application_terminal.rb`

### 6. Stdin Buffer Loss (Bug, not security)

**Before:** `handleStdin` created a new `bufio.NewReader(os.Stdin)` on every `read_line` call, discarding buffered data. This broke piped input (e.g., scripted imports).

**Fix:** Single `bufio.Reader` stored on the `Client` struct, reused across all stdin reads.

**Files:** `client.go`, `resources.go`

## Security Model After Fixes

```
Server request
     │
     ▼
┌─────────────────────────────┐
│ Resource type?              │
├─────────────────────────────┤
│ file/directory              │
│  → isAllowedPath()?         │──▶ NO  → "access denied"
│  → size < 10MB?             │──▶ NO  → "file too large"
│  → proceed                  │
├─────────────────────────────┤
│ environment_variable        │
│  → in allowlist?            │──▶ NO  → "access denied"
├─────────────────────────────┤
│ browser                     │
│  → http/https scheme?       │──▶ NO  → "access denied"
├─────────────────────────────┤
│ stdout/stderr/stdin         │
│  → always allowed (I/O)     │
└─────────────────────────────┘
```

## Remaining Considerations

| Area | Status | Notes |
|------|--------|-------|
| TLS enforcement | Mitigated | `wss://` used by default; `MICEPAD_URL` override could use `ws://` for local dev |
| WebSocket origin validation | Server-side | Ensure server validates WebSocket origin headers |
| Rate limiting | Server-side | Protect against rapid reconnect or request flooding |
| JWT token storage | Acceptable | `session.jwt` stored in `~/.micepad/` with user-only permissions |
| `MICEPAD_URL` spoofing | User responsibility | A malicious URL connects to a rogue server; user must trust the URL they set |
| CSV/XLSX injection | Server-side | Server should sanitize imported data before rendering in web UI |

## Recommendations

1. **Consider `ws://` warning:** Print a warning when `MICEPAD_URL` uses `ws://` (unencrypted) outside of `localhost`.
2. **Pin server certificate:** For production use, consider certificate pinning to prevent MITM even with valid TLS.
3. **Audit log:** Log all file access operations client-side for forensic review.
