package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bicosteve/callory-tracker/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) RegisterUser(username, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_ = hash

	stm := `INSERT INTO users (username,email,hashed_password,created_at, updated_at)
			VALUES (?,?,?,NOW(),NOW())`

	_, err = u.DB.Exec(stm, strings.Title(username), email, string(hash))
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "Duplicate entry") {
			return models.ErrDuplicateEmail
		}

		_ = ok
	}

	return err
}

func (u *UserModel) LoginUser(email, password string) (int, error) {
	var id int
	var hash []byte

	stm := `SELECT id, hashed_password FROM users WHERE email = ?`
	row := u.DB.QueryRow(stm, email)
	err := row.Scan(&id, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, models.ErrorInvalidCredentials
	}

	if err != nil {
		return 0, err
	}

	// compare provided password and hashed password
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, models.ErrorInvalidCredentials
	}

	if err != nil {
		return 0, err
	}

	// Match is correct
	return int(id), nil
}

func (u *UserModel) GetUserDetails(id int) (*models.User, error) {
	user := &models.User{}
	stm := "SELECT id,username,email,created_at FROM users WHERE id = ?"
	err := u.DB.QueryRow(stm, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNoRecord
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

//func (u *UserModel) Exists(id int) (bool, error) {
//	var exists bool
//
//	stm := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
//	err := u.DB.QueryRow(stm, id).Scan(&exists)
//	return exists, err
//}
