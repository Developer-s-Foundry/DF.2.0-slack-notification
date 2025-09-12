package seed

import (
	"fmt"
	"log"
	"time"

	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/repository/postgres"
	"github.com/Developer-s-Foundry/DF.2.0-slack-notification/utils"
)

func data() []postgres.Task {
	dummyTasks := []postgres.Task{
		{
			ID:          utils.Uuid(),
			Name:        "Write project proposal",
			Status:      "pending",
			Description: "Draft and submit the initial project proposal document for review.",
			AssignedTo:  "Alice",
			Expires_at:  time.Now().Add(24 * time.Hour),
		},
		{
			ID:          utils.Uuid(),
			Name:        "Design database schema",
			Status:      "in_progress",
			Description: "Create ER diagrams and define the PostgreSQL schema for the application.",
			AssignedTo:  "Bob",
			Expires_at:  time.Now().Add(48 * time.Hour),
		},
		{
			ID:          utils.Uuid(),
			Name:        "Implement authentication",
			Status:      "completed",
			Description: "Add user authentication with JWT and password hashing.",
			AssignedTo:  "Charlie",
			Expires_at:  time.Now().Add(72 * time.Hour),
		},
		{
			ID:          utils.Uuid(),
			Name:        "Integrate payment gateway",
			Status:      "pending",
			Description: "Set up Paystack API integration to handle user payments securely.",
			AssignedTo:  "David",
			Expires_at:  time.Now().Add(36 * time.Hour),
		},
		{
			ID:          utils.Uuid(),
			Name:        "Deploy to staging",
			Status:      "pending",
			Description: "Deploy the latest build to the staging environment on Koyeb for QA testing.",
			AssignedTo:  "Eve",
			Expires_at:  time.Now().Add(60 * time.Hour),
		},
	}

	return dummyTasks
}

func SeedTasks(p *postgres.PostgresConn) error {
	seeded, err := p.CheckDataExists()

	if err != nil {
		return err
	}

	if seeded {
		log.Println("no data seeding performed data already exist in DB")
		return nil
	}
	for _, task := range data() {
		if err := p.Insert(&task); err != nil {
			return fmt.Errorf("failed to insert task %s: %w", task.Name, err)
		}
	}

	log.Printf("Seeded %d dummy tasks successfully", len(data()))
	return nil
}
