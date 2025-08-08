package store

type User struct {
	ID          string  `json:"id"`
	FullName    string  `json:"full_name"`
	DisplayName *string `json:"display_name"`
	Email       string  `json:"email"`
}

type CreateUserParams struct {
	ID          string
	FullName    string
	DisplayName *string
	Email       string
}

type UpdateUserParams struct {
	ID          **string
	Email       **string
	FullName    **string
	DisplayName **string
}
