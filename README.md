# Developer Foundry 2.0 – Team H - Slack Notification integration

## Project Overview

This repository contains the backend API for **Team H** of the Developer Foundry 2.0 Bootcamp.
The project is written purely in **Golang** and consumes events on a Redis queue and then publishes events onto a Slack channel:

---

## Features

*  **Consume Events** – Consume change events from placed on Redis queues
*  **Notify Slack Channels** – Publish change event into a Slack channel
---

## 🛠️ Tech Stack

* **Backend Framework/API Layer**: Golang stdlib
* **Database**: PostgreSQL (recommended) / SQLite (development)
---

## 👨‍💻 Team Members
* **Oluwadarasimi Temitope Shina-kelani** – Backend Developer (Golang)
* **Stephen Basoah Dankyi** – Backend Developer (Golang)

---

## ⚙️ Setup Instructions

### 1️⃣ Clone the Repository

```bash
git clone git@github.com:Developer-s-Foundry/DF.2.0-task-mgt-authentication.git
cd DF.2.0-task-mgt-authentication
```

### 2️⃣ Import Packages

```bash
go mod vendor
go mod tidy
```

### 3️⃣ Build Repository

```bash
go build
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Create a Pull Request


---
