# Overview

 This is a simple implementation for PoW algorithm, can be used for some study purposes, but not for production code. A lot of required things for production code are missed. Project contains:
* PoW server, it returns one random quote if client resolves challenge.
* PoW Client library
* Sample client using pow client library

# Connection Schema

1        ----request service---> 
2        <---challenge---------- 
3 Client ----solved challenge--> Server
4        <---qoute--------------
5        <---close connection---

## Schema description
1. Client creates a connection to the server.
2. Server generates a challenge and sends it back to the client. Server adds difficulty value to the end of response.
3. Client finds a solution for provided challenge and send to the server.
4. Server checks the solution and if it's correct return random quote.
5. Server closes the connection.

# Protocol description
1. Challenge packet consist three parts - total packet length (one byte, this byte isn't used in calculation), challenge (some random string) and difficulty (one byte, how many leading zeros is accepted by server)

| Name | Length | Sample |
| ---  |:------:|-------:|
| Packet size | 1 | `14` |
| Challenge | various | `[]byte("hello, world!")` |
| Difficulty | 1 | `20` |

2. Solution packet has three logical blocks: packet length, solution, suffix. Solution equals to `sha(challenge+prefix)` start with at least `difficulty` zeros in binary interpretation

| Name | Length | Sample |
| ---  |:------:|-------:|
| Packet size | 1 | `24` |
| Solution | 20 | `[20]byte{...}` |
| Suffix | various | `[]byte{0,5,0,0}` |

# Client library using

Client library can be used only in Go lang project. To create PoW connection use `client.GetPoWClient().Dial("tcp", "%SERVER_HOST%:%SERVER_PORT%")`, connection will be ready to use after. It's a standart `net.Conn`.

# PoW Server overview

Simple UUID v4 is used for random string generation, but it can be updated to other algorithm without client changes. Also text file `data/quotes.txt` is used as quotes storage for response after success verification process.

# Project structure
* cmd/client - client main file, it's just a sample of using of PoW client library
* cmd/server - server main file
* data/quotes - sample quotes storage file
* internal/server - server codebase
* internal/utils - some helper functions used in client and server
* pkg/client - client library, can be used in other projects
* docker-compose

# How to run

```
make server-linux-build
```
Compile server for Docker

```
make client-linux-build
```
Compile client for Docker

```
make start-server
```
Compile and start server in a Docker container

```
make start-client
```
Compile and start client in a Docker container

```
make test
```
Run unit tests

