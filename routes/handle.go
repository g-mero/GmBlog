package routes

import (
	"gmeroblog/model"
	"gmeroblog/utils/static"
	"html/template"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取全局设置项
func GetSet(k string) string {
	return static.Get("BASE_" + k)
}

// 获取一条Mood
func GetAMood() string {
	return model.GetMoodRand()
}

// 获取主题设置项
func GetThemeSet(k string) string {
	return static.Get("THEME_" + model.ThemeName + "_" + k)
}

// 首页
func Index(c *gin.Context) {
	pn, _ := strconv.Atoi(c.Param("pn"))
	title := GetSet("site_name")
	if pn > 0 {
		title = "第" + strconv.Itoa(pn) + "页-" + title
	} else {
		pn = 1
	}
	arts, code, total := model.GetArt(8, pn, 0)
	if code != 200 {
		Error(c)
		return
	}
	c.HTML(200, "index.tmpl", gin.H{
		"title": title,
		"total": total,
		"arts":  arts,
		"pn":    pn,
	})
}

// 文章页
func Article(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		Error(c)
		return
	}
	art, code := model.GetSingleArt(int(id))
	if code != 200 {
		Error(c)
		return
	}
	c.HTML(200, "article.tmpl", gin.H{
		"title":   art.Title + " - " + GetSet("site_name"),
		"article": art,
	})
}

// 分类页
func Category(c *gin.Context) {
	slug := c.Param("slug")
	pn, _ := strconv.Atoi(c.Query("pagenum"))
	var cate model.Category
	var code int
	cate, code = model.GetCateBySlug(slug)

	if code != 200 {
		Error(c)
		return
	}
	if pn <= 0 {
		pn = 1
	}
	arts, code, total := model.GetArt(8, pn, cate.ID)

	if code != 200 {
		Error(c)
		return
	}

	c.HTML(200, "category.tmpl", gin.H{
		"title": "第" + strconv.Itoa(pn) + "页 - " + cate.Name + " - " + GetSet("site_name"),
		"pn":    pn,
		"cate":  cate,
		"arts":  arts,
		"total": total,
	})

}

// 错误页
func Error(c *gin.Context) {
	c.HTML(200, "404.tmpl", gin.H{
		"title": "404 - " + GetSet("site_name"),
	})
}

// 管理后台
func Admin(c *gin.Context) {
	c.HTML(200, "admin", gin.H{
		"title": "后台---" + GetSet("site_name"),
	})
}

// 一下是一些模板函数

// 生成分页
func GenPagination(current int, total int64) template.HTML {
	const pageSize = 8
	totalPages := int(total / pageSize)
	if totalPages <= 1 {
		return ""
	}
	const maxPages = 7
	const centerPages = maxPages - 4
	const avaPages = centerPages + 2
	const offsetPages = centerPages / 2

	const dots = `<li class="page-item"><span class="page-link">...</span></li>`

	var genPagiBynum = func(end int, start int) []string {
		var res []string

		for i := start; i <= end; i++ {
			active := ""
			inner := ""
			pageNum := strconv.Itoa(i)
			if current == i {
				active = "active"
				inner = `<span class="page-link">` + pageNum + `</span>`
			} else {
				inner = `<a class="page-link" href="/page/` +
					pageNum + `">` + pageNum + `</a>`
			}
			var tmp = `<li class="page-item ` + active + `">` + inner + `</li>`
			res = append(res, tmp)
		}

		return res
	}

	firstPage := genPagiBynum(1, 1)[0]
	finalPage := genPagiBynum(totalPages, totalPages)[0]

	var prevPage string
	var nextPage string

	if current <= 1 {
		prevPage = `<li class="page-item disabled"><span class="page-link">&lt;</span></li>`
		nextPage = `<li class="page-item"><a class="page-link" href="/page/2">&gt;</a></li>`
	} else {
		prevPage = `<li class="page-item"><a class="page-link" href="/page/` + strconv.Itoa(current-1) + `">&lt;</a></li>`
		if current >= totalPages {
			nextPage = `<li class="page-item disabled"><span class="page-link">&gt;</span></li>`
		} else {
			nextPage = `<li class="page-item"><a class="page-link" href="/page/` + strconv.Itoa(current+1) + `">&gt;</a></li>`
		}
	}

	result := `<nav aria-label="Page navigation" class="mt-3">
    <ul class="pagination justify-content-center pagination-sm">`

	result += prevPage

	if totalPages <= maxPages {
		for _, v := range genPagiBynum(totalPages, 1) {
			result += v
		}
	} else if current < avaPages {
		for _, v := range genPagiBynum(avaPages, 1) {
			result += v
		}
		result += dots
		result += finalPage
	} else if current > totalPages-avaPages+1 {
		result += firstPage
		result += dots
		for _, v := range genPagiBynum(totalPages, totalPages-avaPages+1) {
			result += v
		}
	} else {
		result += firstPage
		result += dots
		for _, v := range genPagiBynum(current+offsetPages, current-(centerPages-offsetPages-1)) {
			result += v
		}
		result += dots
		result += finalPage
	}
	result += nextPage
	result += `</ul></nav>`
	return template.HTML(result)
}

// string转成html
func Str2Html(text string) template.HTML {
	return template.HTML(text)
}

func FormatTime(time time.Time) string {
	return time.Format("2006-01-02")
}

func GetCateTree() []model.Catetree {
	return model.GetCates()
}

func String2Int(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}

func Divide(a int, b int) int {
	if a%b == 0 {
		return a / b
	} else {
		return a/b + 1
	}
}

func IntArray(a int) []int {
	var res []int
	for i := 0; i < a; i++ {
		res = append(res, i)
	}
	return res
}

func Add(a int, b int) int {
	return a + b
}

// 获取请求头
func GetHeader() string {
	var c *gin.Context
	return c.GetHeader("X-PJAX")
}

func RundInt(n int) int {
	return rand.Intn(n)
}
