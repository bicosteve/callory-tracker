package mock

import (
	"time"

	"github.com/bicosteve/callory-tracker/pkg/models"
)

// mockFood is a canned food record returned by the mock for known IDs.
var mockFood = &models.Food{
	ID:            1,
	Meal:          "Lunch",
	Name:          "Rice",
	Protein:       10,
	Carbohydrates: 40,
	Fat:           5,
	Calories:      220,
	CreatedAt:     time.Now(),
	UpdatedAt:     time.Now(),
	UserID:        1,
}

// FoodModel is a mock implementation of models.FoodModelInterface used in tests.
type FoodModel struct{}

func (m *FoodModel) InsertFood(meal string, name string, protein int, carbohydrate int, fat int, calories int, userId int) (int, error) {
	return 2, nil
}

func (m *FoodModel) GetFood(foodId, userId int) (*models.Food, error) {
	switch foodId {
	case 1:
		return mockFood, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *FoodModel) GetFoodTotal(userId int, createdAt string) (*models.Food, error) {
	return &models.Food{
		Protein:       10,
		Carbohydrates: 40,
		Fat:           5,
		Calories:      220,
	}, nil
}

func (m *FoodModel) GetFoods(userId int) ([]*models.Food, error) {
	return []*models.Food{mockFood}, nil
}

func (m *FoodModel) UpdateFood(meal string, name string, protein, cabs, fat, calory, foodId, userId int) (int, error) {
	return 1, nil
}

func (m *FoodModel) DeleteFood(foodId, userId int) (int, error) {
	return 1, nil
}
