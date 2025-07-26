package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ListID uint64

func ParseListID(s string) (ListID, error) {
	parsed, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return 0, err
	}

	return ListID(parsed), nil
}

func (id ListID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type List struct {
	ID    ListID
	Title string
}

type ListCreate struct {
	Title string
}

type ListUpdate struct {
	Title *string
}

func CreateList(list ListCreate) (ListID, error) {
	var listId, err = CreateLists([]ListCreate{list})
	if err != nil {
		return 0, err
	}

	if len(listId) == 0 {
		return 0, sql.ErrNoRows
	}

	return listId[0], nil
}

func CreateLists(lists []ListCreate) ([]ListID, error) {
	if len(lists) == 0 {
		return []ListID{}, nil
	}

	// Build query with placeholders
	query := `
		INSERT INTO lists (title)
		VALUES
	`
	args := []interface{}{}
	placeholderIndex := 1

	for i, list := range lists {
		if i > 0 {
			query += ",\n"
		}
		query += fmt.Sprintf("($%d)", placeholderIndex)
		args = append(args, list.Title)
		placeholderIndex += 1
	}

	query += "\nRETURNING id"

	// Run the query
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []ListID
	for rows.Next() {
		var id ListID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func GetList(id ListID) (*List, error) {
	var lists, err = GetLists([]ListID{id})

	if err != nil {
		return nil, err
	}

	if len(lists) == 0 {
		return nil, sql.ErrNoRows
	}

	return &lists[0], err
}

func GetLists(ids []ListID) ([]List, error) {

	if len(ids) == 0 {
		return []List{}, nil
	}

	var builder strings.Builder

	builder.WriteString(
		`
		select id, title, description, priority, due_date, assignee, labels
		from lists
		where id in (
		`,
	)

	for i, id := range ids {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("'")
		builder.WriteString(id.String())
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

	var lists []List = make([]List, 0, len(ids))
	for rows.Next() {
		var list List
		rows.Scan(&list.ID, &list.Title)
		lists = append(lists, list)
	}

	return lists, nil
}

func UpdateList(id ListID, update ListUpdate) error {
	return UpdateLists([]ListID{id}, update)
}

func UpdateLists(ids []ListID, update ListUpdate) error {

	if len(ids) == 0 {
		return nil
	}

	var builder strings.Builder
	builder.WriteString("update lists set ")

	var args []interface{}
	var argIndex int = 1
	for field, value := range map[string]interface{}{
		"title": update.Title,
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

	builder.WriteString(" where id in (")
	for i, id := range ids {
		if i > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString("'")
		builder.WriteString(id.String())
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

func DeleteList(id ListID) error {
	return DeleteLists([]ListID{id})
}

func DeleteLists(ids []ListID) error {

	if len(ids) == 0 {
		return nil
	}

	var builder strings.Builder
	builder.WriteString(
		`
		delete from lists
		where id in (
		`,
	)

	for i, id := range ids {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("'")
		builder.WriteString(id.String())
		builder.WriteString("'")
	}

	builder.WriteString(")")

	_, err := DB.Exec(builder.String())

	if err != nil {
		return err
	}

	return nil
}
