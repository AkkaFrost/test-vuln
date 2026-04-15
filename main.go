package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	_ "github.com/lib/pq"
)

// Hardcoded secrets and credentials
const (
	DBHost     = "prod-db.internal.company.com"
	DBUser     = "admin"
	DBPassword = "SuperSecret123!"
	DBName     = "production"

	AWSAccessKeyID     = "AKIAIOSFODNN7EXAMPLE"
	AWSSecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	JWTSecret     = "my-jwt-secret-do-not-share"
	APIKey        = "sk-proj-abc123def456ghi789jkl012mno345"
	SlackWebhook  = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
	PrivateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0Z3VS5JJcds3xfn/ygWyF3PBWLEBNhTzPIGKJB0WkVOmqNiI
bMIEpAIBAAKCAQEA0Z3VS5JJcds3xfn0ygWep8BNhTzPIGKJB0WkVOmqNiIEXAMPLE
-----END RSA PRIVATE KEY-----`
)

var db *sql.DB

func init() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBUser, DBPassword, DBName)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

// SQL Injection — user input concatenated directly into query
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	query := fmt.Sprintf("SELECT id, email, password_hash FROM users WHERE username = '%s'", username)
	row := db.QueryRow(query)

	var id int
	var email, passwordHash string
	row.Scan(&id, &email, &passwordHash)

	// Sensitive data exposure — returning password hash to client
	fmt.Fprintf(w, "ID: %d, Email: %s, PasswordHash: %s", id, email, passwordHash)
}

// Command Injection — unsanitized user input passed to shell
func pingHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	cmd := exec.Command("sh", "-c", "ping -c 1 "+host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(output)
}

// Cross-Site Scripting (XSS) — reflected user input without escaping
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Search results for: %s</h1>", query)
}

// Path Traversal — user-controlled file path with no validation
func readFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	data, err := ioutil.ReadFile("/var/data/" + filename)
	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}
	w.Write(data)
}

// Weak cryptography — MD5 for password hashing
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// Insecure deserialization / mass assignment — trusting all client fields
func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	role := r.FormValue("role") // user can set themselves to "admin"
	email := r.FormValue("email")

	query := fmt.Sprintf("UPDATE users SET email='%s', role='%s' WHERE id=1", email, role)
	db.Exec(query)

	fmt.Fprintf(w, "Profile updated")
}

// Open redirect
func loginHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := r.URL.Query().Get("redirect")
	// No validation — attacker can redirect to a phishing site
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// SSRF — fetching arbitrary URLs from user input
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")
	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Failed to fetch", 500)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)
}

// Verbose error exposure — leaking stack traces and internals
func debugHandler(w http.ResponseWriter, r *http.Request) {
	_, err := db.Query("SELECT * FROM nonexistent_table")
	if err != nil {
		fmt.Fprintf(w, "Database error: %v\nConnection: host=%s user=%s password=%s",
			err, DBHost, DBUser, DBPassword)
	}
}

func main() {
	http.HandleFunc("/user", getUserHandler)         // SQL injection
	http.HandleFunc("/ping", pingHandler)             // Command injection
	http.HandleFunc("/search", searchHandler)         // XSS
	http.HandleFunc("/read", readFileHandler)         // Path traversal
	http.HandleFunc("/profile", updateProfileHandler) // Mass assignment + SQLi
	http.HandleFunc("/login", loginHandler)           // Open redirect
	http.HandleFunc("/proxy", proxyHandler)           // SSRF
	http.HandleFunc("/debug", debugHandler)           // Info disclosure

	// Listening on all interfaces, no TLS
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
