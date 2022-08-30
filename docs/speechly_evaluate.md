## speechly evaluate

Compute accuracy between annotated examples (given by 'speechly annotate') and ground truth.

```
speechly evaluate [<app_id>] [<input_file>] [flags]
```

### Examples

```
speechly evaluate --input output.txt --ground-truth ground-truth.txt
```

### Options

```
  -t, --ground-truth string   Manually verified ground-truths for annotated examples.
  -h, --help                  help for evaluate
  -i, --input string          SAL annotated utterances, as given by 'speechly annotate' command.
```

### SEE ALSO

* [speechly](speechly.md)	 - Speechly CLI

