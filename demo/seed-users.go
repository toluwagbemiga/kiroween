package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DemoUser struct {
	Email    string
	Password string
	Name     string
	Role     string
}

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "haunted_saas")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to database successfully!")

	// Demo users with different roles
	demoUsers := []DemoUser{
		{Email: "admin@haunted.dev", Password: "Admin123!", Name: "Admin User", Role: "admin"},
		{Email: "member@haunted.dev", Password: "Member123!", Name: "Member User", Role: "member"},
		{Email: "john.doe@haunted.dev", Password: "John123!", Name: "John Doe", Role: "member"},
		{Email: "jane.smith@haunted.dev", Password: "Jane123!", Name: "Jane Smith", Role: "member"},
		{Email: "viewer@haunted.dev", Password: "Viewer123!", Name: "Viewer User", Role: "viewer"},
		{Email: "guest@haunted.dev", Password: "Guest123!", Name: "Guest User", Role: "viewer"},
	}

	fmt.Println("\n=== Creating Demo Users ===")
	for _, user := range demoUsers {
		if err := createUser(db, user); err != nil {
			log.Printf("Warning: Failed to create user %s: %v", user.Email, err)
		} else {
			fmt.Printf("✓ Created user: %s (Role: %s, Password: %s)\n", user.Email, user.Role, user.Password)
		}
	}

	fmt.Println("\n=== Demo Users Created Successfully ===")
	fmt.Println("\nYou can now login with any of these accounts:")
	fmt.Println("\nADMIN (Full Access):")
	fmt.Println("  Email: admin@haunted.dev")
	fmt.Println("  Password: Admin123!")
	fmt.Println("\nMEMBERS (Standard Access):")
	fmt.Println("  Email: member@haunted.dev | Password: Member123!")
	fmt.Println("  Email: john.doe@haunted.dev | Password: John123!")
	fmt.Println("  Email: jane.smith@haunted.dev | Password: Jane123!")
	fmt.Println("\nVIEWERS (Read-Only Access):")
	fmt.Println("  Email: viewer@haunted.dev | Password: Viewer123!")
	fmt.Println("  Email: guest@haunted.dev | Password: Guest123!")
}

func createUser(db *sql.DB, user DemoUser) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user
	var userID int64
	err = db.QueryRow(`
		INSERT INTO users (email, password_hash, name, is_active, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, true, true, NOW(), NOW())
		ON CONFLICT (email) DO UPDATE 
		SET password_hash = EXCLUDED.password_hash, name = EXCLUDED.name, updated_at = NOW()
		RETURNING id
	`, user.Email, string(hashedPassword), user.Name).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Get role ID
	var roleID int64
	err = db.QueryRow("SELECT id FROM roles WHERE name = $1", user.Role).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Assign role to user
	_, err = db.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
