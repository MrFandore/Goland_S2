package http

import (
	"encoding/json"
	"net/http"

	"Prac1/services/auth/internal/service"
	"github.com/gorilla/mux"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type VerifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error,omitempty"`
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/v1/auth/login", loginHandler).Methods("POST")
	r.HandleFunc("/v1/auth/verify", verifyHandler).Methods("GET")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := service.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := LoginResponse{
		AccessToken: token,
		TokenType:   "bearer",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("authorization")
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "bearer " {
		token = authHeader[7:]
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(VerifyResponse{Valid: false, Error: "missing or invalid authorization header"})
		return
	}

	valid, subject := service.VerifyToken(token)
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(VerifyResponse{Valid: false, Error: "unauthorized"})
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(VerifyResponse{Valid: true, Subject: subject})
}
