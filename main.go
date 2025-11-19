package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	_ "top20-back/docs"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Top20 API
// @version 1.0
// @description API para gerenciar submissões de top 20 jogadores
// @contact.name API Support
// @contact.email support@top20.com
// @host localhost:3000
// @BasePath /api

type Player struct {
	Position int    `json:"position"`
	Name     string `json:"name"`
}

type Submission struct {
	ID          int      `json:"id"`
	Players     []Player `json:"players"`
	SubmittedBy string   `json:"submittedBy"`
	IPAddress   string   `json:"ipAddress"`
	CreatedAt   string   `json:"createdAt"`
}

// PlayerStats representa as estatísticas de um jogador
type PlayerStats struct {
	PlayerName       string             `json:"playerName"`
	TotalSubmissions int                `json:"totalSubmissions"`
	PositionBreakdown []PositionCount   `json:"positionBreakdown"`
}

// PositionCount representa quantas vezes um jogador apareceu em cada posição
type PositionCount struct {
	Position int `json:"position"`
	Count    int `json:"count"`
}

var db *sql.DB

func waitForDB(database *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	backoff := time.Second
	maxBackoff := 10 * time.Second

	for {
		err := database.PingContext(ctx)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for database: %w", err)
		case <-time.After(backoff):
			log.Printf("Waiting for database (retrying in %v)...", backoff)
			// Exponential backoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}

func main() {
	var err error
	
	// Check if DATABASE_URL is provided (Fly.io, Heroku, etc)
	databaseURL := os.Getenv("DATABASE_URL")
	
	var port string
	if databaseURL != "" {
		// Use DATABASE_URL directly
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			log.Fatal("Error connecting to database:", err)
		}
		
		port = os.Getenv("API_PORT")
		if port == "" {
			port = "8080"
		}
	} else {
		// Use individual environment variables (for local development)
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		port = os.Getenv("API_PORT")

		if port == "" {
			port = "3000"
		}

		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)

		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal("Error connecting to database:", err)
		}
	}
	defer db.Close()

	// Wait for database to be ready with exponential backoff
	if err := waitForDB(db, 60*time.Second); err != nil {
		log.Fatal("Database not available:", err)
	}

	log.Println("Database connected successfully")

	// Create table
	createTable := `
	CREATE TABLE IF NOT EXISTS submissions (
		id SERIAL PRIMARY KEY,
		players JSONB NOT NULL,
		submitted_by VARCHAR(255) NOT NULL,
		ip_address VARCHAR(45) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	http.HandleFunc("/api/health", handleHealth)
	http.HandleFunc("/api/submissions", handleSubmissions)
	http.HandleFunc("/api/players/stats", handlePlayerStats)
	http.HandleFunc("/api/docs/", httpSwagger.WrapHandler)

	serverPort := ":" + port
	fmt.Printf("Server running on http://localhost%s\n", serverPort)
	fmt.Printf("API Docs available at http://localhost%s/api/docs/\n", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}

// @Summary Health check
// @Description Verifica o status da API e conexão com banco de dados
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/health [get]
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse representa uma resposta de sucesso na submissão
type SuccessResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
}

// SubmissionRequest representa o corpo da requisição de submissão
type SubmissionRequest struct {
	Players     []Player `json:"players"`
	SubmittedBy string   `json:"submittedBy"`
}

// handleSubmissions gerencia tanto POST quanto GET para /api/submissions
func handleSubmissions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleSubmit(w, r)
	case http.MethodGet:
		handleGetSubmissions(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary Submeter top 20 jogadores
// @Description Submete uma lista com exatamente 20 jogadores. Apenas uma submissão por IP é permitida.
// @Tags submissions
// @Accept json
// @Produce json
// @Param submission body SubmissionRequest true "Lista de 20 jogadores"
// @Success 201 "Submissão criada com sucesso"
// @Failure 400 {object} ErrorResponse "Número incorreto de jogadores ou campos faltando"
// @Failure 409 {object} ErrorResponse "IP já submeteu anteriormente"
// @Failure 500 {object} ErrorResponse "Erro interno do servidor"
// @Router /api/submissions [post]
func handleSubmit(w http.ResponseWriter, r *http.Request) {
	var submission SubmissionRequest

	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate that exactly 20 players are provided
	if len(submission.Players) != 20 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Invalid number of players",
			Message: fmt.Sprintf("Exactly 20 players are required, got %d", len(submission.Players)),
		})
		return
	}

	// Validate submittedBy is not empty
	if submission.SubmittedBy == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Missing required field",
			Message: "submittedBy is required",
		})
		return
	}

	// Get IP address from request
	ipAddress := getIPAddress(r)

	// Check if IP already submitted
	var existingID int
	err := db.QueryRow("SELECT id FROM submissions WHERE ip_address = $1", ipAddress).Scan(&existingID)
	if err == nil {
		// IP already exists
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "IP address already submitted",
			Message: "Only one submission per IP address is allowed",
		})
		return
	} else if err != sql.ErrNoRows {
		// Database error
		log.Println("Error checking IP:", err)
		http.Error(w, "Error checking submission", http.StatusInternalServerError)
		return
	}

	playersJSON, err := json.Marshal(submission.Players)
	if err != nil {
		http.Error(w, "Error processing players", http.StatusInternalServerError)
		return
	}

	var id int
	err = db.QueryRow(
		"INSERT INTO submissions (players, submitted_by, ip_address) VALUES ($1, $2, $3) RETURNING id",
		playersJSON, submission.SubmittedBy, ipAddress,
	).Scan(&id)

	if err != nil {
		log.Println("Error inserting submission:", err)
		http.Error(w, "Error saving submission", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// getIPAddress extracts the real IP address from the request
// considering proxies and load balancers
func getIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header (common with proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		for idx := 0; idx < len(forwarded); idx++ {
			if forwarded[idx] == ',' {
				return forwarded[:idx]
			}
		}
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i] == ':' {
			return ip[:i]
		}
	}
	return ip
}

// @Summary Listar submissões
// @Description Retorna todas as submissões ou filtra por submittedBy
// @Tags submissions
// @Produce json
// @Param submittedBy query string false "Filtrar por quem submeteu"
// @Success 200 {array} Submission
// @Failure 500 {object} ErrorResponse
// @Router /api/submissions [get]
func handleGetSubmissions(w http.ResponseWriter, r *http.Request) {
	submittedBy := r.URL.Query().Get("submittedBy")

	var rows *sql.Rows
	var err error

	if submittedBy == "" {
		rows, err = db.Query("SELECT id, players, submitted_by, ip_address, created_at FROM submissions ORDER BY created_at DESC")
	} else {
		rows, err = db.Query("SELECT id, players, submitted_by, ip_address, created_at FROM submissions WHERE submitted_by = $1 ORDER BY created_at DESC", submittedBy)
	}

	if err != nil {
		log.Println("Error querying submissions:", err)
		http.Error(w, "Error fetching submissions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []Submission
	for rows.Next() {
		var sub Submission
		var playersJSON []byte
		var createdAt time.Time

		err := rows.Scan(&sub.ID, &playersJSON, &sub.SubmittedBy, &sub.IPAddress, &createdAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		err = json.Unmarshal(playersJSON, &sub.Players)
		if err != nil {
			log.Println("Error unmarshaling players:", err)
			continue
		}

		sub.CreatedAt = createdAt.Format(time.RFC3339)
		result = append(result, sub)
	}

	if result == nil {
		result = []Submission{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// @Summary Estatísticas de um jogador
// @Description Retorna estatísticas de quantas vezes um jogador apareceu em cada posição
// @Tags players
// @Produce json
// @Param name query string true "Nome do jogador"
// @Success 200 {object} PlayerStats
// @Failure 400 {object} ErrorResponse "Nome do jogador não fornecido"
// @Failure 404 {object} ErrorResponse "Jogador não encontrado em nenhuma submissão"
// @Failure 500 {object} ErrorResponse "Erro interno do servidor"
// @Router /api/players/stats [get]
func handlePlayerStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	playerName := r.URL.Query().Get("name")
	if playerName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Missing required parameter",
			Message: "Player name is required",
		})
		return
	}

	// Query all submissions
	rows, err := db.Query("SELECT id, players FROM submissions ORDER BY created_at DESC")
	if err != nil {
		log.Println("Error querying submissions:", err)
		http.Error(w, "Error fetching submissions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Map to count positions
	positionCounts := make(map[int]int)
	totalSubmissions := 0

	for rows.Next() {
		var id int
		var playersJSON []byte

		err := rows.Scan(&id, &playersJSON)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		var players []Player
		err = json.Unmarshal(playersJSON, &players)
		if err != nil {
			log.Println("Error unmarshaling players:", err)
			continue
		}

		// Search for the player in this submission
		for _, player := range players {
			// Case-insensitive comparison
			if equalsCaseInsensitive(player.Name, playerName) {
				positionCounts[player.Position]++
				totalSubmissions++
				break // Found in this submission, move to next
			}
		}
	}

	// If player not found in any submission
	if totalSubmissions == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Player not found",
			Message: fmt.Sprintf("Player '%s' was not found in any submission", playerName),
		})
		return
	}

	// Build response
	positionBreakdown := make([]PositionCount, 0, len(positionCounts))
	for position, count := range positionCounts {
		positionBreakdown = append(positionBreakdown, PositionCount{
			Position: position,
			Count:    count,
		})
	}

	// Sort by position
	sort.Slice(positionBreakdown, func(i, j int) bool {
		return positionBreakdown[i].Position < positionBreakdown[j].Position
	})

	stats := PlayerStats{
		PlayerName:        playerName,
		TotalSubmissions:  totalSubmissions,
		PositionBreakdown: positionBreakdown,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// equalsCaseInsensitive compares two strings case-insensitively
func equalsCaseInsensitive(a, b string) bool {
	return len(a) == len(b) && strings.EqualFold(a, b)
}

