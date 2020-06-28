package controllers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"sipemanager/dao"
	"time"
)

var secretKey = "Svhv*Zv5g&Wecr61BpTh&"

func (this *Controller) GetUser(c *gin.Context) (*dao.User, error) {
	token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot convert claim to mapclaim")
	}
	username := claim["username"].(string)
	user, err := this.dao.GetUserByUsername(username)
	if !ok {
		return nil, errors.New("username does not exists")
	}
	return user, nil
}

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//用户注册
// @Summary 用户注册
// @Tags user
// @Accept  json
// @Produce  json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} JSONResult{data=int} "{"code":0,"data":"用户Id"}"
// @Router /user/register [post]
func (this *Controller) Register(c *gin.Context) {
	var user User
	if err := c.Bind(&user); err != nil {
		this.echoError(c, err)
		return
	}
	loadUser, err := this.dao.GetUserByUsername(user.Username)
	if err == nil {
		this.echoError(c, errors.New("user exists"))
		return
	}
	loadUser = &dao.User{
		Username: user.Username,
		Password: user.Password,
	}
	id, err := this.dao.CreateUser(loadUser)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, id)
}

//用户登录
// @Summary 用户登录
// @Tags user
// @Accept  json
// @Produce  json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} JSONResult{data=object} "{"token": token, "user_id": loadUser.ID}"
// @Router /user/login [post]
func (this *Controller) Login(c *gin.Context) {
	var user User
	if err := c.Bind(&user); err != nil {
		this.echoError(c, err)
		return
	}
	loadUser, err := this.dao.GetUserByUsername(user.Username)
	if err != nil {
		this.echoError(c, err)
		return
	}
	if !dao.CheckPassword(loadUser.Password, user.Password) {
		this.echoError(c, errors.New("password is invalid"))
		return
	}
	token, err := createToken(loadUser.ID, loadUser.Username)
	if err != nil {
		this.echoError(c, err)
		return
	}
	this.echoResult(c, map[string]interface{}{"token": token, "user_id": loadUser.ID})
}

func createToken(userId uint, username string) (string, error) {
	claim := jwt.MapClaims{
		"id":       userId,
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(2)).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(secretKey))
}

func ValidateTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "Unauthorized access to this resource",
				"code": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "Token is not valid",
				"code": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
