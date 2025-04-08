# cueitup's web interface

Local Development
---

Prequisites

- gleam

```sh
# start local development server
# from project root
go run . serve <PROFILE>

cd client
# replace window.location() in src/effects.gleam with http://127.0.0.1:<PORT>
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
