# annotate

Create SAL annotations for a list of examples using Speechly

### Usage

```
speechly annotate [flags]
```

To evaluate already deployed Speechly app, you need a set of evaluation examples that users of your application might say.

### Flags

* `--app` `-a` _(string)_ - Application to evaluate. Can be given as the second positional argument.
* `--de-annotate` `-d` _(bool)_ - Instead of adding annotation, remove annotations from output.
* `--evaluate` `-e` _(bool)_ - Print evaluation stats instead of the annotated output.
* `--help` `-h` _(bool)_ - help for annotate
* `--input` `-i` _(string)_ - Evaluation utterances, separated by newline, if not provided, read from stdin. Can be given as the first positional argument.
* `--output` `-o` _(string)_ - Where to store annotated utterances, if not provided, print to stdout.
* `--reference-date` `-r` _(string)_ - Reference date in YYYY-MM-DD format, if not provided use current date.

### Examples

```
speechly annotate input.txt <app_id>
speechly annotate --app <app_id> --input input.txt > output.txt
speechly annotate --app <app_id> --reference-date 2021-01-20 --input input.txt > output.txt
```
