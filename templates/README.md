# Configuration File Format

Information about the structure and usage of our JSON configuration files.

## Example config

```json
{
  "id": "example",
  "name": "Example",
  "output": "example.txt",
  "request": {
    "method": "GET",
    "headers": {
      "Accept": "text/plain"
    },
    "timeout": 7
  },
  "match": {
    "from": "body",
    "type": "regex",
    "value": "^[A-Z_]+=[^\\n]+"
  },
  "paths": ["/example"]
}
```

## Field Information

| Field           | Type    | Description                         |
| --------------- | ------- | ----------------------------------- |
| id              | string  | Identifier for the configuration    |
| name            | string  | Name of the configuration           |
| output          | string  | Name of the output file             |
| request         | object  | Request configuration               |
| request.method  | string  | HTTP method to use for the request  |
| request.headers | object  | Headers to include in the request   |
| request.timeout | integer | Timeout for the request in seconds  |
| match           | object  | Match configuration                 |
| match.from      | string  | Location to match (body or headers) |
| match.type      | string  | Type of match (regex or string)     |
| match.value     | string  | Value to match (regex or string)    |
| paths           | array   | Paths to scan for the configuration |
