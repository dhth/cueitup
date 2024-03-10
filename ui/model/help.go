package model

var (
	HelpText = `
TUI Reference Manual

cueitup has two sections:
- Message List View
- Message Value View

Keyboard Shortcuts:

General
    <tab>       Switch focus to next section
    <s-tab>     Switch focus to previous section

List View
    h/<Up>      Move cursor up
    k/<Down>    Move cursor down
    n           Fetch the next message from the queue
    N           Fetch up to 10 more messages from the queue
    }           Fetch up to 100 more messages from the queue
    d           Toggle deletion mode; cueitup will delete messages
               after reading them
    <ctrl+s>    Toggle contextual search prompt
    <ctrl+f>    Toggle contextual filtering ON/OFF
    <ctrl+p>    Toggle queue message count polling ON/OFF; ON by default
    p           Toggle persist mode (cueitup will start persisting
                   messages, at the location
                   messages/<topic-name>/<timestamp-when-cueitup-started>/<unix-epoch>-<message-id>.md
    s           Toggle skipping mode; cueitup will consume messages,
                   but not populate its internal list, effectively
                   skipping over them

Message Value View   
    f           Toggle focussed section between full screen and
                   regular mode
    1           Maximize message value view
    q           Minimize section, and return focus to list view
    [           Show details for the previous entry in the list
    ]           Show details for the next entry in the list
`
)
