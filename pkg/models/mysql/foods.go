package mysql

import (
	"database/sql"

	"github.com/bicosteve/callory-tracker/pkg/models"
)

type FoodModel struct {
	DB *sql.DB
}

// InsertFood(): insert food into db
func (f *FoodModel) InsertFood(meal string, name string, protein int, carbohydrate int, fat int, calories int, userId int) (int, error) {
	stm := `INSERT INTO foods
				(meal, name, protein, carbohydrate,fat,calories,created_at,updated_at,userId) 
			VALUES (?,?,?,?,?,?,NOW(),NOW(),?)`

	result, err := f.DB.Exec(stm, meal, name,
		protein, carbohydrate, fat, calories, userId,
	)
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

	err := row.Scan(&food.ID, &food.Name, &food.Protein,
		&food.Carbohydrates, &food.Fat, &food.Calories,
		&food.CreatedAt, &food.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}

	if err != nil {
		return nil, err
	}

	return food, nil

}

func (f *FoodModel) GetFoods(userid int) ([]*models.Food, error) {
	stm := "SELECT * FROM foods WHERE userId = ? ORDER BY created_at DESC LIMIT 100"
	rows, err := f.DB.Query(stm, userid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	foods := []*models.Food{}

	for rows.Next() {
		f := &models.Food{}

		err = rows.Scan(
			&f.ID, &f.Name, &f.Protein, &f.Carbohydrates,
			&f.Fat, &f.Calories, &f.CreatedAt, &f.UpdatedAt,
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

func (f *FoodModel) UpdateFood(foodId, userId int) (int, error) {
	stm := `UPDATE foods SET name = ?, protein = ?, carbohydrates = ?, fat = ?, calory = ?, updated_at = UTC_TIMESTAMP()  WHERE id = ? AND userId = ?`
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