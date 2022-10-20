# deploy

Send the contents of a local directory to training

## Synopsis

The contents of the directory given as argument is sent to the
API and validated. Then, a new model is trained and automatically deployed
as the active model for the application.

```
speechly deploy [<app_id>] <directory> [flags]
```

## Examples

```
speechly deploy <app_id> /path/to/config
speechly deploy -a <app_id> .
```

## Options

```
  -a, --app string        Application to deploy the files to. Can be given as the first positional argument.
  -h, --help              help for deploy
      --skip-validation   Skip the validation step. If there are validation issues, they will not be shown, the deploy will fail silently.
  -w, --watch             Wait for training to be finished.
```

## See also

* [speechly](README.md)	 - Speechly CLI

