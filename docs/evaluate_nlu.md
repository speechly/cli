# evaluate nlu

Evaluate the NLU accuracy of the given application model

### Usage

```
speechly evaluate nlu [flags]
```

To run NLU evaluation, you need a set of ground truth annotations. Use the `annotate` command to get started.

### Flags

* `--help` `-h` _(bool)_ - help for nlu
* `--reference-date` `-r` _(string)_ - Reference date in YYYY-MM-DD format, if not provided use current date.
* `--relax` _(bool)_ - Ignore normalized entity values and casing in matching.

### Examples

```
speechly evaluate nlu <app_id> ground-truths.txt
speechly evaluate nlu <app_id> ground-truths.txt --reference-date 2021-01-20
```
