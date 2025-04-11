<p align="center">
  <h1 align="center">cueitup</h1>
  <p align="center">
    <a href="https://github.com/dhth/cueitup/actions/workflows/build.yml"><img alt="GitHub release" src="https://img.shields.io/github/actions/workflow/status/dhth/cueitup/build.yml?style=flat-square"></a>
    <a href="https://github.com/dhth/cueitup/releases/latest"><img alt="Latest release" src="https://img.shields.io/github/release/dhth/cueitup.svg?style=flat-square"></a>
    <a href="https://github.com/dhth/cueitup/releases"><img alt="Commits since latest release" src="https://img.shields.io/github/commits-since/dhth/cueitup/latest?style=flat-square"></a>
  </p>
</p>

`cueitup` lets you inspect messages in an AWS SQS queue in a simple and
deliberate manner. It was built to simplify the process of investigating the
contents of messages being pushed to an SNS topic. You can pull one or more
messages on demand, peruse through them in a list, and, if needed, persist them
to your local filesystem.

![tui](https://github.com/user-attachments/assets/c1727615-46c1-483c-9e08-ffdb9c9e1fb6)
![web interface](https://github.com/user-attachments/assets/848592e1-0e24-4eb0-a3c0-db12b4d34116)

Install
---

**homebrew**:

```sh
brew install dhth/tap/cueitup
```

**go**:

```sh
go install github.com/dhth/cueitup@latest
```

🛠️ Configuration
---

Create a configuration file that looks like the following. By default,
`outtasync` will look for this file at `~/.config/outtasync.yml`.

```yaml
profiles:
    # a name for a profile; you refer to it when running cueitup
  - name: profile-a #

    # the SQS queue URL
    queue_url: https://sqs.eu-central-1.amazonaws.com/000000000000/queue-a

    # use this to leverage a profile contained in the shared AWS config and credentials files
    # https://docs.aws.amazon.com/sdkref/latest/guide/file-format.html
    aws_config_source: "profile:local-profile"

    # the format of the message body; possible values: [json, none]
    format: json

  - name: profile-b
    queue_url: https://sqs.eu-central-1.amazonaws.com/000000000000/queue-b
    aws_config_source: env
    format: none

  - name: profile-c
    queue_url: https://sqs.eu-central-1.amazonaws.com/000000000000/queue-c
    aws_config_source: env
    format: json

    # to only show the contents of a nested object
    subset_key: Message

    # cueitup will display the value of this key in its list
    context_key: aggregateId
```

⚡️ Usage
---

`cueitup` can display messages via two interfaces: a TUI or a webpage 

```text
$ cueitup tui --help

open cueitup's TUI

Usage:
  cueitup tui <PROFILE> [flags]

Flags:
  -c, --config-path string   location of cueitup's config file (default "/Users/dhruvthakur/Library/Application Support/cueitup/cueitup.yml")
  -d, --debug                whether to only display config picked up by cueitup
  -D, --delete-messages      whether to start the TUI with the setting "delete messages" ON (default true)
  -h, --help                 help for tui
  -P, --persist-messages     whether to start the TUI with the setting "persist messages" ON
  -M, --show-message-count   whether to start the TUI with the setting "show message count" ON (default true)
  -S, --skip-messages        whether to start the TUI with the setting "skip messages" ON
```

<video src="https://github.com/user-attachments/assets/738a5797-89f8-4717-9639-3a0fe72715d8"></video>

```text
open cueitup's web interface

Usage:
  cueitup serve <PROFILE> [flags]

Flags:
  -c, --config-path string   location of cueitup's config file (default "/Users/dhruvthakur/Library/Application Support/cueitup/cueitup.yml")
  -d, --debug                whether to only display config picked up by cueitup
  -D, --delete-messages      whether to start the web interface with the setting "delete messages" ON (default true)
  -h, --help                 help for serve
  -o, --open                 whether to open web interface in browser automatically
  -S, --select-on-hover      whether to start the web interface with the setting "select on hover" ON
  -M, --show-message-count   whether to start the web interface with the setting "show message count" ON (default true)
```

<video src="https://github.com/user-attachments/assets/e11e2d02-c5a4-4379-b6f2-ee498094e122"></video>

Using subset and context keys
---

Say the messages in your SQS queue look like this.

```json
{
  "browserInfo": {
    "browserName": "Firefox",
    "browserVersion": 118,
    "deviceType": "Desktop",
    "platform": "Linux"
  },
  "isBot": true,
  "metadata": "{\"aggregateId\":\"00000000-0000-0000-0000-000000012363\",\"sequenceNr\":347}",
  "sessionId": "987e6543-b21a-34c5-d678-123456789abc",
  "transactionId": "123e4567-e89b-12d3-a456-426614174000"
}
```

If you want to only see the nested object under `browserInfo`, you'd configure a
profile like this:

```yaml
- name: sample-profile
  queue_url: ...
  aws_config_source: ...
  format: json
  subset_key: browserInfo
```

Now, if you want `cueitup` to display the value of the key `platform` under
`browserInfo`, you'd configure `context_key` as well:

```yaml
- name: sample-profile
  queue_url: ...
  aws_config_source: ...
  format: json
  subset_key: browserInfo
  context_key: platform
```

TUI Reference Manual
---

```
cueitup has 3 views:
- Message List View
- Message Value View
- Help View (this one)

Keyboard Shortcuts

General

    <tab>                          Switch focus to next section
    <s-tab>                        Switch focus to previous section
    ?                              Show help view
    q                              Go back or quit

Message List View

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
                                       messages/<topic-name>/<timestamp>-<message-id>.md
    s                              Toggle skipping mode; cueitup will consume messages,
                                       but not populate its internal list, effectively
                                       skipping over them

Message Value View

    [,h                            Show details for the previous entry in the list
    ],l                            Show details for the next entry in the list
```

Acknowledgements
---

`cueitup` is built using the TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
