# annotate

Create SAL annotations for a list of examples using Speechly.

```
speechly annotate [<input file>] [<app id>] [flags]
```

## Examples

```
speechly annotate -a <app_id> --input input.txt
speechly annotate -a <app_id> --input input.txt > output.txt
speechly annotate -a <app_id> --reference-date 2021-01-20 --input input.txt > output.txt

To evaluate already deployed Speechly app, you need a set of evaluation examples that users of your application might say.
```

## Options

```
  -a, --app string              App ID of the application to evaluate. Can alternatively be given as the first positional argument.
  -d, --de-annotate             Instead of adding annotation, remove annotations from output.
  -e, --evaluate                Print evaluation stats instead of the annotated output.
  -h, --help                    help for annotate
  -i, --input string            Evaluation utterances, separated by newline, if not provided, read from stdin. Can alternatively be given as the first positional argument.
  -o, --output string           Where to store annotated utterances, if not provided, print to stdout.
  -r, --reference-date string   Reference date in YYYY-MM-DD format, if not provided use current date.
```

## See also

* [speechly](README.md)	 - Speechly CLI

