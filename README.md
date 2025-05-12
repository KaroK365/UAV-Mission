Got it! Here's the revised **README** that focuses **only on the backend**, removing all frontend mentions:

---

# 🛰️ UAV Mission Planning – Backend (Go + Docker)

This is the backend component of a UAV (Unmanned Aerial Vehicle) Mission Planning System. It provides REST APIs to calculate total distance, fuel usage, and time required based on coordinates sent from a client application.

## 📌 Features

* 🧮 Calculate total mission distance using latitude and longitude waypoints
* ⛽ Estimate fuel consumption based on distance
* ⏱ Calculate total estimated mission time
* 🐳 Dockerized for portability and ease of deployment

---

## 🧠 Tech Stack

| Component | Technology  |
| --------- | ----------- |
| Language  | Go (Golang) |
| APIs      | REST        |
| Deploy    | Docker      |

---

Sure! Here's a **small and clean section** for the frontend preview you can include at the end of the README:

---

## 🌐 Frontend Preview

Although this repository only contains the backend, here’s a quick glimpse of how the frontend looks when integrated with this API:

| Mission Planner UI                                                                                          |
| ----------------------------------------------------------------------------------------------------------- |
| ![UAV Frontend Screenshot](https://github.com/user-attachments/assets/87a84f21-9788-4581-9ca8-ef9f233d8117) |
| ![Waypoint Selection](https://github.com/user-attachments/assets/23bb4e1a-1829-4535-9581-2230c4036ecc)      |
| ![Distance & Time Display](https://github.com/user-attachments/assets/c5c61f26-3a41-4058-93aa-bf99188f9bb3) |
| ![UI Components](https://github.com/user-attachments/assets/d53aacf2-ba97-4dcf-9ee9-85cefb5126f7)           |

> Frontend is developed separately using React and Leaflet.
> 🔗 Frontend : https://github.com/KaroK365/UAV-mission-frontend

## 🐳 Running with Docker

1. Clone the repo:

   ```bash
   git clone <repo-url>
   cd uav-mission-backend
   ```

2. Build the Docker image:

   ```bash
   docker build -t uav-backend .
   ```

3. Run the container:

   ```bash
   docker run -p 8080:8080 uav-backend
   ```

4. Access the API:

   * `http://localhost:8080/api/calculate`

---

## 🤖 AI Assistance Disclosure

> Boilerplate code and calculation logic were initially generated using AI tools. Final logic integration, Dockerization, and testing were done manually.

---

## 📄 License

This project is intended for educational purposes only.
