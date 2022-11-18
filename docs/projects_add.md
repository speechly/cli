# projects add

Add access to a pre-existing project

### Usage

```
speechly projects add [flags]
```

### Flags

* `--apikey` _(string)_ - API token, created in Speechly Dashboard. Can also be given as the sole positional argument.
* `--help` `-h` _(bool)_ - help for add
* `--host` _(string)_ - API address
* `--name` _(string)_ - An unique name for the project. If not given the project name configured in Speechly Dashboard will be used.
* `--skip-online-validation` _(bool)_ - Skips validating the API token against the host.

### Examples

```
speechly projects add <api_token>
speechly projects add --apikey <api_token>
```
