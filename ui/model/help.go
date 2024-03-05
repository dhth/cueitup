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
           n           Fetch the next message from the queue
           N           Fetch the next 10 messages from the queue
           s           Toggle skipping mode; cueitup will consume messages,
                           but not populate its internal list, effectively
                           skipping over them

           Message Value View   
           f           Toggle focussed section between full screen and
                           regular mode
           1           Maximize message value view
           q           Minimize section, and return focus to list view
`
)
