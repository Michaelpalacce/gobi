# Diagram 

Client ----  Connects                        ----> Server
Client ----  Vault name                      ----> Server
Client ----  Sync strategy to use            ----> Server
Client ----  Sync Message                    ----> Server
Client <---- Replay events since last sync   ----  Server
Server: Subscribes for events and acts upon them. Check [Event Subscription](#event-subscription)
Client <---- Starts sending new events       ----  Server
Client: Replays events
Client ----  Requests Item(s)                ----> Server
Client <---- Sends Item(s)                   ----  Server
Client ----  Notifies of sync end            ----> Server
Client <---- Sends the sha of every file     ----  Server
Client: Starts watching for changes
Client ----  Detects changes based on sha    ----> Server
Client ----  Starts sending events to server ----> Server
Server: Publishes events to the PubSub
    All Servers act based on the strategy outlined in [Event Subscription](#event-subscription)

## Event Subscription

- Replays the event with the Storage Driver.
    - Storage Driver has to decide if each server needs to do modifications.
        - Local: Each has to do it.
        - Remote: Acquire lock, one does it, others wait, then they continue without doing anyting.
