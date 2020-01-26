package users

import (
	"fmt"

	"github.com/imteyazraja/bookstore_users_api/utils/mysql_utils"

	"github.com/imteyazraja/bookstore_users_api/utils/date_utils"

	"github.com/imteyazraja/bookstore_users_api/datasources/mysql/users_db"

	"github.com/imteyazraja/bookstore_users_api/utils/errors"
)

const (
	statusActive        = "active"
	indexUniqueEmail    = "email_UNIQUE"
	errorNoRows         = "no rows in result set"
	queryInsertUser     = "INSERT INTO users (first_name, last_name,email,date_created,password,status) VALUES (?,?,?,?,?,?);"
	queryGetUser        = "SELECT id,first_name,last_name,email,date_created,status FROM users WHERE id=?"
	queryUpdateUser     = "UPDATE users SET first_name=?,last_name=?,email=? WHERE id=?"
	queryDeleteUser     = "DELETE FROM users WHERE id=?"
	queryFindUserStatus = "SELECT id,first_name,last_name,email,date_created,status FROM users WHERE status=?"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		/*if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundErr(fmt.Sprintf("user %d doesn't exist", user.Id))
		}
		return errors.NewInternalServerErr(fmt.Sprintf("error while trying to get user %d : %s", user.Id, err.Error()))*/
		return mysql_utils.ParseError(getErr)
	}

	/*if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}
	result := userDB[user.Id]
	if result == nil {
		return errors.NewNotFoundErr(fmt.Sprintf("user %d not found", user.Id))
	}
	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated*/
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Status = statusActive
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status)
	/*if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestErr(fmt.Sprintf("error email %s already exists", user.Email))
		}
		return errors.NewInternalServerErr(fmt.Sprintf("error while trying to save user %s", err.Error()))
	}*/
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		//return errors.NewInternalServerErr(fmt.Sprintf("error while trying to save user %s", err.Error()))
		return mysql_utils.ParseError(err)
	}
	user.Id = userId
	/*current := userDB[user.Id]
	if current != nil {
		if current.Email == user.Email {
			return errors.NewBadRequestErr(fmt.Sprintf("email %s already exists", user.Email))
		}
		return errors.NewBadRequestErr(fmt.Sprintf("user %d already exists", user.Id))

	}
	user.DateCreated = date_utils.GetNowString()
	userDB[user.Id] = user*/
	return nil
}
func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()
	_, updtErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)

	if updtErr != nil {
		return mysql_utils.ParseError(updtErr)
	}
	return nil
}
func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()

	if _, delErr := stmt.Exec(user.Id); delErr != nil {
		return mysql_utils.ParseError(delErr)
	}
	return nil
}
func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserStatus)
	if err != nil {
		return nil, errors.NewInternalServerErr(err.Error())
	}
	defer stmt.Close()
	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerErr(err.Error())
	}
	defer rows.Close()
	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundErr(fmt.Sprintf("no user matching status %s", status))
	}
	return results, nil
}
