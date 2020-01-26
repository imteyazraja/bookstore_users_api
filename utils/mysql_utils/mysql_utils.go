package mysql_utils

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/imteyazraja/bookstore_users_api/utils/errors"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundErr("no record found")
		}
		return errors.NewInternalServerErr(fmt.Sprintf("error parsing database response %s", err.Error()))
	}
	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestErr(fmt.Sprintf("Invalid Data %s", err.Error()))
	}
	return errors.NewInternalServerErr(fmt.Sprintf("error processing request %s", err.Error()))
}
