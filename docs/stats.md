# stats

Get utterance statistics for the current project or an application in it

```
speechly stats [<app_id>] [flags]
```

# Examples

```
speechly stats [<app_id>]
speechly stats -a APP_ID
speechly stats > output.csv
speechly stats --start-date 2021-03-01 --end-date 2021-04-01
```

# Options

```
  -a, --app string          Application to get the statistics for. Can alternatively be given as the sole positional argument.
      --end-date string     End date for statistics, not included in results.
      --export              Print report as CSV
  -h, --help                help for stats
      --start-date string   Start date for statistics.
```

# See also

* [speechly](README.md)	 - Speechly CLI

