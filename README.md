### How the protocol works

The protocol works in a simple way, there is a format to the request, and another format to the response.
A request expect at least two values, a operation, a key and optionaly a value. Ex: `OPERATION KEY VALUE`
Valid operations are:
  - GET
    - expects a KEY
  - SET
    - expects a KEY and a VALUE

Two values are expected in a response, a status and a message. Ex: `STATUS MESSAGE`
Valid statuses are:
  - OK
  - ERROR

### To Do:

- [x] TCP server
- [x] In memory store
- [ ] ttl
- [ ] Go client
- [ ] Gracefull shutdown
- [ ] Persistence
