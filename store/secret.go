package store

type UserSecret struct {
	ID   string
	Pass string
}

type CreateUserSecretParams struct {
	ID   string
	Pass string
}

type UpdateUserSecretParams struct {
	Pass string
}
