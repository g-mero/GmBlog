package middleware

import (
	"encoding/json"
	"gmeroblog/utils/cache"
	"gmeroblog/utils/config"
	"gmeroblog/utils/errmsg"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var Jwtkey = []byte(config.Jwtkey)
var code int

type MyClaims struct {
	UserID int `json:"user_id"`
	Role   int `json:"role"`
	jwt.StandardClaims
}

// generate token
func SetToken(uid int, role int) (string, int) {
	// token 7天过期
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	setClaims := MyClaims{
		uid,
		role,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gmeroblog",
		},
	}
	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, setClaims)
	token, err := reqClaim.SignedString(Jwtkey)
	if err != nil {
		return "", errmsg.ERROR
	}
	if cacheValue, err := json.Marshal(setClaims); err == nil {
		cache.Set(token, cacheValue)
	} else {
		return token, errmsg.ERROR_JSON_DECODE
	}
	return token, errmsg.SUCCES
}

// verify token
func CheckToken(token string) (*MyClaims, int) {
	// 修改为从cache中读取，不存在则返回错误
	v := cache.Get(token)
	if v == nil {
		return nil, errmsg.ERROR
	}
	var value MyClaims
	if err := json.Unmarshal(v, &value); err != nil {
		return nil, errmsg.ERROR_JSON_DECODE
	}
	return &value, errmsg.SUCCES
	/* 	setToken, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
	   		return Jwtkey, nil
	   	})
	   	if key, _ := setToken.Claims.(*MyClaims); setToken.Valid {
	   		return key, errmsg.SUCCES
	   	} else {
	   		return nil, errmsg.ERROR
	   	} */
}

// jwt middleware
func JwtToken(opt int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")

		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		checkToken := strings.SplitN(tokenHeader, " ", 2)
		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		// 获取token
		token := checkToken[1]
		key, tCode := CheckToken(token)
		if tCode != errmsg.SUCCES {
			// 在cache中没有找到token，输出token错误的信息
			code = errmsg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		// token过期
		if time.Now().Unix() > key.ExpiresAt {
			// 在缓存中删除token
			cache.Del(token)
			code = errmsg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		// 管理员鉴权
		if key.Role != 1 && opt == 1 {
			code = errmsg.ERROR_USER_NOT_RIGHT
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		c.Set("userID", key.UserID)
		c.Set("userRole", key.Role)
		c.Set("loginToken", token)
	}
}
