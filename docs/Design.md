# Design

`gobi` is an API framework for syncing multiple client "vaults" with the help of Websockets. 

## Database

- MongoDB will be utilized for a database.
- The free tier of mongodb will be used on the cloud to simplify the deployment.
- The connection string will be passed as an environment variable.

### Scalability

- MongoDB supports sharding and can horizontally scale if needed

### Performance

- The data stored in MongoDB will be files metadata.
- Limited growth is expected even with millions of files.

### Hardware Requirements

- Since I want gobi to be deployed even in hardware challenged environments, we need the database to not have a big footprint.

## Web Framework

- `https://github.com/gin-gonic/gin` will be used since it's simple, fast and well established

## Authentication

- `https://github.com/gin-gonic/gin/blob/master/docs/doc.md#using-basicauth-middleware` will be used for now

## Sync

- To achieve bi-directional syncing, we'll implement a versioning system.
  - Store changes to a file as metadata
  - Store the versions of the file as well

## Development 

- Notifications
	- Will use pub/sub to send notifications with the help of Redis
- Websockets 
	- Connection pools
- Server 
	- No state stored. 
- Client
	- Client will store Last Sync locally 
		- This will be sent to the server
	- The server will check for events since last sync
	- The server will send back what the client needs to do. 
	- Client will store the changes they need to do and check what has changed locally. We rely on the file change timestamp for information 

## Scenarios 

- First time client 
	- Time stamp is 0 and everything is downloaded 
- Client who has added a file offline 
	- The sync happens, server has no changes since last sync, the client does 
	- Stores what needs to happen locally
	- Sends the file 
- Server has a new file on connection 
	- The check happens the server sends that there is a new file since last sync
	- The client pulls it 
- Client has deleted a file offline 
	- The check happens. 
	- Client says that it doesn't have the file 
	- The server does 
	- The file gets downloaded 
- Client has deleted a file offline with the client running 
	- The change should be added to the list of operations
- Server has deleted a file since last sync 
	- We see there is tombstone metadata for a file in the database. 
	- The file gets deleted on the client. 
	- After one year, garbage collection will delete expired data and files.
- Server/Client has renamed the file 
	- One deletion and one addition
- Both server and client have changes on a single file 
	- Show a conflict resolution window 
- General notes 
	- The server always relies on the timestamp in the database 
	- The client always relies on the os for timestamps 
	- The client queue of what needs to happen is a simple file 
