# Matchmaking Engine System Design

![System Design](./Matchmaking%20engine%20design.png)

## Folder Structure

```
└── 📁Matchmaking-Engine
    └── 📁cmd
        └── 📁api
            └── server.go
        └── main.go
        └── 📁matchmaker
            └── run.go
        └── 📁ws_client
            └── main.go
    └── 📁examples
        └── services.md
        └── user flow.md
    └── 📁internal
        └── 📁clientSim
            └── join_request.go
        └── 📁matchmaking
            └── elo.go
            └── matchmaking.go
        └── 📁player
            └── handler.go
            └── player.go
        └── 📁ws
            └── handler.go
            └── pubSub.go
    └── 📁pkg
        └── 📁clients
            └── manager.go
    └── 📁utils
        └── redis.go
    └── .env
    └── .gitignore
    └── go.mod
    └── go.sum
    └── Matchmaking engine design.png
    └── README.md
```
