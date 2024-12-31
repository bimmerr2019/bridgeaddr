// types.go
package main

type LNURLResponse struct {
    Status         string `json:"status"`
    Tag            string `json:"tag"`
    Callback       string `json:"callback"`
    MaxSendable    int64  `json:"maxSendable"`
    MinSendable    int64  `json:"minSendable"`
    Metadata       string `json:"metadata"`
    CommentAllowed int    `json:"commentAllowed"`
}

type SuccessAction struct {
    Tag     string `json:"tag"`
    Message string `json:"message"`
}

type InvoiceResponse struct {
    Status        string         `json:"status"`
    PR            string         `json:"pr"`
    Routes        []interface{}  `json:"routes"`
    Disposable    bool          `json:"disposable,omitempty"`
    SuccessAction *SuccessAction `json:"successAction,omitempty"`
}

