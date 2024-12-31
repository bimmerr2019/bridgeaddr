// makeinvoice.go
package main

import (
    "encoding/json"
    "fmt"
    "net"
    "strings"
    "time"
    
    "github.com/jb55/lnsocket/go"
    "github.com/rs/zerolog/log"
)

func makeInvoice(username, domain string, msat int) (bolt11 string, err error) {
    log.Debug().
        Str("username", username).
        Str("domain", domain).
        Int("msat", msat).
        Msg("creating invoice")

    // grab all the necessary data from DNS
    var (
        kind     string
        host     string
        nodeid   string
        rune_    string
    )

    // Try DNS lookups with detailed logging
    kindRecord := "_kind." + domain
    log.Debug().Str("looking up", kindRecord).Msg("DNS lookup")
    if v, err := net.LookupTXT(kindRecord); err == nil && len(v) > 0 {
        kind = strings.Trim(v[0], "\"")
        log.Debug().Str("kind", kind).Msg("found kind")
    } else {
        return "", fmt.Errorf("missing kind for %s: %v", kindRecord, err)
    }

    hostRecord := "_host." + domain
    log.Debug().Str("looking up", hostRecord).Msg("DNS lookup")
    if v, err := net.LookupTXT(hostRecord); err == nil && len(v) > 0 {
        host = strings.Trim(v[0], "\"")
        log.Debug().Str("host", host).Msg("found host")
    } else {
        return "", fmt.Errorf("missing host for %s: %v", hostRecord, err)
    }

    if kind != "commando" {
        return "", fmt.Errorf("unsupported backend kind: %s", kind)
    }

    nodeidRecord := "_nodeid." + domain
    log.Debug().Str("looking up", nodeidRecord).Msg("DNS lookup")
    if v, err := net.LookupTXT(nodeidRecord); err == nil && len(v) > 0 {
        nodeid = strings.Trim(v[0], "\"")
        log.Debug().Str("nodeid", nodeid).Msg("found nodeid")
    } else {
        return "", fmt.Errorf("missing nodeid for %s: %v", nodeidRecord, err)
    }

    runeRecord := "_rune." + domain
    log.Debug().Str("looking up", runeRecord).Msg("DNS lookup")
    if v, err := net.LookupTXT(runeRecord); err == nil && len(v) > 0 {
        rune_ = strings.Trim(v[0], "\"")
        log.Debug().Msg("found rune")  // Don't log the actual rune
    } else {
        return "", fmt.Errorf("missing rune for %s: %v", runeRecord, err)
    }

    // Create lnsocket connection
    ln := &lnsocket.LNSocket{}
    ln.GenKey()
    
    log.Debug().
        Str("host", host).
        Str("nodeid", nodeid).
        Msg("connecting to node")
    
    err = ln.ConnectAndInit(host, nodeid)
    if err != nil {
        return "", fmt.Errorf("connect error: %v", err)
    }
    defer ln.Disconnect()

    // Prepare the invoice parameters with more unique label
    label := fmt.Sprintf("bridgeaddr/%s/%d", username, time.Now().Unix())
    description := fmt.Sprintf("Payment to %s@%s", username, domain)

    // Create command with proper JSON structure
    cmdMap := map[string]interface{}{
        "amount_msat": msat,
        "label":      label,
        "description": description,
    }
    
    cmdBytes, err := json.Marshal(cmdMap)
    if err != nil {
        return "", fmt.Errorf("failed to create command JSON: %v", err)
    }
    
    log.Debug().
        Str("cmd", string(cmdBytes)).
        Msg("sending command to node")
    
    response, err := ln.Rpc(rune_, "invoice", string(cmdBytes))
    if err != nil {
        return "", fmt.Errorf("rpc error: %v", err)
    }

    log.Debug().
        Str("response", response).
        Msg("received response from node")

    // Parse the response to get bolt11
    var resp struct {
        Result struct {
            Bolt11 string `json:"bolt11"`
        } `json:"result"`
        Error *struct {
            Message string `json:"message"`
        } `json:"error,omitempty"`
    }
    
    err = json.Unmarshal([]byte(response), &resp)
    if err != nil {
        return "", fmt.Errorf("failed to parse response: %v, raw: %s", err, response)
    }

    // Check for error in response
    if resp.Error != nil {
        return "", fmt.Errorf("node error: %s", resp.Error.Message)
    }

    // Validate bolt11
    if resp.Result.Bolt11 == "" {
        return "", fmt.Errorf("empty bolt11 in response: %s", response)
    }

    return resp.Result.Bolt11, nil
}

