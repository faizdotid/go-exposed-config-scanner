# Configuration File Format

Information about the structure and usage of our JSON configuration files.

## Example config

```json
{
  "id": "example",
  "name": "Example",
  "output": "example.txt",
  "request": {
    "method": "GET", // POST, PUT , HEAD, DELETE, UPDATE, PATCH
    "headers": {
      "Accept": "text/plain"
    },
    "timeout": 7
  },
  "match": {
    "status_code": 0, // can be "200,404" or "all" or just 200
    "from": "body", // headers
    "type": "regex", // words, word, binary
    "value": "^[A-Z_]+=[^\\n]+" // binary := 616263
  },
  "paths": ["/example"]
}
```

## Field Information

| Field             | Type      | Description                                    |
| ---------------   | -------   | ---------------------------------------------- |
| id                | string    | Identifier for the configuration               |
| name              | string    | Name of the configuration                      |
| output            | string    | Name of the output file                        |
| request           | object    | Request configuration                          |
| request.method    | string    | HTTP method to use for the request             |
| request.headers   | object    | Headers to include in the request              |
| request.timeout   | integer   | Timeout for the request in seconds             |
| match             | object    | Match configuration                            |
| match.status_code | string    | Status code to be matched                      |
| match.from        | string    | Location to match (body or headers)            |
| match.type        | string    | Type of match (regex or string, binary)        |
| match.value       | string    | Value to match (regex or string)               |
| paths             | array     | Paths to scan for the configuration            |
