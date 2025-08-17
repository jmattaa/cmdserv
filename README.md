# cmdserv

create a `endpoints.json` next to the binary of the file.

# example endpoints.json

```json
{
    "endpoints": [
        {
            "endpoint": "/",
            "command": ["say", "Hello World!"]
        },
        {
            "endpoint": "/spotify",
            "command": ["open", "/Applications/Spotify.app"]
        }
    ]
}
```

send a http request from any devince in the network and the command will run.
