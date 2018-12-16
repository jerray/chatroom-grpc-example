# Chatroom

gRPC bidirectional streaming example.

## Build

```sh
make build
```

## Run Server

```sh
./chatroom server -p 3000
```

## Run Client

```sh
./chatroom client -s 127.0.0.1:3000
```

Then type `#login [name]` to register a user name on server.

Open another terminal tab and start another client. Register another user.
Then users can chat with each other. Type something and press Enter in the
console, another user will receive the message.

### Client Commands

* `#login [name]` Register a user name
* `@[name] [message]` Send a message to a user
* `[message]` Send a message to all users in the chatroom

## Licence 

MIT
