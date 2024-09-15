import { check, sleep } from "k6";
import { Counter, Trend } from "k6/metrics";
import ws from "k6/ws";

// Define custom metrics for performance tracking
let messageDuration = new Trend("websocket_message_duration");
let messageCount = new Counter("websocket_message_count");
let connectionErrors = new Counter("websocket_connection_errors");
let rateLimitedUsers = new Counter("websocket_rate_limited_users");

export let options = {
  stages: [
    { duration: "1m", target: 50 }, // Ramp-up to 50 users over 1 minute
    { duration: "3m", target: 50 }, // Stay at 50 users for 3 minutes
    { duration: "1m", target: 100 }, // Ramp-up to 100 users
    { duration: "3m", target: 100 }, // Stay at 100 users for 3 minutes
    { duration: "1m", target: 0 }, // Ramp-down to 0 users
  ],
};

export default function () {
  let url = "ws://nginx:8080/ws"; // WebSocket URL for Nginx
  let params = { tags: { my_tag: "load_test" } };

  const res = ws.connect(url, params, function (socket) {
    socket.on("open", function () {
      console.log("Connected to WebSocket server");

      let start = new Date();
      // Simulate periodic message sending every second for 5 seconds
      for (let i = 0; i < 10; i++) {
        socket.send(
          JSON.stringify({ type: "message", content: `Hello Server! (message ${i + 1})` })
        );
        let responseTime = new Date() - start;
        messageDuration.add(responseTime); // Track message round-trip duration
        messageCount.add(1); // Increment message count metric

        // Pause for 1 second between messages
        sleep(1);
      }

      socket.on("message", function (data) {
        let messageData = JSON.parse(data);
        check(messageData, {
          "received message": (msg) => msg.content === "Hello Client!",
        });

        // Check for rate limiting
        if (messageData.content === "rate_limited") {
          console.warn("User rate-limited");
          rateLimitedUsers.add(1); // Track rate-limited users
        }

        console.log("Message received: ", data);
      });

      socket.setTimeout(function () {
        console.log("Closing connection");
        socket.close();
      }, 5000); // Close connection after 5 seconds
    });

    socket.on("close", function () {
      console.log("Connection closed");
    });

    socket.on("error", function (e) {
      console.error("WebSocket connection error: ", e.error());
      connectionErrors.add(1); // Track connection errors
    });
  });

  check(res, { "status is 101": (r) => r && r.status === 101 }); // Ensure successful connection (status code 101)
  sleep(1);
}
