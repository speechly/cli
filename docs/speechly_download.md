## speechly download

Download the active configuration or model of the given app.

### Synopsis

Fetches the currently stored configuration or model. This command does not check for validity of the stored configuration, but downloads the latest version.

```
speechly download [<app_id>] <directory> [flags]
```

### Examples

```
speechly download <app_id> /path/to/config
speechly download -a <app_id> .
speechly download -a <app_id> . --model tflite
```

### Options

```
  -a, --app string     Which application's configuration or model to download. Can be given as the first positional argument.
  -h, --help           help for download
      --model string   Specify the machine learning framework of the model to download. Available options are: ort, tflite, coreml and all. This feature is available on Enterprise plans (https://speechly.com/pricing)
```

### SEE ALSO

* [speechly](speechly.md)	 - Speechly CLI

