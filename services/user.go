package services

import (
	"context"
	"errors"
	"tasks/errorcodes"
	"tasks/option"
	"tasks/store"

	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptUserEncryptionCost = bcrypt.DefaultCost
)

var ST store.Store

type ServiceContext struct {
	ctx    context.Context
	string string
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

func CreateUser(ctx ServiceContext, args CreateUserParams) error {

	var serr *store.StoreError
	serr = ST.CreateUser(ctx.ctx, store.CreateUserParams{
		ID:          args.ID,
		Email:       args.Email,
		FullName:    args.FullName,
		DisplayName: option.FromPtr(args.DisplayName),
	})

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserIDNull,
			errorcodes.UserIDAlreadyExists,
			errorcodes.UserIDFormat,
			errorcodes.UserEmailNull,
			errorcodes.UserEmailAlreadyExists,
			errorcodes.UserEmailFormat,
			errorcodes.UserFullNameFormat,
			errorcodes.UserDisplayNameFormat:
			return errors.New(serr.Msg)

		default:
			return errors.New("internal server error")
		}
	}

	// Encrypt the password
	var err error
	var hashedPass []byte
	hashedPass, err = bcrypt.GenerateFromPassword(
		[]byte(args.Pass), BcryptUserEncryptionCost)

	if err != nil {

		return errors.New(err.Error())
	}

	// Create user secrets
	serr = ST.CreateUserSecret(ctx.ctx, store.CreateUserSecretParams{
		ID:   args.ID,
		Pass: string(hashedPass),
	})

	if serr != nil {

		return errors.New(serr.Msg)
	}

	return nil
}

func GetUser(ctx ServiceContext, id string) (*User, error) {

	var user *store.User
	var serr *store.StoreError
	if user, serr = ST.GetUser(ctx.ctx, id); serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound:
			return nil, errors.New(serr.Msg)

		default:
			return nil, errors.New("internal server error")
		}
	}

	return &User{
		ID:          user.ID,
		FullName:    user.FullName,
		DisplayName: user.DisplayName.Ptr(),
	}, nil
}

func UpdateUser(ctx ServiceContext, id string, args UpdateUserParams) error {

	var serr *store.StoreError = ST.UpdateUser(ctx.ctx, id, store.UpdateUserParams{
		ID:          option.Some(args.ID),
		Email:       option.Some(args.Email),
		FullName:    option.Some(args.FullName),
		DisplayName: option.FromPtr(args.DisplayName),
	})

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound,
			errorcodes.UserIDNull,
			errorcodes.UserIDAlreadyExists,
			errorcodes.UserIDFormat,
			errorcodes.UserEmailNull,
			errorcodes.UserEmailAlreadyExists,
			errorcodes.UserEmailFormat,
			errorcodes.UserFullNameFormat,
			errorcodes.UserDisplayNameFormat:

			return errors.New(serr.Msg)
		}

		return errors.New(serr.Msg)
	}

	return nil
}

func DeleteUser(ctx ServiceContext, id string) error {

	var serr *store.StoreError
	if serr = ST.DeleteUser(ctx.ctx, id); serr != nil {

		return errors.New(serr.Msg)
	}

	if serr = ST.DeleteUserSecret(ctx.ctx, id); serr != nil {

		return errors.New(serr.Msg)
	}

	return nil
}
