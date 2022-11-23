# deploy

Send the contents of a local directory to training

### Usage

```
speechly deploy [flags]
```

The contents of the directory given as argument is sent to the API and validated. Then, a new model is trained and automatically deployed as the active model for the application.

### Flags

* `--app` `-a` _(string)_ - Application to deploy the files to. Can be given as the first positional argument.
* `--help` `-h` _(bool)_ - help for deploy
* `--skip-validation` _(bool)_ - Skip the validation step. If there are validation issues, they will not be shown, the deploy will fail silently.
* `--watch` `-w` _(bool)_ - Wait for training to be finished.

### Examples

```
speechly deploy <app_id> /path/to/config
speechly deploy --app <app_id> .
speechly deploy --watch --app <app_id> .
```
