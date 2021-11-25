# MulticastSDCCProject
Once installed Docker, run ```docker-compose up``` to run the application.

### Algorithm's type
- B   : basic multicast
- SQ : totally ordered centralized with sequencer
- SC : totally ordered distributed with scalar clock
- VC  : casually ordered with vector clock

# MULTICAST API
- Basepath of api : multicast/v1

| PATH | METHOD | SUMMARY | 
| ---- | ---------- | -------- |
| /groups | POST | Create Multicast Group |
|  /groups/{mId} | GET | Get information about a group|
|  /groups/{mId} | DELETE | Close an existing group |
|  /groups/{mId} | PUT | Start a multicast group |
|  /messages/{mId} | GET | Get messages of a group |
|  /messages/{mId} | POST | Send a message to a group |
| /deliverQueue/{mId} | GET | Get deliver queue of a group |
