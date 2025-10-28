# 🚖 Saarathi — A Ride-Sharing Platform Built with Go Microservices
Saarathi is a distributed ride-sharing platform designed with scalability, event-driven communication, and clean architecture in mind.  
It simulates the full lifecycle of a ride request — from fare estimation to trip completion — powered by modern backend technologies.


## 🧩 Architecture Overview

Saarathi is composed of multiple independent **Go microservices**, communicating via **NATS JetStream** and following a **Hexagonal Architecture** pattern.

### Core Services

- **🧍 Users Service**  
  Manages all user-related data (drivers and riders). Handles registration, authentication, and profile management.

- **🚗 Trips Service**  
  Handles the **ride lifecycle**:
  - Generates fare estimates using **OSRM (Open Source Routing Machine)**.  
  - Creates new trips with `PENDING` state once a rider accepts the price.  
  - Listens for offer confirmations to update trip status.

- **📡 Rider Service**  
  Establishes **Server-Sent Events (SSE)** connections to deliver real-time trip updates to riders.  
  Riders wait for responses after accepting fare estimates.

- **🎯 Offers Service**  
  Listens for `TripCreated` events from the Trips Service and begins **driver matching**:
  - Sends match requests via events to the **RMS Service** to find nearby drivers.  
  - Manages offer locks and expirations using **PostgreSQL** as the source of truth (Redis as cache).  
  - Publishes offer events to the Driver-State Service.

- **🗺️ RMS (Ride Matching Service)**  
  Provides the list of **nearest drivers** based on geolocation data from Redis.

- **⚡ Driver-State Service**  
  Maintains **WebSocket** connections with all online drivers:
  - Receives offer events from Offers Service.  
  - Forwards them to the correct driver’s WebSocket connection.  
  - On driver acceptance, sends confirmation back to Offers Service.

## 🧠 Technical Highlights

- **Language:** Go (Golang)  
- **Message Broker:** NATS JetStream  
- **Database:** PostgreSQL  
- **Caching & Realtime Data:** Redis  
- **External Routing:** OSRM  
- **Protocols:** gRPC, WebSockets, SSE  
- **Containerization:** Kubernetes & Docker Compose  

## 🚀 Future Improvements
- Cheorographed SAGA is yet to be implemented with compensating event 
- Improved fault tolerance and retry mechanisms.  

**Built with ❤️ and Go.**

