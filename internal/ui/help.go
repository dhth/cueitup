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
      ?                              Show help view
      q                              Go back or quit
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
      M                              Toggle polling for message count in queue
      p                              Toggle persist mode (cueitup will start persisting
                                         messages, at the location
                                         messages/<topic-name>/<timestamp>-<message-id>.(json|txt)
      s                              Toggle skipping mode; cueitup will consume messages,
                                         but not populate its internal list, effectively
                                         skipping over them
`),
	helpHeaderStyle.Render("Message Value View   "),
	helpSectionStyle.Render(`
      [,h                            Show details for the previous entry in the list
      ],l                            Show details for the next entry in the list
`),
)
