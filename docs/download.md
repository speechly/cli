# download

Download the active configuration or model bundle of the given app.

## Synopsis

Fetches the currently stored configuration or model bundle. This command does not check for validity of the stored configuration, but downloads the latest version.

```
speechly download [<app_id>] <directory> [flags]
```

## Examples

```
speechly download <app_id> /path/to/config
speechly download -a <app_id> .
speechly download -a <app_id> . --model tflite
```

## Options

```
  -a, --app string     Application which configuration or model bundle to download. Can be given as the first positional argument.
  -h, --help           help for download
  -m, --model string   Model bundle machine learning framework. Available options are: ort, tflite, coreml and all. This feature is available on Enterprise plans (https://speechly.com/pricing)
```

## See also

* [speechly](README.md)	 - Speechly CLI

