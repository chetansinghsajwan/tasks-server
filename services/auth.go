package services

import (
	"fmt"
	"tasks/errorcodes"
	"tasks/store"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	JwtSigningKey             = ""
	AccessTokenCookiePath     = ""
	AccessTokenCookieDomain   = ""
	AccessTokenCookieLifetime = 3600
)

type AuthToken struct {
	UserID string
}

type LoginParams struct {
	UserID   string
	Password string
}

func AuthenticateToken(tokenStr string) (*AuthToken, *ServiceError) {

	var err error
	var token *jwt.Token
	token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		// Make sure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return JwtSigningKey, nil
	})

	if err != nil || !token.Valid {

		return nil, &ServiceError{
			Code: errorcodes.InvalidToken,
		}
	}

	// Extract username if needed
	var userID string
	if userID, err = token.Claims.GetSubject(); err != nil {

		return nil, &ServiceError{
			Code: errorcodes.InvalidToken,
		}
	}

	if userID == "" {

		return nil, &ServiceError{
			Code: errorcodes.InvalidToken,
		}
	}

	return &AuthToken{
		UserID: userID,
	}, nil
}

func GenerateToken(userID string) (string, error) {

	var claims = jwt.NewWithClaims(
		jwt.SigningMethodPS256.SigningMethodRSA,
		jwt.MapClaims{
			"id":         userID,
			"expires_at": time.Now().Add(time.Hour * 2),
		},
	)

	var token, err = claims.SignedString(JwtSigningKey)

	if err != nil {
		return "", err
	}

	return token, nil
}

func LoginUser(ctx ServiceContext, args LoginParams) (string, *ServiceError) {

	var serr *store.StoreError
	var secret *store.UserSecret
	secret, serr = ST.GetUserSecret(ctx.Ctx, args.UserID)

	if serr != nil {

		switch serr.Code {

		case errorcodes.UserNotFound:
			return "", &ServiceError{
				Code: serr.Code,
			}

		default:
			return "", &ServiceError{
				Code: errorcodes.Internal,
			}
		}
	}

	var err error = bcrypt.CompareHashAndPassword(
		[]byte(secret.Pass), []byte(args.Password))

	if err != nil {

		return "", &ServiceError{
			Code: errorcodes.AuthMatchFail,
		}
	}

	var claims = jwt.NewWithClaims(
		jwt.SigningMethodPS256.SigningMethodRSA,
		jwt.MapClaims{
			"id":         args.UserID,
			"expires_at": time.Now().Add(time.Hour * 2),
		},
	)

	var token string
	token, err = claims.SignedString(JwtSigningKey)

	if err != nil {

		return "", &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return token, nil
}
