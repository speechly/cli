# validate

Validate the given configuration for syntax errors

### Usage

```
speechly validate [flags]
```

The contents of the directory given as argument is sent to the API and validated. Possible errors are printed to stdout.

### Flags

* `--app` `-a` _(string)_ - Application to validate the files for. Can be given as the first positional argument.
* `--help` `-h` _(bool)_ - help for validate

### Examples

```
speechly validate <app_id> .
speechly validate --app <app_id> /path/to/config
```
