package main

import (
	"log"
	"net/http"

	"careerness/api/internal/handler"
	"careerness/api/internal/session"
)

func main() {
	mux := http.NewServeMux()

	// session attachment（session → workspace_root の束縛）を共有する in-memory store。
	// attach と apply-patch が同じ束縛を参照する（ADR-006）。
	store := session.NewStore()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/api/v1/conversations/message", handler.PostMessage)
	mux.HandleFunc("/api/v1/patches/validate", handler.PostValidatePatch)
	mux.HandleFunc("/api/v1/extract", handler.PostExtract)
	mux.HandleFunc("/api/v1/workspace/attach", handler.PostAttach(store))
	mux.HandleFunc("/api/v1/apply-patch", handler.PostApplyPatch)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsMiddleware(mux),
	}

	log.Println("CareerNess API starting on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// corsMiddleware は開発環境向けに Vite dev server からのリクエストを許可する。
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://localhost:4173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
