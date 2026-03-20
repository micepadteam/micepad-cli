# Fish completion for micepad CLI
# Add to ~/.config/fish/completions/micepad.fish (symlink or copy)

# Disable file completions by default
complete -c micepad -f

# Helper: check if we're completing at a specific subcommand level
function __micepad_no_subcommand
    set -l cmd (commandline -opc)
    test (count $cmd) -eq 1
end

function __micepad_using_command
    set -l cmd (commandline -opc)
    test (count $cmd) -ge 2; and test "$cmd[2]" = "$argv[1]"
end

function __micepad_using_subcommand
    set -l cmd (commandline -opc)
    test (count $cmd) -ge 3; and test "$cmd[2]" = "$argv[1]"; and test "$cmd[3]" = "$argv[2]"
end

function __micepad_needs_subcmd
    set -l cmd (commandline -opc)
    test (count $cmd) -eq 2; and test "$cmd[2]" = "$argv[1]"
end

# Root commands
complete -c micepad -n __micepad_no_subcommand -a login -d 'Login to your account'
complete -c micepad -n __micepad_no_subcommand -a logout -d 'Logout of your account'
complete -c micepad -n __micepad_no_subcommand -a whoami -d 'Show current user and account'
complete -c micepad -n __micepad_no_subcommand -a events -d 'Manage events'
complete -c micepad -n __micepad_no_subcommand -a pax -d 'Manage participants'
complete -c micepad -n __micepad_no_subcommand -a campaigns -d 'Manage campaigns'
complete -c micepad -n __micepad_no_subcommand -a templates -d 'Manage templates'
complete -c micepad -n __micepad_no_subcommand -a checkins -d 'Check-in operations'
complete -c micepad -n __micepad_no_subcommand -a groups -d 'Manage groups'
complete -c micepad -n __micepad_no_subcommand -a plans -d 'Manage one-time plans for this event'
complete -c micepad -n __micepad_no_subcommand -a admin -d 'Admin operations (super admin only)'

# --- events ---
complete -c micepad -n '__micepad_needs_subcmd events' -a list -d 'List your events'
complete -c micepad -n '__micepad_needs_subcmd events' -a use -d 'Set default event context'
complete -c micepad -n '__micepad_needs_subcmd events' -a current -d 'Show current event details'
complete -c micepad -n '__micepad_needs_subcmd events' -a stats -d 'Show event dashboard statistics'

# --- pax ---
complete -c micepad -n '__micepad_needs_subcmd pax' -a list -d 'List participants'
complete -c micepad -n '__micepad_needs_subcmd pax' -a show -d 'Show participant details'
complete -c micepad -n '__micepad_needs_subcmd pax' -a add -d 'Add a new participant'
complete -c micepad -n '__micepad_needs_subcmd pax' -a checkin -d 'Check in a participant'
complete -c micepad -n '__micepad_needs_subcmd pax' -a checkout -d 'Check out a participant'
complete -c micepad -n '__micepad_needs_subcmd pax' -a count -d 'Count participants by status'
complete -c micepad -n '__micepad_needs_subcmd pax' -a export -d 'Export participants to CSV'

# pax list options
complete -c micepad -n '__micepad_using_subcommand pax list' -l status -d 'Filter by RSVP status' -xa 'confirmed unconfirmed declined waitlisted pending_approval'
complete -c micepad -n '__micepad_using_subcommand pax list' -l checkin -d 'Filter by check-in status' -xa 'checked_in not_checked_in checked_out'
complete -c micepad -n '__micepad_using_subcommand pax list' -l group -d 'Filter by group name'
complete -c micepad -n '__micepad_using_subcommand pax list' -l search -d 'Search by name or email'
complete -c micepad -n '__micepad_using_subcommand pax list' -l limit -d 'Max rows to show'

# pax add options
complete -c micepad -n '__micepad_using_subcommand pax add' -l email -d 'Email address (required)'
complete -c micepad -n '__micepad_using_subcommand pax add' -l first_name -d 'First name (required)'
complete -c micepad -n '__micepad_using_subcommand pax add' -l last_name -d 'Last name'
complete -c micepad -n '__micepad_using_subcommand pax add' -l company -d 'Company name'
complete -c micepad -n '__micepad_using_subcommand pax add' -l job_title -d 'Job title'

# pax count options
complete -c micepad -n '__micepad_using_subcommand pax count' -l by -d 'Group by' -xa 'rsvp checkin group'

# pax export options
complete -c micepad -n '__micepad_using_subcommand pax export' -l status -d 'Filter by RSVP status' -xa 'confirmed unconfirmed declined waitlisted pending_approval'
complete -c micepad -n '__micepad_using_subcommand pax export' -l output -d 'Output file path'

# --- campaigns ---
complete -c micepad -n '__micepad_needs_subcmd campaigns' -a list -d 'List campaigns'
complete -c micepad -n '__micepad_needs_subcmd campaigns' -a create -d 'Create a new campaign'
complete -c micepad -n '__micepad_needs_subcmd campaigns' -a show -d 'Show campaign details'
complete -c micepad -n '__micepad_needs_subcmd campaigns' -a send -d 'Send a campaign'
complete -c micepad -n '__micepad_needs_subcmd campaigns' -a cancel -d 'Cancel a scheduled campaign'
complete -c micepad -n '__micepad_needs_subcmd campaigns' -a stats -d 'Show campaign statistics'

# campaigns list options
complete -c micepad -n '__micepad_using_subcommand campaigns list' -l type -d 'Filter by type' -xa 'email whatsapp'

# campaigns create options
complete -c micepad -n '__micepad_using_subcommand campaigns create' -l type -d 'Campaign type (required)' -xa 'email whatsapp'
complete -c micepad -n '__micepad_using_subcommand campaigns create' -l name -d 'Campaign name'
complete -c micepad -n '__micepad_using_subcommand campaigns create' -l template -d 'Template ID'

# campaigns send options
complete -c micepad -n '__micepad_using_subcommand campaigns send' -l confirm -d 'Skip confirmation prompt'

# campaigns stats options
complete -c micepad -n '__micepad_using_subcommand campaigns stats' -l watch -d 'Live-refresh stats'

# --- templates ---
complete -c micepad -n '__micepad_needs_subcmd templates' -a list -d 'List templates'

# templates list options
complete -c micepad -n '__micepad_using_subcommand templates list' -l type -d 'Filter by type' -xa 'email whatsapp'

# --- checkins ---
complete -c micepad -n '__micepad_needs_subcmd checkins' -a stats -d 'Show check-in statistics'
complete -c micepad -n '__micepad_needs_subcmd checkins' -a recent -d 'Show recent check-in activity'

# checkins stats options
complete -c micepad -n '__micepad_using_subcommand checkins stats' -l watch -d 'Live-refresh stats'

# checkins recent options
complete -c micepad -n '__micepad_using_subcommand checkins recent' -l limit -d 'Number of recent check-ins'

# --- groups ---
complete -c micepad -n '__micepad_needs_subcmd groups' -a list -d 'List groups'
complete -c micepad -n '__micepad_needs_subcmd groups' -a show -d 'Show group details'

# --- plans ---
complete -c micepad -n '__micepad_needs_subcmd plans' -a current -d 'Show current plan for this event'
complete -c micepad -n '__micepad_needs_subcmd plans' -a list -d 'List available one-time plans'
complete -c micepad -n '__micepad_needs_subcmd plans' -a subscribe -d 'Subscribe this event to a one-time plan'
complete -c micepad -n '__micepad_needs_subcmd plans' -a usage -d 'Show plan usage for this event'
complete -c micepad -n '__micepad_needs_subcmd plans' -a add_ons -d 'List available add-ons'

# plans list options
complete -c micepad -n '__micepad_using_subcommand plans list' -l group -d 'Filter by group' -xa 'free registration onsite bundle starter pro'

# plans subscribe options
complete -c micepad -n '__micepad_using_subcommand plans subscribe' -l confirm -d 'Skip confirmation prompt'

# --- admin ---
complete -c micepad -n '__micepad_needs_subcmd admin' -a dashboard -d 'Show platform-wide statistics'
complete -c micepad -n '__micepad_needs_subcmd admin' -a accounts -d 'List accounts'
complete -c micepad -n '__micepad_needs_subcmd admin' -a users -d 'List users'
complete -c micepad -n '__micepad_needs_subcmd admin' -a gatherings -d 'List all gatherings'
complete -c micepad -n '__micepad_needs_subcmd admin' -a subscriptions -d 'List active subscriptions'

# admin accounts options
complete -c micepad -n '__micepad_using_subcommand admin accounts' -l limit -d 'Max rows'
complete -c micepad -n '__micepad_using_subcommand admin accounts' -l search -d 'Search by name or email'

# admin users options
complete -c micepad -n '__micepad_using_subcommand admin users' -l limit -d 'Max rows'
complete -c micepad -n '__micepad_using_subcommand admin users' -l search -d 'Search by name or email'

# admin gatherings options
complete -c micepad -n '__micepad_using_subcommand admin gatherings' -l limit -d 'Max rows'
complete -c micepad -n '__micepad_using_subcommand admin gatherings' -l status -d 'Filter by status'

# admin subscriptions options
complete -c micepad -n '__micepad_using_subcommand admin subscriptions' -l limit -d 'Max rows'
