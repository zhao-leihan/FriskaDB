# FriskaDB 💝

**Your Friendly Query Language Database - Now with Server Mode!**

FriskaDB is a lightweight database with a unique query language that uses **FRIS**-prefixed keywords. Built with Go for performance and concurrent access support.

**✨ NEW: Server Mode** - Use FriskaDB as a database server like MongoDB or PostgreSQL!

## 🚀 Features

- 🎯 **Unique Friska Syntax**: FRIS-prefixed keywords instead of SQL
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
git clone https://github.com/yourusername/friskadb
cd FriskaDB

# Install dependencies
go get

# Build server
go build -o friskadb-server.exe cmd/friskadb-server/main.go

# Build REPL client
go build -o friskadb.exe cmd/friskadb/main.go

# Build example client
go build -o basic_client.exe examples/basic_client/main.go
```

---

## 🌐 Server Mode (NEW!)

### Starting the Server

```bash
# Start with defaults (port 7171, admin/friska123)
./friskadb-server

# Custom configuration
./friskadb-server -host 0.0.0.0 -port 8080 -user myuser -pass mypass -db production
```

**Server Options:**
- `-host` - Server host (default: `0.0.0.0`)
- `-port` - Server port (default: `7171`)
- `-db` - Database name (default: `mydb`)
- `-dir` - Data directory (default: `~/.friskadb`)
- `-user` - Admin username (default: `admin`)
- `-pass` - Admin password (default: `friska123`)

### Using Client Library

```go
package main

import (
    "fmt"
    "log"
    "friskadb/pkg/client"
)

func main() {
    // Connect to server
    db, err := client.Connect("localhost:7171", "admin", "friska123")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create table
    _, err = db.Exec(`
        FRISRATE FRISKABLE users (
            name TEXT,
            email TEXT,
            age NUMBER
        );
    `)

    // Insert data
    msg, err := db.Exec(`
        FRISERT FRISINTO users (name, email, age)
        FRISVALUES ('Friska', 'friska@db.com', 25);
    `)
    fmt.Println(msg) // "✅ Saved successfully! Total rows: 1"

    // Query data
    rows, err := db.Query("FRISLECT * FRISFROM users;")
    for _, row := range rows {
        fmt.Printf("User: %v\n", row)
    }
}
```

---

## 💻 REPL Mode (Interactive)

```bash
# Start interactive REPL
./friskadb

# REPL commands available:
# - help: Show help
# - save: Manually save database
# - exit/quit: Save and exit
```

**Important:** All queries must end with a semicolon `;` for multi-line support!

```friska
friska> FRISRATE FRISKABLE users (
     ->     name TEXT,
     ->     age NUMBER
     -> );
✨ Table 'users' created successfully!

friska> FRISERT FRISINTO users (name, age) 
     -> FRISVALUES ('Friska', 25);
✅ Saved successfully! Total rows: 1

friska> FRISLECT * FRISFROM users;
🎉 Found 1 row(s)! Here they are:
...
```

---

## 🌟 Friska Query Language

### Keywords Dictionary

| SQL | Friska | Description |
|-----|--------|-------------|
| CREATE | FRISRATE | Create table |
| SELECT | FRISLECT | Select data |
| INSERT | FRISERT | Insert data |
| UPDATE | FRISDATE | Update data |
| DELETE | FRISLETE | Delete data |
| DROP | FRISDROP | Drop table |
| DESC | FRISC | Describe table |
| SHOW | FRISSHOW | Show tables |
| TABLE | FRISKABLE | Table keyword |
| FROM | FRISFROM | From clause |
| WHERE | FRISWHERE | Where clause |
| AND | FRISAND | Logical AND |
| OR | FRISOR | Logical OR |

### Operators

- **Comparison**: `ABOVE` (>), `BELOW` (<), `ATLEAST` (>=), `ATMOST` (<=)
- **Special**: `FRISXIST` (NOT NULL), `FRISMISS` (NULL), `FRISLOVE` (LIKE)
- **Logical**: `FRISAND` (AND), `FRISOR` (OR)

---

## 📚 Query Examples

### Create Table
```friska
FRISRATE FRISKABLE users (
    name TEXT,
    age NUMBER,
    email TEXT,
    active BOOLEAN
);
```

### Insert Data
```friska
FRISERT FRISINTO users (name, age, email, active) 
FRISVALUES ('Friska', 25, 'friska@example.com', true);
```

### Select Data
```friska
-- All columns
FRISLECT * FRISFROM users;

-- Specific columns
FRISLECT name, age FRISFROM users;

-- With conditions
FRISLECT * FRISFROM users FRISWHERE age ABOVE 18;

-- Pattern matching
FRISLECT * FRISFROM users FRISWHERE email FRISLOVE '%@gmail.com';

-- Multiple conditions
FRISLECT * FRISFROM users 
FRISWHERE age ABOVE 18 FRISAND active FRISXIST;
```

### Update Data
```friska
FRISDATE users FRISSET age = 26 FRISWHERE name = 'Friska';
```

### Delete Data
```friska
FRISLETE FRISFROM users FRISWHERE age BELOW 18;
```

### Other Commands
```friska
-- Describe table
FRISC users;

-- Show all tables
FRISSHOW FRISKABLES;

-- Drop table
FRISDROP FRISKABLE old_users;
```

---

## 🏗️ Architecture

```
FriskaDB/
├── cmd/
│   ├── friskadb/           # REPL client
│   └── friskadb-server/    # TCP server
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
db, _ := friskadb.Connect("localhost:7171", "admin", "friska123")
defer db.Close()

// Use Friska queries
rows, _ := db.Query("FRISLECT * FRISFROM products;")
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

Built with ❤️ by the FriskaDB team. Special thanks to the Go community!

---

**Happy querying with Friska! 🎉**
