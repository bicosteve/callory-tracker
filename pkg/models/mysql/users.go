package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bicosteve/callory-tracker/pkg/logger"
	"github.com/bicosteve/callory-tracker/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB  *sql.DB
	Log *logger.Logger
}

// log returns the model's logger, falling back to a discarding logger when one
// has not been wired up (for example in tests). This keeps every DB method
// nil-safe while still emitting useful output in production.
func (u *UserModel) log() *logger.Logger {
	if u.Log == nil {
		return logger.Discard()
	}
	return u.Log
}

func (u *UserModel) RegisterUser(username, email, password string) error {
	u.log().Info.Printf("RegisterUser: registering email=%q", email)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.log().Error.Printf("RegisterUser: hashing password failed for email=%q: %v", email, err)
		return err
	}

	stm := `INSERT INTO users (username,email,hashed_password,created_at, updated_at)
		VALUES (?,?,?,NOW(),NOW())`

	_, err = u.DB.Exec(stm, strings.Title(username), email, string(hash))
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "Duplicate entry") {
			u.log().Warning.Printf("RegisterUser: duplicate email=%q", email)
			return models.ErrDuplicateEmail
		}

		u.log().Error.Printf("RegisterUser: exec failed for email=%q: %v", email, err)
		return err
	}

	u.log().Info.Printf("RegisterUser: registered email=%q", email)
	return nil
}

func (u *UserModel) LoginUser(email, password string) (int, error) {
	var id int
	var hash []byte

	u.log().Info.Printf("LoginUser: authenticating email=%q", email)

	stm := `SELECT id, hashed_password FROM users WHERE email = ?`
	row := u.DB.QueryRow(stm, email)
	err := row.Scan(&id, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		u.log().Warning.Printf("LoginUser: no user for email=%q", email)
		return 0, models.ErrorInvalidCredentials
	}

	if err != nil {
		u.log().Error.Printf("LoginUser: scan failed for email=%q: %v", email, err)
		return 0, err
	}

	// compare provided password and hashed password
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		u.log().Warning.Printf("LoginUser: invalid password for email=%q", email)
		return 0, models.ErrorInvalidCredentials
	}

	if err != nil {
		u.log().Error.Printf("LoginUser: password compare failed for email=%q: %v", email, err)
		return 0, err
	}

	// Match is correct
	u.log().Info.Printf("LoginUser: authenticated userId=%d email=%q", id, email)
	return int(id), nil
}

func (u *UserModel) GetUserDetails(id int) (*models.User, error) {
	user := &models.User{}

	u.log().Info.Printf("GetUserDetails: fetching userId=%d", id)

	stm := "SELECT id,username,email,created_at FROM users WHERE id = ?"
	err := u.DB.QueryRow(stm, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		u.log().Warning.Printf("GetUserDetails: no record for userId=%d", id)
		return nil, models.ErrNoRecord
	}

	if err != nil {
		u.log().Error.Printf("GetUserDetails: scan failed for userId=%d: %v", id, err)
		return nil, err
	}

	u.log().Info.Printf("GetUserDetails: found userId=%d", id)
	return user, nil
}
