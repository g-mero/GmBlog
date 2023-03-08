package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gmeroblog/utils/errmsg"
	"gmeroblog/utils/rand"
	"gmeroblog/utils/static"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          int    `gorm:"primarykey" json:"id"`
	Username    string `json:"username" validate:"required,min=4,max=30"`
	Password    string `json:"password" validate:"required,min=6,max=20"`
	Nickname    string `gorm:"DEFAULT:'unknown user'" json:"nickname" validate:"required,min=4,max=20"`
	GithubId    string `json:"github_id"`
	GithubToken string `json:"github_token"`
	AvatarUrl   string `json:"avatar_url"`
	Email       string `json:"email" validate:"required,email"`
	Role        int    `gorm:"DEFAULT:2" json:"role" validate:"required,lte=2"`
}

func InitUser() int {
	var user User
	err := db.Where("role = ?", 1).First(&user).Error
	code := errmsg.SUCCES
	// 创建初始管理员账号
	if err == gorm.ErrRecordNotFound {
		user.Username = "admin"
		user.Password = "admin"
		user.Nickname = "admin"
		user.Role = 1

		if static.Get("admin_github_id") != "" {
			user.GithubId = static.Get("admin_github_id")
		}

		code = CreateUser(&user)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return errmsg.ERROR
	}
	return code
}

// is user exist
func CheckUser(name string) (code int) {
	var users User
	err := db.Select("id").Where("username = ?", name).First(&users).Error
	if err == gorm.ErrRecordNotFound {
		return errmsg.SUCCES
	}
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USED
	}
	return errmsg.ERROR
}

// create user
func CreateUser(data *User) int {
	data.Password = ScryptPw(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

func GetUserByGithubId(id string) (User, int) {
	var user User
	err := db.Where("github_id = ?", id).First(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCES
}

func GetUserByUsername(name string) (User, int) {
	var user User
	err := db.Where("username = ?", name).First(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCES
}

func GetUserByUserID(id int) (User, int) {
	var user User
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCES
}

// get user list
func GetUsers(pageSize int, pageNum int) []User {
	var users []User
	err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	return users
}

// edit user
func EditUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["nickname"] = data.Nickname
	maps["password"] = ScryptPw(data.Password)
	err = db.Model(&user).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES

}

// delete user
func DeleteUser(id int) int {
	var user User
	err = db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

// password scrypt
func ScryptPw(password string) string {
	const Keylen = 10
	salt := []byte{2, 33, 44, 55, 3, 11, 23, 98}

	HashPw, err := scrypt.Key([]byte(password), salt, 8192, 8, 1, Keylen)
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
}

// verify login
func CheckLogin(username string, password string) (User, int) {
	var user User
	db.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}
	if ScryptPw(password) != user.Password {
		return user, errmsg.ERROR_PASSWORD_WRONG
	}
	if user.Role != 1 {
		return user, errmsg.ERROR_USER_NOT_RIGHT
	}
	return user, errmsg.SUCCES
}

func request_get(url string, access_token string) io.ReadCloser {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+access_token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	if resp.StatusCode != 200 {
		return nil
	}
	return resp.Body
}

// github oauth
func GithubLogin(access_token string) (User, int) {
	var user User

	// 获取用户基本信息
	body := request_get("https://api.github.com/user", access_token)

	if body == nil {
		return user, errmsg.ERROR_REQUEST
	}

	var gituser struct {
		Login     string `json:"login"`      // github用户名，这是可以修改的
		Id        int    `json:"id"`         // githubID , 唯一的
		AvatarUrl string `json:"avatar_url"` // github头像地址
		Email     string `json:"email"`
	}

	if err := json.NewDecoder(body).Decode(&gituser); err != nil {
		fmt.Println(err)
		return user, errmsg.ERROR_JSON_DECODE
	}
	body.Close()

	// 获取用户邮箱
	nbody := request_get("https://api.github.com/user/emails", access_token)
	if nbody == nil {
		fmt.Println("1")
		return user, errmsg.ERROR_REQUEST
	}

	var emailStruct []struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(nbody).Decode(&emailStruct); err != nil {
		fmt.Println(err)
		return user, errmsg.ERROR_JSON_DECODE
	}
	nbody.Close()

	// 取获取到的第一个邮箱
	gituser.Email = emailStruct[0].Email

	fmt.Println(gituser)

	user.GithubId = strconv.Itoa(gituser.Id)
	user.Nickname = gituser.Login
	user.AvatarUrl = gituser.AvatarUrl
	user.Email = gituser.Email
	user.GithubToken = access_token
	// 新建用户
	if !github_id_exist(user.GithubId) {
		user.Username = "github-" + user.GithubId
		user.Password = rand.String(10)
		user.Role = 3
		code := CreateUser(&user)
		if code != errmsg.SUCCES {
			return user, errmsg.ERROR
		}
		return user, errmsg.SUCCES
	} else {
		// 老用户登录，更新数据
		old_user, code := GetUserByGithubId(user.GithubId)
		old_user.Nickname = user.Nickname
		old_user.AvatarUrl = user.AvatarUrl
		old_user.Email = user.Email
		old_user.GithubToken = access_token
		if code != errmsg.SUCCES {
			return user, errmsg.ERROR
		}
		db.Save(&old_user)
		return old_user, errmsg.SUCCES
	}
}

func github_id_exist(gid string) bool {
	err := db.Where("github_id = ?", gid).First(&User{}).Error
	return err != gorm.ErrRecordNotFound
}
