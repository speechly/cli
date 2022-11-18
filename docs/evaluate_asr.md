# evaluate asr

Evaluate the ASR accuracy of the given application model.

### Usage

```
speechly evaluate asr [flags]
```

To run ASR evaluation, you need a set of ground truth transcripts. Use the `transcribe` command to get started.

### Flags

* `--help` `-h` _(bool)_ - help for asr
* `--streaming` _(bool)_ - Use the Streaming API instead of the Batch API.

### Examples

```
speechly evaluate asr <app_id> ground-truths.jsonl
speechly evaluate asr <app_id> ground-truths.jsonl --streaming
```
