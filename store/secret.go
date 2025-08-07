package store

type UserSecret struct {
	ID   UserID
	Pass string
}

type CreateUserSecretParams struct {
	ID   UserID
	Pass string
}

type UpdateUserSecretParams struct {
	Pass string
}
