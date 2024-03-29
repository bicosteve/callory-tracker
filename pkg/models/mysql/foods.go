package mysql

import (
	"database/sql"
	"errors"
	"github.com/bicosteve/callory-tracker/pkg/models"
	"strings"
)

type FoodModel struct {
	DB *sql.DB
}

// InsertFood(): insert food into db
func (f *FoodModel) InsertFood(
	meal string, name string, protein int, carbohydrate int,
	fat int, calories int, userId int,
) (int, error) {
	stm := `INSERT INTO foods
				(meal, name, protein, carbohydrate,fat,calories,created_at,updated_at,userId) 
			VALUES (?,?,?,?,?,?,NOW(),NOW(),?)`

	result, err := f.DB.Exec(stm, strings.Title(meal), strings.Title(name),
		protein, carbohydrate, fat, calories, userId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	// will return the last inserted id on the table
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (f *FoodModel) GetFood(foodId, userId int) (*models.Food, error) {
	stm := `SELECT * FROM foods WHERE id = ? AND userId = ? LIMIT 1`
	row := f.DB.QueryRow(stm, foodId, userId)

	food := &models.Food{}

	err := row.Scan(&food.ID, &food.Meal, &food.Name, &food.Protein,
		&food.Carbohydrates, &food.Fat, &food.Calories,
		&food.CreatedAt, &food.UpdatedAt, &food.UserID)

	if errors.Is(err, sql.ErrNoRows) {

		return nil, models.ErrNoRecord
	}

	if err != nil {

		return nil, err
	}

	return food, nil
}

func (f *FoodModel) GetFoodTotal(
	userId int, createdAt string,
) (*models.Food, error) {
	total := &models.Food{}
	defer f.DB.Close()
	stm := `SELECT SUM(protein), SUM(carbohydrate), SUM(fat), SUM(calories) 
			FROM foods WHERE userId = ? AND created_at LIKE CONCAT('%',?)`
	row := f.DB.QueryRow(stm, userId, createdAt)
	err := row.Scan(&total.Protein, &total.Carbohydrates, &total.Fat, &total.Calories)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNoRecord
	}

	if err != nil {
		return nil, err
	}

	return total, nil
}

func (f *FoodModel) GetFoods(userId int) ([]*models.Food, error) {
	stm := "SELECT * FROM foods WHERE userId = ? ORDER BY created_at DESC LIMIT 100"
	rows, err := f.DB.Query(stm, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var foods []*models.Food

	for rows.Next() {
		f := &models.Food{}

		err = rows.Scan(
			&f.ID, &f.Meal, &f.Name, &f.Protein, &f.Carbohydrates,
			&f.Fat, &f.Calories, &f.CreatedAt, &f.UpdatedAt, &f.UserID,
		)

		if err != nil {
			return nil, err
		}

		foods = append(foods, f)
	}

	// retries any rows error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return foods, nil
}

func (f *FoodModel) UpdateFood(
	meal string, name string, protein, cabs, fat, calory, foodId, userId int,
) (int, error) {
	stm := `UPDATE foods SET meal = ?, name = ?, protein = ?, carbohydrate = ?, fat = ?, calories = ?, updated_at = UTC_TIMESTAMP()  WHERE id = ? AND userId = ?`
	result, err := f.DB.Exec(stm, meal, name, protein, cabs, fat, calory, foodId, userId)
	if err != nil {
		return 0, err
	}

	id, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (f *FoodModel) DeleteFood(foodId, userId int) (int, error) {
	stm := "DELETE FROM foods WHERE id = ? AND userId = ?"
	result, err := f.DB.Exec(stm, foodId, userId)
	if err != nil {
		return 0, err
	}

	id, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
