package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=123456 dbname=cal_salary sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	fmt.Println("--- USERS ---")
	var users []struct {
		ID       string         `db:"id"`
		Username string         `db:"username"`
		Email    sql.NullString `db:"email"`
	}
	err = db.Select(&users, "SELECT id, username, email FROM users")
	if err != nil {
		log.Println("Error:", err)
	}
	for _, u := range users {
		fmt.Printf("User: %s, ID: %s, Email: %v\n", u.Username, u.ID, u.Email.String)
	}

	fmt.Println("\n--- ROLES ---")
	var roles []struct {
		ID   string `db:"id"`
		Slug string `db:"slug"`
		Name string `db:"name"`
	}
	err = db.Select(&roles, "SELECT id, slug, name FROM roles")
	if err != nil {
		log.Println("Error:", err)
	}
	for _, r := range roles {
		fmt.Printf("Role: %s, Slug: %s, ID: %s\n", r.Name, r.Slug, r.ID)
	}

	fmt.Println("\n--- USER ROLES MAPPING ---")
	var userRoles []struct {
		UserID string `db:"user_id"`
		RoleID string `db:"role_id"`
	}
	err = db.Select(&userRoles, "SELECT user_id, role_id FROM user_roles")
	if err != nil {
		log.Println("Error:", err)
	}
	for _, ur := range userRoles {
		fmt.Printf("UserID: %s -> RoleID: %s\n", ur.UserID, ur.RoleID)
	}

	fmt.Println("\n--- USER PROFILES ---")
	var profiles []struct {
		ID       string         `db:"id"`
		UserID   sql.NullString `db:"user_id"`
		Code     string         `db:"code"`
		FullName string         `db:"full_name"`
	}
	err = db.Select(&profiles, "SELECT id, user_id, code, full_name FROM user_profiles")
	if err != nil {
		log.Println("Error:", err)
	}
	for _, p := range profiles {
		fmt.Printf("Profile: %s, Code: %s, UserID: %v, ProfileID: %s\n", p.FullName, p.Code, p.UserID.String, p.ID)
	}
}
