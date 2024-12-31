package main

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
)

func makeMetadata(username, domain string) string {
    // Create metadata array for LNURL-pay
    metadata := [][]string{
        {"text/identifier", username + "@" + domain},
        {"text/plain", "Satoshis to " + username + "@" + domain + "."},
    }

    // Convert to JSON string
    metadataJSON, err := json.Marshal(metadata)
    if err != nil {
        return "[]"
    }
    return string(metadataJSON)
}

func base64ImageFromURL(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 {
        return "", errors.New("image returned status " + strconv.Itoa(resp.StatusCode))
    }

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read image from %s: %w", url, err)
    }

    return base64.StdEncoding.EncodeToString(data), nil
}

