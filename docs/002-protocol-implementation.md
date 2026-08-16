# Protocol Implementation

The Answer Protocol (TAP) is strictly designed to follow the RFC 42TAP specifications. Our implementation ensures robust, concurrent handling of TCP connections with full compliance to the required message formats and error codes.

## Handshake and Greeting
The server handshake and initial greeting strictly adhere to the RFC 42TAP specifications upon any client connection (whether CLI or GUI). Immediately after a connection is established, the server sends the expected protocol version greeting before awaiting the `CONNECT` command.

Example greeting:
```text
S: OK hello proto=1
```

*(Any further deviations or specific implementation choices regarding ABNF syntax parsing or event broadcasting will be documented here.)*
