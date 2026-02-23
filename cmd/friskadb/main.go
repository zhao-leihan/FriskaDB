package main

import (
	"bufio"
	"fmt"
	"RayhanDB/pkg/core"
	"RayhanDB/pkg/parser"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var (
	colorGreen  = color.New(color.FgGreen, color.Bold)
	colorYellow = color.New(color.FgYellow, color.Bold)
	colorRed    = color.New(color.FgRed, color.Bold)
	colorCyan   = color.New(color.FgCyan, color.Bold)
	colorBlue   = color.New(color.FgBlue, color.Bold)
)

func main() {
	dbName := "mydb"
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".RayhanDB")

	// Parse simple args
	if len(os.Args) > 1 {
		dbName = os.Args[1]
	}

	storage := core.NewStorage(dataDir)
	db, err := storage.Load(dbName)
	if err != nil {
		db = core.NewDatabase(dbName)
	}

	executor := parser.NewExecutor(db)

	printWelcome()

	scanner := bufio.NewScanner(os.Stdin)
	running := true
	var queryBuffer strings.Builder
	inMultiLine := false

	for running {
		// Show different prompt for multi-line
		if inMultiLine {
			fmt.Print("     -> ")
		} else {
			fmt.Print("friska> ")
		}

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		// Build query buffer
		if queryBuffer.Len() > 0 {
			queryBuffer.WriteString(" ")
		}
		queryBuffer.WriteString(line)

		// Check if query ends with semicolon or is a special command
		queryText := strings.TrimSpace(queryBuffer.String())

		// Check for query terminator (semicolon) or single-line special commands
		isComplete := strings.HasSuffix(queryText, ";") ||
			strings.EqualFold(queryText, "exit") ||
			strings.EqualFold(queryText, "quit") ||
			strings.EqualFold(queryText, "help") ||
			strings.EqualFold(queryText, "save")

		// If query not complete, continue reading
		if !isComplete && queryText != "" {
			inMultiLine = true
			continue
		}

		// Reset for next query
		inMultiLine = false
		input := strings.TrimSuffix(queryText, ";") // Remove semicolon
		input = strings.TrimSpace(input)
		queryBuffer.Reset()

		if input == "" {
			continue
		}

		switch strings.ToLower(input) {
		case "exit", "quit":
			colorCyan.Println("\n💾 Saving database...")
			if err := storage.Save(db); err != nil {
				colorRed.Printf("❌ Error: %v\n", err)
			} else {
				colorGreen.Println("✅ Database saved successfully!")
			}
			colorCyan.Println("👋 Goodbye! See you soon!\n")
			running = false
			continue
		case "help":
			printHelp()
			continue
		case "save":
			if err := storage.Save(db); err != nil {
				colorRed.Printf("❌ Error: %v\n", err)
			} else {
				colorGreen.Println("✅ Database saved successfully!")
			}
			continue
		}

		// Parse and execute query
		p := parser.NewParser(input)
		query, err := p.Parse()
		if err != nil {
			colorRed.Printf("❌ Error: %v\n", err)
			continue
		}

		result, err := executor.Execute(query)
		if err != nil {
			colorRed.Printf("❌ Error: %v\n", err)
			continue
		}

		// Format output
		switch query.Type {
		case parser.QuerySelect:
			if rows, ok := result.([]core.Row); ok {
				printTable(rows)
			}
		case parser.QueryDescribe:
			if table, ok := result.(*core.Table); ok {
				printTableSchema(table)
			}
		case parser.QueryShowTables:
			if tables, ok := result.([]string); ok {
				printTableList(tables)
			}
		default:
			if msg, ok := result.(string); ok {
				colorGreen.Println(msg)
			}
		}
	}
}

func printWelcome() {
	banner := `
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║   ███████╗██████╗ ██╗███████╗██╗  ██╗ █████╗        ║
║   ██╔════╝██╔══██╗██║██╔════╝██║ ██╔╝██╔══██╗       ║
║   █████╗  ██████╔╝██║███████╗█████╔╝ ███████║       ║
║   ██╔══╝  ██╔══██╗██║╚════██║██╔═██╗ ██╔══██║       ║
║   ██║     ██║  ██║██║███████║██║  ██╗██║  ██║       ║
║   ╚═╝     ╚═╝  ╚═╝╚═╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝       ║
║                                                       ║
║        Your Friendly Query Language Database! 💝      ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝
`
	colorCyan.Println(banner)
	colorYellow.Println("Type 'help' for available commands, 'exit' to quit.\n")
}

func printHelp() {
	help := `
🌟 Friska Query Language Commands:

📝 CREATE TABLE:
   FRISRATE FRISKABLE table_name (col1 TYPE, col2 TYPE, ...)
   Example: FRISRATE FRISKABLE users (name TEXT, age NUMBER)

➕ INSERT:
   FRISERT FRISINTO table (col1, col2) FRISVALUES (val1, val2)
   Example: FRISERT FRISINTO users (name, age) FRISVALUES ('Rayhan', 25)

🔍 SELECT:
   FRISLECT columns FRISFROM table [FRISWHERE condition]
   Example: FRISLECT * FRISFROM users
   Example: FRISLECT name FRISFROM users FRISWHERE age ABOVE 18

✏️ UPDATE:
   FRISDATE table FRISSET col=val [FRISWHERE condition]
   Example: FRISDATE users FRISSET age=26 FRISWHERE name='Rayhan'

🗑️ DELETE:
   FRISLETE FRISFROM table [FRISWHERE condition]
   Example: FRISLETE FRISFROM users FRISWHERE age BELOW 18

💥 DROP TABLE:
   FRISDROP FRISKABLE table
   Example: FRISDROP FRISKABLE old_users

📋 DESCRIBE:
   FRISC table
   Example: FRISC users

📚 SHOW TABLES:
   FRISSHOW FRISKABLES

🎯 Operators:
   =, !=, ABOVE (>), BELOW (<), ATLEAST (>=), ATMOST (<=)
   FRISLOVE (LIKE), FRISAMONG (IN)
   FRISXIST (NOT NULL), FRISMISS (NULL)
   FRISAND (AND), FRISOR (OR)

💡 Special Commands:
   help  - Show this help
   exit  - Exit RayhanDB
   save  - Save database to disk
`
	colorBlue.Println(help)
}

func printTable(rows []core.Row) {
	if len(rows) == 0 {
		colorYellow.Println("😅 No rows found!")
		return
	}

	var columns []string
	for col := range rows[0] {
		columns = append(columns, col)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)

	for _, row := range rows {
		var rowData []string
		for _, col := range columns {
			rowData = append(rowData, fmt.Sprintf("%v", row[col]))
		}
		table.Append(rowData)
	}

	colorGreen.Printf("\n🎉 Found %d row(s)! Here they are:\n\n", len(rows))
	table.Render()
	fmt.Println()
}

func printTableSchema(tbl *core.Table) {
	colorCyan.Printf("\n📋 Table: %s\n\n", tbl.Name)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Column", "Type"})

	for _, col := range tbl.Schema.Columns {
		table.Append([]string{col.Name, string(col.Type)})
	}

	table.Render()
	fmt.Printf("\nTotal rows: %d\n\n", tbl.Count())
}

func printTableList(tables []string) {
	if len(tables) == 0 {
		colorYellow.Println("😅 No tables found! Create one with FRISRATE FRISKABLE.")
		return
	}

	colorCyan.Printf("\n📚 Available Tables (%d):\n\n", len(tables))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Table Name"})

	for i, name := range tables {
		table.Append([]string{fmt.Sprintf("%d", i+1), name})
	}

	table.Render()
	fmt.Println()
}
