package main

import (
	"fmt"
	"golang_api/config"
	"golang_api/database"
	"golang_api/handlers"
	"log"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {
	cfg := config.Load()

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: cfg.SuperTokensConnectionURI,
			APIKey:        cfg.SuperTokensAPIKey,
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "golang_api",
			APIDomain:     fmt.Sprintf("http://localhost:%s", cfg.Port),
			WebsiteDomain: fmt.Sprintf("http://localhost:%s", cfg.Port),
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				CookieSameSite: &[]string{"lax"}[0],
			}),
		},
	})

	if err != nil {
		log.Fatalf("Failed to initialize SuperTokens: %v", err)
	}

	db, err := database.New(cfg.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	authHandler := handlers.NewAuthHandler(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", authHandler.Register)
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/logout", session.VerifySession(nil, authHandler.Logout))
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})

	handler := supertokens.Middleware(corsMiddleware(mux))

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, st-auth-mode")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
