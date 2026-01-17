# Reddis - A Redis Clone in Go

A lightweight Redis-like server implementation written in Go that supports basic key-value operations.

## Features

- TCP server implementation
- Basic Redis protocol (RESP) support
- Key-value storage operations (GET/SET)
- Concurrent client handling
- **TTL (Time To Live)** - Set expiration times for keys
- **Data Persistence** - Save and load data from disk
- Simple in-memory database

## Getting Started

### Prerequisites

- Go 1.x or higher
- Basic understanding of Redis protocols

### Installation

1. Clone the repository
2. Navigate to the project directory
3. Run the server:

```bash
go run main.go
```

By default, the server listens on port 5000. You can specify a custom port using the `-listenAddress` flag:

```bash
go run main.go -listenAddress=":6379"
```

## Supported Commands

- `SET key value` - Store a key-value pair
- `GET key` - Retrieve a value by key
- `HELLO` - Server greeting command
- `CLIENT INFO` - Get client information

## Key Features

### TTL (Time To Live)
Keys can be set with expiration times. The server automatically manages key expiration using an internal expiration map that tracks timestamps and durations for each key.

### Data Persistence
The server supports data persistence, allowing you to save the database state to disk and restore it on server restart. This ensures your data survives server restarts and system crashes.

## Architecture

The server implements:
- Concurrent client handling using goroutines
- Channel-based communication between components
- RESP (Redis Serialization Protocol) for client-server communication
- In-memory key-value storage

## Implementation Details

- Uses the `tidwall/resp` package for RESP protocol handling
- Implements a peer-based connection management system
- Features non-blocking message handling via channels
- Supports graceful connection handling and error management

## Limitations

- Basic implementation with limited Redis commands
- Simple error handling

## Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements.