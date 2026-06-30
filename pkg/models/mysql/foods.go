package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/bicosteve/callory-tracker/pkg/logger"
	"github.com/bicosteve/callory-tracker/pkg/models"
)

type FoodModel struct {
	DB  *sql.DB
	Log *logger.Logger
}

// log returns the model's logger, falling back to a discarding logger when one
// has not been wired up (for example in tests). This keeps every DB method
// nil-safe while still emitting useful output in production.
func (f *FoodModel) log() *logger.Logger {
	if f.Log == nil {
		return logger.Discard()
	}
	return f.Log
}

// InsertFood(): insert food into db
func (f *FoodModel) InsertFood(
	meal string, name string, protein int, carbohydrate int,
	fat int, calories int, userId int,
) (int, error) {
	stm := `INSERT INTO foods
		(meal, name, protein, carbohydrate,fat,calories,created_at,updated_at,userId)
		VALUES (?,?,?,?,?,?,NOW(),NOW(),?)`

	f.log().Info.Printf("InsertFood: inserting meal=%q name=%q for userId=%d", meal, name, userId)

	result, err := f.DB.Exec(stm, strings.Title(meal), strings.Title(name),
		protein, carbohydrate, fat, calories, userId)
	if err != nil {
		f.log().Error.Printf("InsertFood: exec failed for userId=%d: %v", userId, err)
		return 0, err
	}

	id, err := result.LastInsertId()
	// will return the last inserted id on the table
	if err != nil {
		f.log().Error.Printf("InsertFood: LastInsertId failed for userId=%d: %v", userId, err)
		return 0, err
	}

	f.log().Info.Printf("InsertFood: inserted food id=%d for userId=%d", id, userId)
	return int(id), nil
}

func (f *FoodModel) GetFood(foodId, userId int) (*models.Food, error) {
	stm := `SELECT * FROM foods WHERE id = ? AND userId = ? LIMIT 1`
	f.log().Info.Printf("GetFood: fetching foodId=%d for userId=%d", foodId, userId)
	row := f.DB.QueryRow(stm, foodId, userId)

	food := &models.Food{}

	err := row.Scan(&food.ID, &food.Meal, &food.Name, &food.Protein,
		&food.Carbohydrates, &food.Fat, &food.Calories,
		&food.CreatedAt, &food.UpdatedAt, &food.UserID)

	if errors.Is(err, sql.ErrNoRows) {
		f.log().Warning.Printf("GetFood: no record for foodId=%d userId=%d", foodId, userId)
		return nil, models.ErrNoRecord
	}

	if err != nil {
		f.log().Error.Printf("GetFood: scan failed for foodId=%d userId=%d: %v", foodId, userId, err)
		return nil, err
	}

	f.log().Info.Printf("GetFood: found foodId=%d for userId=%d", foodId, userId)
	return food, nil
}

func (f *FoodModel) GetFoodTotal(
	userId int, createdAt string,
) (*models.Food, error) {
	total := &models.Food{}
	defer f.DB.Close()
	stm := `SELECT SUM(protein), SUM(carbohydrate), SUM(fat), SUM(calories)
		FROM foods WHERE userId = ? AND created_at LIKE CONCAT('%',?)`

	f.log().Info.Printf("GetFoodTotal: aggregating totals for userId=%d on %q", userId, createdAt)

	row := f.DB.QueryRow(stm, userId, createdAt)
	err := row.Scan(&total.Protein, &total.Carbohydrates, &total.Fat, &total.Calories)

	if errors.Is(err, sql.ErrNoRows) {
		f.log().Warning.Printf("GetFoodTotal: no records for userId=%d on %q", userId, createdAt)
		return nil, models.ErrNoRecord
	}

	if err != nil {
		f.log().Error.Printf("GetFoodTotal: scan failed for userId=%d: %v", userId, err)
		return nil, err
	}

	f.log().Info.Printf("GetFoodTotal: computed totals for userId=%d", userId)
	return total, nil
}

func (f *FoodModel) GetFoods(userId int) ([]*models.Food, error) {
	stm := "SELECT * FROM foods WHERE userId = ? ORDER BY created_at DESC LIMIT 100"

	f.log().Info.Printf("GetFoods: listing foods for userId=%d", userId)

	rows, err := f.DB.Query(stm, userId)
	if err != nil {
		f.log().Error.Printf("GetFoods: query failed for userId=%d: %v", userId, err)
		return nil, err
	}

	defer rows.Close()

	var foods []*models.Food

	for rows.Next() {
		food := &models.Food{}

		err = rows.Scan(
			&food.ID, &food.Meal, &food.Name, &food.Protein, &food.Carbohydrates,
			&food.Fat, &food.Calories, &food.CreatedAt, &food.UpdatedAt, &food.UserID,
		)

		if err != nil {
			f.log().Error.Printf("GetFoods: scan failed for userId=%d: %v", userId, err)
			return nil, err
		}

		foods = append(foods, food)
	}

	// retries any rows error encountered during iteration
	err = rows.Err()
	if err != nil {
		f.log().Error.Printf("GetFoods: rows iteration error for userId=%d: %v", userId, err)
		return nil, err
	}

	f.log().Info.Printf("GetFoods: returned %d foods for userId=%d", len(foods), userId)
	return foods, nil
}

func (f *FoodModel) UpdateFood(
	meal string, name string, protein, cabs, fat, calory, foodId, userId int,
) (int, error) {
	stm := `UPDATE foods SET meal = ?, name = ?, protein = ?, carbohydrate = ?, fat = ?, calories = ?, updated_at = UTC_TIMESTAMP() WHERE id = ? AND userId = ?`

	f.log().Info.Printf("UpdateFood: updating foodId=%d for userId=%d", foodId, userId)

	result, err := f.DB.Exec(stm, meal, name, protein, cabs, fat, calory, foodId, userId)
	if err != nil {
		f.log().Error.Printf("UpdateFood: exec failed for foodId=%d userId=%d: %v", foodId, userId, err)
		return 0, err
	}

	id, err := result.RowsAffected()
	if err != nil {
		f.log().Error.Printf("UpdateFood: RowsAffected failed for foodId=%d userId=%d: %v", foodId, userId, err)
		return 0, err
	}

	if id == 0 {
		f.log().Warning.Printf("UpdateFood: no rows affected for foodId=%d userId=%d", foodId, userId)
	} else {
		f.log().Info.Printf("UpdateFood: updated %d row(s) for foodId=%d userId=%d", id, foodId, userId)
	}

	return int(id), nil
}

func (f *FoodModel) DeleteFood(foodId, userId int) (int, error) {
	stm := "DELETE FROM foods WHERE id = ? AND userId = ?"

	f.log().Info.Printf("DeleteFood: deleting foodId=%d for userId=%d", foodId, userId)

	result, err := f.DB.Exec(stm, foodId, userId)
	if err != nil {
		f.log().Error.Printf("DeleteFood: exec failed for foodId=%d userId=%d: %v", foodId, userId, err)
		return 0, err
	}

	id, err := result.RowsAffected()
	if err != nil {
		f.log().Error.Printf("DeleteFood: RowsAffected failed for foodId=%d userId=%d: %v", foodId, userId, err)
		return 0, err
	}

	if id == 0 {
		f.log().Warning.Printf("DeleteFood: no rows affected for foodId=%d userId=%d", foodId, userId)
	} else {
		f.log().Info.Printf("DeleteFood: deleted %d row(s) for foodId=%d userId=%d", id, foodId, userId)
	}

	return int(id), nil
}
