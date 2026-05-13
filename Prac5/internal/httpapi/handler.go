package httpapi

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MrFandore/Go_S2/internal/student"
)

type Handler struct {
	repo        *student.Repo
	stmtByID    *sql.Stmt
	stmtByEmail *sql.Stmt // для prepared statement по email
}

// NewHandler принимает подготовленные statements
func NewHandler(repo *student.Repo, stmtByID *sql.Stmt, stmtByEmail *sql.Stmt) *Handler {
	return &Handler{
		repo:        repo,
		stmtByID:    stmtByID,
		stmtByEmail: stmtByEmail,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"scheme": "https",
	})
}

func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawID := r.URL.Query().Get("id")
	if rawID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var st student.Student
	err = h.stmtByID.QueryRow(id).Scan(&st.ID, &st.FullName, &st.StudyGroup, &st.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}

// Дополнительное задание 3: получить студента по email
func (h *Handler) GetStudentByEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	var st student.Student
	err := h.stmtByEmail.QueryRow(email).Scan(&st.ID, &st.FullName, &st.StudyGroup, &st.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}
