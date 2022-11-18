# download

Download the active configuration or model bundle of the given app.

### Usage

```
speechly download [flags]
```

Fetches the currently stored configuration or model bundle. This command does not check for validity of the stored configuration, but downloads the latest version.

### Flags

* `--app` `-a` _(string)_ - Application which configuration or model bundle to download. Can be given as the first positional argument.
* `--help` `-h` _(bool)_ - help for download
* `--model` `-m` _(string)_ - Model bundle machine learning framework. Available options are: ort, tflite, coreml and all. This feature is available on Enterprise plans (https://speechly.com/pricing)

### Examples

```
speechly download <app_id> /path/to/config
speechly download --app <app_id> . --model tflite
```
