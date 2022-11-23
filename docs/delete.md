# delete

Delete an existing application

### Usage

```
speechly delete [flags]
```

### Flags

* `--app` `-a` _(string)_ - Application to delete. Can be given as the sole positional argument.
* `--dry-run` `-d` _(bool)_ - Don't perform the deletion.
* `--force` `-f` _(bool)_ - Skip confirmation prompt.
* `--help` `-h` _(bool)_ - help for delete

### Examples

```
speechly delete <app_id>
speechly delete --app <app_id> --force
```
