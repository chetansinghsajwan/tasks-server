package services

import (
	"context"
	"log"
	"tasks/errorcodes"
	"tasks/store"

	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptUserEncryptionCost = bcrypt.DefaultCost
)

type ServiceError struct {
	Code errorcodes.Code
}

var ST store.Store

type ServiceContext struct {
	Ctx    context.Context
	UserID string
}

type User struct {
	ID          string
	Email       string
	FullName    string
	DisplayName *string
}

type CreateUserParams struct {
	ID          string
	Email       string
	FullName    string
	DisplayName *string
	Pass        string
}

type UpdateUserParams struct {
	ID          *string
	Email       *string
	FullName    *string
	DisplayName **string
}

func CreateUser(ctx ServiceContext, args CreateUserParams) *ServiceError {

	var serr *store.StoreError
	serr = ST.CreateUser(ctx.Ctx, store.CreateUserParams{
		ID:          args.ID,
		Email:       args.Email,
		FullName:    args.FullName,
		DisplayName: args.DisplayName,
	})

	if serr != nil {

		switch serr.Code {

		case errorcodes.UserIDNull,
			errorcodes.UserIDAlreadyExists,
			errorcodes.InvalidUserIDFormat,
			errorcodes.UserEmailNull,
			errorcodes.UserEmailAlreadyExists,
			errorcodes.InvalidUserEmailFormat,
			errorcodes.InvalidUserFullNameFormat,
			errorcodes.InvalidUserDisplayNameFormat:

			return &ServiceError{
				Code: serr.Code,
			}

		default:
			log.Printf("SERVICE: CreateUser: unexpected error: %v", serr)

			return &ServiceError{
				Code: errorcodes.Internal,
			}
		}
	}

	// Encrypt the password
	var err error
	var hashedPass []byte
	hashedPass, err = bcrypt.GenerateFromPassword(
		[]byte(args.Pass), BcryptUserEncryptionCost)

	if err != nil {

		return &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	// Create user secrets
	serr = ST.CreateUserSecret(ctx.Ctx, store.CreateUserSecretParams{
		ID:   args.ID,
		Pass: string(hashedPass),
	})

	if serr != nil {

		switch serr.Code {

		case errorcodes.InvalidUserSecretPassFormat:
			return &ServiceError{
				Code: serr.Code,
			}

		default:

			log.Printf("SERVICE: CreateUser: unexpected error: %v", serr)

			return &ServiceError{
				Code: errorcodes.Internal,
			}
		}
	}

	return nil
}

func GetUser(ctx ServiceContext, id string) (*User, *ServiceError) {

	var user *store.User
	var serr *store.StoreError
	if user, serr = ST.GetUser(ctx.Ctx, id); serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound:
			return nil, &ServiceError{
				Code: serr.Code,
			}

		default:
			return nil, &ServiceError{
				Code: errorcodes.Internal,
			}
		}
	}

	return &User{
		ID:          user.ID,
		FullName:    user.FullName,
		DisplayName: user.DisplayName,
	}, nil
}

func UpdateUser(ctx ServiceContext, id string, args UpdateUserParams) *ServiceError {

	var serr *store.StoreError = ST.UpdateUser(ctx.Ctx, id, store.UpdateUserParams{
		ID:          &args.ID,
		Email:       &args.Email,
		FullName:    &args.FullName,
		DisplayName: args.DisplayName,
	})

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound,
			errorcodes.UserIDNull,
			errorcodes.UserIDAlreadyExists,
			errorcodes.InvalidUserIDFormat,
			errorcodes.UserEmailNull,
			errorcodes.UserEmailAlreadyExists,
			errorcodes.InvalidUserEmailFormat,
			errorcodes.InvalidUserFullNameFormat,
			errorcodes.InvalidUserDisplayNameFormat:

			return &ServiceError{
				Code: serr.Code,
			}
		}

		return &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return nil
}

func DeleteUser(ctx ServiceContext, id string) *ServiceError {

	var serr *store.StoreError
	if serr = ST.DeleteUser(ctx.Ctx, id); serr != nil {

		return &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	if serr = ST.DeleteUserSecret(ctx.Ctx, id); serr != nil {

		return &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return nil
}
