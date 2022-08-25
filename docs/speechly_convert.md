## speechly convert

Converts an Alexa Interaction Model in JSON format to a Speechly configuration

```
speechly convert [-l language] <input_file> [flags]
```

### Examples

```
speechly convert my-alexa-skill.json
speechly convert -l en-US my-alexa-skill.json
```

### Options

```
  -h, --help              help for convert
  -l, --language string   Language of input (default "en-US")
```

### SEE ALSO

* [speechly](speechly.md)	 - Speechly API Client

