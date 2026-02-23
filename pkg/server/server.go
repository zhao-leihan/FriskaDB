package server

import (
	"bufio"
	"fmt"
	"RayhanDB/pkg/auth"
	"RayhanDB/pkg/core"
	"RayhanDB/pkg/parser"
	"RayhanDB/pkg/protocol"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

// Server represents the RayhanDB TCP server
type Server struct {
	host          string
	port          int
	db            *core.Database
	storage       *core.Storage
	executor      *parser.Executor
	authenticator *auth.Authenticator
	listener      net.Listener
	running       atomic.Bool
	connections   sync.WaitGroup
	mu            sync.Mutex
}

// Config holds server configuration
type Config struct {
	Host          string
	Port          int
	DatabaseName  string
	DataDir       string
	AdminUser     string
	AdminPassword string
}

// NewServer creates a new RayhanDB server
func NewServer(cfg *Config) (*Server, error) {
	// Initialize storage
	storage := core.NewStorage(cfg.DataDir)

	// Load or create database
	db, err := storage.Load(cfg.DatabaseName)
	if err != nil {
		db = core.NewDatabase(cfg.DatabaseName)
	}

	// Initialize authenticator
	authenticator := auth.NewAuthenticator()
	if err := authenticator.AddUser(cfg.AdminUser, cfg.AdminPassword); err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	return &Server{
		host:          cfg.Host,
		port:          cfg.Port,
		db:            db,
		storage:       storage,
		executor:      parser.NewExecutor(db),
		authenticator: authenticator,
	}, nil
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.listener = listener
	s.running.Store(true)

	log.Printf("🚀 RayhanDB Server started on %s", addr)
	log.Printf("💝 Ready to accept connections!")

	// Accept connections
	for s.running.Load() {
		conn, err := listener.Accept()
		if err != nil {
			if !s.running.Load() {
				break
			}
			log.Printf("❌ Error accepting connection: %v", err)
			continue
		}

		s.connections.Add(1)
		go s.handleConnection(conn)
	}

	return nil
}

// Stop stops the server gracefully
func (s *Server) Stop() error {
	s.running.Store(false)

	if s.listener != nil {
		s.listener.Close()
	}

	// Wait for all connections to finish
	s.connections.Wait()

	// Save database
	if err := s.storage.Save(s.db); err != nil {
		return fmt.Errorf("failed to save database: %w", err)
	}

	log.Println("👋 Server stopped successfully")
	return nil
}

// handleConnection handles a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer s.connections.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("✅ New connection from %s", clientAddr)

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	authenticated := false

	for scanner.Scan() {
		line := scanner.Bytes()

		// Decode request
		req, err := protocol.DecodeRequest(line)
		if err != nil {
			s.sendError(writer, "", fmt.Errorf("invalid request format: %w", err))
			continue
		}

		// Authenticate if not already authenticated
		if !authenticated {
			if err := s.authenticator.Authenticate(req.Auth.Username, req.Auth.Password); err != nil {
				s.sendError(writer, req.ID, fmt.Errorf("authentication failed: %w", err))
				conn.Close()
				return
			}
			authenticated = true
			log.Printf("🔐 User '%s' authenticated from %s", req.Auth.Username, clientAddr)
		}

		// Execute query
		s.executeQuery(writer, req)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("❌ Connection error from %s: %v", clientAddr, err)
	}

	log.Printf("👋 Connection closed from %s", clientAddr)
}

// executeQuery executes a Friska query and sends response
func (s *Server) executeQuery(writer *bufio.Writer, req *protocol.Request) {
	// Special handling for FRISREGISTER command (user registration)
	if strings.HasPrefix(strings.TrimSpace(req.Query), "FRISREGISTER") {
		s.handleRegistration(writer, req)
		return
	}

	// Parse query
	p := parser.NewParser(req.Query)
	query, err := p.Parse()
	if err != nil {
		s.sendError(writer, req.ID, err)
		return
	}

	// Execute query
	result, err := s.executor.Execute(query)
	if err != nil {
		s.sendError(writer, req.ID, err)
		return
	}

	// Format response based on query type
	var resp *protocol.Response

	switch query.Type {
	case parser.QuerySelect:
		if rows, ok := result.([]core.Row); ok {
			resp = protocol.NewSuccessResponse(req.ID, protocol.ConvertRows(rows), fmt.Sprintf("Found %d row(s)", len(rows)))
		}
	case parser.QueryDescribe:
		if table, ok := result.(*core.Table); ok {
			resp = protocol.NewSuccessResponse(req.ID, protocol.ConvertTableInfo(table), "Table schema")
		}
	case parser.QueryShowTables:
		if tables, ok := result.([]string); ok {
			resp = protocol.NewSuccessResponse(req.ID, tables, fmt.Sprintf("Found %d table(s)", len(tables)))
		}
	default:
		// For CREATE, INSERT, UPDATE, DELETE, DROP
		if msg, ok := result.(string); ok {
			resp = protocol.NewSuccessResponse(req.ID, nil, msg)
		}
	}

	if resp == nil {
		s.sendError(writer, req.ID, fmt.Errorf("unexpected result type"))
		return
	}

	s.sendResponse(writer, resp)
}

// handleRegistration handles user registration requests
func (s *Server) handleRegistration(writer *bufio.Writer, req *protocol.Request) {
	// Parse FRISREGISTER user:username pass:password
	parts := strings.Fields(req.Query)
	if len(parts) != 3 {
		s.sendError(writer, req.ID, fmt.Errorf("invalid registration format, expected: FRISREGISTER user:username pass:password"))
		return
	}

	var username, password string
	for _, part := range parts[1:] {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "user":
			username = kv[1]
		case "pass":
			password = kv[1]
		}
	}

	if username == "" || password == "" {
		s.sendError(writer, req.ID, fmt.Errorf("missing username or password"))
		return
	}

	// Check if user already exists
	if s.authenticator.UserExists(username) {
		s.sendError(writer, req.ID, fmt.Errorf("user '%s' already exists", username))
		return
	}

	// Add user
	if err := s.authenticator.AddUser(username, password); err != nil {
		s.sendError(writer, req.ID, fmt.Errorf("failed to register user: %w", err))
		return
	}

	log.Printf("👤 New user '%s' registered", username)
	resp := protocol.NewSuccessResponse(req.ID, nil, fmt.Sprintf("User '%s' registered successfully", username))
	s.sendResponse(writer, resp)
}

// sendResponse sends a response to the client
func (s *Server) sendResponse(writer *bufio.Writer, resp *protocol.Response) {
	data, err := protocol.EncodeResponse(resp)
	if err != nil {
		log.Printf("❌ Failed to encode response: %v", err)
		return
	}

	writer.Write(data)
	writer.WriteByte('\n')
	writer.Flush()
}

// sendError sends an error response to the client
func (s *Server) sendError(writer *bufio.Writer, requestID string, err error) {
	resp := protocol.NewErrorResponse(requestID, err)
	s.sendResponse(writer, resp)
}

// AddUser adds a new user to the server
func (s *Server) AddUser(username, password string) error {
	return s.authenticator.AddUser(username, password)
}

// SaveDatabase manually saves the database
func (s *Server) SaveDatabase() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage.Save(s.db)
}
