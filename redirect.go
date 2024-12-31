// redirect.go
package main

import (
    "fmt"
    "net/http"
)

func handleRedirect(w http.ResponseWriter, r *http.Request) {
    // Get the domain from Host header
    domain := r.Host
    
    // Decode and validate the domain
    if !validateDomain(domain) {
        http.Error(w, "Invalid domain", http.StatusBadRequest)
        return
    }

    // Build the redirection URL
    url := fmt.Sprintf("https://%s/.well-known/lnurlp/", domain)
    
    // Redirect to the LNURL endpoint
    http.Redirect(w, r, url, http.StatusFound)
}

