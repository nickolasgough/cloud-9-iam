package userrepository

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
)

func (r *repository) CreateUser(user *usermodel.User, password string) (*usermodel.User, error) {
	user.ID = uuid.NewString()
	now := time.Now().UTC()
	user.Created = now
	user.Updated = now
	sqlStr := insertUserStatement()
	_, err := r.database.Exec(
		sqlStr,
		user.ID,                           // 1
		user.FirstName,                    // 2
		user.LastName,                     // 3
		user.DisplayName,                  // 4
		user.DisplayImageURL,              // 5
		user.Email,                        // 6
		password,                          // 7
		user.Created.Format(time.RFC3339), // 8
		user.Updated.Format(time.RFC3339), // 9
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func insertUserStatement() string {
	statement := `
		INSERT INTO %[1]s 
		(
			id,
			firstName,
			lastName,
			displayName,
			displayImageURL,
			email,
			password,
			created,
			updated
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9
		);
`
	formattedStatement := fmt.Sprintf(
		statement,
		usermodel.UsersTable, // 1
	)
	fmt.Println(formattedStatement)
	return formattedStatement
}
