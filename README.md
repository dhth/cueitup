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

TODO
---

- [ ] Add ability to only save records with a chosen set of keys

Acknowledgements
---

`cueitup` is built using the awesome TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
