# Milestones

## Milestone 1: Core Setup & Auth

- Setup backend repo and services in Golang
- PostgreSQL schema migration setup
- User registration and login with JWT
- Basic frontend with auth flow

## Milestone 2: WebSocket + Matchmaking

- Implement authenticated WebSocket gateway
- Create matchmaking queue in Redis
- Notify matched users via WebSocket
- Session management in Redis

## Milestone 3: WebRTC Signaling

- Forward SDP offers/answers and ICE candidates
- Implement frontend WebRTC logic
- Establish peer-to-peer video + text chat
- Handle disconnects gracefully

## Milestone 4: Moderation & Reporting

- Implement abuse report handling
- Persist reports in DB
- Block abusive users from matchmaking temporarily

## Milestone 5: MVP Hardening & Testing

- Add rate limits and error handling
- Integration tests for match-disconnect-reconnect loop
- Deploy using Docker Compose
- Manual testing of full user journey

## Milestone 6: Launch

- Monitor Redis & DB performance
- Final frontend polish
- Enable limited user rollout