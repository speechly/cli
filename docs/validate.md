# validate

Validate the given configuration for syntax errors

# Synopsis

The contents of the directory given as argument is sent to the
API and validated. Possible errors are printed to stdout.

```
speechly validate [<app_id>] <directory> [flags]
```

# Examples

```
speechly validate -a <app_id> .
speechly validate <app_id> /path/to/config
```

# Options

```
  -a, --app string   Application to validate the files for. Can be given as the first positional argument.
  -h, --help         help for validate
```

# See also

* [speechly](README.md)	 - Speechly CLI

