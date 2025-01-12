package userrepository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nickolasgough/cloud-9-iam/internal/shared/ierrors"
	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
)

func (r *repository) GetUser(userID string) (*usermodel.User, error) {
	sqlStr := getUserBaseStatement("id")
	return r.getUser(sqlStr, userID)
}

func (r *repository) GetUserByEmail(email string) (*usermodel.User, error) {
	sqlStr := getUserBaseStatement("email")
	return r.getUser(sqlStr, email)
}

func (r *repository) getUser(sqlStr string, identifier string) (*usermodel.User, error) {
	var id, firstName, lastName, displayName, displayImageURL, email string
	var created, updated time.Time
	row := r.database.QueryRow(sqlStr, identifier)
	err := row.Scan(&id, &firstName, &lastName, &displayName, &displayImageURL, &email, &created, &updated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ierrors.NewError(ierrors.NotFound, err)
		}
		return nil, err
	}
	user := &usermodel.User{
		ID:              id,
		FirstName:       firstName,
		LastName:        lastName,
		DisplayName:     displayName,
		DisplayImageURL: displayImageURL,
		Email:           email,
		Created:         created,
		Updated:         updated,
	}
	return user, nil
}

func getUserBaseStatement(field string) string {
	stmnt := `
		SELECT
			id,
			firstName,
			lastName,
			displayName,
			displayImageURL,
			email,
			created,
			updated
		FROM %[1]s
		WHERE
			%[2]s = $1
	`
	formattedStmnt := fmt.Sprintf(stmnt, usermodel.UsersTable, field)
	return formattedStmnt
}

func (r *repository) GetUserAndPasswordByEmail(email string) (*usermodel.User, string, error) {
	sqlStr := getUserAndPasswordByEmail()
	user, password, err := r.getUserAndPassword(sqlStr, email)
	if err != nil {
		return nil, "", err
	}
	// TODO: set a random password for a user created by Google sign-in so this isn't necessary
	if password == "" {
		return nil, "", ierrors.NewError(ierrors.InvalidArgument, errors.New("user not verified"))
	}
	return user, password, nil
}

func (r *repository) getUserAndPassword(sqlStr string, identifier string) (*usermodel.User, string, error) {
	var id, firstName, lastName, displayName, displayImageURL, email string
	var created, updated time.Time
	var password string
	row := r.database.QueryRow(sqlStr, identifier)
	err := row.Scan(&id, &firstName, &lastName, &displayName, &displayImageURL, &email, &created, &updated, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", ierrors.NewError(ierrors.NotFound, err)
		}
		return nil, "", err
	}
	user := &usermodel.User{
		ID:              id,
		FirstName:       firstName,
		LastName:        lastName,
		DisplayName:     displayName,
		DisplayImageURL: displayImageURL,
		Email:           email,
		Created:         created,
		Updated:         updated,
	}
	return user, "", nil
}

func getUserAndPasswordByEmail() string {
	statement := `
		SELECT
		id,
			firstName,
			lastName,
			displayName,
			displayImageURL,
			email,
			created,
			updated,
			password
		FROM %[1]s
		WHERE
			email = $1
	`
	formattedStatement := fmt.Sprintf(statement, usermodel.UsersTable)
	return formattedStatement
}
