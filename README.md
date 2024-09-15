# WebSocket Optimization Project

## Features
WebSocket Server: Handles WebSocket connections with support for rate limiting, heartbeat, and session management.
Rate Limiting: Uses Redis to control the rate of messages per client.
Message Prioritization: Publishes messages to RabbitMQ for further processing.
Monitoring: Integrated with Prometheus and Grafana for real-time monitoring and visualization.
Load Testing: Includes a k6 script for performance testing.


## Prerequisites
Docker and Docker Compose installed.
Basic knowledge of WebSockets, Redis, RabbitMQ, and Grafana.


# Getting Started
1. Clone the Repository
```bash
git clone https://github.com/yourusername/websocket-optimization.git
cd websocket-optimization
```
2. Set Up Docker Containers

Run the following command to start all services:
```bash
docker-compose up --build
```
or 
```bash
docker compose up --build --scale websocket-app=3
```
This command will start the WebSocket server, Redis, RabbitMQ, Nginx, InfluxDB, Grafana, and k6.

3. Configure Nginx

Ensure your nginx.conf file is properly configured to proxy WebSocket connections. Example configuration:
```
server {
    listen 8080;

    location /ws {
        proxy_pass http://websocket-app:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```
5. Monitoring
## Grafana
Grafana URL: http://localhost:3000
Default Credentials: Username: admin, Password: admin

### Add Prometheus Data Source:
Go to Configuration -> Data Sources -> Add Data Source.
Select InfluxDB and set the URL to http://InfluxDB:9090.


6. Accessing Services
- WebSocket Server: ws://localhost:8080/ws
- Nginx: http://localhost:8080
- Grafana: http://localhost:3000
- InfluxDB: http://localhost:8086


# License
This project is licensed under the MIT License
