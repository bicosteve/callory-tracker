package mock

import (
	"time"

	"github.com/bicosteve/callory-tracker/pkg/models"
)

var mockUser = &models.User{
	ID:        1,
	Username:  "Alice",
	Email:     "alice@example.com",
	Password:  "password",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

type UserModel struct{}

func (m *UserModel) RegisterUser(username, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) LoginUser(email, password string) (int, error) {
	if email == "alice@example.com" && password == "password" {
		return 1, nil
	}
	return 0, models.ErrorInvalidCredentials
}

func (m *UserModel) GetUserDetails(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
