### How the protocol works

The protocol works in a simple way, there is a format to the request, and another format to the response.
A request expects at least two values: an operation, a key and optionally a value. Example: `OPERATION KEY VALUE`
Valid operations are:
  - **GET**
    - retrieve a key from the store
    - expects a KEY
  - **SET**
    - store a key-value pair
    - expects a KEY and a VALUE
  - **EXP**
    - set an expiration date to a key
    - expects a KEY and a Unix timestamp

A response is expected to include two values: a status and a message. Example: `STATUS MESSAGE`
Valid statuses are:
  - **OK**
  - **ERROR**

### To Do:

- [x] TCP server
- [x] In memory store
- [x] ttl
- [x] Gracefull shutdown
- [ ] Go client
- [ ] Persistence
