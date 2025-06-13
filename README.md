<p align="center">
  <h1 align="center">cueitup</h1>
  <p align="center">
    <a href="https://github.com/dhth/cueitup/actions/workflows/main.yml"><img alt="GitHub release" src="https://img.shields.io/github/actions/workflow/status/dhth/cueitup/main.yml?style=flat-square"></a>
    <a href="https://github.com/dhth/cueitup/releases/latest"><img alt="Latest release" src="https://img.shields.io/github/release/dhth/cueitup.svg?style=flat-square"></a>
    <a href="https://github.com/dhth/cueitup/releases"><img alt="Commits since latest release" src="https://img.shields.io/github/commits-since/dhth/cueitup/latest?style=flat-square"></a>
  </p>
</p>

`cueitup` lets you inspect messages in an AWS SQS queue in a simple and
deliberate manner. You can pull one or more messages on demand, peruse through
them in a list, and, if needed, persist them to your local filesystem. `cueitup`
offers both a terminal UI and a web interface.

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

Or get the binaries directly from a
[release](https://github.com/dhth/cueitup/releases). Read more about verifying
the authenticity of released artifacts [here](#-verifying-release-artifacts).

🛠️ Configuration
---

Create a YAML configuration file that looks like the following. The location of
this file depends on your operating system, and can be determined by running
`cueitup
-h`.

```yaml
profiles:
    # a name for a profile; you refer to it when running cueitup
  - name: profile-a #

    # the SQS queue URL
    queue_url: https://sqs.eu-central-1.amazonaws.com/000000000000/queue-a

    # use this to leverage a profile contained in the shared AWS config and credentials files
    # https://docs.aws.amazon.com/sdkref/latest/guide/file-format.html
    aws_config_source: profile:local-profile

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

    # cueitup will display this key value pair as "context" in its list
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
  -c, --config-path string   location of cueitup's config file (default "/Users/user/Library/Application Support/cueitup/cueitup.yml")
  -d, --debug                whether to only display config picked up by cueitup
  -D, --delete-messages      whether to start the TUI with the setting "delete messages" ON (default true)
  -h, --help                 help for tui
  -P, --persist-messages     whether to start the TUI with the setting "persist messages" ON
  -M, --show-message-count   whether to start the TUI with the setting "show message count" ON (default true)
  -S, --skip-messages        whether to start the TUI with the setting "skip messages" ON
```

<video src="https://github.com/user-attachments/assets/738a5797-89f8-4717-9639-3a0fe72715d8"></video>

```text
$ cueitup serve --help

open cueitup's web interface

Usage:
  cueitup serve <PROFILE> [flags]

Flags:
  -c, --config-path string   location of cueitup's config file (default "/Users/user/Library/Application Support/cueitup/cueitup.yml")
  -d, --debug                whether to only display config picked up by cueitup
  -D, --delete-messages      whether to start the web interface with the setting "delete messages" ON (default true)
  -h, --help                 help for serve
  -o, --open                 whether to open web interface in browser automatically
  -S, --select-on-hover      whether to start the web interface with the setting "select on hover" ON
  -M, --show-message-count   whether to start the web interface with the setting "show message count" ON (default true)
```

<video src="https://github.com/user-attachments/assets/e11e2d02-c5a4-4379-b6f2-ee498094e122"></video>

Various ways to display JSON messages
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

You can configure `cueitup` to show a specific key value pair as "context" in
its UI. You do do via the `context_key` configuration property.

```yaml
- name: sample-profile
  queue_url: ...
  aws_config_source: ...
  format: json
  context_key: transactionId
```
![](https://github.com/user-attachments/assets/8da47c9a-d883-40e0-b75c-c15810e7807c)

If you want to only see the nested object under `browserInfo`, you'd configure a
profile like this:

```yaml
- name: sample-profile
  queue_url: ...
  aws_config_source: ...
  format: json
  subset_key: browserInfo
  context_key: platform
```

![](https://github.com/user-attachments/assets/04027abf-25ea-4bd3-8b4d-c192e4e5bacc)

`cueitup` can also work with stringified JSON.

```yaml
- name: sample-profile
  queue_url: ...
  aws_config_source: ...
  format: json
  subset_key: metadata
  context_key: aggregateId
```

![](https://github.com/user-attachments/assets/1f2d93f7-5d91-40ea-82e6-9eed28ac99c6)

TUI Keyboard shortcuts
---

### General

| Keymap    | Description                      |
|-----------|----------------------------------|
| `<tab>`   | Switch focus to next section     |
| `<s-tab>` | Switch focus to previous section |
| `?`       | Show help view                   |
| `q`       | Go back or quit                  |

### Message List Pane

| Keymap     | Description                                                                  |
|------------|------------------------------------------------------------------------------|
| `h/<Up>`   | Move cursor up                                                               |
| `k/<Down>` | Move cursor down                                                             |
| `n`        | Fetch the next message from the queue                                        |
| `N`        | Fetch up to 10 more messages from the queue                                  |
| `}`        | Fetch up to 100 more messages from the queue                                 |
| `d`        | Toggle deletion mode; cueitup will delete messages after reading them        |
| `M`        | Toggle polling for message count in queue                                    |
| `p`        | Toggle persist mode (messages will be saved to a specific location)          |
| `s`        | Toggle skipping mode (consume messages without populating the internal list) |

### Message Value Pane

| Keymap   | Description                                     |
|----------|-------------------------------------------------|
| `[`, `h` | Show details for the previous entry in the list |
| `]`, `l` | Show details for the next entry in the list     |

🔐 Verifying release artifacts
---

In case you get the `cueitup` binary directly from a
[release](https://github.com/dhth/cueitup/releases), you may want to verify its
authenticity. Checksums are applied to all released artifacts, and the resulting
checksum file is signed using
[cosign](https://docs.sigstore.dev/cosign/installation/).

Steps to verify (replace `A.B.C` in the commands listed below with the version
you want):

1. Download the following files from the release:

    - cueitup_A.B.C_checksums.txt
    - cueitup_A.B.C_checksums.txt.pem
    - cueitup_A.B.C_checksums.txt.sig

2. Verify the signature:

   ```shell
   cosign verify-blob cueitup_A.B.C_checksums.txt \
       --certificate cueitup_A.B.C_checksums.txt.pem \
       --signature cueitup_A.B.C_checksums.txt.sig \
       --certificate-identity-regexp 'https://github\.com/dhth/cueitup/\.github/workflows/.+' \
       --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
   ```

3. Download the compressed archive you want, and validate its checksum:

   ```shell
   curl -sSLO https://github.com/dhth/cueitup/releases/download/vA.B.C/cueitup_A.B.C_linux_amd64.tar.gz
   sha256sum --ignore-missing -c cueitup_A.B.C_checksums.txt
   ```

3. If checksum validation goes through, uncompress the archive:

   ```shell
   tar -xzf cueitup_A.B.C_linux_amd64.tar.gz
   ./cueitup
   # profit!
   ```

Acknowledgements
---

`cueitup` is built using the TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
