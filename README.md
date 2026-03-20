# Micepad CLI

Command-line interface for the [Micepad](https://micepad.co) event management platform. Manage events, participants, check-ins, campaigns, and more — right from your terminal.

## How It Works

Micepad CLI is a thin client powered by [Terminalwire](https://terminalwire.com). All commands execute on the Micepad server via a WebSocket connection — the CLI itself contains no business logic. This means you always get the latest features without updating the gem.

## Installation

**Requirements:** Ruby 3.0+

### From RubyGems

```bash
gem install micepad-cli
```

### From GitHub

```bash
gem install specific_install
gem specific_install https://github.com/micepadteam/micepad-cli
```

### From source

```bash
git clone https://github.com/micepadteam/micepad-cli.git
cd micepad-cli
gem build micepad-cli.gemspec
gem install ./micepad-cli-0.1.0.gem
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
micepad terminal tree
```

## Configuration

By default, the CLI connects to `wss://alpha.micepad.co/terminal`.

Override the server URL with the `MICEPAD_URL` environment variable:

```bash
# Connect to production
MICEPAD_URL="wss://studio.micepad.co/terminal" micepad login

# Connect to local development server
MICEPAD_URL="ws://localhost:3000/terminal" micepad login

# Export for persistent use
export MICEPAD_URL="wss://studio.micepad.co/terminal"
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
| `micepad events use SLUG` | Set the active event context |
| `micepad events current` | Show details of the active event |
| `micepad events stats` | Show event dashboard statistics |

### Participants

| Command | Description |
|---------|-------------|
| `micepad pax list` | List participants (with filters) |
| `micepad pax show ID` | Show participant details |
| `micepad pax add` | Add a new participant |
| `micepad pax checkin ID` | Check in a participant |
| `micepad pax checkout ID` | Check out a participant |
| `micepad pax count` | Count participants by status |
| `micepad pax export` | Export participants to CSV |
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

```bash
# Show storage directory path
micepad pax import --storage-path

# Download import template
micepad pax import --template
micepad pax import --template --format xlsx

# Import from CSV (interactive — prompts for group, mappings, confirmation)
cp attendees.csv $(micepad pax import --storage-path)/
micepad pax import attendees.csv

# Import with options
micepad pax import attendees.csv --group "VIP" --action add --yes

# Dry run (validate only, no changes)
micepad pax import attendees.csv --dry-run

# Export failed rows
micepad pax import attendees.csv --errors-out errors.csv
```

### Check-ins

| Command | Description |
|---------|-------------|
| `micepad checkins stats` | Show check-in statistics with velocity |
| `micepad checkins stats --watch` | Live-refresh stats every 2 seconds |
| `micepad checkins recent` | Show recent check-in activity |

### Campaigns

| Command | Description |
|---------|-------------|
| `micepad campaigns list` | List campaigns |
| `micepad campaigns create` | Create a new campaign |
| `micepad campaigns show ID` | Show campaign details |
| `micepad campaigns send ID` | Send a campaign |
| `micepad campaigns cancel ID` | Cancel a scheduled campaign |
| `micepad campaigns stats ID` | Show campaign delivery statistics |

```bash
# List email campaigns only
micepad campaigns list --type email

# Create campaign from template
micepad campaigns create --type email --name "Welcome" --template tpl_xxxxx

# Watch delivery stats in real-time
micepad campaigns stats cmp_xxxxx --watch
```

### Templates

| Command | Description |
|---------|-------------|
| `micepad templates list` | List all templates |
| `micepad templates list --type email` | Filter by type (email, whatsapp) |

### Groups

| Command | Description |
|---------|-------------|
| `micepad groups list` | List all groups |
| `micepad groups show NAME` | Show group details with RSVP breakdown |

### Admin (super admin only)

| Command | Description |
|---------|-------------|
| `micepad admin dashboard` | Platform-wide statistics (users, DAU, MAU) |
| `micepad admin accounts` | List accounts |
| `micepad admin users` | List users |
| `micepad admin gatherings` | List all gatherings |
| `micepad admin subscriptions` | List active subscriptions |

```bash
# Search accounts
micepad admin accounts --search "acme"

# Filter gatherings by status
micepad admin gatherings --status published --limit 50
```

### Meta

| Command | Description |
|---------|-------------|
| `micepad terminal tree` | Show full command tree |
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
