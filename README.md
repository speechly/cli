# Speechly CLI

A command line tool to:

- list apps
- upload config for apps
- download app configs

# Installation

There are binary releases for macOS, Linux and Windows, see [releases](releases/).

If you are using [Homebrew](https://brew.sh), you can install the `speechly` command with:

- `brew tap speechly/tap`
- `brew install speechly` to get the latest release

After tapping, normal brew updates will include any new versions of `speechly`.

# Usage

You need an API key to be able to access the API. After creating one in the
[Speechly dashboard](https://www.speechly.com/dashboard/), create a
configuration file for the CLI by adding a configuration context:

    speechly config add --name default --apikey APIKEY [--host api.speechly.com]

The latest context added will be used as the current context. See help for config
command to discover other uses.

After configuration, it is possible to:

- `list` apps in project
- `describe` apps to get their status
- `download` to get application's training data files as tar
- `upload` to upload a new package of training files and start training

The versioning of the configuration files should be done properly, ie. keep them in a version control system. Consider the upload/download functionality to be a tool for the training pipeline instead of collaboration or versioning.

# Develop and debug the tool

`grpcurl` is a nice tool to access gRPC APIs.

## Raw API Access with grpcurl

    grpcurl -d '{}' -H 'authorization: Bearer APIKEY' \
        -v -authority api.speechly.com \
        -protoset config.protoset \
        api.speechly.com:443 speechly.config.v1.ConfigAPI/GetProject

## protoset generation

    protoc --proto_path=protos --descriptor_set_out=config.protoset --include_imports protos/speechly/config/v1/config_api.proto
