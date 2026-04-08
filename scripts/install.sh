#!/usr/bin/env bash
# install.sh — Install micepad CLI
#
# Usage:
#   curl -fsSL https://micepad.co/install-cli | bash
#
# Options (via environment):
#   MICEPAD_BIN_DIR     Where to install binary (default: ~/.local/bin)
#   MICEPAD_VERSION     Specific version to install (default: latest)

set -euo pipefail

REPO="micepad/micepad-cli"
BIN_DIR="${MICEPAD_BIN_DIR:-$HOME/.local/bin}"
VERSION="${MICEPAD_VERSION:-}"

# Color helpers — respect NO_COLOR (https://no-color.org)
if [[ -z "${NO_COLOR:-}" ]] && [[ -t 1 ]]; then
  bold()  { printf '\033[1m%s\033[0m' "$1"; }
  green() { printf '\033[32m%s\033[0m' "$1"; }
  red()   { printf '\033[31m%s\033[0m' "$1"; }
  dim()   { printf '\033[2m%s\033[0m' "$1"; }
else
  bold()  { printf '%s' "$1"; }
  green() { printf '%s' "$1"; }
  red()   { printf '%s' "$1"; }
  dim()   { printf '%s' "$1"; }
fi

info()  { echo "  $(green "✓") $1"; }
step()  { echo "  $(bold "→") $1"; }
error() { echo "  $(red "✗ ERROR:") $1" >&2; exit 1; }

find_sha256_cmd() {
  if command -v sha256sum &>/dev/null; then
    echo "sha256sum"
  elif command -v shasum &>/dev/null; then
    echo "shasum -a 256"
  else
    error "No SHA256 tool found (need sha256sum or shasum)"
  fi
}

detect_platform() {
  local os arch

  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$os" in
    darwin) os="darwin" ;;
    linux)  os="linux" ;;
    *) error "Unsupported OS: $os" ;;
  esac

  arch=$(uname -m)
  case "$arch" in
    x86_64|amd64)   arch="amd64" ;;
    aarch64|arm64)   arch="arm64" ;;
    *) error "Unsupported architecture: $arch" ;;
  esac

  echo "${os}_${arch}"
}

get_latest_version() {
  local url version
  url=$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/${REPO}/releases/latest" 2>/dev/null) || true
  version="${url##*/}"
  version="${version#v}"
  if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    error "Could not determine latest version (resolved '${version:-<empty>}' from '${url:-<no URL>}'). Check your network connection or https://github.com/${REPO}/releases"
  fi
  echo "$version"
}

verify_checksums() {
  local version="$1" tmp_dir="$2" archive_name="$3"
  local base_url="https://github.com/${REPO}/releases/download/v${version}"

  step "Verifying checksums..."

  if ! curl -fsSL "${base_url}/checksums.txt" -o "${tmp_dir}/checksums.txt"; then
    error "Failed to download checksums.txt"
  fi

  local expected actual
  expected=$(awk -v f="$archive_name" '$2 == f || $2 == ("*" f) {print $1; exit}' "${tmp_dir}/checksums.txt")
  actual=$(cd "$tmp_dir" && $(find_sha256_cmd) "$archive_name" | awk '{print $1}')
  [[ -n "$expected" && "$expected" == "$actual" ]] \
    || error "Checksum verification failed for $archive_name"

  info "Checksum verified"
}

download_binary() {
  local version="$1" platform="$2" tmp_dir="$3"
  local archive_name="micepad_${version}_${platform}.tar.gz"
  local url="https://github.com/${REPO}/releases/download/v${version}/${archive_name}"

  step "Downloading micepad v${version} for ${platform}..."

  if ! curl -fsSL "$url" -o "${tmp_dir}/${archive_name}"; then
    error "Failed to download from $url"
  fi

  verify_checksums "$version" "$tmp_dir" "$archive_name"

  step "Extracting..."
  tar -xzf "${tmp_dir}/${archive_name}" -C "$tmp_dir"

  if [[ ! -f "${tmp_dir}/micepad" ]]; then
    error "Binary not found in archive"
  fi

  mkdir -p "$BIN_DIR"
  mv "${tmp_dir}/micepad" "$BIN_DIR/"
  chmod +x "$BIN_DIR/micepad"

  info "Installed micepad to $BIN_DIR/micepad"
}

setup_path() {
  if [[ ":$PATH:" == *":$BIN_DIR:"* ]]; then
    return 0
  fi

  step "Adding $BIN_DIR to PATH..."

  local shell_rc=""
  case "${SHELL:-}" in
    */zsh)  shell_rc="$HOME/.zshrc" ;;
    */bash) shell_rc="$HOME/.bashrc" ;;
    *)      shell_rc="$HOME/.profile" ;;
  esac

  local path_line="export PATH=\"$BIN_DIR:\$PATH\""

  if [[ -f "$shell_rc" ]] && grep -qF "$BIN_DIR" "$shell_rc" 2>/dev/null; then
    info "PATH already configured in $shell_rc"
  else
    echo "" >> "$shell_rc"
    echo "# Added by micepad installer" >> "$shell_rc"
    echo "$path_line" >> "$shell_rc"
    info "Added to $shell_rc"
    info "Run: source $shell_rc"
  fi
}

verify_install() {
  local installed_version
  if installed_version=$("$BIN_DIR/micepad" --version 2>/dev/null); then
    info "$(green "${installed_version} installed")"
    return 0
  fi
  error "Installation failed — micepad not working"
}

install_skill() {
  if ! command -v npx &>/dev/null; then
    dim "  Skipped skill install (npx not found)"
    return
  fi

  # If skills already installed, check for updates instead of re-adding
  if [[ -d "$HOME/.claude/skills/micepad" ]]; then
    local check_output
    check_output=$(npx -y skills check 2>/dev/null || true)
    if echo "$check_output" | grep -q "update(s) available"; then
      step "Updating Micepad skills for Claude Code..."
      if npx -y skills update 2>/dev/null; then
        info "Micepad skills updated"
      else
        dim "  Skill update skipped (non-critical)"
      fi
    else
      info "Micepad skills already up to date"
    fi
    return
  fi

  step "Installing Micepad skill for Claude Code..."
  # Use < /dev/tty so the interactive agent-selection prompt works
  # even when this script is piped from curl (curl ... | bash)
  if npx -y skills add micepad/skills -g < /dev/tty 2>/dev/null; then
    # Remove find-skills that gets bundled by the skills CLI
    rm -rf "$HOME/.claude/skills/find-skills" 2>/dev/null
    info "Micepad skills installed (use /micepad or /micepad-admin in Claude Code)"
  else
    dim "  Skipped skill install (non-critical)"
  fi
}

show_banner() {
  local cols
  cols=$(tput cols 2>/dev/null || echo 80)

  if [[ "$cols" -ge 40 ]]; then
    local c1="" c2="" c3="" c4="" c5="" c6="" b="" r=""
    if [[ -z "${NO_COLOR:-}" ]] && [[ -t 1 ]]; then
      c1=$'\033[38;2;255;107;107m'  # Red
      c2=$'\033[38;2;255;179;71m'   # Orange
      c3=$'\033[38;2;255;217;61m'   # Yellow
      c4=$'\033[38;2;72;219;153m'   # Green
      c5=$'\033[38;2;99;102;241m'   # Indigo
      c6=$'\033[38;2;186;103;255m'  # Purple
      b=$'\033[1m'
      r=$'\033[0m'
    fi

    echo ""
    echo "${c1}    __  ____                           __${r}"
    echo "${c2}   /  |/  (_)_______  ____  ____ _____/ /${r}"
    echo "${c3}  / /|_/ / / ___/ _ \\/ __ \\/ __ \`/ __  /${r}"
    echo "${c4} / /  / / / /__/  __/ /_/ / /_/ / /_/ /${r}"
    echo "${c5}/_/  /_/_/\\___/\\___/ .___/\\__,_/\\__,_/${r}"
    echo "${c6}                  /_/            ${b}CLI${r}"
    echo ""
  else
    echo ""
    echo "  $(bold "Micepad CLI")"
    echo ""
  fi
}

main() {
  show_banner

  if ! command -v curl &>/dev/null; then
    error "curl is required but not installed"
  fi

  local platform version tmp_dir
  platform=$(detect_platform)

  if [[ -n "$VERSION" ]]; then
    version="$VERSION"
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]]; then
      error "Invalid version '${version}'. Expected semver format (e.g. 1.2.3 or 1.2.3-rc.1)."
    fi
  else
    version=$(get_latest_version)
  fi

  tmp_dir=$(mktemp -d)
  trap "rm -rf '${tmp_dir}'" EXIT

  # Check if already installed and up to date
  local current_version=""
  if command -v micepad &>/dev/null; then
    current_version=$(micepad --version 2>/dev/null | awk '{print $2}' || true)
  fi

  if [[ -n "$current_version" && "$current_version" == "$version" ]]; then
    info "Already up to date (v${version})"
    setup_path
    install_skill
    return
  fi

  if [[ -n "$current_version" ]]; then
    step "Updating v${current_version} → v${version}"
  fi

  download_binary "$version" "$platform" "$tmp_dir"
  setup_path
  verify_install
  install_skill

  echo ""
  echo "  Next steps:"
  echo "    $(bold "micepad login")          Authenticate with Micepad"
  echo "    $(bold "micepad help")           See all commands"
  echo "    $(bold "micepad events list")    List your events"
  echo "    $(bold "micepad env")            Switch environments (prod/alpha/dev)"
  echo "    $(bold "micepad update")         Update to the latest version"
  echo ""
}

main "$@"
