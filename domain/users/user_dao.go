package users

import (
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/tamihyo/bookstore_users-api/datasources/mysql/users_db"
	"github.com/tamihyo/bookstore_utils-go/logger"

	"github.com/tamihyo/bookstore_users-api/utils/mysql_utils"
	"github.com/tamihyo/bookstore_utils-go/rest_errors"
)

const (
	errorNoRows                 = "now rows in result set"
	queryInsertUser             = ("INSERT INTO users (first_name,last_name,email,date_created,password,status) VALUES(?,?,?,?,?,?);")
	queryGetUser                = "SELECT id,first_name,last_name, email,date_created,status FROM users WHERE id = ?;"
	queryUpdateUser             = "UPDATE users SET first_name =?,last_name=?,email=? WHERE id = ?;"
	queryDeleteUser             = "DELETE FROM users WHERE id = ?;"
	queryFindUserByStatus       = "SELECT id, first_name,last_name, email,date_created, status FROM users WHERE status = ?; "
	queryFindByEmailAndPassword = "SELECT id ,first_name,last_name,email,date_created,status FROM users WHERE email= ? AND password=?"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepate get user statement", err) //logging to system
		// return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
		return rest_errors.NewInternalServerError("datatabase error", errors.New("database error")) //logging to user
	}

	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		return mysql_utils.ParseError(getErr)
	}

	return nil
}

func (user *User) Save() *rest_errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryInsertUser)

	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("datatabase error", errors.New("database error"))
	}

	defer stmt.Close() //this is important if using Prepare statement, will be executed just before function return

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status)

	if saveErr != nil {
		logger.Error("error when trying to save user statement", saveErr)
		return rest_errors.NewInternalServerError("datatabase error", errors.New("database error"))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating new user", err)
		return rest_errors.NewInternalServerError("datatabase error", errors.New("database error"))
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)

	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("datatabase error", errors.New("database error"))
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("datatabase error", errors.New("database error"))
	}
	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)

	if err != nil {
		logger.Error("error when trying to prepare delte user statement", err)
		return rest_errors.NewInternalServerError("error when trying to prepare delte user statement", errors.New("database error"))
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	if err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user by status", err)
		return nil, rest_errors.NewInternalServerError("error when trying to prepare find user by status", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)

	if err != nil {
		logger.Error("error when trying to find user by status", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find user by status", errors.New("database error"))
	}
	defer rows.Close() //put defer after return error

	results := make([]User, 0)
	for rows.Next() {
		var user User
		/*
			pointer always need to be passed inside Scan,
			so it does not return empty value
			and return whatever it get  from database
		*/
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when scan user row into user struct to find users by status", err)
			return nil, rest_errors.NewInternalServerError("error when scan user row into user struct to find users by status", errors.New("database error"))

		}
		results = append(results, user)
	}

	if len(results) == 0 {
		logger.Error("error when trying to find user by status", err)
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no user matching status %s", status))
	}

	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepate get user by email & passwor statement", err) //logging to system
		// return rest_errors.NewInternalServerError(rest_errors.NewError("database error"))
		return rest_errors.NewInternalServerError("error when trying to prepate get user by email & passwor statement", errors.New("database error")) //logging to user
	}

	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invali user credential")
		}
		logger.Error("error when trying to prepate get user by email & passwor statement", getErr) //logging to system
		return mysql_utils.ParseError(getErr)
	}

	return nil
}
