# Distributed Locking System using Redis with Gin

This project implements a distributed locking system using Redis and Gin framework in Go. It provides endpoints to lock, check the lock status, and release locks for templates identified by unique IDs.

## Features

- Locking of templates to prevent concurrent access or modification.
- Checking the status of a template lock.
- Releasing locks on templates.

## Prerequisites

- Go programming language installed.
- Redis server running with proper configuration.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your_username/distributed-locking.git
   ```

2. Navigate to the project directory:

   ```bash
   cd distributed-locking
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

## Configuration

Ensure that Redis server is running and accessible from the application. Modify the Redis connection details in the `init()` function of `main.go` if needed.

```go
rdb = redis.NewClient(&redis.Options{
    Addr:     "redis-12200.c330.asia-south1-1.gce.cloud.redislabs.com:12200",
    Password: "td0g9FgW67e7BZx1RMx5UVNceFSvVkKa",
    DB:       0,
})
```

## Usage

1. Start the application:

   ```bash
   go run main.go
   ```

2. Access the endpoints using a REST client or curl:

   - `POST /api/locktemplate/:id`: Lock a template identified by `:id`.
   - `GET /api/checklocktemplate/:id`: Check the lock status of the template identified by `:id`.
   - `DELETE /api/releaselocktemplate/:id`: Release the lock on the template identified by `:id`.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.