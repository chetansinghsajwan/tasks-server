package errorcodes

type Code uint32

const (
	Unknown Code = iota
	Internal
	TxCreate
	TxCommit

	// Auth related errors
	InvalidToken
	AuthMatchFail

	// User related errors
	UserNotFound
	UserIDNull
	UserIDAlreadyExists
	InvalidUserIDFormat
	UserEmailNull
	UserEmailAlreadyExists
	InvalidUserEmailFormat
	InvalidUserFullNameFormat
	InvalidUserDisplayNameFormat

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
