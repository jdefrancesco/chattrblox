# Chatrblox

The goal is to create a minimal clone of Chatroulette/Omegle where users are randomly paired in 1-on-1 video and text chat sessions. The product supports anonymous-style interactions but requires users to log in, enabling moderation and basic user history tracking. Users can initiate a chat, get matched with a random online user, communicate through video and text, and either continue or disconnect to be matched with someone new. The backend is implemented in Golang to leverage its high-concurrency strengths and performance for real-time communication.

## Requirements

### Must Have

* User authentication (email/password or social login)
* Real-time 1-on-1 random pairing of users
* WebRTC-based video and text chat support
* Backend matching engine (in Golang) to pair online users
* Disconnect and "Next" functionality to match with another user
* Moderation hooks (report user, basic abuse flagging)
* Should Have
* User profile (minimal - avatar, age, interests)
* Basic chat logs (text only) with time-limited retention
* WebSocket signaling server for WebRTC setup
* Rate limiting to prevent abuse (e.g. max skips per minute)
* Could Have
* Interest-based matching
* Reconnection to previous users (if both agree)
* Admin dashboard for moderation
* Won’t Have
* Group chat or more than 1:1 matching
* Native mobile app (initial version is web-based)

## Method

### High-Level Architecture

The system is composed of the following components:

* Frontend Client (Web)
* Connects via WebRTC for video + text.
* Uses WebSocket for signaling and communication with backend.
* Handles UI for login, matching, and chat.
* Authentication Service (Golang)
* Manages user login, registration, and JWT token issuance.
* Interfaces with a PostgreSQL database for user accounts.
* Uses GORM ORM for data access and bcrypt for password security.
* Matchmaking Engine (Golang)
* Maintains pool of connected users via in-memory queue. (redis)
* Matches random available users for chat.
* Communicates match data over WebSocket.
* Signaling Server (Golang)
* Handles WebRTC signaling (SDP offers/answers, ICE candidates).
* Routes messages between paired clients until peer-to-peer connection is established.
* Moderation Service (Golang)
* Handles reports, abuse flagging, and user blocking.
* Database (PostgreSQL)
* Stores user data, basic profile info, and limited chat logs.