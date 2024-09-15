package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ghulamazad/websocket-optimization/internal/messaging"
	"github.com/ghulamazad/websocket-optimization/internal/session"
	"github.com/ghulamazad/websocket-optimization/internal/throttle"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

var redisThrottler = throttle.NewRedisClient("redis:6379")

// WebSocketHandler handles WebSocket connections and rate limiting.
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to upgrade connection:", err)
        return
    }
    defer conn.Close()

    clientIp := r.RemoteAddr
    sessionId := session.StoreSession(clientIp)

    // Set pong handler and read deadline
    conn.SetPongHandler(func(appData string) error {
        log.Println("Pong received from client:", clientIp)
        conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    // Start heartbeat mechanism
    go handlerHeartbeat(conn, 30*time.Second)

    for {
        conn.SetReadDeadline(time.Now().Add(60 * time.Second))

        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            conn.Close()
            break
        }

        // Check if the client is allowed to send messages based on rate limiting
        allowed := redisThrottler.AllowConnection(clientIp, 5, 10*time.Second)
        if !allowed {
            log.Println("Client rate limited:", clientIp)
            err = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","content":"rate_limited"}`))
            if err != nil {
                log.Println("Failed to send rate limited message:", err)
            }
            continue // Skip processing this message and wait for the next one
        }

        log.Printf("Message from client (%s): %s", clientIp, string(message))

        // Handle message prioritization via RabbitMQ
        err = messaging.PublishMessage(5, string(message))
        if err != nil {
            log.Println("Error publishing message:", err)
        }
    }

    // Remove session on disconnect
    session.RemoveSession(sessionId)
}

// handlerHeartbeat sends periodic ping messages to the client.
func handlerHeartbeat(conn *websocket.Conn, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                log.Println("Heartbeat failed:", err)
                conn.Close() // Close the connection on ping failure
                return
            }
        }
    }
}
