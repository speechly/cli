<div align="center" markdown="1">
<a href="https://www.speechly.com/#gh-light-mode-only">
   <img src="https://d33wubrfki0l68.cloudfront.net/f15fc952956e1952d6bd23661b7a7ee6b775faaa/c1b30/img/speechly-logo-duo-black.svg" height="48" />
</a>
<a href="https://www.speechly.com/#gh-dark-mode-only">
   <img src="https://d33wubrfki0l68.cloudfront.net/5622420d87a4aad61e39418e6be5024c56d4cd1d/94452/img/speechly-logo-duo-white.svg" height="48" />
</a>

### Real-time automatic speech recognition and natural language understanding tools in one flexible API

[Website](https://www.speechly.com/)
&ensp;|&ensp;
[Docs](https://docs.speechly.com/)
&ensp;|&ensp;
[Discussions](https://github.com/speechly/speechly/discussions)
&ensp;|&ensp;
[Blog](https://www.speechly.com/blog/)
&ensp;|&ensp;
[Podcast](https://anchor.fm/the-speechly-podcast)

---
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
