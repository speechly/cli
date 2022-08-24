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

To install Speechly CLI, run these commands from any directory in your terminal:

### Homebrew

```bash
brew tap speechly/tap
brew install speechly
```

### Scoop

```bash
scoop bucket add speechly https://github.com/speechly/scoop-bucket
scoop install speechly
```

There are binary releases for macOS, Linux and Windows, see [releases](https://github.com/speechly/cli/releases). Also a [Docker image](https://hub.docker.com/repository/docker/speechly/cli) is built and published.

## Usage

After installing the Speechly CLI, you can run the `speechly` command.

```bash
speechly [command]
```

Speechly CLI follows an approach similar to git or docker, where different functionalities of the tool are accessed by specifying a command followed by arguments to this command.


## Documentation

To learn how to log in to Speechly and start deploying configurations, visit the [Speechly Documentation](https://docs.speechly.com/dev-tools/command-line-tool/)

For a full command reference, se the list below or, view [the CLI Reference](docs)

### Commands

[`annotate`](docs/speechly_annotate.md)
Create SAL annotations for a list of examples using Speechly.

[`completion`](docs/speechly_completion.md)
Generate the autocompletion script for the specified shell

[`convert`](docs/speechly_convert.md)
Converts an Alexa Interaction Model in JSON format to a Speechly configuration

[`create`](docs/speechly_create.md)
Create a new application in the current context (project)

[`delete`](docs/speechly_delete.md)
Delete an existing application

[`deploy`](docs/speechly_deploy.md)
Send the contents of a local directory to training

[`describe`](docs/speechly_describe.md)
Print details about an application

[`download`](docs/speechly_download.md)
Download the active configuration of the given app.

[`edit`](docs/speechly_edit.md)
Edit an existing application

[`evaluate`](docs/speechly_evaluate.md)
Compute accuracy between annotated examples (given by 'speechly annotate') and ground truth.

[`help`](docs/speechly_help.md)
Help about any command

[`list`](docs/speechly_list.md)
List applications in the current context (project)

[`projects`](docs/speechly_projects.md)
Manage API access to Speechly projects

[`sample`](docs/speechly_sample.md)
Sample a set of examples from the given SAL configuration

[`stats`](docs/speechly_stats.md)
Get utterance statistics for the current project or an application in it

[`transcribe`](docs/speechly_transcribe.md)
Transcribe the given jsonlines file

[`utterances`](docs/speechly_utterances.md)
Get a sample of recent utterances.

[`validate`](docs/speechly_validate.md)
Validate the given configuration for syntax errors

[`version`](docs/speechly_version.md)
Print the version number

## Develop and debug the tool

### Compile and run tests

There are github actions for CI/CD, and locally you can run `make test` to run tests and `make lint` to run golangci-lint for the code.

### Speechly API access

See the [Speechly API](https://github.com/speechly/api) for more information about the API and how to access it, as well as documentation.
