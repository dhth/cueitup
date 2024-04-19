# cueitup

✨ Overview
---

`cueitup` lets you inspect messages in an AWS SQS queue in a simple and
deliberate manner. It was built to simplify the process of investigating the
contents of messages being pushed to an SNS topic. You can pull one or more
messages on demand, peruse through them in a list, and, if needed, persist them
to your local filesystem.

<p align="center">
  <img src="./assets/cueitup.gif?raw=true" alt="Usage" />
</p>

Demo video:

[![Demo Video](https://img.youtube.com/vi/95HsXNUL4J4/0.jpg)](https://www.youtube.com/watch?v=95HsXNUL4J4)

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

⚡️ Usage
---

### Consuming JSON messages

#### Basic usage

```bash

cueitup \
    -aws-profile="<PROFILE>" \
    -aws-region="<REGION>" \
    -queue-url="https://sqs.eu-central-1.amazonaws.com/<ABC>/<XYX>" \
    -msg-format=json
```

#### Viewing a subset of the full payload

To only view the nested object with the key 'Message' in the JSON
payload below, use 👇

```json
{
  "Type": "Notification",
  "MessageId": "f7bbec51-1cd1-4630-8eb3-7b124de6d6f4",
  "TopicArn": "arn:aws:sns:eu-central-1:123:queue-name",
  "Message": {
    "companyId": "af8e74b2-82db-4349-b861-c1d9d1a3033f",
    "resourceId": "611a709e-2b96-41e3-9274-8bbd4e191334",
    "aggregateId": "93422d4d-90ec-4a20-a794-3f835d7605cf",
    "sequenceNr": 59,
    "dateTime": "b5692ca5-e060-4318-8a40-e2b806a4018a",
    "type": "com.some.kind.of.event",
    "version": 1
  },
  "Timestamp": "5f0244d6-e640-43f4-86d0-5b9aa639c7df",
  "SignatureVersion": "1",
  "Signature": "XYZ",
  "SigningCertURL": "https://sns.eu-central-1.amazonaws.com/SimpleNotificationService-ABC",
  "UnsubscribeURL": "https://sns.eu-central-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:eu-central-1:XYZ"
}
```

```bash
cueitup \
    -aws-profile="<PROFILE>" \
    -aws-region="<REGION>" \
    -queue-url="https://sqs.eu-central-1.amazonaws.com/<ABC>/<XYX>" \
    -msg-format='json' \
    -subset-key='Message'
```

#### Adding Context to Your List

You can provide a key, whose value will be shown as for context in the list.

```bash
cueitup \
    -aws-profile="<PROFILE>" \
    -aws-region="<REGION>" \
    -queue-url="https://sqs.eu-central-1.amazonaws.com/<ABC>/<XYX>" \
    -msg-format='json' \
    -subset-key='Message' \
    -context-key='resourceId'

```

Reference Manual
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

List View

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

Message Value View

   f                              Toggle focussed section between full screen and
                                      regular mode
   1                              Maximize message value view
   q                              Minimize section, and return focus to list view
   [                              Show details for the previous entry in the list
   ]                              Show details for the next entry in the list
```

TODO
---

- [ ] Add ability to only save records with a chosen set of keys

Acknowledgements
---

`cueitup` is built using the awesome TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
