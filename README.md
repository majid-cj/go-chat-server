# Go Chat Server - README

## Overview

This project aims to create a robust chat server that facilitates real-time communication using WebSockets. It is built using the Iris web framework and incorporates several components for managing WebSocket sessions, logging, authentication, and database interactions. `Go Chat Server` serves as the centralized configuration manager, ensuring seamless integration of these components.

---

## Features

- **Real-Time Communication**: Provides WebSocket session management via the `melody` package for real-time messaging.
- **IP Information**: Includes an `ipinfo.Client` for geolocation and IP-based information retrieval.
- **Logging**: Utilizes `zap` for structured logging with timestamps.
- **Application Context**: Manages global application state using `context.Context`.
- **Persistence**: Connects to the database via `persistence.Repository` for storing user and chat data.
- **Authentication**: Facilitates authentication through `auth.DBAuth` and `auth.Token`.
- **File Uploads**: Supports file upload operations using `fileupload.UploadFile` for sharing media in chats.
- **Session Management**: Manages WebSocket sessions in a thread-safe manner to ensure reliable communication.
- **Error Handling**: Centralized error channel (`ErrChan`) for capturing application errors.
- **CORS Support**: Configured with flexible cross-origin resource sharing settings for secure API access.

---

## Architecture

### Application Flow

1. **Server Initialization**:
   - Loads environment variables using `godotenv`.
   - Sets up logging, persistence, authentication, and WebSocket session handling.
   - Configures middleware for logging, recovery, and CORS.
2. **WebSocket Management**:
   - Handles WebSocket connections and sessions using `melody`.
3. **Routing**:
   - Defines API routes using the `router` package.
4. **Graceful Shutdown**:
   - Listens for termination signals to clean up resources and close active connections.

---

## Components

### Struct Fields

- **`Melody`**: Handles WebSocket communications for the chat server.
- **`IPInfo`**: Manages IP-based geolocation using `github.com/ipinfo/go/v2/ipinfo`.
- **`Log`**: Logging utility configured with `zap`.
- **`AppContext`**: Base context for the application.
- **`Wg`**: WaitGroup for goroutine synchronization.
- **`ErrChan`**: Error channel to capture and propagate errors.
- **`App`**: Iris web application instance.
- **`Persistence`**: Repository layer for database interactions, including user and message storage.
- **`Auth`**: Handles user authentication.
- **`Token`**: Token generation and validation utility.
- **`Upload`**: File upload utility for media sharing in chats.
- **`Session`**: Thread-safe session store for WebSocket connections.

---

## Usage

### Server Setup

The main entry point is `main.go`, which initializes the server and configures all components. Key parts include:

1. **Environment Variables**:
   Ensure the following environment variables are set:

   - `PORT`: Port on which the server will run.
   - `IP_INFO`: API key for IPInfo.
   - `UPLOADS`: Directory path for uploaded files.

2. **Server Initialization**:

   ```go
   appConfig, err := config.NewAppConfig()
   if err != nil {
       log.Fatalf("Failed to initialize AppConfig: %v", err)
   }
   defer appConfig.Persistence.Client.Disconnect(appConfig.AppContext)
   defer appConfig.Auth.DB.Close()
   ```

3. **CORS Middleware**:

   ```go
   CORS := cors.New(cors.Options{
       AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
       AllowedHeaders:   []string{"Content-Type", "Accept", "Authorization"},
       AllowedOrigins:   []string{"*"},
       AllowCredentials: true,
       MaxAge:           5000,
   })
   appConfig.App.UseRouter(CORS)
   ```

4. **Routing**:
   Routes are defined in the `router` package:

   ```go
   router.APIVersionOne(appConfig)
   ```

5. **Start Server**:
   ```go
   err := appConfig.App.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")), iris.WithOptimizations)
   if err != nil {
       appConfig.Log.Errorf("Error starting server: %v", err)
   }
   ```

### Graceful Shutdown

The server listens for system interrupts to shut down gracefully:

```go
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
sig := <-c
appConfig.Log.Info("Received signal:", sig)
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
appConfig.App.Shutdown(ctx)
```

---

## Dependencies

The project relies on the following external packages:

- `github.com/olahol/melody`: WebSocket management.
- `github.com/ipinfo/go/v2/ipinfo`: IP information services.
- `go.uber.org/zap`: Structured logging.
- `github.com/kataras/iris/v12`: Web framework.
- `github.com/iris-contrib/middleware/cors`: Cross-origin resource sharing middleware.
- `github.com/joho/godotenv`: Environment variable loader.
- `github.com/majid-cj/go-chat-server/infrastructure/auth`: Authentication utilities.
- `github.com/majid-cj/go-chat-server/infrastructure/persistence`: Database interaction layer.
- `github.com/majid-cj/go-chat-server/util/fileupload`: File upload utilities.

---

## Contributing

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature-name`.
3. Commit your changes: `git commit -m 'Add some feature'`.
4. Push to the branch: `git push origin feature-name`.
5. Open a pull request.
