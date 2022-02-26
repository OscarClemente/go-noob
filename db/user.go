package db

import (
	"database/sql"

	"github.com/OscarClemente/go-noob/models"
)

func (db Database) GetAllUsers() (*models.UserList, error) {
	list := &models.UserList{}
	rows, err := db.Conn.Query("SELECT * FROM user")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Friends)
		if err != nil {
			return list, err
		}
		list.Users = append(list.Users, user)
	}
	return list, nil
}

func (db Database) AddUser(user *models.User) error {
	var id int
	query := `INSERT INTO users (name, email, friends) VALUES ($1, $2, $3) RETURNING id`
	err := db.Conn.QueryRow(query, user.Name, user.Email, user.Friends).Scan(&id)
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (db Database) GetUserById(userId int) (models.User, error) {
	user := models.User{}
	query := `SELECT * FROM users WHERE id = $1;`
	row := db.Conn.QueryRow(query, userId)
	switch err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Friends); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}

func (db Database) DeleteUser(userId int) error {
	query := `DELETE FROM users WHERE id = $1;`
	_, err := db.Conn.Exec(query, userId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateUser(userId int, userData models.User) (models.User, error) {
	user := models.User{}
	query := `UPDATE users SET name=$1, email=$2, friends=$3 WHERE id=$5 RETURNING id, game, title, content, created_at;`
	err := db.Conn.QueryRow(query, userData.Name, userData.Email, userData.Friends, userId).Scan(&user.ID, &user.Name, &user.Email, &user.Friends)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
		return user, err
	}
	return user, nil
}
