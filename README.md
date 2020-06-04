# Speechly CLI

A command line tool to:

- list apps
- upload config for apps
- download app configs

You need an API key to be able to access the API. Add it to a configuration
file named `.speechly.yaml` in the working directory:

```
current-context: default
contexts:
- name: default
  host: staging.speechly.com
  apikey: APIKEY

```

## Raw API Access with grpcurl

    grpcurl -d '{}' -H 'authorization: Bearer APIKEY' \
        -v -authority staging.speechly.com \
        -protoset config.protoset \
        staging.speechly.com:443 speechly.config.v1.ConfigAPI/GetProject


## protoset generation

    protoc --proto_path=protos --descriptor_set_out=config.protoset --include_imports protos/speechly/config/v1/config_api.proto
