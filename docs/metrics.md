# Metrics

To evaluate whether the MVP meets its objectives and performs reliably, the following methods will be used:

- **Session Completion Rate**  
  Track the number of successful pairings and completed sessions (not immediately disconnected). A high bounce rate may indicate issues in video connection or poor UX.

- **Session Duration Metrics**  
  Use the `sessions` table to analyze average session times. Helps understand user engagement and session quality.

- **Error Logging & Monitoring**  
  Review logs from the signaling and matchmaking services for unexpected disconnects, matching delays, or signaling failures.

- **WebRTC Failure Rate**  
  Track how many connections fail to complete the handshake. Helps in debugging ICE/STUN issues.

- **Moderation Reports**  
  Volume and frequency of reports give insight into user behavior and potential abuse hotspots.

- **User Feedback (Optional)**  
  Lightweight post-session thumbs-up/thumbs-down can help collect subjective satisfaction data.

- **Performance Metrics**

  - Redis queue latency
  - WebSocket throughput
  - DB write latency (sessions, reports)