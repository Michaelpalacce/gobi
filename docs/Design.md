# Design

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

- https://github.com/gin-gonic/gin/blob/master/docs/doc.md#using-basicauth-middleware can be used potentially

## Sync

- To achieve bi-directional syncing, we'll implement a versioning system.
  - Store changes to a file as metadata
  - Store the versions of the file as well
