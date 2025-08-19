package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"effective-golang/internal/auth"
	"effective-golang/internal/game"
	"effective-golang/internal/leaderboard"
	"effective-golang/internal/models"
	"effective-golang/pkg/utils"
)

// Application represents the main application
type Application struct {
	server           *http.Server
	authService      *auth.AuthService
	gameService      *game.GameService
	leaderboardSvc   *leaderboard.LeaderboardService
	unitOfWork       models.UnitOfWork
	
	// Graceful shutdown
	shutdownCh       chan os.Signal
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewApplication creates a new application instance
func NewApplication() (*Application, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
	
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	
	// Initialize repositories (in real app, these would be database implementations)
	// For this learning project, we'll use in-memory implementations
	unitOfWork := utils.NewInMemoryUnitOfWork()
	
	// Initialize services
	authService := auth.NewAuthService(
		unitOfWork.UserRepository(),
		unitOfWork.CacheRepository(),
	)
	
	gameService := game.NewGameService(
		unitOfWork.GameRepository(),
		unitOfWork.UserRepository(),
		unitOfWork.LeaderboardRepository(),
		unitOfWork.CacheRepository(),
		10, // max workers
		100, // queue size
	)
	
	leaderboardSvc := leaderboard.NewLeaderboardService(
		unitOfWork.LeaderboardRepository(),
		unitOfWork.UserRepository(),
		unitOfWork.CacheRepository(),
		3600, // cache TTL in seconds
	)
	
	// Create router
	router := mux.NewRouter()
	
	// Setup middleware
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)
	
	// Setup routes
	setupRoutes(router, authService, gameService, leaderboardSvc)
	
	// Create HTTP server
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	app := &Application{
		server:         server,
		authService:    authService,
		gameService:    gameService,
		leaderboardSvc: leaderboardSvc,
		unitOfWork:     unitOfWork,
		shutdownCh:     make(chan os.Signal, 1),
		ctx:            ctx,
		cancel:         cancel,
	}
	
	// Setup graceful shutdown
	signal.Notify(app.shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	
	return app, nil
}

// Start starts the application
func (app *Application) Start() error {
	log.Printf("Starting server on port %s", app.server.Addr)
	
	// Start server in a goroutine
	go func() {
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()
	
	// Wait for shutdown signal
	<-app.shutdownCh
	log.Println("Shutdown signal received")
	
	return app.Shutdown()
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() error {
	log.Println("Shutting down application...")
	
	// Cancel context
	app.cancel()
	
	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Shutdown HTTP server
	if err := app.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	
	// Close services
	if err := app.gameService.Close(); err != nil {
		log.Printf("Game service shutdown error: %v", err)
	}
	
	app.leaderboardSvc.Close()
	
	// Close unit of work
	if err := app.unitOfWork.Close(); err != nil {
		log.Printf("Unit of work shutdown error: %v", err)
	}
	
	log.Println("Application shutdown complete")
	return nil
}

// setupRoutes sets up all application routes
func setupRoutes(
	router *mux.Router,
	authService *auth.AuthService,
	gameService *game.GameService,
	leaderboardSvc *leaderboard.LeaderboardService,
) {
	// Health check
	router.HandleFunc("/health", healthHandler).Methods("GET")
	
	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	
	// Auth routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", registerHandler(authService)).Methods("POST")
	auth.HandleFunc("/login", loginHandler(authService)).Methods("POST")
	auth.HandleFunc("/logout", logoutHandler(authService)).Methods("POST")
	
	// Game routes
	games := api.PathPrefix("/games").Subrouter()
	games.HandleFunc("", createGameHandler(gameService)).Methods("POST")
	games.HandleFunc("/{gameID}/start", startGameHandler(gameService)).Methods("POST")
	games.HandleFunc("/{gameID}/score", updateScoreHandler(gameService)).Methods("PUT")
	games.HandleFunc("/{gameID}/end", endGameHandler(gameService)).Methods("POST")
	games.HandleFunc("/active", getActiveGamesHandler(gameService)).Methods("GET")
	games.HandleFunc("/{gameID}", getGameHandler(gameService)).Methods("GET")
	
	// Leaderboard routes
	leaderboards := api.PathPrefix("/leaderboards").Subrouter()
	leaderboards.HandleFunc("", createLeaderboardHandler(leaderboardSvc)).Methods("POST")
	leaderboards.HandleFunc("/{leaderboardID}/scores", addScoreHandler(leaderboardSvc)).Methods("POST")
	leaderboards.HandleFunc("/{leaderboardID}/top", getTopEntriesHandler(leaderboardSvc)).Methods("GET")
	leaderboards.HandleFunc("/{leaderboardID}/rank/{userID}", getUserRankHandler(leaderboardSvc)).Methods("GET")
	leaderboards.HandleFunc("/{leaderboardID}/stats", getLeaderboardStatsHandler(leaderboardSvc)).Methods("GET")
	leaderboards.HandleFunc("/{leaderboardID}", getLeaderboardHandler(leaderboardSvc)).Methods("GET")
	
	// User routes
	users := api.PathPrefix("/users").Subrouter()
	users.HandleFunc("/{userID}/stats", getUserStatsHandler(authService)).Methods("GET")
}

// Middleware functions

// loggingMiddleware logs all HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Call next handler
		next.ServeHTTP(w, r)
		
		// Log request
		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Handler functions

// healthHandler handles health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}

// Helper functions

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// main function
func main() {
	// Create application
	app, err := NewApplication()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}
	
	// Start application
	if err := app.Start(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
