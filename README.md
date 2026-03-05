# RayhanDB 💝

**Your Friendly Query Language Database - Now with Server Mode!**

RayhanDB is a lightweight database with a unique query language that uses **FRIS**-prefixed keywords. Built with Go for performance and concurrent access support.

**✨ NEW: Server Mode** - Use RayhanDB as a database server like MongoDB or PostgreSQL!

## 🚀 Features

- 🎯 **Unique Rayhan Syntax**: Ray-prefixed keywords instead of SQL
- 🌐 **Server Mode**: TCP server with client-server architecture
- 📚 **Client Library**: Connect from any Go application
- 🔒 **Thread-Safe**: Concurrent access support with mutex locks
- 🔐 **Authentication**: Bcrypt password hashing
- 💾 **JSON Persistence**: Automatic save/load
- 🎨 **Pretty Output**: Colorized tables in REPL mode
- 📝 **Full CRUD**: Complete database operations

---

## 📦 Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/RayhanDB
cd RayhanDB

# Install dependencies
go get

# Build server
go build -o RayhanDB-server.exe cmd/RayhanDB-server/main.go

# Build REPL client
go build -o RayhanDB.exe cmd/RayhanDB/main.go

# Build example client
go build -o basic_client.exe examples/basic_client/main.go
```

---

## 🌐 Server Mode (NEW!)

### Starting the Server

```bash
# Start with defaults (port 7171, admin/rayhan123)
./RayhanDB-server

# Custom configuration
./RayhanDB-server -host 0.0.0.0 -port 8080 -user myuser -pass mypass -db production
```

**Server Options:**
- `-host` - Server host (default: `0.0.0.0`)
- `-port` - Server port (default: `7171`)
- `-db` - Database name (default: `mydb`)
- `-dir` - Data directory (default: `~/.RayhanDB`)
- `-user` - Admin username (default: `admin`)
- `-pass` - Admin password (default: `rayhan123`)

### Using Client Library

```go
package main

import (
    "fmt"
    "log"
    "RayhanDB/pkg/client"
)

func main() {
    // Connect to server
    db, err := client.Connect("localhost:7171", "admin", "rayhan123")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create table
    _, err = db.Exec(`
        RAYRATE RAYTABLE users (
            name TEXT,
            email TEXT,
            age NUMBER
        );
    `)

    // Insert data
    msg, err := db.Exec(`
        RAYERT RAYINTO users (name, email, age)
        RAYVALUES ('Rayhan', 'rayhan@db.com', 25);
    `)
    fmt.Println(msg) // "✅ Saved successfully! Total rows: 1"

    // Query data
    rows, err := db.Query("RAYLECT * RAYFROM users;")
    for _, row := range rows {
        fmt.Printf("User: %v\n", row)
    }
}
```

---

## 💻 REPL Mode (Interactive)

```bash
# Start interactive REPL
./RayhanDB

# REPL commands available:
# - help: Show help
# - save: Manually save database
# - exit/quit: Save and exit
```

**Important:** All queries must end with a semicolon `;` for multi-line support!

```Rayhan
rayhan> RAYRATE RAYTABLE users (
     ->     name TEXT,
     ->     age NUMBER
     -> );
✨ Table 'users' created successfully!

rayhan> RAYERT RAYINTO users (name, age) 
     -> RAYVALUES ('Rayhan', 25);
✅ Saved successfully! Total rows: 1

rayhan> RAYLECT * RAYFROM users;
🎉 Found 1 row(s)! Here they are:
...
```

---

## 🌟 Rayhan Query Language

### Keywords Dictionary

| SQL | Rayhan | Description |
|-----|--------|-------------|
| CREATE | RAYRATE | Create table |
| SELECT | RAYLECT | Select data |
| INSERT | RAYERT | Insert data |
| UPDATE | RAYDATE | Update data |
| DELETE | RAYLETE | Delete data |
| DROP | RAYDROP | Drop table |
| DESC | RAYC | Describe table |
| SHOW | RAYSHOW | Show tables |
| TABLE | RAYTABLE | Table keyword |
| FROM | RAYFROM | From clause |
| WHERE | RAYWHERE | Where clause |
| AND | RAYAND | Logical AND |
| OR | RAYOR | Logical OR |

### Operators

- **Comparison**: `ABOVE` (>), `BELOW` (<), `ATLEAST` (>=), `ATMOST` (<=)
- **Special**: `RAYXIST` (NOT NULL), `RAYMISS` (NULL), `RAYLOVE` (LIKE)
- **Logical**: `RAYAND` (AND), `RAYOR` (OR)

---

## 📚 Query Examples

### Create Table
```rayhan
RAYRATE RAYTABLE users (
    name TEXT,
    age NUMBER,
    email TEXT,
    active BOOLEAN
);
```

### Insert Data
```rayhan
RAYERT RAYINTO users (name, age, email, active) 
RAYVALUES ('Rayhan', 25, 'rayhan@example.com', true);
```

### Select Data
```rayhan
-- All columns
RAYLECT * RAYFROM users;

-- Specific columns
RAYLECT name, age RAYFROM users;

-- With conditions
RAYLECT * RAYFROM users RAYWHERE age ABOVE 18;

-- Pattern matching
RAYLECT * RAYFROM users RAYWHERE email RAYLOVE '%@gmail.com';

-- Multiple conditions
RAYLECT * RAYFROM users 
RAYWHERE age ABOVE 18 RAYAND active RAYXIST;
```

### Update Data
```rayhan
RAYDATE users RAYSET age = 26 RAYWHERE name = 'Rayhan';
```

### Delete Data
```rayhan
RAYLETE RAYFROM users RAYWHERE age BELOW 18;
```

### Other Commands
```rayhan
-- Describe table
RAYC users;

-- Show all tables
RAYSHOW RAYTABLES;

-- Drop table
RAYDROP RAYTABLE old_users;
```

---

## 🏗️ Architecture

```
RayhanDB/
├── cmd/
│   ├── RayhanDB/           # REPL client
│   └── RayhanDB-server/    # TCP server
├── pkg/
│   ├── core/               # Database engine
│   ├── parser/             # Query parser
│   ├── protocol/           # Network protocol
│   ├── auth/               # Authentication
│   ├── server/             # TCP server
│   └── client/             # Client library
└── examples/
    └── basic_client/       # Example usage
```

---

## 🔧 Technical Details

### Core Stack
- **Language**: Go 1.21+
- **Protocol**: JSON over TCP
- **Port**: 7171 (default)
- **Storage**: JSON files
- **Concurrency**: Thread-safe with `sync.RWMutex`

### Dependencies
- `github.com/fatih/color` - Colored output
- `github.com/olekukonko/tablewriter` - Pretty tables
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/google/uuid` - Request IDs

---

## 🎯 Use Cases

### Development & Prototyping
```go
// Quick database for your Go app
db, _ := RayhanDB.Connect("localhost:7171", "admin", "rayhan123")
defer db.Close()

// Use rayhan queries
rows, _ := db.Query("RAYLECT * RAYFROM products;")
```

### Learning Database Concepts
- Friendly syntax with FRIS keywords
- Interactive REPL for experimentation
- Clear error messages with emoji

### Small Applications
- Embedded database for Go apps
- Personal projects
- Microservices data store

---

## 📊 Performance

- **Concurrent Clients**: Supports multiple simultaneous connections
- **Thread-Safe**: All operations protected by mutexes
- **In-Memory**: Fast query execution
- **Persistent**: Auto-save to JSON on shutdown

---

## 🤝 Contributing

Contributions welcome! Feel free to submit issues and pull requests.

---

## 📝 License

MIT License - use however you'd like!

---

## 🙏 Acknowledgments

Built with ❤️ by the RayhanDB team. Special thanks to the Go community!

---

**Happy querying with Friska! 🎉**
