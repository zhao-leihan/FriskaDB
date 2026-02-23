# RayhanDB Server Mode - Quick Start Guide

## 🚀 Running the Server

### 1. Start Server
Open first terminal:
```bash
cd c:\Users\Rayhan\Music\RayhanDB
.\RayhanDB-server.exe
```

Expected output:
```
🚀 RayhanDB Server started on 0.0.0.0:7171
💝 Ready to accept connections!
```

### 2. Run Example Client
Open second terminal:
```bash
cd c:\Users\Rayhan\Music\RayhanDB
.\basic_client.exe
```

Expected output:
```
🔌 Connecting to RayhanDB server...
✅ Connected successfully!

📝 Creating table...
✨ Table 'users' created successfully!

➕ Inserting data...
✅ Saved successfully! Total rows: 1
✅ Saved successfully! Total rows: 2

🔍 Querying data...
🎉 Found 2 row(s):
  1. Name: Friska, Email: rayhan@db.com, Age: 25
  2. Name: Alice, Email: alice@example.com, Age: 30

💝 Demo completed successfully!
```

---

## 📚 Using in Your Go Code

### Example Application

```go
package main

import (
    "fmt"
    "log"
    "RayhanDB/pkg/client"
)

func main() {
    // Connect
    db, err := client.Connect("localhost:7171", "admin", "rayhan123")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create products table
    db.Exec(`
        RAYRATE RAYTABLE products (
            id NUMBER,
            name TEXT,
            price NUMBER,
            stock NUMBER
        );
    `)

    // Insert products
    db.Exec(`RAYERT RAYINTO products (id, name, price, stock) 
             RAYVALUES (1, 'Laptop', 15000000, 10);`)
    
    db.Exec(`RAYERT RAYINTO products (id, name, price, stock) 
             RAYVALUES (2, 'Mouse', 250000, 50);`)

    // Query in-stock products
    rows, _ := db.Query(`
        RAYLECT name, price, stock 
        RAYFROM products 
        RAYWHERE stock ABOVE 0;
    `)

    for _, row := range rows {
        fmt.Printf("%v - Rp %v (Stock: %v)\n", 
            row["name"], row["price"], row["stock"])
    }
}
```

---

## 🔧 Server Configuration

### CLI Flags

```bash
./RayhanDB-server -h

# Common options:
-host string    Server host (default "0.0.0.0")
-port int       Server port (default 7171)
-db string      Database name (default "mydb")
-dir string     Data directory (default "~/.RayhanDB")
-user string    Admin username (default "admin")
-pass string    Admin password (default "rayhan123")
```

### Custom Server

```bash
# Production setup
./RayhanDB-server \
  -host 0.0.0.0 \
  -port 8080 \
  -db production \
  -user prodadmin \
  -pass SecurePass123 \
  -dir /var/lib/RayhanDB
```

---

## 🌐 Client Library API

### Connection
```go
db, err := client.Connect(address, username, password)
defer db.Close()
```

### Execute Query (SELECT)
```go
rows, err := db.Query("RAYLECT * RAYFROM users;")
// Returns: []core.Row
```

### Execute Command (INSERT, UPDATE, DELETE, CREATE, DROP)
```go
msg, err := db.Exec("RAYERT RAYINTO users (name) RAYVALUES ('Alice');")
// Returns: success message string
```

---

## 📊 Database Operations

### Tables
```go
// Create
db.Exec(`RAYRATE RAYTABLE orders (id NUMBER, total NUMBER);`)

// Describe
db.Query(`RAYC orders;`)

// List all
db.Query(`RAYSHOW RAYTABLES;`)

// Drop
db.Exec(`RAYDROP RAYTABLE old_table;`)
```

### Data Manipulation
```go
// Insert
db.Exec(`RAYERT RAYINTO orders (id, total) RAYVALUES (1, 150000);`)

// Update
db.Exec(`RAYDATE orders RAYSET total=200000 RAYWHERE id=1;`)

// Delete
db.Exec(`RAYLETE RAYFROM orders RAYWHERE total BELOW 100000;`)

// Select
rows, _ := db.Query(`RAYLECT * RAYFROM orders;`)
```

---

## ⚡ Advanced Usage

### Concurrent Clients

RayhanDB supports multiple concurrent connections:

```go
// Client 1
go func() {
    db1, _ := client.Connect("localhost:7171", "admin", "rayhan123")
    defer db1.Close()
    db1.Exec(`RAYERT RAYINTO logs (msg) RAYVALUES ('Client 1');`)
}()

// Client 2
go func() {
    db2, _ := client.Connect("localhost:7171", "admin", "rayhan123")
    defer db2.Close()
    db2.Exec(`RAYERT RAYINTO logs (msg) RAYVALUES ('Client 2');`)
}()
```

### Error Handling

```go
msg, err := db.Exec(`RAYERT RAYINTO nonexistent (x) RAYVALUES (1);`)
if err != nil {
    // Handle error: "query failed: table 'nonexistent' not found 😢"
    fmt.Println(err)
}
```

---

## 🎯 Tips & Best Practices

1. **Always close connections**: Use `defer db.Close()`
2. **Handle errors**: Check all returned errors
3. **Use semicolons**: End all queries with `;`
4. **Authentication**: Change default password in production
5. **Concurrent-safe**: Server handles multiple clients automatically

---

## 🐛 Troubleshooting

### Connection refused
```
Error: failed to connect: dial tcp 127.0.0.1:7171: connection refused
```
**Solution**: Make sure server is running with `./RayhanDB-server`

### Authentication failed
```
Error: authentication failed: invalid username or password
```
**Solution**: Check username/password match server configuration

### Query syntax error
```
Error: query failed: unexpected token: ...
```
**Solution**: Check Friska syntax, ensure semicolon at end

---

**Enjoy using RayhanDB Server! 💝**
