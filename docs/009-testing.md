# Testing

## Verifying Protocol Handshake

To verify that the initial greeting matches the RFC 42TAP specification:
1. Start the server in one terminal: 
   ```bash
   make run-server
   ```
2. In a new terminal, start the CLI client: 
   ```bash
   make run-client
   ```
3. In another terminal, start the GUI client: 
   ```bash
   make run-client-gui
   ```
4. Check the logs/output of both clients to verify that the very first message received is exactly:
   ```text
   S: OK hello proto=1
   ```
