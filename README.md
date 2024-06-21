# gobi

Go Bi-Directional Sync With API

## Development

### Setting up the environment

```
docker-compose up -d
```

### Running the server 

```bash
go run cmd/gobi/main.go
```

### Registering a user

```bash
curl --location 'http://localhost:8080/api/v1/users/' \
--header 'Content-Type: application/json' \
--data '{
    "username": "root",
    "password": "toor"
}'
```

## Thoughts

NOTE: Went in the wrong direction for using websockets for file transfer as well as communication. Need to refactor big parts of the project
to accomodate these changes. Lessons learned: always start with a design.

### Problems with multiple concurrent requests.

1. Server wants to send 5 files to the client.
2. The server sends 5 requests to the client at the same time
3. Everything breaks apart
4. I think the correct thing to do would be to use normal REST communication in this case.

If I do use rest, I can queue up fetches and pushes on the client easilly. The Server does not need to queue anything,
rather the server just needs to notify the clients. The rest part of the server will publish to redis and the websockets will
listen for events.

- [ ] Implement REST interface for sending and receiving items
- [ ] Refactor syncstrategies
- [ ] Refactor websocket to not send and receive file anymore, create a separate helper
- [ ] Refactor the processors to send the files with REST
- [ ] Remove unnecessary messages.

## Roadmap

- [x] Database
- [x] Project design and structure
- [x] Basic Authentication
- [x] Storage Driver Interface
- [x] Local Storage Driver
- [x] File Uploading
- [x] File Pushing
- [x] Conflict resolution
- [x] Better server interrupts handling ( send data first and then stop )
- [ ] Docker Compose For Mongo And Redis
- [ ] Encryption at rest
- [ ] Multiple Targets
- [ ] Bi-Directional Syncing
- [ ] Versioning
- [ ] Deletion resolution

## Principles

- No single point of failure
  - Multiple API instances
- Scalable
  - The API can be horizontally scaled
  - Database can be horizontally scaled
  - Different storage drivers
- Resilient
- Secure

## Requirements

### Functional Requirements

### Non-Functional Requirements

## Design

### Websocket Handling

The client will use websockets as a way of communication, to receive notifications from the server about any changes done by other connected
clients to the same vault. 

Any form of metadata will be transferred via websockets, to speed up the overall application.

Files will never be sent out via websockets and instead will be handled by the [Data Transmission Use Case](#data-transmission). However,
renaming, deletion and other similar actions are accepted.

It is possible for renaming and deletion to result in a notification to the client that the operation could not be completed. If this is so,
the server must notify what actions the client should take to catch up and be in sync with the rest. 

The communication will be Bi-Directional between the client and the server, where each party can notify the other for any changes that they
wish to do to the method of communication. The server and client must be able to respond to any such changes effective from receiving that
message.

#### Notifying the client for errors from the requested operation

In case when the server denied the request, that means that conflict resolution happend and the server decided that the change was stale. At
this point, the server must notify the client as he would normally of what kind of change is needed.

Example would be the client changes a file, but the server receives a delete signal before the client has a chance to react and the
sync strategy determined that the server has the correct record.

### Data Transmission

The client will use REST communication for sending and receiving files.

Upon receiving a file if the file is changed from what the server has locally, the receiving server should notify all others that the 
file has been changed by way of Redis channels

### Sync Strategies Abstraction

Sync starategies will be used to hold different conflict resolution methods. They are an abstraction that is supposed to make an automated
or non-automated decision what should happen in case of a sync conflict. 

### Bi-Directional Syncing

The server will store a copy of events from a variable amount of time. By default this will be set to 1 year.

When a client connects, the client will notify the Server when the last event received was and the server will send all the 
events that the client needs to repeat. Upon receiving this list, the client will enqueue the events that it needs to replay 
and start executing on them. Whenever files need to be fetched from the server, the [Data Transmission](#data-transmission)
section will be used to handle that. While the process is ongoing, the server may notify the client of any additional changes that have happened. Any changes that have
ocurred, will be added at the end of the client's queue.

#### Conflict Resolution

In case of a conflict while executing the strategy, the client will be prompted to make a decision. 
If the client's changes are accepted, any events sent to the queue later on regarding that file will be ignored. 
This also includes additional events that may have been sent after. Once the queue is cleared, and the initial sync is marked as complete,
then events will be taken on a case by case bassis.

### Storage Abstraction

The storage layer of the application will be abstracted, allowing different drivers to be created in the future.

- [ ] Local
- [ ] AWS
- [ ] NFS

### Concurrency

### Bi-Directional Syncing

