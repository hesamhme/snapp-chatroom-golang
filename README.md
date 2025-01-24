# Snapp Chatroom Golang

[![Go](https://img.shields.io/badge/Go-1.23-blue)](https://golang.org)
[![NATS](https://img.shields.io/badge/NATS-2.10.24-blue)](https://nats.io)
[![Docker](https://img.shields.io/badge/Docker-20.10.7-blue)](https://www.docker.com)
[![Redis](https://img.shields.io/badge/Redis-7.4.2-red)](https://redis.io)

This project is a chatroom application built using Go (Golang) with TCP, Redis, and NATS messaging system.


## Features

- User can join/leave chatrooms
- Broadcast messages across chatrooms using NATS
- Store chatroom user data in Redis
- Interactive CLI client for users
- Logging capabilities for monitoring

## Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Nats
- Redis

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/hesamhme/snapp-chatroom-golang.git
   cd snapp-chatroom-golang
   ```
2. Build and run the application using Docker Compose:

   ```
   docker-compose up --build
   ```
3. Build server and client binary file:
   ```
   go build -o bin/server cmd/server/main.go
   go build -o bin/client cmd/client/main.go
   ```
4. Run server:
   ```
   $ bin/server 
   ```
   you should see:
   ```
   <date> <time> Connected to Redis
   ```
   means you are all set and can run as many as client you want in other terminals.

5. Run a client in a seperate terminal:
   ```
   bin/client

   ```
   you should see:
   ```
   Enter your username:
   Enter chatroom name: 

   ```

  Enter your name and then enter a chatroom name. If the chatroom already exists, you will join it; if not, don't worry—the chatroom will be created, and you will join it!

## Runing Tests
1. To run unit tests:
   ```
   go test ./tests -v

   ```

## Usage
1. Start the chatroom server and client.
2. Enter a username and chatroom name when prompted.
3. Send messages by typing and pressing Enter.
4. Use commands like ```#users``` to list chatroom users and ```#rooms``` to list active chatrooms.
5. Exit the chatroom by typing ```#exit```.

## Environment Variables
 The following environment variables can be configured:

 - ```NATS_URL``` - NATS server address (default: ```nats://localhost:4222```)
 - ```REDIS_URL``` - Redis server address (default: ```localhost:6379```)

## Logs

The application logs important events and errors to help with debugging and monitoring. By default, logs are stored in the `log` directory within the project root.

### Log Levels
The application uses structured logging with various levels, including:

- `INFO` – General information about application flow.
- `WARN` – Non-critical issues that may require attention.
- `ERROR` – Serious issues that need immediate attention.
- `FATAL` – Critical issues causing the application to terminate.

#### Default Log Location

Logs are stored in the following file:
 ```
 log/app.log
 ```


#### Configuring Logging
To customize log settings, you can:

1. Modify the `logger.go` file inside `internal/infrastructure/logger/`.
2. Change the log output directory or level via environment variables.

#### Viewing Logs

To view logs in real-time, run:

```
tail -f log/app.log
```
For Windows users:

```
Get-Content log/app.log -Wait

```
### Example Log Output
```
{"level":"info","time":"2025-01-23T22:25:58+03:30","msg":"Chat server starting..."}
{"level":"info","time":"2025-01-23T22:26:08+03:30","msg":"New client connected: 127.0.0.1:42678"}
{"level":"error","time":"2025-01-23T22:27:28+03:30","msg":"Failed to connect to Redis"}
```

## Contributions
 Contributions are welcome! Feel free to open an issue or submit a pull request.

## License
 This project is licensed under the MIT License.
