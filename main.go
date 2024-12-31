// main.go
package main

import (
    "encoding/json"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func init() {
    // Initialize zerolog
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    
    // Set log level based on environment variable
    if os.Getenv("DEBUG") != "" {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    } else {
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    }
}

func handleCanIssue(w http.ResponseWriter, r *http.Request) {
    domain := r.URL.Query().Get("domain")
    log.Debug().
        Str("domain", domain).
        Msg("certificate issuance request")
    
    // Allow certificates for all domains
    w.WriteHeader(http.StatusOK)
}

func sendError(w http.ResponseWriter, msg string) {
    log.Error().Msg(msg)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "ERROR",
        "reason": msg,
    })
}

func main() {
    r := mux.NewRouter()
    
    // Add CORS middleware
    corsMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
    
    r.Use(corsMiddleware)
    
    // Add the can-issue endpoint for Caddy
    r.HandleFunc("/can-issue", handleCanIssue)
    
    // LNURL endpoints
    r.HandleFunc("/.well-known/lnurlp/{username}", handleLNURL)
    
    // Start server
    addr := "0.0.0.0:12345"
    log.Info().Str("addr", addr).Msg("starting server")
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatal().Err(err).Msg("server failed")
    }
}

