# Configuration File Format

Information about the structure and usage of our JSON configuration files.

## Example config

```json
{
    "id": "example",
    "name": "Example",
    "output": "example.txt",
    "match_from": "body",
    "match": {
        "type": "regex",
        "value": "^[A-Z_]+=[^\\n]+"
    },
    "paths": [
        "/example"
    ]
}
```

## Field Information
| Field       | Type   | Description                                                    |
|-------------|--------|----------------------------------------------------------------|
| id          | string | Identifier for the configuration                               |
| name        | string | Name of the configuration                                      |
| output      | string | Name of the output file                                        |
| match_from  | string | Specifies where to apply the match (either "body" or "header") |
| match       | object | Defines the matching criteria                                  |
| match.type  | string | Type of match to perform (either "regex" or "words")           |
| match.value | string | The pattern or words to match                                  |
| paths       | array  | List of URL paths to scan                                      |