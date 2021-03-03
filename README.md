# TickTackGo

> a local network tcp tick tack toe game, written in go!

## Usage

1. Start a server
    ```bash
    > go run cmd/ticktackgo/server/main.go
    ```
2. Start two clients
3. In the client that was started first, provide a user name
4. When prompted for a game code, in the first client, press enter
5. In the client that was started second, provide a username, and then provide the code from the first client
6. Play tick-tack-go!
7. Each user will get a turn, as indicated. The user may input 0-8 and press enter. 0-8 represent squares on the tick-tack-go board as follows
    ```bash
      0  |  1  |  2  
    ----------------
      3  |  4  |  5
    ----------------
      6  |  7  |  8  
    ```
8. Each client must wait for the other client to take their turn.
9. Once a winning combination occurs, the game will complete!

### Todo's / Improvements

- [ ] Handle stalemate scendarios better
    - [ ] if 9 turns taken (stored on board), claim stalemate and close game, or reset board
- [ ] Allow for internet connectivity to the game server
    - [ ] Create certificates for secure connection
- [ ] Better error handling
    - [ ] in cases of placing in existing sqaures
    - [ ] in case of server crash or client crash
    - [ ] better graceful handling of client exit
- [ ] Push message handling for send and receive into channels for better read/write performance
    - [ ] ideally these channels could either be on the game, board, or clients themselves. Likely clients.
- [ ] Do not block on name-setting. Upon new connection, push connection into async chan handler.
