# Speechly CLI

A command line tool to:

- list apps
- upload config for apps
- download app configs

You need an API key to be able to access the API. After creating one in the
[Speechly dashboard](https://www.speechly.com/dashboard/), create a
configuration file for the CLI by adding a configuration context:

    speechly config add --name default --apikey APIKEY [--host api.speechly.com]

The latest context added will be used as the current context. See help for config
command to discover other uses.

After configuration, it is possible to:

- `list` apps in project
- `describe` apps to get their status
- `download` to get application's training data file
- `upload` to upload a new version of the training file and start training


## Raw API Access with grpcurl

    grpcurl -d '{}' -H 'authorization: Bearer APIKEY' \
        -v -authority staging.speechly.com \
        -protoset config.protoset \
        staging.speechly.com:443 speechly.config.v1.ConfigAPI/GetProject


## protoset generation

    protoc --proto_path=protos --descriptor_set_out=config.protoset --include_imports protos/speechly/config/v1/config_api.proto
