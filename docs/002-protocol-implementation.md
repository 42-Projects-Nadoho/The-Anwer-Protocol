# Protocol Implementation

The Answer Protocol (TAP) is strictly designed to follow the RFC 42TAP specifications. Our implementation ensures robust, concurrent handling of TCP connections with full compliance to the required message formats and error codes.

## Handshake and Greeting
The server handshake and initial greeting strictly adhere to the RFC 42TAP specifications upon any client connection (whether CLI or GUI). Immediately after a connection is established, the server sends the expected protocol version greeting before awaiting the `CONNECT` command.

Example greeting:
```text
S: OK hello proto=1
```

## ABNF Syntax Compliance

The server's message parser is strictly built against the ABNF syntax definitions provided in the RFC:
- All incoming client commands are validated against their expected ABNF structure.
- Malformed commands, missing arguments, or invalid syntax immediately trigger protocol-compliant error responses (e.g., `ERR`) rather than causing server instability.
- All outbound events and responses from the server are formatted strictly according to the RFC ABNF rules, ensuring perfect compatibility with any compliant third-party clients.
