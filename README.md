# Go Key-Value Store

This project implements a simple key-value store in Go, using the Redis Serialization Protocol (RESP) for communication. It supports basic operations such as `GET`, `SET`, `EXISTS`, `DEL`, and advanced functionalities like handling Time-To-Live (TTL), data persistence, and RESP message parsing and serialization.

## Features

- **SET with TTL**: Store keys with an optional expiration time. Keys expire and are automatically deleted from the store after the TTL passes.
- **GET**: Retrieve the value of a stored key. Returns an error if the key has expired or does not exist.
- **EXISTS**: Check and return the count of specified keys that currently exist in the store.
- **DEL**: Delete specified keys and return the count of successfully deleted entries.
- **SAVE**: Persist the current state of the key-value store to disk.
- **RESP Parsing and Serialization**: Handles conversion of client-server communication data using RESP.
- **Error Handling**: Manages parsing errors and connection issues efficiently.

## Architecture Overview

### Server

- **Functionality**: Listens for TCP connections, handling incoming RESP commands and sending responses formatted in RESP.

### RESP Parser

- **Functionality**: Parses incoming RESP data into structured `Payload`.

### Command Handlers in store.go

- **Functionality**: Processes commands like `GET`, `SET`,`EXSIST`, `DELETE` and manages TTL expirations.

### Data Persistence in diskstore.go

- **Functionality**: Manages saving and loading of the key-value store data from disk.

## Installation

Clone the repository and build the project:

```bash
git clone https://example.com/your-repo.git
cd your-repo
go build

Using the Server
Interact with the server using any Redis-compatible client, such as redis-cli. Below are some example commands:


redis-cli -p 6379
127.0.0.1:6379> SET mykey myvalue EX 60
OK
127.0.0.1:6379> GET mykey
"myvalue"
127.0.0.1:6379> EXISTS mykey
(integer) 1
127.0.0.1:6379> DEL mykey
(integer) 1
127.0.0.1:6379> SAVE
OK
```
