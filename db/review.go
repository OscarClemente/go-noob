package db

import (
	"database/sql"

	"github.com/OscarClemente/go-noob/models"
)

func (db Database) GetAllReviews() (*models.ReviewList, error) {
	list := &models.ReviewList{}
	rows, err := db.Conn.Query("SELECT * FROM reviews")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var review models.Review
		err := rows.Scan(&review.ID, &review.Game, &review.Title, &review.Content, &review.Rating, &review.UserID, &review.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Reviews = append(list.Reviews, review)
	}
	return list, nil
}

func (db Database) AddReview(review *models.Review) error {
	var id int
	var createdAt string
	query := `INSERT INTO reviews (game, title, content, rating, author) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err := db.Conn.QueryRow(query, review.Game, review.Title, review.Content, review.Rating, review.UserID).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	review.ID = id
	review.CreatedAt = createdAt
	return nil
}

func (db Database) GetReviewById(reviewId int) (models.Review, error) {
	review := models.Review{}
	query := `SELECT * FROM reviews WHERE id = $1;`
	row := db.Conn.QueryRow(query, reviewId)
	switch err := row.Scan(&review.ID, &review.Game, &review.Title, &review.Content, &review.Rating, &review.UserID, &review.CreatedAt); err {
	case sql.ErrNoRows:
		return review, ErrNoMatch
	default:
		return review, err
	}
}

func (db Database) DeleteReview(reviewId int) error {
	query := `DELETE FROM reviews WHERE id = $1;`
	_, err := db.Conn.Exec(query, reviewId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateReview(reviewId int, reviewData models.Review) (models.Review, error) {
	review := models.Review{}
	query := `UPDATE reviews SET game=$1, title=$2, content=$3, rating=$4  WHERE id=$5 RETURNING id, game, title, content, created_at;`
	err := db.Conn.QueryRow(query, reviewData.Game, reviewData.Title, reviewData.Content, reviewData.Rating, reviewId).Scan(&review.ID, &review.Game, &review.Title, &review.Content, &review.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return review, ErrNoMatch
		}
		return review, err
	}
	return review, nil
}
