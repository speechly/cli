## speechly sample

Sample a set of examples from the given SAL configuration

### Synopsis

The contents of the directory given as argument is sent to the
API and compiled. If configuration is valid, a set of examples are printed to stdout.

```
speechly sample [<app_id>] <directory> [flags]
```

### Examples

```
speechly sample -a <app_id> .
speechly sample -a <app_id> /path/to/config
speechly sample <app_id> /path/to/config --stats
```

### Options

```
  -a, --app string                 Application to sample the files from. Can alternatively be given as the first positional argument.
      --batch-size int             How many examples to return. Must be between 32 and 10000. (default 100)
      --seed int                   Random seed to use when initializing the sampler.
      --stats                      Print intent and entity distributions to the output.
      --advanced-stats             Print entity type, value and value pair distributions to the output.
      --advanced-stats-limit int   Line limit for advanced_stats. The lines are ordered by count. (default 10)
  -h, --help                       help for sample
```

### SEE ALSO

* [speechly](speechly.md)	 - Speechly API Client

