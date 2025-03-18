# Signaling Server

Responsibilities:
- Handle WebSocket connections from users.
- Relay signaling messages (SDP offer/answer, ICE candidates) between matched peers.
- Cleanup on disconnect.

Each user is identified via JWT token during WebSocket handshake. Session and peer information is stored in Redis temporarily.

WebSocket Message Types:

- `signal_offer`: sent from initiator with SDP offer.
- `signal_answer`: sent from responder with SDP answer.
- `ice_candidate`: for exchanging ICE data.
- `disconnect`: notify backend and partner.
- `report`: send a report about current peer.
