<div align="center" markdown="1">
<br/>

![speechly-logo-duo-black](https://user-images.githubusercontent.com/2579244/193574443-130d16d6-76f1-4401-90f2-0ed753b39bc0.svg#gh-light-mode-only)
![speechly-logo-duo-white](https://user-images.githubusercontent.com/2579244/193574464-e682b2ce-dd48-4f70-88d4-a1fc5241fa47.svg#gh-dark-mode-only)

[Website](https://www.speechly.com/)
&ensp;&middot;&ensp;
[Docs](https://docs.speechly.com/)
&ensp;&middot;&ensp;
[Support](https://github.com/speechly/speechly/discussions)
&ensp;&middot;&ensp;
[Blog](https://www.speechly.com/blog/)
&ensp;&middot;&ensp;
[Login](https://api.speechly.com/dashboard/)

<br/>
</div>

# Speechly CLI

Speechly's Command Line Interface lets you manage your projects and applications, deploy new versions, download configurations, evaluate accuracy and more.

## Installation

To install Speechly CLI, open your terminal and run:

```bash
# Using Homebrew
brew tap speechly/tap
brew install speechly

# Using Scoop
scoop bucket add speechly https://github.com/speechly/scoop-bucket
scoop install speechly
```

There are binary releases for macOS, Linux and Windows, see [GitHub Releases](https://github.com/speechly/cli/releases). Also a [Docker image](https://hub.docker.com/repository/docker/speechly/cli) is built and published.

## Usage

After installing and [authenticating](https://docs.speechly.com/features/cli#authentication) Speechly CLI, you can get a list of commands by running:

```bash
speechly
```

To get a list of available sub-commands, arguments & flags run:

```bash
speechly [command]
```

Speechly CLI follows an approach similar to git or docker, where different functionalities of the tool are accessed by specifying a command followed by arguments to this command.

## Documentation

See [Using Speechly CLI](https://docs.speechly.com/features/cli) to learn more about how to use the tool.

For a full command reference, see [Speechly CLI reference](https://docs.speechly.com/reference/cli/) or the [`docs`](docs) folder.

## Generate documentation

Docs are generated when running `make`. 

You can run the generate script separately:

```bash
go run docs/generate.go docs
# or if inside /docs directory
go run generate.go .
```

## Compile and run tests

There are github actions for CI/CD, and locally you can run `make test` to run tests and `make lint` to run golangci-lint for the code.

## Speechly API access

See the [Speechly API](https://github.com/speechly/api) for more information about the API and how to access it, as well as documentation.
