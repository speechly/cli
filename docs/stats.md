# stats

Get utterance statistics for the current project or an application in it

### Usage

```
speechly stats [flags]
```

### Flags

* `--app` `-a` _(string)_ - Application to get the statistics for. Can be given as the sole positional argument.
* `--end-date` _(string)_ - End date for statistics, not included in results.
* `--export` _(bool)_ - Print report as CSV
* `--help` `-h` _(bool)_ - help for stats
* `--start-date` _(string)_ - Start date for statistics.

### Examples

```
speechly stats <app_id>
speechly stats --app <app_id>
speechly stats > output.csv
speechly stats --start-date 2021-03-01 --end-date 2021-04-01
```
