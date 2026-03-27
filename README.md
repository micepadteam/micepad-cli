# Micepad CLI

Command-line interface for the [Micepad](https://micepad.co) event management platform. Manage events, participants, check-ins, campaigns, and more — right from your terminal.

## How It Works

Micepad CLI is a thin client powered by [Terminalwire](https://terminalwire.com). All commands execute on the Micepad server via a WebSocket connection — the CLI itself contains no business logic. This means you always get the latest features without updating the binary.

## Installation

```bash
curl -fsSL https://github.com/micepadteam/micepad-cli/releases/latest/download/install.sh | bash
```

This downloads the prebuilt binary for your platform (macOS/Linux, amd64/arm64), verifies checksums, and adds it to your PATH.

To update:

```bash
micepad update
```

## Quick Start

```bash
# Authenticate (opens browser for authorization)
micepad login

# List your events
micepad events list

# Select an event to work with
micepad events use my-event-slug

# Check who you are and what's selected
micepad whoami

# See all available commands
micepad tree
```

## Configuration

By default, the CLI connects to `wss://studio.micepad.co/terminal`.

Override the server URL:

```bash
# Set via configure command
micepad configure --url "wss://studio.micepad.co/terminal"

# Or via environment variable
export MICEPAD_URL="wss://studio.micepad.co/terminal"

# Connect to local development server
MICEPAD_URL="ws://localhost:3000/terminal" micepad login
```

## Commands

### Authentication

| Command | Description |
|---------|-------------|
| `micepad login` | Authenticate via browser-based authorization |
| `micepad logout` | Clear session |
| `micepad whoami` | Show current user, account, and event |

### Events

| Command | Description |
|---------|-------------|
| `micepad events list` | List all events in your account |
| `micepad events create` | Create a new event |
| `micepad events use SLUG` | Set the active event context |
| `micepad events current` | Show details of the active event |
| `micepad events stats` | Show event dashboard statistics |

### Participants

| Command | Description |
|---------|-------------|
| `micepad pax list` | List participants (with filters) |
| `micepad pax show ID` | Show participant details |
| `micepad pax add` | Add a new participant |
| `micepad pax update ID` | Update a participant |
| `micepad pax checkin ID` | Check in a participant |
| `micepad pax checkout ID` | Check out a participant |
| `micepad pax count` | Count participants by status |
| `micepad pax export` | Export participants to CSV/XLSX |
| `micepad pax import FILE` | Import participants from CSV/Excel |

#### Filtering participants

```bash
# Filter by RSVP status
micepad pax list --status confirmed

# Filter by check-in status
micepad pax list --checkin checked_in

# Filter by group
micepad pax list --group "VIP"

# Search by name or email
micepad pax list --search "john"

# Combine filters
micepad pax list --status confirmed --checkin not_checked_in --limit 100
```

#### Importing participants

The CLI automatically copies local files to its sandboxed storage — just pass the file path directly:

```bash
# Import from CSV/Excel (interactive wizard)
micepad pax import ~/Downloads/attendees.csv

# Import with options
micepad pax import attendees.xlsx --group "Speakers" --action add --yes

# Dry run (validate only, no changes)
micepad pax import attendees.csv --dry-run

# Download import template
micepad pax import --template
micepad pax import --template --format xlsx
```

### Groups

| Command | Description |
|---------|-------------|
| `micepad groups list` | List all groups |
| `micepad groups create` | Create a new group |
| `micepad groups show NAME` | Show group details with RSVP breakdown |

### Registration Types

| Command | Description |
|---------|-------------|
| `micepad regtypes list` | List registration types with capacity |
| `micepad regtypes create` | Create a registration type |

### Forms

| Command | Description |
|---------|-------------|
| `micepad forms list` | List forms |
| `micepad forms fields ID` | List all fields including hidden ones |
| `micepad forms add-field ID` | Add a field |
| `micepad forms update-field ID SLUG` | Update a field |
| `micepad forms reorder ID` | Reorder form fields |
| `micepad forms update ID` | Update form settings |
| `micepad forms publish ID` | Publish form (makes it live) |
| `micepad forms unpublish ID` | Close registration |
| `micepad forms url ID` | Get the public registration URL |

### Check-ins

| Command | Description |
|---------|-------------|
| `micepad checkins stats` | Show check-in statistics with velocity |
| `micepad checkins stats --watch` | Live-refresh stats every 2 seconds |
| `micepad checkins recent` | Show recent check-in activity |
| `micepad checkins add-staff` | Add check-in staff |
| `micepad checkins remove-staff EMAIL` | Remove staff member |
| `micepad checkins staff` | List check-in staff |
| `micepad checkins staff-activity` | Staff performance stats |

### Campaigns

| Command | Description |
|---------|-------------|
| `micepad campaigns list` | List campaigns |
| `micepad campaigns create` | Create a new campaign |
| `micepad campaigns show ID` | Show campaign details |
| `micepad campaigns update ID` | Update campaign settings |
| `micepad campaigns add-section ID` | Add a content section |
| `micepad campaigns sections ID` | List all sections |
| `micepad campaigns add-recipients ID` | Add recipients |
| `micepad campaigns send ID` | Send a campaign |
| `micepad campaigns cancel ID` | Cancel a scheduled campaign |
| `micepad campaigns stats ID` | Show delivery statistics |

### Badges

| Command | Description |
|---------|-------------|
| `micepad badges list` | List badge templates |
| `micepad badges create` | Create a badge template |
| `micepad badges show ID` | Show badge template with fields |
| `micepad badges add-field ID` | Add a field to badge template |

### QR Login Tokens (Kiosks)

| Command | Description |
|---------|-------------|
| `micepad qrlogin generate` | Create a kiosk access token |
| `micepad qrlogin list` | List active tokens |
| `micepad qrlogin revoke ID` | Revoke a token |

### Templates

| Command | Description |
|---------|-------------|
| `micepad templates list` | List all templates |

### Meta

| Command | Description |
|---------|-------------|
| `micepad tree` | Show full command tree |
| `micepad help` | Show top-level help |
| `micepad events help` | Show subcommand help |

## Participant Identifiers

Several commands accept a participant identifier. You can use any of:

- **Prefix ID**: `pax_abc123`
- **Email**: `john@example.com`
- **Code**: the participant's registration code
- **QR Code**: the participant's QR code value

```bash
micepad pax show pax_abc123
micepad pax show john@example.com
micepad pax checkin pax_abc123
```

## License

MIT
