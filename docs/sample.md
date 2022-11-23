# sample

Sample a set of examples from the given SAL configuration

### Usage

```
speechly sample [flags]
```

The contents of the directory given as argument is sent to the API and compiled. If configuration is valid, a set of examples are printed to stdout.

### Flags

* `--app` `-a` _(string)_ - Application to sample the files from. Can be given as the first positional argument.
* `--batch-size` _(int)_ - How many examples to return. Must be between 32 and 10000.
* `--seed` _(int)_ - Random seed to use when initializing the sampler.
* `--stats` _(bool)_ - Print intent and entity distributions to the output.
* `--advanced-stats` _(bool)_ - Print entity type, value and value pair distributions to the output.
* `--advanced-stats-limit` _(int)_ - Line limit for advanced_stats. The lines are ordered by count.
* `--help` `-h` _(bool)_ - help for sample

### Examples

```
speechly sample <app_id> .
speechly sample --app <app_id> /path/to/config
speechly sample <app_id> /path/to/config --stats
```
