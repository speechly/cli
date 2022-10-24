# create

Create a new application in the current project

## Synopsis

Creates a new application in the current project and a config file in the current working directory.

```
speechly create [<application name>] [flags]
```

## Examples

```
speechly create "My app"
speechly create --name "My app" --output-dir /foo/bar

```

## Options

```
  -h, --help                help for create
  -l, --language string     Application language. Available options are 'en-US' and 'fi-FI'. (default "en-US")
  -n, --name string         Application name. Can be given as the sole positional argument.
  -o, --output-dir string   Output directory for the config file.
```

## See also

* [speechly](README.md)	 - Speechly CLI

