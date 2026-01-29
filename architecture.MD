/cmd
    main.go

/internal
     /config
        initdb.go

/users
    domain.go
    repository.go
    service.go
    http_handler.go
    sqlite_repository.go
    /migrations
        001_create_users_table.sql

/auth
    authenticator.go
    session_store.go
    middleware.go

/posts
    domain.go
    repository.go
    service.go
    http_handler.go
    sqlite_repository.go
    /migrations
        001_create_posts_table.sql
        002_add_indexes.sql

/comments
    domain.go
    repository.go
    service.go
    http_handler.go
    sqlite_repository.go
    /migrations
        001_create_comments_table.sql

/categories
    domain.go
    repository.go
    service.go
    http_handler.go
    sqlite_repository.go
    /migrations
        001_create_categories_table.sql

/messages
    domain.go
    repository.go
    service.go
    websocket_handler.go
    websocket_notifier.go
    hub.go
    client.go
    sqlite_repository.go
    /migrations
        001_create_messages_table.sql

/server
    server.go

/shared
    errors.go
    utils.go

/frontend
    index.html
    app.js
    style.css

/vault
    forum.db

.env
.gitignore
go.mod
go.sum