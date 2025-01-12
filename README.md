## Cacher

Cacher is an in-memory database server. It's not meant for production usage; I've created this for learning purposes. 
The project has two programs, a server you can make TCP requests to and a CLI that can make requests to the server. 
Both are written in Go.

The server stores key-value pairs in memory using a [hash table](https://en.wikipedia.org/wiki/Hash_table) (Golang's built-in map). 
You can get, set, delete, and expire keys using four supported operations: `GET`, `SET`, `DEL` and `EXP`.

### Running

You can run the server by compiling it or by running a Docker container:
  - Building
    - `make run/server`
  - Docker
    - `make build/docker`
    - `make run/docker`

### Making requests

You can make requests to the server by opening a TCP connection to `:8595` using any tool you prefer or the CLI.
  - Building
    - `make build/cli`
    - `./bin/cli -operation SET -key foo -value bar`
  - Docker
    - `make build/docker`
    - `make up/docker`
    - `docker exec -it cacher /usr/local/bin/cacher/cli -operation SET -key foo -value bar`

### How the protocol works

The protocol works in a simple way, there is a format to the request, and another format to the response.
A request expects at least two values: an operation, a key and optionally a value. Example: `SET foo bar`
Valid operations are:
  - **GET**
    - retrieve a key from the store
    - expects a KEY
  - **SET**
    - store a key-value pair
    - expects a KEY and a VALUE
  - **DEL**
    - delete a key from the store
    - expects a KEY
  - **EXP**
    - set an expiration date to a key
    - expects a KEY and a Unix timestamp

A response is expected to include two values: a status and a message. Example: `ERROR should provide a value when operation is SET`
Valid statuses are:
  - **OK**
  - **ERROR**

### To Do:

- [x] TCP server
- [x] In memory store
- [x] ttl
- [x] Gracefull shutdown
- [x] CLI client
- [x] Docker
- [ ] Persistence
