### `annotate`

Create SAL annotations for a list of examples using Speechly.

```
speechly annotate [<input file>] [<app id>] [flags]
```

##### Examples

```
speechly annotate -a <app_id> --input input.txt
speechly annotate -a <app_id> --input input.txt > output.txt
speechly annotate -a <app_id> --reference-date 2021-01-20 --input input.txt > output.txt

To evaluate already deployed Speechly app, you need a set of evaluation examples that users of your application might say.
```

##### Options

```
  -a, --app string              App ID of the application to evaluate. Can alternatively be given as the first positional argument.
  -d, --de-annotate             Instead of adding annotation, remove annotations from output.
  -e, --evaluate                Print evaluation stats instead of the annotated output.
  -h, --help                    help for annotate
  -i, --input string            Evaluation utterances, separated by newline, if not provided, read from stdin. Can alternatively be given as the first positional argument.
  -o, --output string           Where to store annotated utterances, if not provided, print to stdout.
  -r, --reference-date string   Reference date in YYYY-MM-DD format, if not provided use current date.
```

### `convert`

Converts an Alexa Interaction Model in JSON format to a Speechly configuration

```
speechly convert [-l language] <input_file> [flags]
```

##### Examples

```
speechly convert my-alexa-skill.json
speechly convert -l en-US my-alexa-skill.json
```

##### Options

```
  -h, --help              help for convert
  -l, --language string   Language of input (default "en-US")
```

### `create`

Create a new application in the current context (project)

```
speechly create [<application name>] [flags]
```

##### Options

```
  -h, --help              help for create
  -l, --language string   Application language. Current only 'en-US' and 'fi-FI' are supported. (default "en-US")
  -n, --name string       Application name, can alternatively be given as the sole positional argument.
```

### `delete`

Delete an existing application

```
speechly delete [<app_id>] [flags]
```

##### Options

```
  -a, --app string   Application ID to delete. Can alternatively be given as the sole positional argument.
  -d, --dry-run      Don't perform the deletion.
  -f, --force        Skip confirmation prompt.
  -h, --help         help for delete
```

### `deploy`

Send the contents of a local directory to training

##### Synopsis

The contents of the directory given as argument is sent to the
API and validated. Then, a new model is trained and automatically deployed
as the active model for the application.

```
speechly deploy [<app_id>] <directory> [flags]
```

##### Examples

```
speechly deploy <app_id> /path/to/config
speechly deploy -a <app_id> .
```

##### Options

```
  -a, --app string        Application to deploy the files to. Can be given as the first positional argument.
  -h, --help              help for deploy
      --skip-validation   Skip the validation step. If there are validation issues, they will not be shown, the deploy will fail silently.
  -w, --watch             Wait for training to be finished.
```

### `describe`

Print details about an application

```
speechly describe [<app_id>] [flags]
```

##### Options

```
  -a, --app string   Application ID to describe. Can alternatively be given as the sole positional argument.
  -h, --help         help for describe
  -w, --watch        If app status is training, wait until it is finished.
```

### `download`

Download the active configuration of the given app.

##### Synopsis

Fetches the currently stored configuration from the API. This command
does not check for validity of the stored configuration, but downloads the latest
version.

```
speechly download [<app_id>] <directory> [flags]
```

##### Examples

```
speechly download <app_id> /path/to/config
speechly download -a <app_id> .
```

##### Options

```
  -a, --app string   Which application's configuration to download. Can be given as the first positional argument.
  -h, --help         help for download
```

### `edit`

Edit an existing application

```
speechly edit [flags]
```

##### Options

```
  -a, --app string    Application ID
  -h, --help          help for edit
  -n, --name string   Application name
```

### `evaluate`

Compute accuracy between annotated examples (given by 'speechly annotate') and ground truth.

```
speechly evaluate [<app_id>] [<input_file>] [flags]
```

##### Examples

```
speechly evaluate --input output.txt --ground-truth ground-truth.txt
```

##### Options

```
  -t, --ground-truth string   Manually verified ground-truths for annotated examples.
  -h, --help                  help for evaluate
  -i, --input string          SAL annotated utterances, as given by 'speechly annotate' command.
```

### `list`

List applications in the current context (project)

```
speechly list [flags]
```

##### Options

```
  -h, --help   help for list
```

### `projects add`

Add access to a pre-existing project

```
speechly projects add [apikey] [flags]
```

##### Options

```
      --apikey string            API key, created in Speechly Dashboard. Can also be given as the sole positional argument.
  -h, --help                     help for add
      --host string              API address (default "api.speechly.com")
      --name string              An unique name for the project. If not given the project name configured in Speechly Dashboard will be used.
      --skip-online-validation   Skips validating the API key against the host.
```

### `projects remove`

Remove access to a project

```
speechly projects remove [flags]
```

##### Options

```
  -h, --help          help for remove
      --name string   The name for the project for which access is to be removed.
```

### `projects use`

Select the default project used

```
speechly projects use [flags]
```

##### Options

```
  -h, --help          help for use
      --name string   An unique name for the project.
```

### `sample`

Sample a set of examples from the given SAL configuration

##### Synopsis

The contents of the directory given as argument is sent to the
API and compiled. If configuration is valid, a set of examples are printed to stdout.

```
speechly sample [<app_id>] <directory> [flags]
```

##### Examples

```
speechly sample -a <app_id> .
speechly sample -a <app_id> /path/to/config
speechly sample <app_id> /path/to/config --stats
```

##### Options

```
  -a, --app string                 Application to sample the files from. Can alternatively be given as the first positional argument.
      --batch-size int             How many examples to return. Must be between 32 and 10000. (default 100)
      --seed int                   Random seed to use when initializing the sampler.
      --stats                      Print intent and entity distributions to the output.
      --advanced-stats             Print entity type, value and value pair distributions to the output.
      --advanced-stats-limit int   Line limit for advanced_stats. The lines are ordered by count. (default 10)
  -h, --help                       help for sample
```

### `stats`

Get utterance statistics for the current project or an application in it

```
speechly stats [<app_id>] [flags]
```

##### Examples

```
speechly stats [<app_id>]
speechly stats -a APP_ID
speechly stats > output.csv
speechly stats --start-date 2021-03-01 --end-date 2021-04-01
```

##### Options

```
  -a, --app string          Application to get the statistics for. Can alternatively be given as the sole positional argument.
      --end-date string     End date for statistics, not included in results.
      --export              Print report as CSV
  -h, --help                help for stats
      --start-date string   Start date for statistics.
```

### `transcribe`

Transcribe the given jsonlines file

```
speechly transcribe <app_id> <input_file> [flags]
```

##### Examples

```
speechly transcribe <app_id> <input_file>
```

##### Options

```
  -h, --help           help for transcribe
  -m, --model string   On-device model file paths as a comma separated list.
```

### `utterances`

Get a sample of recent utterances.

##### Synopsis

Fetches a sample of recent utterances and their SAL-annotated transcript.

```
speechly utterances <app_id> [flags]
```

##### Options

```
  -h, --help   help for utterances
```

### `validate`

Validate the given configuration for syntax errors

##### Synopsis

The contents of the directory given as argument is sent to the
API and validated. Possible errors are printed to stdout.

```
speechly validate [<app_id>] <directory> [flags]
```

##### Examples

```
speechly validate -a <app_id> .
speechly validate <app_id> /path/to/config
```

##### Options

```
  -a, --app string   Application to validate the files for. Can alternatively be given as the first positional argument.
  -h, --help         help for validate
```

### `version`

Print the version number

```
speechly version [flags]
```

##### Options

```
  -h, --help   help for version
```

