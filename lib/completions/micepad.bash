# Bash completion for micepad CLI
# Add to ~/.bashrc: source /path/to/lib/completions/micepad.bash

_micepad() {
  local cur prev words cword
  _init_completion || return

  local root_commands="login logout whoami events pax campaigns templates checkins groups plans admin"

  # Subcommands per group
  local events_cmds="list use current stats"
  local pax_cmds="list show add checkin checkout count export"
  local campaigns_cmds="list create show send cancel stats"
  local templates_cmds="list"
  local checkins_cmds="stats recent"
  local groups_cmds="list show"
  local plans_cmds="current list subscribe usage add_ons"
  local admin_cmds="dashboard accounts users gatherings subscriptions"

  case "${words[1]}" in
    events)
      if [[ $cword -eq 2 ]]; then
        COMPREPLY=($(compgen -W "$events_cmds" -- "$cur"))
      fi
      ;;
    pax)
      if [[ $cword -eq 2 ]]; then
        COMPREPLY=($(compgen -W "$pax_cmds" -- "$cur"))
      else
        case "${words[2]}" in
          list)
            COMPREPLY=($(compgen -W "--status --checkin --group --search --limit" -- "$cur"))
            case "$prev" in
              --status) COMPREPLY=($(compgen -W "confirmed unconfirmed declined waitlisted pending_approval" -- "$cur")) ;;
              --checkin) COMPREPLY=($(compgen -W "checked_in not_checked_in checked_out" -- "$cur")) ;;
            esac
            ;;
          add)
            COMPREPLY=($(compgen -W "--email --first_name --last_name --company --job_title" -- "$cur"))
            ;;
          count)
            COMPREPLY=($(compgen -W "--by" -- "$cur"))
            [[ "$prev" == "--by" ]] && COMPREPLY=($(compgen -W "rsvp checkin group" -- "$cur"))
            ;;
          export)
            COMPREPLY=($(compgen -W "--status --output" -- "$cur"))
            case "$prev" in
              --status) COMPREPLY=($(compgen -W "confirmed unconfirmed declined waitlisted pending_approval" -- "$cur")) ;;
              --output) _filedir csv ;;
            esac
            ;;
        esac
      fi
      ;;
    campaigns)
      if [[ $cword -eq 2 ]]; then
        COMPREPLY=($(compgen -W "$campaigns_cmds" -- "$cur"))
      else
        case "${words[2]}" in
          list)
            COMPREPLY=($(compgen -W "--type" -- "$cur"))
            [[ "$prev" == "--type" ]] && COMPREPLY=($(compgen -W "email whatsapp" -- "$cur"))
            ;;
          create)
            COMPREPLY=($(compgen -W "--type --name --template" -- "$cur"))
            [[ "$prev" == "--type" ]] && COMPREPLY=($(compgen -W "email whatsapp" -- "$cur"))
            ;;
          send)
            COMPREPLY=($(compgen -W "--confirm" -- "$cur"))
            ;;
          stats)
            COMPREPLY=($(compgen -W "--watch" -- "$cur"))
            ;;
        esac
      fi
      ;;
    templates)
      if [[ $cword -eq 2 ]]; then
        COMPREPLY=($(compgen -W "$templates_cmds" -- "$cur"))
      else
        case "${words[2]}" in
          list)
            COMPREPLY=($(compgen -W "--type" -- "$cur"))
            [[ "$prev" == "--type" ]] && COMPREPLY=($(compgen -W "email whatsapp" -- "$cur"))
            ;;
        esac
      fi
      ;;
    checkins)
      if [[ $cword -eq 2 ]]; then
        COMPREPLY=($(compgen -W "$checkins_cmds" -- "$cur"))
      else
        case "${words[2]}" in
          stats) COMPREPLY=($(compgen -W "--watch" -- "$cur")) ;;
          recent) COMPREPLY=($(compgen -W "--limit" -- "$cur")) ;;
        esac
      fi
      ;;
    groups)
      if [[ $cword -eq 2 ]]; then
        COMPREPLY=($(compgen -W "$groups_cmds" -- "$cur"))
      fi
      ;;
    plans)
      if [[ $cword -eq 2 ]]; then
        COMPREPLY=($(compgen -W "$plans_cmds" -- "$cur"))
      else
        case "${words[2]}" in
          list)
            COMPREPLY=($(compgen -W "--group" -- "$cur"))
            [[ "$prev" == "--group" ]] && COMPREPLY=($(compgen -W "free registration onsite bundle starter pro" -- "$cur"))
            ;;
          subscribe)
            COMPREPLY=($(compgen -W "--confirm" -- "$cur"))
            ;;
        esac
      fi
      ;;
    admin)
      if [[ $cword -eq 2 ]]; then
        COMPREPLY=($(compgen -W "$admin_cmds" -- "$cur"))
      else
        case "${words[2]}" in
          accounts|users)
            COMPREPLY=($(compgen -W "--limit --search" -- "$cur"))
            ;;
          gatherings)
            COMPREPLY=($(compgen -W "--limit --status" -- "$cur"))
            ;;
          subscriptions)
            COMPREPLY=($(compgen -W "--limit" -- "$cur"))
            ;;
        esac
      fi
      ;;
    *)
      if [[ $cword -eq 1 ]]; then
        COMPREPLY=($(compgen -W "$root_commands" -- "$cur"))
      fi
      ;;
  esac
}

complete -F _micepad micepad
