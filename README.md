# gobi

Go Bi-Directional Sync With API

## Roadmap

- [ ] Database
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
  - Cassandra Cluster
- Scalable
  - The API can be horizontally scaled
  - Cassandra can be horizontally scaled
  - Different storage drivers
- Resilient
- Secure
