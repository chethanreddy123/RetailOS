package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SuperAdminHandler struct {
	pool       *pgxpool.Pool
	jwtSecret  string
	adminKey   string
}

func NewSuperAdminHandler(pool *pgxpool.Pool, jwtSecret, adminKey string) *SuperAdminHandler {
	return &SuperAdminHandler{pool: pool, jwtSecret: jwtSecret, adminKey: adminKey}
}

// POST /superadmin/seed — creates the first super admin (protected by X-Admin-Key)
func (h *SuperAdminHandler) Seed(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	_, err = h.pool.Exec(r.Context(),
		`INSERT INTO super_admins (email, hashed_password) VALUES ($1, $2)`,
		req.Email, string(hashed),
	)
	if err != nil {
		writeError(w, http.StatusConflict, "super admin already exists or db error: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "super admin created"})
}

// POST /superadmin/login — email + password → JWT with role=superadmin
func (h *SuperAdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var id, hashedPw string
	err := h.pool.QueryRow(r.Context(),
		`SELECT id, hashed_password FROM super_admins WHERE email = $1`,
		req.Email,
	).Scan(&id, &hashedPw)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  id,
		"role": "superadmin",
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not sign token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": signed})
}
