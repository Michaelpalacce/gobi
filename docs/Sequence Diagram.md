# Diagram 

Client ---- Connects ----> Server
Client ---- Vault name ----> Server
Client ---- Sync strategy to use ----> Server
Client ---- Sync Message ----> Server
Client <---- Replay events since last sync ---- Server
Server: Subscribes for events
Client: Replays events


## Notes

