// lnurl.go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/rs/zerolog/log"
)

func handleLNURL(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]
    
    // Get the amount from query string if present
    amount := r.URL.Query().Get("amount")
    
    // Extract the domain from the Host header
    domain := r.Host
    
    log.Debug().
        Str("username", username).
        Str("domain", domain).
        Str("amount", amount).
        Msg("handling LNURL request")

    if amount == "" {
        // First step: return the LNURL-pay parameters
        metadata := makeMetadata(username, domain)
        
        resp := LNURLResponse{
            Status:        "OK",
            Tag:          "payRequest",
            Callback:     fmt.Sprintf("https://%s/.well-known/lnurlp/%s", domain, username),
            MinSendable:  1000,          // 1 sat minimum
            MaxSendable:  100000000,     // 1000 sats maximum
            Metadata:     metadata,
            CommentAllowed: 0,
        }
        
        json.NewEncoder(w).Encode(resp)
        return
    }

    // Second step: create an invoice
    msats, err := strconv.Atoi(amount)
    if err != nil {
        sendError(w, "invalid amount")
        return
    }

    bolt11, err := makeInvoice(username, domain, msats)
    if err != nil {
        sendError(w, "failed to create invoice: "+err.Error())
        return
    }

    resp := InvoiceResponse{
        Status: "OK",
        PR:     bolt11,
        Routes: []interface{}{},  // Changed from []string{}
        SuccessAction: &SuccessAction{
            Tag:     "message",
            Message: "Payment received!",
        },
    }

    json.NewEncoder(w).Encode(resp)
}

