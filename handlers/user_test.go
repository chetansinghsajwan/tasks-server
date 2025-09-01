package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"os"
	"tasks/db"
	"tasks/routes"
	"tasks/services"
	"tasks/store/pg"
	"testing"
)

type CreateUserRequest struct {
	ID          string  `json:"id"`
	Pass        string  `json:"pass"`
	Email       string  `json:"email"`
	FullName    string  `json:"full_name"`
	DisplayName *string `json:"display_name"`
}

type CreateUserResponse struct {
	User *struct {
		ID string `json:"id"`
	} `json:"user,omitempty"`

	Error *struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Hint string `json:"hint,omitempty"`
	} `json:"error,omitempty"`
}

type GetUserResponse struct {
	User *struct {
		ID          string  `json:"id"`
		Email       string  `json:"email"`
		FullName    string  `json:"full_name"`
		DisplayName *string `json:"display_name"`
	} `json:"user,omitempty"`

	Error *struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Hint string `json:"hint,omitempty"`
	} `json:"error,omitempty"`
}

const (
	validUserID_0          = "johndoe-0"
	validUserID_1          = "johndoe-1"
	validUserPass_0        = "password123"
	validUserEmail_0       = "myemail@domain.com"
	validUserFullName_0    = "John Doe"
	validUserDisplayName_0 = "John"
)

func TestUserHandler(t *testing.T) {

	os.Setenv("DATABASE_URL", "postgres://devuser:devpass@database:5432/testdb?sslmode=disable")
	db.Init()

	db.Pool.Exec(context.Background(), "select reinitialize_schema()")

	services.ST = pg.PostgresStore{
		Pool: db.Pool,
	}

	var router = routes.SetupRouter()

	t.Run("Create user with invalid ids", func(t *testing.T) {

		var invalidIDs = []string{
			"",  // empty
			" ", // whitespace only
			"a", // too short

			// We need to fix this validation check, right now i think that
			// the long value is being truncated by the database
			// "this-id-is-way-too-long-for-the-validation-check",

			"-abc",    // starts with hyphen
			"abc-",    // ends with hyphen
			"Abc",     // uppercase not allowed
			"abc@123", // special char @
			"abc#def", // special char #
			"abc.def", // dot not allowed
			"abc_def", // underscore not allowed
			"abc😀",    // emoji
			"हिंदी",   // non-latin
			"12--34",  // double hyphen
			"abc ",    // trailing space
			" abc",    // leading space
			"a b c",   // internal whitespace
		}

		for _, invalidID := range invalidIDs {

			// Capture the range var
			var invalidID = invalidID

			t.Run(fmt.Sprintf("Create user with invalid id '%s'", invalidID), func(t *testing.T) {

				t.Parallel()

				var reqBody = fmt.Sprintf(
					`{
						"id": "%s",
						"pass": "%s",
						"email": "%s",
						"full_name": "%s",
						"display_name": "%s"
					}`,
					invalidID,
					validUserPass_0,
					validUserEmail_0,
					validUserFullName_0,
					validUserDisplayName_0,
				)

				fmt.Printf("Request body: %s\n", reqBody)

				var req = httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")

				var rec = httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				if rec.Code != http.StatusBadRequest {
					t.Fatalf("Expected status '%v', got '%d'", http.StatusBadRequest, rec.Code)
				}

				var resBody CreateUserResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
					t.Fatalf("Failed to parse response body, err: %v, body: %s", err, rec.Body.String())
				}

				if resBody.Error == nil {
					t.Fatalf("Expected error object, got none")
				}

				if resBody.Error.Code != "InvalidUserIDFormat" {
					t.Fatalf("Expected error code '%s', got '%s'", "InvalidUserIDFormat", resBody.Error.Code)
				}
			})
		}
	})

	t.Run("Create user with invalid emails", func(t *testing.T) {

		var invalidEmails = []string{
			"",                      // empty
			"   ",                   // spaces only
			"plainaddress",          // missing @
			"username.domain.com",   // missing @
			"user@com",              // missing dot
			"user@.com",             // starts with dot
			"user@example.",         // ends with dot
			"user@example.c",        // TLD too short
			"user@@example.com",     // double @
			"user@exa_mple.com",     // underscore in domain
			"user@exam$ple.com",     // invalid char in domain
			"user name@example.com", // space in local
			"user@exam😀ple.com",     // emoji
			".user@example.com",     // starts with dot
			"user.@example.com",     // ends with dot before @
			"USER@EXAMPLE.COM",      // uppercase (passes in your case because of ~*, but worth noting)
		}

		for _, invalidEmail := range invalidEmails {

			// Capture the range var
			var invalidEmail = invalidEmail

			t.Run(fmt.Sprintf("Create user with invalid email '%s'", invalidEmail), func(t *testing.T) {

				t.Parallel()

				var reqBody = fmt.Sprintf(
					`{
						"id": "%s",
						"pass": "%s",
						"email": "%s",
						"full_name": "%s",
						"display_name": "%s"
					}`,
					validUserID_0,
					validUserPass_0,
					invalidEmail,
					validUserFullName_0,
					validUserDisplayName_0,
				)

				fmt.Printf("Request body: %s\n", reqBody)

				var req = httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")

				var rec = httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				if rec.Code != http.StatusBadRequest {
					t.Fatalf("Expected status '%v', got '%d'", http.StatusBadRequest, rec.Code)
				}

				var resBody CreateUserResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
					t.Fatalf("Failed to parse response body, err: %v, body: %s", err, rec.Body.String())
				}

				if resBody.Error == nil {
					t.Fatalf("Expected error object, got none")
				}

				if resBody.Error.Code != "InvalidUserEmailFormat" {
					t.Fatalf("Expected error code '%s', got '%s'", "InvalidUserEmailFormat", resBody.Error.Code)
				}
			})
		}
	})

	t.Run("Create user with invalid full names", func(t *testing.T) {

		var invalidFullNames = []string{
			"",          // empty
			" ",         // single space
			"   ",       // multiple spaces
			" John Doe", // leading space
			"John Doe ", // trailing space
			"  Alice",   // multiple leading spaces
			"Bob  ",     // trailing spaces
		}

		for _, invalidFullName := range invalidFullNames {

			// Capture the range var
			var invalidFullName = invalidFullName

			t.Run(fmt.Sprintf("Create user with invalid full name '%s'", invalidFullName), func(t *testing.T) {

				t.Parallel()

				var reqBody = fmt.Sprintf(
					`{
						"id": "%s",
						"pass": "%s",
						"email": "%s",
						"full_name": "%s",
						"display_name": "%s"
					}`,
					validUserID_0,
					validUserPass_0,
					validUserEmail_0,
					invalidFullName,
					validUserDisplayName_0,
				)

				fmt.Printf("Request body: %s\n", reqBody)

				var req = httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")

				var rec = httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				if rec.Code != http.StatusBadRequest {
					t.Fatalf("Expected status '%v', got '%d'", http.StatusBadRequest, rec.Code)
				}

				var resBody CreateUserResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
					t.Fatalf("Failed to parse response body, err: %v, body: %s", err, rec.Body.String())
				}

				if resBody.Error == nil {
					t.Fatalf("Expected error object, got none")
				}

				if resBody.Error.Code != "InvalidUserFullNameFormat" {
					t.Fatalf("Expected error code '%s', got '%s'", "InvalidUserFullNameFormat", resBody.Error.Code)
				}
			})
		}
	})

	t.Run("Create user with invalid display names", func(t *testing.T) {

		var invalidDisplayNames = []string{
			"",        // empty
			" ",       // single space
			"   ",     // multiple spaces
			" John",   // leading space
			"John ",   // trailing space
			"  Alice", // multiple leading spaces
			"Bob  ",   // multiple trailing spaces
		}

		for _, invalidDisplayName := range invalidDisplayNames {

			// Capture the range var
			var invalidDisplayName = invalidDisplayName

			t.Run(fmt.Sprintf("Create user with invalid display name '%s'", invalidDisplayName), func(t *testing.T) {

				t.Parallel()

				var reqBody = fmt.Sprintf(
					`{
						"id": "%s",
						"pass": "%s",
						"email": "%s",
						"full_name": "%s",
						"display_name": "%s"
					}`,
					validUserID_0,
					validUserPass_0,
					validUserEmail_0,
					validUserFullName_0,
					invalidDisplayName,
				)

				fmt.Printf("Request body: %s\n", reqBody)

				var req = httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")

				var rec = httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				if rec.Code != http.StatusBadRequest {
					t.Fatalf("Expected status '%v', got '%d'", http.StatusBadRequest, rec.Code)
				}

				var resBody CreateUserResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
					t.Fatalf("Failed to parse response body, err: %v, body: %s", err, rec.Body.String())
				}

				if resBody.Error == nil {
					t.Fatalf("Expected error object, got none")
				}

				if resBody.Error.Code != "InvalidUserDisplayNameFormat" {
					t.Fatalf("Expected error code '%s', got '%s'", "InvalidUserDisplayNameFormat", resBody.Error.Code)
				}
			})
		}
	})

	t.Run("Create user with invalid password", func(t *testing.T) {

		// TODO: We need to first implement password validation logic (in services maybe?)
		var invalidPasswords = []string{
			// "",
			// "  ",
		}

		for _, invalidPassword := range invalidPasswords {

			// Capture the range var
			var invalidPassword = invalidPassword

			t.Run(fmt.Sprintf("Create user with invalid password '%s'", invalidPassword), func(t *testing.T) {

				t.Parallel()

				var reqBody = fmt.Sprintf(
					`{
						"id": "%s",
						"pass": "%s",
						"email": "%s",
						"full_name": "%s",
						"display_name": "%s"
					}`,
					validUserID_0,
					invalidPassword,
					validUserEmail_0,
					validUserFullName_0,
					validUserDisplayName_0,
				)

				fmt.Printf("Request body: %s\n", reqBody)

				var req = httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(reqBody))
				req.Header.Set("Content-Type", "application/json")

				var rec = httptest.NewRecorder()

				router.ServeHTTP(rec, req)

				if rec.Code != http.StatusBadRequest {
					t.Fatalf("Expected status '%v', got '%d'", http.StatusBadRequest, rec.Code)
				}

				var resBody CreateUserResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
					t.Fatalf("Failed to parse response body, err: %v, body: %s", err, rec.Body.String())
				}

				if resBody.Error == nil {
					t.Fatalf("Expected error object, got none")
				}

				if resBody.Error.Code != "InvalidUserSecretPassFormat" {
					t.Fatalf("Expected error code '%s', got '%s'", "InvalidUserSecretPassFormat", resBody.Error.Code)
				}
			})
		}
	})

	t.Run("Create user with valid data", func(t *testing.T) {

		var reqBody = fmt.Sprintf(
			`{
				"id": "%s",
				"pass": "%s",
				"email": "%s",
				"full_name": "%s",
				"display_name": "%s"
			}`,
			validUserID_0,
			validUserPass_0,
			validUserEmail_0,
			validUserFullName_0,
			validUserDisplayName_0,
		)

		fmt.Printf("Request body: %s\n", reqBody)

		var req = httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		var rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("Expected status '%v', got '%d'", http.StatusCreated, rec.Code)
		}
	})

	t.Run("Create user with duplicate id", func(t *testing.T) {

		var reqBody = fmt.Sprintf(
			`{
				"id": "%s",
				"pass": "%s",
				"email": "%s",
				"full_name": "%s",
				"display_name": "%s"
			}`,
			validUserID_0,
			validUserPass_0,
			validUserEmail_0,
			validUserFullName_0,
			validUserDisplayName_0,
		)

		fmt.Printf("Request body: %s\n", reqBody)

		var req = httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		var rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("Expected status '%v', got '%d'", http.StatusBadRequest, rec.Code)
		}

		var resBody CreateUserResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
			t.Fatalf("Failed to parse response body, err: %v, body: %s", err, rec.Body.String())
		}

		if resBody.Error == nil {
			t.Fatalf("Expected error object, got none")
		}

		if resBody.Error.Code != "UserIDAlreadyExists" {
			t.Fatalf("Expected error code '%s', got '%s'", "UserIDAlreadyExists", resBody.Error.Code)
		}
	})

	t.Run("Get user that exists", func(t *testing.T) {

		var req = httptest.NewRequest("GET", fmt.Sprintf("/api/users/%s", validUserID_0), nil)
		var rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status '%v', got '%d'", http.StatusOK, rec.Code)
		}

		var resBody GetUserResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
			t.Fatalf("Failed to parse response body, err: %v, body: %s", err, rec.Body.String())
		}

		if resBody.Error != nil {
			t.Fatalf("Expected no error, got Code: %s, Msg: %s", resBody.Error.Code, resBody.Error.Msg)
		}

		if resBody.User == nil {
			t.Fatalf("Expected user got nil")
		}

		if resBody.User.ID != validUserID_0 {
			t.Fatalf("Expected user id '%s', got '%s'", validUserID_0, resBody.User.ID)
		}

		if resBody.User.Email != validUserEmail_0 {
			t.Fatalf("Expected user email '%s', got '%s'", validUserEmail_0, resBody.User.Email)
		}

		if resBody.User.FullName != validUserFullName_0 {
			t.Fatalf("Expected user full name '%s', got '%s'", validUserFullName_0, resBody.User.FullName)
		}

		if resBody.User.DisplayName == nil {
			t.Fatalf("Expected user display name '%s', got 'nil'", validUserDisplayName_0)
		}

		if *resBody.User.DisplayName != validUserDisplayName_0 {
			t.Fatalf("Expected user display name '%s', got '%s'", validUserDisplayName_0, *resBody.User.DisplayName)
		}
	})

	t.Run("Get user that doesn't exist", func(t *testing.T) {

		var req = httptest.NewRequest("GET", fmt.Sprintf("/api/users/%s", validUserID_1), nil)
		var rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("Expected status '%v', got '%d'", http.StatusNotFound, rec.Code)
		}

		var resBody GetUserResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
			t.Fatalf("Failed to parse response body, err: %v, body: %s", err, rec.Body.String())
		}

		if resBody.Error == nil {
			t.Fatalf("Expected error object, got none")
		}

		if resBody.Error.Code != "UserNotFound" {
			t.Fatalf("Expected error code '%s', got '%s'", "UserNotFound", resBody.Error.Code)
		}
	})

	t.Run("Delete user that exists", func(t *testing.T) {

		var req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/users/%s", validUserID_0), nil)
		var rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("Expected status '%v', got '%d'", http.StatusOK, rec.Code)
		}
	})

	t.Run("Delete user that doesn't exist", func(t *testing.T) {

		var req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/users/%s", validUserID_0), nil)
		var rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("Expected status '%v', got '%d'", http.StatusNotFound, rec.Code)
		}

		var resBody GetUserResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
			t.Fatalf("Failed to parse response body, err: %v, body: %s", err, rec.Body.String())
		}

		if resBody.Error == nil {
			t.Fatalf("Expected error object, got none")
		}

		if resBody.Error.Code != "UserNotFound" {
			t.Fatalf("Expected error code '%s', got '%s'", "UserNotFound", resBody.Error.Code)
		}
	})
}
