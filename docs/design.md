# Implementation

## Auth Service (Golang)

- Setup HTTP server with `chi` or `gin`.
- Implement user registration and login endpoints.
- Use bcrypt for password hashing.
- Return JWTs for authenticated sessions.
- Connect to PostgreSQL using `pgx` or `gorm`.

## WebSocket Gateway (Golang)

- Implement an authenticated WebSocket handler.
- Handle incoming "ready", "signal_offer", "signal_answer", "ice_candidate", "disconnect", and "report" messages.
- Use Redis for mapping userID -> partnerID in sessions.

## Matchmaking Engine

- Build a simple Redis-backed queue.
- On "ready", push user ID into `matchmaking_queue`.
- Pop 2 users, pair them, notify both via WebSocket.
- Store session info temporarily in Redis and persist to DB after session ends.

## WebRTC Signaling (within WebSocket Server)

- Forward SDP/ICE messages between paired clients.
- Support reconnect/retry logic on failure.
- Disconnect should clean up Redis state and notify partner.

## Moderation API

- Accept abuse reports via WebSocket or HTTP.
- Store in PostgreSQL with reference to reporter/reported.

## Frontend (Web - React or similar)

- Login and Register UI.
- WebSocket client for signaling.
- WebRTC connection setup using standard `RTCPeerConnection`.
- In-call UI with video, text chat, and “Next” button.

##  Infrastructure

- Use Docker for all services.
- Deploy using Docker Compose or Kubernetes.
- Use environment variables for secrets/config.
- Add logging and error monitoring (e.g., Sentry or Logtail).
- Use DigitalOcean