"# swiftgem_go_apis" 

social-app/
├── cmd/
│   └── main.go              # Application entry point (wires dependencies, starts server)
├── internal/
│   ├── config/
│   │   └── config.go        # Loads configs from env/YAML (DB, JWT, ports)
│   ├── db/
│   │   └── db.go            # Database connection pooling and migration setup
│   ├── models/
│   │   ├── user.go          # User entity (ID, email, profile, etc.)
│   │   ├── post.go          # Post entity (content, media, timestamps)
│   │   ├── feed.go          # Feed item entity (aggregated posts)
│   │   └── notification.go  # Notification entity (type, recipient, payload)
│   ├── repositories/
│   │   ├── user_repo.go     # User CRUD operations (CreateUser, GetByEmail)
│   │   ├── post_repo.go     # Post queries (CreatePost, GetByID)
│   │   ├── feed_repo.go     # Feed aggregation (GetUserFeed)
│   │   └── notification_repo.go # Notification storage/retrieval
│   ├── services/
│   │   ├── auth_service.go  # Auth logic (hashing, token generation)
│   │   ├── user_service.go  # User business rules (profile updates, follows)
│   │   ├── post_service.go  # Post handling (validation, interactions)
│   │   ├── feed_service.go  # Feed curation and personalization
│   │   └── notification_service.go # Notification generation/dispatch
│   ├── handlers/
│   │   ├── auth_handler.go  # HTTP endpoints for auth (signup, login)
│   │   ├── user_handler.go  # User profile/follow endpoints
│   │   ├── post_handler.go  # Post CRUD and like/comment handlers
│   │   ├── feed_handler.go  # Feed retrieval endpoints
│   │   └── notification_handler.go # Notification listing/mark-as-read
│   ├── middlewares/
│   │   ├── jwt.go           # JWT authentication and authorization
│   │   └── logging.go       # Request/response logging middleware
│   └── routes/
│       └── routes.go        # Defines API routes/groups (v1/auth, v1/posts)
└── pkg/                     # Shared utilities
    ├── errors/              # Custom error types
    ├── helpers/             # Utility functions (e.g., pagination)
    └── constants/           # App-wide constants (e.g., error codes)