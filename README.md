"# swiftgem_go_apis" 

├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Loads .env (uses PORT, DB_DSN, etc.)
│   ├── db/
│   │   └── db.go            # DB connection with DB_DSN
│   ├── models/
│   │   ├── user.go          # User model
│   │   ├── post.go          # Post model
│   │   ├── feed.go          # Feed model
│   │   ├── notification.go  # Notification model
│   │   └── chat.go          # Add for chats module
│   ├── repositories/
│   │   ├── user_repo.go     # User CRUD
│   │   ├── post_repo.go     # Post CRUD
│   │   ├── feed_repo.go     # Feed queries
│   │   ├── notification_repo.go # Notification queries
│   │   └── chat_repo.go     # Add for chats module
│   ├── services/
│   │   ├── auth_service.go  # Auth logic (uses JWT_EXPIRATION_MIN)
│   │   ├── user_service.go  # Profile logic
│   │   ├── post_service.go  # Post logic
│   │   ├── feed_service.go  # Feed logic
│   │   ├── notification_service.go # Notification logic
│   │   └── chat_service.go  # Add for chats module
│   ├── handlers/
│   │   ├── auth_handler.go  # Auth endpoints (signup, login, verify OTP)
│   │   ├── user_handler.go  # Profile endpoints
│   │   ├── post_handler.go  # Post endpoints
│   │   ├── feed_handler.go  # Feed endpoints
│   │   ├── notification_handler.go # Notification endpoints
│   │   └── chat_handler.go  # Add for chats module
│   ├── middlewares/
│   │   ├── jwt.go           # JWT middleware
│   │   └── logging.go       # Logging middleware (optional)
│   └── routes/
│       └── routes.go        # API routes (add chats routes)
└── pkg/
    ├── errors/              # Custom errors
    ├── helpers/             # Utilities (e.g., pagination)
    └── constants/           # Constants