package main

import (
	"fmt"
	"log"

	"RayhanDB/pkg/client"
)

func main() {
	// Connect to RayhanDB server
	fmt.Println("🔌 Connecting to RayhanDB...")
	cfg := client.ClientConfig{
		Host:     "localhost",
		Port:     7171,
		Username: "admin",
		Password: "rayhan123",
	}

	c, err := client.NewClient(cfg)
	if err != nil {
		log.Fatalf("❌ Connection failed: %v", err)
	}
	defer c.Close()
	fmt.Println("✅ Connected!\n")

	// Demo 1: Create Products Table
	fmt.Println("📦 Creating 'products' table...")
	_, err = c.Exec(`
		RAYCREATE RAYTABLE products (
			id NUMBER,
			name TEXT,
			price NUMBER,
			stock NUMBER,
			category TEXT
		);
	`)
	if err != nil {
		log.Printf("⚠️  Table might already exist: %v\n", err)
	} else {
		fmt.Println("✅ Table created!\n")
	}

	// Demo 2: Insert Sample Data
	fmt.Println("➕ Inserting products...")
	products := []struct {
		id       int
		name     string
		price    int
		stock    int
		category string
	}{
		{1, "Laptop Gaming ROG", 25000000, 3, "Laptop"},
		{2, "MacBook Pro M3", 35000000, 2, "Laptop"},
		{3, "Mouse Logitech G Pro", 450000, 15, "Accessories"},
		{4, "Keyboard Keychron K2", 1200000, 8, "Accessories"},
		{5, "Monitor LG UltraGear", 5500000, 5, "Monitor"},
		{6, "Headset SteelSeries", 2500000, 10, "Accessories"},
	}

	for _, p := range products {
		query := fmt.Sprintf(
			"RAYERT RAYINTO products (id, name, price, stock, category) RAYVALUES (%d, '%s', %d, %d, '%s');",
			p.id, p.name, p.price, p.stock, p.category,
		)
		_, err = c.Exec(query)
		if err != nil {
			log.Printf("⚠️  Insert error: %v\n", err)
		} else {
			fmt.Printf("  ✅ Added: %s\n", p.name)
		}
	}
	fmt.Println()

	// Demo 3: Select All
	fmt.Println("📊 All Products:")
	rows, err := c.Query("RAYLECT * RAYFROM products;")
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}
	printRows(rows)

	// Demo 4: Filter by Price
	fmt.Println("\n💰 Products under 10 million:")
	rows, err = c.Query("RAYLECT name, price, stock RAYFROM products RAYWHERE price BELOW 10000000;")
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}
	printRows(rows)

	// Demo 5: Update Stock
	fmt.Println("\n📦 Updating stock for 'Mouse Logitech G Pro'...")
	_, err = c.Exec("RAYDATE products RAYSET stock = 20 RAYWHERE name = 'Mouse Logitech G Pro';")
	if err != nil {
		log.Fatalf("❌ Update failed: %v", err)
	}
	fmt.Println("✅ Stock updated!")

	// Demo 6: Verify Update
	fmt.Println("\n🔍 Checking updated stock:")
	rows, err = c.Query("RAYLECT name, stock RAYFROM products RAYWHERE name = 'Mouse Logitech G Pro';")
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}
	printRows(rows)

	// Demo 7: Delete Expensive Items
	fmt.Println("\n🗑️  Deleting products over 30 million...")
	_, err = c.Exec("RAYLETE RAYFROM products RAYWHERE price ABOVE 30000000;")
	if err != nil {
		log.Fatalf("❌ Delete failed: %v", err)
	}
	fmt.Println("✅ Deleted!")

	// Demo 8: Final Check
	fmt.Println("\n📊 Remaining Products:")
	rows, err = c.Query("RAYLECT * RAYFROM products;")
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}
	printRows(rows)

	// Demo 9: Show All Tables
	fmt.Println("\n📋 All Tables in Database:")
	rows, err = c.Query("RAYSHOW RAYTABLES;")
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}
	for _, row := range rows {
		for _, v := range row {
			fmt.Printf("  - %v\n", v)
		}
	}

	fmt.Println("\n✨ Demo completed successfully!")
}

func printRows(rows []map[string]interface{}) {
	if len(rows) == 0 {
		fmt.Println("  (No results)")
		return
	}

	for i, row := range rows {
		fmt.Printf("  %d. ", i+1)
		first := true
		for k, v := range row {
			if !first {
				fmt.Print(", ")
			}
			fmt.Printf("%s: %v", k, v)
			first = false
		}
		fmt.Println()
	}
}
