package db

import (
	"database/sql"

	"github.com/OscarClemente/go-noob/models"
)

func (db Database) GetAllReviews() (*models.ReviewList, error) {
	list := &models.ReviewList{}
	rows, err := db.Conn.Query("SELECT * FROM review")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.Game, &review.Title, &review.Content, &review.Rating, &review.User, &review.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Reviews = append(list.Reviews, review)
	}
	return list, nil
}

func (db Database) AddReview(item *models.Item) error {
	var id int
	var createdAt string
	query := `INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, item.Name, item.Description).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	item.ID = id
	item.CreatedAt = createdAt
	return nil
}

func (db Database) GetReviewById(itemId int) (models.Item, error) {
	item := models.Item{}
	query := `SELECT * FROM items WHERE id = $1;`
	row := db.Conn.QueryRow(query, itemId)
	switch err := row.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt); err {
	case sql.ErrNoRows:
		return item, ErrNoMatch
	default:
		return item, err
	}
}

func (db Database) DeleteReview(itemId int) error {
	query := `DELETE FROM items WHERE id = $1;`
	_, err := db.Conn.Exec(query, itemId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateReview(itemId int, itemData models.Item) (models.Item, error) {
	item := models.Item{}
	query := `UPDATE items SET name=$1, description=$2 WHERE id=$3 RETURNING id, name, description, created_at;`
	err := db.Conn.QueryRow(query, itemData.Name, itemData.Description, itemId).Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return item, ErrNoMatch
		}
		return item, err
	}
	return item, nil
}
