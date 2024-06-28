# gobi

Go Bi-Directional Sync With API

## Development

### Injecting environment variables

Example of environment variables that need to be set in order to run the application.
```bash
export MONGO_CONNECTION_STRING="mongodb://mongo:mongo@127.0.0.1:27017" # This is the default connection string for the docker-compose file
export MONGO_DATABASE="gobi" # This is the name of the database that will be used to store the data
export LOCAL_VAULTS_LOCATION=".dev/vaults/" # This is where the vaults will be stored
```

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


## Roadmap

- [x] Database
- [x] Project design and structure
- [x] Basic Authentication
- [x] Storage Driver Interface
- [x] Local Storage Driver
- [ ] File Uploading
- [ ] File Pushing
- [ ] Conflict resolution
- [ ] Better server interrupts handling ( send data first and then stop )
- [x] Docker Compose For Mongo And Redis
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

Events will be squashed into a single event, to avoid sending multiple events for the same file. This will be done by the server. The
squashing will be done by following the rules:
- If a file is created and then deleted, the file will be ignored.
- If a file is created and then modified, the file will be modified.
- If a file is modified and then deleted, the file will be deleted.
- If a file is modified and then modified again, the file will be modified with the latest changes.
- If a file is deleted and then created, the file will be created.

#### Conflict Resolution

In case of a conflict while executing the strategy, the client will be prompted to make a decision. 
If the client's changes are accepted, any events sent to the queue later on regarding that file will be ignored. 
This also includes additional events that may have been sent after. Once the queue is cleared, and the initial sync is marked as complete,
then events will be taken on a case by case bassis.

Conflicts will be detected by following the following rules:
- If the client has a file that the server does not have, the client's file will be sent to the server.
- If the server has a file that the client does not have, the server's file will be sent to the client.
- If both the client and the server have a file, the file with the latest timestamp will be sent to the other party.


#### Offline Syncing

##### Server

The server will store a copy of the events that have happened in the past year. This will be used to replay events to clients that have been offline for a while.

##### Client

The client will store the operations it would normally send to the server in a queue. This queue will be processed once the client is back online after the 
server has finished sending all the events that the client has missed.

### Storage Abstraction

The storage layer of the application will be abstracted, allowing different drivers to be created in the future.

- [x] Local
- [ ] AWS
- [ ] NFS

### Concurrency

