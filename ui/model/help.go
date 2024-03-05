package model

var (
	helpText = `
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
           N           Fetch the next 10 messages from the queue
           }           Fetch the next 100 messages from the queue
           d           Toggle deletion mode; cueitup will delete messages
                       after reading them
           <ctrl+p     Toggle queue message count polling ON/OFF; ON by default
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
