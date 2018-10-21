// model.go

package main

import (
	"database/sql"
	"fmt"
)

type user struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Points   string `json:"points"`
}

type subject struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// users options

func (u *user) getUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT username, email, points FROM users WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.UserName, &u.Email, &u.Points)
}

func (u *user) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE users SET username='%s', email='%s', password='%s', points='%s'WHERE id=%d", u.UserName, u.Email, u.Password, u.Points, u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *user) deleteUser(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *user) checkEmail(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id FROM users WHERE email='%s'", u.Email)
	return db.QueryRow(statement).Scan(&u.ID)
}

func (u *user) createUser(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO users(username, email, password, points) VALUES('%s', '%s', '%s', '%s')", u.UserName, u.Email, u.Password, u.Points)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	statement := fmt.Sprintf("SELECT id, username, email, points FROM users LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.UserName, &u.Email, &u.Points); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// subject

func (s *subject) getSubject(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name, token FROM subject WHERE id=%d", s.ID)
	return db.QueryRow(statement).Scan(&s.Name, &s.Token)
}

func (s *subject) createSubject(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO subject(name, token) VALUES('%s', '%s')", s.Name, s.Token)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *subject) updateSubject(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE subjects SET name='%s', token='%s' WHERE id=%d", s.Name, s.Token, s.ID)
	_, err := db.Exec(statement)
	return err
}

func (s *subject) deleteSubject(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM subjects WHERE id=%d", s.ID)
	_, err := db.Exec(statement)
	return err
}

func getSubjects(db *sql.DB, start, count int) ([]subject, error) {
	statement := fmt.Sprintf("SELECT id, name, token FROM subject LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	subjects := []subject{}

	for rows.Next() {
		var s subject
		if err := rows.Scan(&s.ID, &s.Name, &s.Token); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}

	return subjects, nil
}

// user auth

func (u *user) authUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, username, email FROM users WHERE email='%s' AND password='%s'", u.Email, u.Password)
	return db.QueryRow(statement).Scan(&u.ID, &u.UserName, &u.Email)
}
