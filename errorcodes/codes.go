package errorcodes

type Code uint32

const (
	Unknown Code = iota
	TxCreate
	TxCommit

	// User related errors
	UserNotFound
	UserIDNull
	UserIDAlreadyExists
	UserIDFormat
	UserEmailNull
	UserEmailAlreadyExists
	UserEmailFormat
	UserFullNameFormat
	UserDisplayNameFormat

	// List related errors
	ListNotFound
	ListIDNull
	ListIDAlreadyExists
	ListIDFormat

	// List Accesss related errors
	ListAccessNotFound
	ListAccessAlreadyExists
	ListAccessOwnerAlreadyExists

	// Task related errors
	TaskNotFoundCode

	// Secret related errors
	UserSecretNotFound
	InvalidSecretIDFormat
	InvalidSecretScope
	InvalidUserSecretPassFormat
)
