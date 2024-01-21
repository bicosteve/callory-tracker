package mysql

import (
	"database/sql"

	"github.com/bicosteve/callory-tracker/pkg/models"
)

type FoodModel struct {
	DB *sql.DB
}

// InsertFood(): insert food into db
func (f *FoodModel) InsertFood(name string, protein int, carbohydrates int, fat int, calories int) (int, error) {
	stm := `INSERT INTO foods (name, protein, carbohydrates,fat,calories,created_at,update_at) VALUES (?,?,?,?,?,UTC_TIMESTAMP(),UTC_TIMESTAMP())`

	result, err := f.DB.Exec(stm, name, protein, carbohydrates, fat, calories)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	// will return the last inserted id on the table
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (f *FoodModel) GetFoods(userid int) ([]*models.Food, error) {
	stm := "SELECT * FROM foods WHERE userId = ? ORDER BY created_at DESC"
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
