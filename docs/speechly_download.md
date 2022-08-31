## speechly download

Download the active configuration of the given app.

### Synopsis

Fetches the currently stored configuration from the API. This command
does not check for validity of the stored configuration, but downloads the latest
version.

```
speechly download [<app_id>] <directory> [flags]
```

### Examples

```
speechly download <app_id> /path/to/config
speechly download -a <app_id> .
```

### Options

```
  -a, --app string   Which application's configuration to download. Can be given as the first positional argument.
  -h, --help         help for download
```

### SEE ALSO

* [speechly](speechly.md)	 - Speechly CLI

