package main

import (
	"fmt"
	"log"

	"RayhanDB/pkg/client"
)

func main() {
	// Connect to RayhanDB server
	fmt.Println("🔌 Connecting to RayhanDB server...")
	db, err := client.Connect("localhost:7171", "admin", "rayhan123")
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer db.Close()

	fmt.Println("✅ Connected successfully!\n")

	// Create table
	fmt.Println("📝 Creating table...")
	msg, err := db.Exec(`
		RAYRATE RAYTABLE users (
			name TEXT,
			email TEXT,
			age NUMBER
		);
	`)
	if err != nil {
		log.Printf("Error creating table: %v", err)
	} else {
		fmt.Printf("✨ %s\n\n", msg)
	}

	// Insert data
	fmt.Println("➕ Inserting data...")
	msg, err = db.Exec(`
		RAYERT RAYINTO users (name, email, age)
		RAYVALUES ('Rayhan', 'rayhan@db.com', 25);
	`)
	if err != nil {
		log.Printf("Error inserting: %v", err)
	} else {
		fmt.Printf("✅ %s\n\n", msg)
	}

	msg, err = db.Exec(`
		RAYERT RAYINTO users (name, email, age)
		RAYVALUES ('Alice', 'alice@example.com', 30);
	`)
	if err != nil {
		log.Printf("Error inserting: %v", err)
	} else {
		fmt.Printf("✅ %s\n\n", msg)
	}

	// Query data
	fmt.Println("🔍 Querying data...")
	rows, err := db.Query("RAYLECT * RAYFROM users;")
	if err != nil {
		log.Fatalf("Error querying: %v", err)
	}

	fmt.Printf("🎉 Found %d row(s):\n", len(rows))
	for i, row := range rows {
		fmt.Printf("  %d. Name: %v, Email: %v, Age: %v\n",
			i+1, row["name"], row["email"], row["age"])
	}

	fmt.Println("\n💝 Demo completed successfully!")
}
