<h1 align="center">
<a href="https://www.speechly.com/?utm_source=github&utm_medium=cli&utm_campaign=header"><img src="https://www.speechly.com/images/logo.png" height="100" alt="Speechly"></a>
</h1>
<h2 align="center">
Speechly is the Fast, Accurate, and Simple Voice Interface API for Web and Mobile Apps
</h2>

[Speechly website](https://www.speechly.com/?utm_source=github&utm_medium=cli&utm_campaign=header)&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;[Docs](https://www.speechly.com/docs/?utm_source=github&utm_medium=cli&utm_campaign=header)&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;[Blog](https://www.speechly.com/blog/?utm_source=github&utm_medium=cli&utm_campaign=header)

# Speechly CLI

A command line tool to:

- list apps
- deploy configurations for Speechly apps
- generate a sample from the configuration to see how the SAL is transformed into text
- see statistics about the usage of Speechly apps

Learn about the [configuration syntax](https://docs.speechly.com/slu-examples/) and other topics in the [Speechly docs](https://docs.speechly.com).

# Installation

There are binary releases for macOS, Linux and Windows, see [releases](https://github.com/speechly/cli/releases).

### Homebrew

If you are using [Homebrew](https://brew.sh), you can install the `speechly` command with:

- `brew tap speechly/tap`
- `brew install speechly` to get the latest release

After tapping, normal brew updates will include any new versions of `speechly`.

### Scoop

[Scoop](https://github.com/lukesampson/scoop) is a package manager for Windows. `speechly` can be installed with scoop with:

- `scoop bucket add speechly https://github.com/speechly/scoop-bucket`
- `scoop install speechly` to install the latest release

You can get updates with `scoop update`.

# Usage

You need an API key to be able to access the API. After creating one in the
[Speechly dashboard](https://www.speechly.com/dashboard/?utm_source=github&utm_medium=cli&utm_campaign=text), create a
configuration file for the CLI by adding a configuration context:

    speechly config add --name default --apikey APIKEY [--host api.speechly.com]

The latest context added will be used as the current context. See help for config
command to discover other uses.

After configuration of the Speechly CLI, it is possible to:

- `create` create a new application in the current context (project)
- `delete` delete an existing application
- `deploy` deploy to upload a directory containing SAL configuration file(s), train a model out of them and take the model into use.
- `describe` describe apps to get their status
- `list` list apps in project
- `sample` sample a set of examples from the given SAL configuration
- `validate` validate the given SAL configuration for syntax errors
- `stats` see statistics about the apps in current context

The versioning of the SAL configuration files should be done properly, ie. keep them in a version control system. Consider the deploy/download functionality to be a tool for the training pipeline instead of collaboration or versioning.

Read our [tutorial](https://www.speechly.com/blog/configure-voice-ui-command-line/) for more information on using the Command Line Tool

# Develop and debug the tool

### Compile and run tests

There are github actions for CI/CD, and locally you can run `make test` to run tests and `make lint` to run golangci-lint for the code.

### Speechly API access

See the [Speechly API](https://github.com/speechly/api) for more information about the API and how to access it, as well as documentation.

# About Speechly

Speechly is a developer tool for building real-time multimodal voice user interfaces. It enables developers and designers to enhance their current touch user interface with voice functionalities for better user experience. Speechly key features:

#### Speechly key features

- Fully streaming API
- Multi modal from the ground up
- Easy to configure for any use case
- Fast to integrate to any touch screen application
- Supports natural corrections such as "Show me red – i mean blue t-shirts"
- Real time visual feedback encourages users to go on with their voice

|                  Example application                  | Description                                                                                                                                                                                                                                                                                                                               |
| :---------------------------------------------------: | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| <img src="https://i.imgur.com/v9o1JHf.gif" width=50%> | Instead of using buttons, input fields and dropdowns, Speechly enables users to interact with the application by using voice. <br />User gets real-time visual feedback on the form as they speak and are encouraged to go on. If there's an error, the user can either correct it by using traditional touch user interface or by voice. |
