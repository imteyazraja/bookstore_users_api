package users

import (
	"net/http"
	"strconv"

	"github.com/imteyazraja/bookstore_users_api/utils/errors"

	"github.com/imteyazraja/bookstore_users_api/services"

	"github.com/imteyazraja/bookstore_users_api/domain/users"

	"github.com/gin-gonic/gin"
)

func GetUserId(userIdParam string) (int64, *errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestErr("user id should be a number")
	}
	return userId, nil
}
func Create(c *gin.Context) {
	var user users.User
	/*bytes, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		// TO DO
		return
	}
	if err := json.Unmarshal(bytes, &user); err != nil {
		// TO DO
		fmt.Println(err.Error())
		return
	}*/
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestErr("invalid json body")
		c.JSON(restErr.Status, restErr)
		//fmt.Println(err.Error())
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	//fmt.Println(user)
	//fmt.Println(string(bytes))
	//fmt.Println(err)
	//c.String(http.StatusNotImplemented, "implement me!")
	c.JSON(http.StatusCreated, result)
}
func Get(c *gin.Context) {
	/*userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestErr("user id should be a number")
		c.JSON(err.Status, err)
		return
	}*/
	userId, idErr := GetUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	user, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user)
	//c.String(http.StatusNotImplemented, "implement me!")
}
func Update(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestErr("user id should be a number")
		c.JSON(err.Status, err)
		return
	}
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestErr("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user.Id = userId
	isPartial := c.Request.Method == http.MethodPatch
	result, err := services.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
func Delete(c *gin.Context) {
	userId, idErr := GetUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	if err := services.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.FindByStatus(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users)
}
