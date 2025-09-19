# Developer Foundry 2.0 – Team H - Slack Notification Integration

This project is a Slack Notification Service built to send real-time updates to Slack whenever users interact with tasks in the system.

---

## Features

* Task Creation Notification – Sends a Slack message when a new task is created.
* Task Update Notification – Sends a message when a task is updated.
* Task Status Change Notification – Notifies the team whenever the status of a task changes (e.g., In Progress → Completed).

---

## Tech Stack

* Go (Golang) – Implemented using the standard library
* PostgreSQL – Database for tasks and notifications
* Redis – Used for caching and queueing
* Koyeb – Hosting platform

---

## System Design

![System Design](./image/slack_notification_service.png)

### How It Works

1. **Client** sends a request to create, update, or change a task.
2. **Cache** stores temporary data for fast retrieval.
3. **Scheduler** manages timing and background processing of notifications.
4. **Database** persists all task information (task id, task name, assigned\_to, due\_at, created\_at, updated\_at).
5. **Redis Queue** holds notification jobs before they are sent.
6. **Slack Server** receives the final payload and delivers the message to the Slack channel.

---

## Environment Variables

The service requires the following environment variables to be set:

```env
DATABASE_URL=postgres://user:password@localhost:5432/slack_service
REDIS_URL=redis://localhost:6379
SLACK_BOT_TOKEN=xxxxx
PORT=8080
```

---

## API Endpoints

| Method | Endpoint      | Description                                     |
| ------ | ------------- | ----------------------------------------------- |
| POST   | `/tasks`      | Create a new task (triggers Slack notification) |
| PATCH  | `/tasks/{id}` | Update task details (triggers notification)     |

---

## Deployment

The service is hosted on **Koyeb**.

---

## 👨‍💻 Team Members
* **Oluwadarasimi Temitope Shina-kelani** – Backend Developer (Golang)
* **Stephen Basoah Dankyi** – Backend Developer (Golang)
