// domain_validate.go
package main

import (
    "net"
    "strings"

    "github.com/rs/zerolog/log"
)

func validateDomain(domain string) bool {
    // Check for known invalid domains
    if strings.HasSuffix(domain, ".local") ||
        strings.HasSuffix(domain, ".localhost") ||
        strings.HasSuffix(domain, ".internal") {
        log.Debug().Str("domain", domain).Msg("rejected internal domain")
        return false
    }

    // Look up the domain
    ips, err := net.LookupIP(domain)
    if err != nil {
        log.Debug().Err(err).Str("domain", domain).Msg("domain lookup failed")
        return false
    }

    // Must have at least one IP
    if len(ips) == 0 {
        log.Debug().Str("domain", domain).Msg("no IPs found for domain")
        return false
    }

    return true
}

