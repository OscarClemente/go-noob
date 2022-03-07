package db

import (
	"database/sql"
	"fmt"

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
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return list, err
		}
		list.Users = append(list.Users, user)
	}
	return list, nil
}

func (db Database) AddUser(user *models.User) error {
	var id int
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	err := db.Conn.QueryRow(query, user.Name, user.Email).Scan(&id)
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (db Database) GetUserById(userId int) (models.User, error) {
	fmt.Println("GetUserById", userId)
	user := models.User{}
	query := `SELECT * FROM users WHERE id = $1;`
	row := db.Conn.QueryRow(query, userId)
	switch err := row.Scan(&user.ID, &user.Name, &user.Email); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}

func (db Database) GetUsersById(userIds []int) ([]*models.User, error) {
	users := []*models.User{}
	joinIds := ""
	for i, userId := range userIds {
		joinIds = joinIds + fmt.Sprint(userId)
		if i < len(userIds)-1 {
			joinIds = joinIds + ","
		}
	}

	query := `SELECT * FROM users WHERE id IN(` + joinIds + `);`
	rows, err := db.Conn.Query(query)
	if err != nil {
		fmt.Println("Error fetching user ids", err)
		return users, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (db Database) AddFriendToUser(userId, friendId int) error {
	var id int
	query := `INSERT INTO friends (username, friend) VALUES ($1, $2) RETURNING id`
	err := db.Conn.QueryRow(query, userId, friendId).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (db Database) GetFriendsOfUserById(userId int) ([]*models.User, error) {
	users := []*models.User{}
	query := `SELECT users.ID, users.Name, users.Email
		FROM friends
		INNER JOIN users ON users.id = friends.friend
		WHERE friends.username = $1;`
	rows, err := db.Conn.Query(query, userId)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	return users, nil
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
	query := `UPDATE users SET name=$1, email=$2 WHERE id=$5 RETURNING id, game, title, content, created_at;`
	err := db.Conn.QueryRow(query, userData.Name, userData.Email, userId).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
		return user, err
	}
	return user, nil
}
