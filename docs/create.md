# create

Create a new application in the current project

### Usage

```
speechly create [flags]
```

Creates a new application in the current project and a config file in the current working directory.

### Flags

* `--help` `-h` _(bool)_ - help for create
* `--language` `-l` _(string)_ - Application language. See docs for available options https://docs.speechly.com/basics/models (default 'en-US')
* `--name` `-n` _(string)_ - Application name. Can be given as the sole positional argument.
* `--output-dir` `-o` _(string)_ - Output directory for the config file.

### Examples

```
speechly create "My App"
speechly create --name "My App" --output-dir /foo/bar
```
