# gobi

Go Bi-Directional Sync With API

## Roadmap

- [x] Database
- [x] Project design and structure
- [x] Basic Authentication
- [ ] Storage Driver Interface
- [ ] Docker Compose For Mongo And Redis
- [ ] Local Storage Driver
- [x] File Uploading
- [x] File Pushing
- [ ] Encryption at rest
- [ ] Multiple Targets
- [ ] Bi-Directional Syncing
- [ ] Versioning
- [ ] Deletion resolution
- [ ] Conflict resolution
- [ ] Better server interrupts handling ( send data first and then stop )

## Principles

- No single point of failure
  - Multiple API instances
- Scalable
  - The API can be horizontally scaled
  - Database can be horizontally scaled
  - Different storage drivers
- Resilient
- Secure

