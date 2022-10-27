# transcribe

Transcribe the given file(s) using on-device or cloud transcription

```
speechly transcribe <input_file> [flags]
```

# Examples

```
speechly transcribe <input_file> --model /path/to/model/bundle
speechly transcribe <input_file> --app <app_id>
```

# Options

```
  -a, --app string     Application ID to use for cloud transcription
  -h, --help           help for transcribe
  -m, --model string   Model bundle file. This feature is available on Enterprise plans (https://speechly.com/pricing)
      --streaming      Use the Streaming API instead of the Batch API.
```

# See also

* [speechly](README.md)	 - Speechly CLI

