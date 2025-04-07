package ui

import "fmt"

var HelpText = fmt.Sprintf(`
  %s
%s
  %s

  %s
%s
  %s
%s
  %s
%s
`,
	helpHeaderStyle.Render("cueitup Reference Manual"),
	helpSectionStyle.Render(`
  (scroll line by line with j/k/arrow keys or by half a page with <c-d>/<c-u>)

  cueitup has 3 views:
  - Message List View
  - Message Value View
  - Help View (this one)
`),
	helpHeaderStyle.Render("Keyboard Shortcuts"),
	helpHeaderStyle.Render("General"),
	helpSectionStyle.Render(`
      <tab>                          Switch focus to next section
      <s-tab>                        Switch focus to previous section
      1                              Maximize message value view
      ?                              Show help view
`),
	helpHeaderStyle.Render("Message List View"),
	helpSectionStyle.Render(`
      h/<Up>                         Move cursor up
      k/<Down>                       Move cursor down
      n                              Fetch the next message from the queue
      N                              Fetch up to 10 more messages from the queue
      }                              Fetch up to 100 more messages from the queue
      d                              Toggle deletion mode; cueitup will delete messages
                                         after reading them
      <ctrl+s>                       Toggle contextual search prompt
      <ctrl+f>                       Toggle contextual filtering ON/OFF
      <ctrl+p>                       Toggle queue message count polling ON/OFF; ON by default
      p                              Toggle persist mode (cueitup will start persisting
                                         messages, at the location
                                         messages/<topic-name>/<timestamp-when-cueitup-started>/<unix-epoch>-<message-id>.md
      s                              Toggle skipping mode; cueitup will consume messages,
                                         but not populate its internal list, effectively
                                         skipping over them
`),
	helpHeaderStyle.Render("Message Value View   "),
	helpSectionStyle.Render(`
      q                              Minimize section, and return focus to list view
      [,h                            Show details for the previous entry in the list
      ],l                            Show details for the next entry in the list
`),
)
