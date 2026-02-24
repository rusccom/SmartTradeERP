package httpx

import (
    "encoding/json"
    "net/http"
)

type Meta struct {
    Page    int `json:"page"`
    PerPage int `json:"per_page"`
    Total   int `json:"total"`
}

type ErrorBody struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

type Envelope struct {
    Data  interface{} `json:"data"`
    Error *ErrorBody  `json:"error"`
    Meta  *Meta       `json:"meta"`
}

func WriteData(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
    body := Envelope{Data: data, Error: nil, Meta: meta}
    writeJSON(w, status, body)
}

func WriteError(w http.ResponseWriter, status int, code, message string, details interface{}) {
    body := Envelope{}
    body.Data = nil
    body.Error = &ErrorBody{Code: code, Message: message, Details: details}
    body.Meta = nil
    writeJSON(w, status, body)
}

func writeJSON(w http.ResponseWriter, status int, body Envelope) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(body)
}
