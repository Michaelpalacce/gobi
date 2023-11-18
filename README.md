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
- The connection string will be passed as an environment variable
