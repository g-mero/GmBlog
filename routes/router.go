package routes

import (
	"fmt"
	v1 "gmeroblog/api/v1"
	"gmeroblog/middleware"
	"gmeroblog/utils/config"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
	if err != nil {
		panic(err.Error())
	}
	includes, err := filepath.Glob(templatesDir + "/includes/*.tmpl")
	if err != nil {
		panic(err.Error())
	}
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		var files []string
		files = append(append(files, include), layoutCopy...) // 确保include在最前面,否则空页面
		r.AddFromFilesFuncs(filepath.Base(include), template.FuncMap{
			"Arts":       GetArtList,
			"Int":        String2Int,
			"Divide":     Divide,
			"IntArray":   IntArray,
			"Add":        Add,
			"CateTree":   GetCateTree,
			"FormatTime": FormatTime,
			"STH":        Str2Html,
			"RandInt":    RundInt,
			"Pagination": GenPagination,
		}, files...)
	}
	if config.LocalAdmin {
		r.AddFromFiles("admin", "./web/admin/dist/index.html") // 前端后台管理界面 solidjs单界面
	}
	r.AddFromFiles("githubAuth", "./web/stuff/github-auth.html")
	return r
}

func InitRouter() {
	gin.DisableConsoleColor()
	// 记录日志
	logFile := &lumberjack.Logger{
		Filename:   "log/server.log",
		MaxSize:    2,     // 文件大小MB
		MaxBackups: 5,     // 最大保留日志文件数量
		MaxAge:     28,    // 保留天数
		Compress:   false, // 是否压缩
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	gin.SetMode(config.AppMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	_ = r.SetTrustedProxies(nil)

	r.HTMLRender = loadTemplates("./web/theme/default/")

	r.Static("/assets", "web/theme/default/assets")
	{
		r.GET("/", Index)
		r.GET("/page/:pn", Index)
		r.GET("/article/:id", Article)
		r.GET("/category/:slug", Category)

		fmt.Println(config.LocalAdmin)

		// 管理界面
		if config.LocalAdmin {
			r.GET("/admin", Admin)
			r.GET("/admin/:func", Admin)
			r.Static("/admin/assets", "web/admin/dist/assets") // 单界面应用的静态文件
		}

		r.NoRoute(Error)
	}
	// 管理员拥有最高读写权限，尽量保证信任的人使用
	admin := r.Group("api/v1")
	admin.Use(middleware.JwtToken(1))
	{
		// User model route
		admin.POST("user/add", v1.AddUser)
		admin.PUT("user/:id", v1.EditUser)
		admin.DELETE("user/:id", v1.DeleteUser)

		admin.POST("category/add", v1.AddCategory)
		admin.PUT("category/:id", v1.EditCate)
		admin.DELETE("category/:id", v1.DeleteCate)

		admin.POST("article/add", v1.AddArt)
		admin.PUT("article/:id", v1.EditArt)
		admin.PUT("update/:id", v1.UpdateHtmlContent)
		admin.DELETE("article/:id", v1.DeleteArt)

		admin.POST("diary/add", v1.AddDiary)
		admin.PUT("diary/:time", v1.EditDiary)
		admin.DELETE("diary/:time", v1.DeleteDiary)
		// 日记
		admin.GET("diary", v1.GetDiary)
		admin.GET("diarydate", v1.GetDiaryDate)
		admin.GET("diary/:time", v1.GetSingleDiary)
		// 设置
		admin.PUT("settings", v1.EditSettings)

	}
	// 普通用户（不相信写入的任何数据，后台要进行严格检查纠错，读权限可放宽，
	// 应该对api使用进行限额）
	user := r.Group("api/v1")
	user.Use(middleware.JwtToken(2))
	{
		user.GET("user/info", v1.GetUserInfo)
		user.POST("user/logout", v1.Logout)

		user.POST("comment/add", v1.AddComment)

		user.GET("comment", v1.GetComments)
		user.GET("comment/reply", v1.GetReplys)

	}
	// 公共区域（不应该开放任何写的权限，读权限应该进行流量控制）
	router := r.Group("api/v1")
	{
		// 鉴权
		router.GET("check", v1.CheckToken)
		// 分类
		router.GET("category", v1.GetCate)
		router.GET("cateTree", v1.GetCateTree)
		router.GET("category/:id", v1.GetSingleCate)
		// 用户
		router.GET("users", v1.GetUsers)
		router.GET("oauth/redirect", v1.GithubLogin)
		router.POST("login", v1.Login)
		// 文章
		router.GET("article", v1.GetArt)
		router.GET("article/:id", v1.GetSingleArt)
		// 评论
		router.GET("comment_limit", v1.GetCommentsLimit)

		router.GET("settings", v1.GetSettings)
	}

	_ = r.Run(config.HttpPort)
}
