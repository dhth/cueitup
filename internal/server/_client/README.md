# cueitup's web interface

Local Development
---

Prerequisites

- gleam

```sh
# start local development server
# from project root
go run . serve <PROFILE>

cd client
# set dev = True in src/effects.gleam
gleam run -m lustre/dev start
```

Before committing code
---

```sh
# ensure local changes in _client/src/effects.gleam are
# reverted
cd client

# compile app to js code
gleam run -m lustre/dev build app
```
