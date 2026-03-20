#compdef micepad
# Zsh completion for micepad CLI
# Add to ~/.zshrc: fpath=(/path/to/lib/completions $fpath) && compinit
# Or: source /path/to/lib/completions/micepad.zsh

_micepad() {
  local -a root_commands
  root_commands=(
    'login:Login to your account'
    'logout:Logout of your account'
    'whoami:Show current user and account'
    'events:Manage events'
    'pax:Manage participants'
    'campaigns:Manage campaigns'
    'templates:Manage templates'
    'checkins:Check-in operations'
    'groups:Manage groups'
    'admin:Admin operations (super admin only)'
  )

  _arguments -C \
    '1:command:->command' \
    '*::arg:->args'

  case $state in
    command)
      _describe 'command' root_commands
      ;;
    args)
      case $words[1] in
        events) _micepad_events ;;
        pax) _micepad_pax ;;
        campaigns) _micepad_campaigns ;;
        templates) _micepad_templates ;;
        checkins) _micepad_checkins ;;
        groups) _micepad_groups ;;
        admin) _micepad_admin ;;
      esac
      ;;
  esac
}

_micepad_events() {
  local -a subcmds
  subcmds=(
    'list:List your events'
    'use:Set default event context'
    'current:Show current event details'
    'stats:Show event dashboard statistics'
  )

  _arguments -C \
    '1:subcommand:->subcmd' \
    '*::arg:->args'

  [[ $state == subcmd ]] && _describe 'subcommand' subcmds
}

_micepad_pax() {
  local -a subcmds
  subcmds=(
    'list:List participants'
    'show:Show participant details'
    'add:Add a new participant'
    'checkin:Check in a participant'
    'checkout:Check out a participant'
    'count:Count participants by status'
    'export:Export participants to CSV'
  )

  _arguments -C \
    '1:subcommand:->subcmd' \
    '*::arg:->args'

  case $state in
    subcmd) _describe 'subcommand' subcmds ;;
    args)
      case $words[1] in
        list)
          _arguments \
            '--status[Filter by RSVP status]:status:(confirmed unconfirmed declined waitlisted pending_approval)' \
            '--checkin[Filter by check-in status]:status:(checked_in not_checked_in checked_out)' \
            '--group[Filter by group name]:group:' \
            '--search[Search by name or email]:query:' \
            '--limit[Max rows to show]:number:'
          ;;
        show) _arguments '1:identifier:' ;;
        add)
          _arguments \
            '--email[Email address (required)]:email:' \
            '--first_name[First name (required)]:name:' \
            '--last_name[Last name]:name:' \
            '--company[Company name]:company:' \
            '--job_title[Job title]:title:'
          ;;
        checkin) _arguments '1:identifier:' ;;
        checkout) _arguments '1:identifier:' ;;
        count)
          _arguments \
            '--by[Group by]:field:(rsvp checkin group)'
          ;;
        export)
          _arguments \
            '--status[Filter by RSVP status]:status:(confirmed unconfirmed declined waitlisted pending_approval)' \
            '--output[Output file path]:file:_files -g "*.csv"'
          ;;
      esac
      ;;
  esac
}

_micepad_campaigns() {
  local -a subcmds
  subcmds=(
    'list:List campaigns'
    'create:Create a new campaign'
    'show:Show campaign details'
    'send:Send a campaign'
    'cancel:Cancel a scheduled campaign'
    'stats:Show campaign statistics'
  )

  _arguments -C \
    '1:subcommand:->subcmd' \
    '*::arg:->args'

  case $state in
    subcmd) _describe 'subcommand' subcmds ;;
    args)
      case $words[1] in
        list)
          _arguments \
            '--type[Filter by type]:type:(email whatsapp)'
          ;;
        create)
          _arguments \
            '--type[Campaign type (required)]:type:(email whatsapp)' \
            '--name[Campaign name]:name:' \
            '--template[Template ID]:template:'
          ;;
        show) _arguments '1:campaign_id:' ;;
        send)
          _arguments \
            '1:campaign_id:' \
            '--confirm[Skip confirmation prompt]'
          ;;
        cancel) _arguments '1:campaign_id:' ;;
        stats)
          _arguments \
            '1:campaign_id:' \
            '--watch[Live-refresh stats]'
          ;;
      esac
      ;;
  esac
}

_micepad_templates() {
  local -a subcmds
  subcmds=(
    'list:List templates'
  )

  _arguments -C \
    '1:subcommand:->subcmd' \
    '*::arg:->args'

  case $state in
    subcmd) _describe 'subcommand' subcmds ;;
    args)
      case $words[1] in
        list)
          _arguments \
            '--type[Filter by type]:type:(email whatsapp)'
          ;;
      esac
      ;;
  esac
}

_micepad_checkins() {
  local -a subcmds
  subcmds=(
    'stats:Show check-in statistics'
    'recent:Show recent check-in activity'
  )

  _arguments -C \
    '1:subcommand:->subcmd' \
    '*::arg:->args'

  case $state in
    subcmd) _describe 'subcommand' subcmds ;;
    args)
      case $words[1] in
        stats) _arguments '--watch[Live-refresh stats]' ;;
        recent) _arguments '--limit[Number of recent check-ins]:number:' ;;
      esac
      ;;
  esac
}

_micepad_groups() {
  local -a subcmds
  subcmds=(
    'list:List groups'
    'show:Show group details'
  )

  _arguments -C \
    '1:subcommand:->subcmd' \
    '*::arg:->args'

  case $state in
    subcmd) _describe 'subcommand' subcmds ;;
    args)
      case $words[1] in
        show) _arguments '1:name_or_id:' ;;
      esac
      ;;
  esac
}

_micepad_admin() {
  local -a subcmds
  subcmds=(
    'dashboard:Show platform-wide statistics'
    'accounts:List accounts'
    'users:List users'
    'gatherings:List all gatherings'
    'subscriptions:List active subscriptions'
  )

  _arguments -C \
    '1:subcommand:->subcmd' \
    '*::arg:->args'

  case $state in
    subcmd) _describe 'subcommand' subcmds ;;
    args)
      case $words[1] in
        accounts|users)
          _arguments \
            '--limit[Max rows]:number:' \
            '--search[Search by name or email]:query:'
          ;;
        gatherings)
          _arguments \
            '--limit[Max rows]:number:' \
            '--status[Filter by status]:status:'
          ;;
        subscriptions)
          _arguments \
            '--limit[Max rows]:number:'
          ;;
      esac
      ;;
  esac
}

_micepad "$@"
