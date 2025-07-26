package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Username string

func NewUsername() Username {
	return ""
}

func ParseUsername(s string) (Username, error) {
	return Username(s), nil
}

func (un Username) String() string {
	return string(un)
}

type User struct {
	Username    Username
	FullName    string
	DisplayName string
	Email       string
}

type UserCreate struct {
	UserName    string
	FullName    string
	DisplayName *string
	Email       string
}

type UserUpdate struct {
	FullName    string
	DisplayName *string
	Email       *string
}

func CreateUser(user UserCreate) (Username, error) {
	var usernames, err = CreateUsers([]UserCreate{user})
	if err != nil {
		return NewUsername(), err
	}

	if len(usernames) == 0 {
		return NewUsername(), sql.ErrNoRows
	}

	return usernames[0], nil
}

func CreateUsers(users []UserCreate) ([]Username, error) {
	if len(users) == 0 {
		return []Username{}, nil
	}

	// Build query with placeholders
	query := `
		INSERT INTO users (username, full_name, display_name, email)
		VALUES
	`
	args := []interface{}{}
	placeholderIndex := 1

	for i, user := range users {
		if i > 0 {
			query += ",\n"
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d)",
			placeholderIndex,
			placeholderIndex+1,
			placeholderIndex+2,
			placeholderIndex+3,
		)
		args = append(args,
			user.UserName,
			user.FullName,
			user.DisplayName,
			user.Email,
		)
		placeholderIndex += 4
	}

	query += "\nRETURNING username"

	// Run the query
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uns []Username
	for rows.Next() {
		var un Username
		if err := rows.Scan(&un); err != nil {
			return nil, err
		}
		uns = append(uns, un)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uns, nil
}

func GetUser(un Username) (*User, error) {
	var users, err = GetUsers([]Username{un})

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return &users[0], err
}

func GetUsers(uns []Username) ([]User, error) {

	if len(uns) == 0 {
		return []User{}, nil
	}

	var builder strings.Builder

	builder.WriteString(
		`
		select username, full_name, display_name, email
		from users
		where username in (
		`,
	)

	for i, un := range uns {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("'")
		builder.WriteString(un.String())
		builder.WriteString("'")
	}

	builder.WriteString(")")

	var rows *sql.Rows
	var err error
	rows, err = DB.Query(builder.String())

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User = make([]User, 0, len(uns))
	for rows.Next() {
		var user User
		rows.Scan(&user.Username, &user.FullName, &user.DisplayName, &user.Email)
		users = append(users, user)
	}

	return users, nil
}

func UpdateUser(un Username, update UserUpdate) error {
	return UpdateUsers([]Username{un}, update)
}

func UpdateUsers(uns []Username, update UserUpdate) error {

	if len(uns) == 0 {
		return nil
	}

	var builder strings.Builder
	builder.WriteString("update users set ")

	var args []interface{}
	var argIndex int = 1
	for field, value := range map[string]any{
		"full_name":    update.FullName,
		"display_name": update.DisplayName,
		"email":        update.Email,
	} {
		if value == nil || reflect.ValueOf(value).IsNil() {
			continue
		}

		if argIndex > 1 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	builder.WriteString(" where username in (")
	for i, un := range uns {
		if i > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString("'")
		builder.WriteString(un.String())
		builder.WriteString("'")
	}
	builder.WriteString(")")

	var result, err = DB.Exec(builder.String(), args...)
	print("result: ")
	print(result)

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return err
}

func DeleteUser(un Username) error {
	return DeleteUsers([]Username{un})
}

func DeleteUsers(uns []Username) error {

	if len(uns) == 0 {
		return nil
	}

	var builder strings.Builder
	builder.WriteString(
		`
		delete from users
		where username in (
		`,
	)

	for i, un := range uns {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("'")
		builder.WriteString(un.String())
		builder.WriteString("'")
	}

	builder.WriteString(")")

	_, err := DB.Exec(builder.String())

	if err != nil {
		return err
	}

	return nil
}
