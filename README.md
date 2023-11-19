# gobi

Go Bi-Directional Sync With API

## Roadmap

- [x] Database
- [ ] Authentication and session management with the help of Redis
- [ ] Encryption at rest
- [ ] Multiple Targets
- [ ] Storage Driver Interface
- [ ] Local Storage Driver
- [ ] File Uploading
- [ ] File Pushing
- [ ] Bi-Directional Syncing
- [ ] Versioning
- [ ] Deletion resolution
- [ ] Conflict Resolution

## Principles

- No single point of failure
  - Multiple API instances
- Scalable
  - The API can be horizontally scaled
  - Database can be horizontally scaled
  - Different storage drivers
- Resilient
- Secure

## Design

### Database

- MongoDB will be utilized for a database.
- The free tier of mongodb will be used on the cloud to simplify the deployment.
- The connection string will be passed as an environment variable.

#### Scalability

- MongoDB supports sharding and can horizontally scale if needed

#### Performance

- The data stored in MongoDB will be files metadata.
- Limited growth is expected even with millions of files.

#### Hardware Requirements

- Since I want gobi to be deployed even in hardware challenged environments, we need the database to not have a big footprint.

### Authentication

- https://github.com/gin-gonic/gin/blob/master/docs/doc.md#using-basicauth-middleware can be used potentially
