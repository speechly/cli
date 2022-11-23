# transcribe

Transcribe the given file(s) using on-device or cloud transcription

### Usage

```
speechly transcribe [flags]
```

To transcribe multiple files, create a JSON Lines file with each audio on their own line using the format `{"audio":"/path/to/file"}`.

### Flags

* `--app` `-a` _(string)_ - Application ID to use for cloud transcription
* `--help` `-h` _(bool)_ - help for transcribe
* `--model` `-m` _(string)_ - Model bundle file. This feature is available on Enterprise plans (https://speechly.com/pricing)
* `--streaming` _(bool)_ - Use the Streaming API instead of the Batch API.

### Examples

```
speechly transcribe file.wav --app <app_id>
speechly transcribe files.jsonl --app <app_id> > output.jsonl
speechly transcribe files.jsonl --model /path/to/model/bundle
```
