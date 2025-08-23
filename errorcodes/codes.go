package errorcodes

type Code string

const (
	Unknown  Code = "Unknown"
	Internal Code = "Internal"
	TxCreate Code = "TxCreate"
	TxCommit Code = "TxCommit"

	// Auth related errors
	InvalidToken  Code = "InvalidToken"
	AuthMatchFail Code = "AuthMatchFail"

	// User related errors
	UserNotFound                 Code = "UserNotFound"
	UserIDNull                   Code = "UserIDNull"
	UserIDAlreadyExists          Code = "UserIDAlreadyExists"
	InvalidUserIDFormat          Code = "InvalidUserIDFormat"
	UserEmailNull                Code = "UserEmailNull"
	UserEmailAlreadyExists       Code = "UserEmailAlreadyExists"
	InvalidUserEmailFormat       Code = "InvalidUserEmailFormat"
	InvalidUserFullNameFormat    Code = "InvalidUserFullNameFormat"
	InvalidUserDisplayNameFormat Code = "InvalidUserDisplayNameFormat"

	// List related errors
	ListNotFound        Code = "ListNotFound"
	ListIDNull          Code = "ListIDNull"
	ListIDAlreadyExists Code = "ListIDAlreadyExists"
	ListIDFormat        Code = "ListIDFormat"

	// List Accesss related errors
	ListAccessNotFound           Code = "ListAccessNotFound"
	ListAccessAlreadyExists      Code = "ListAccessAlreadyExists"
	ListAccessOwnerAlreadyExists Code = "ListAccessOwnerAlreadyExists"

	// Task related errors
	TaskNotFound Code = "TaskNotFound"

	// Secret related errors
	UserSecretNotFound          Code = "UserSecretNotFound"
	InvalidSecretIDFormat       Code = "InvalidSecretIDFormat"
	InvalidSecretScope          Code = "InvalidSecretScope"
	InvalidUserSecretPassFormat Code = "InvalidUserSecretPassFormat"
)
