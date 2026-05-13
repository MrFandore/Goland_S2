package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/MrFandore/Go_S2/internal/config"
	"github.com/MrFandore/Go_S2/internal/httpapi"
	"github.com/MrFandore/Go_S2/internal/student"
)

func main() {
	cfg := config.New()

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	repo := student.NewRepo(db)

	stmtByID, err := repo.PrepareGetByID()
	if err != nil {
		log.Fatal(err)
	}
	defer stmtByID.Close()

	stmtByEmail, err := repo.PrepareGetByEmail()
	if err != nil {
		log.Fatal(err)
	}
	defer stmtByEmail.Close()

	handler := httpapi.NewHandler(repo, stmtByID, stmtByEmail)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/students", handler.GetStudentByID)
	mux.HandleFunc("/students/by-email", handler.GetStudentByEmail)

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := &http.Server{
		Addr:      cfg.Addr,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	log.Printf("HTTPS server started on %s", cfg.Addr)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}
