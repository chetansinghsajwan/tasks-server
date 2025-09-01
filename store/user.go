package store

type User struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	FullName    string  `json:"full_name"`
	DisplayName *string `json:"display_name"`
}

type CreateUserParams struct {
	ID          string
	Email       string
	FullName    string
	DisplayName *string
}

type UpdateUserParams struct {
	ID          *string
	Email       *string
	FullName    *string
	DisplayName **string
}
