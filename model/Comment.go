package model

import (
	"gmeroblog/utils/errmsg"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Comment struct {
	ID          int       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id"`                        // 评论作者id
	ArticleID   int       `json:"article_id"`                     // 文章id
	Content     string    `gorm:"not null;" json:"content"`       // 评论内容
	Status      int8      `gorm:"default:2" json:"status"`        // 评论状态 1 审核通过 2 审核中
	IsTop       bool      `gorm:"default:false" json:"is_top"`    // 是否置顶
	IsEdited    bool      `gorm:"default:false" json:"is_edited"` // 是否修改过
	ToCommentID int       `gorm:"default:0" json:"to_comment_id"` // 回复目标id,0 代表是根评论
	ToUserID    int       `gorm:"default:0" json:"to_user_id"`    // 回复目标用户id,0表示回复的是根评论
	Likes       int       `gorm:"default:0" json:"likes"`         // 点赞数
	Replys      int       `gorm:"default:0" json:"replys"`        // 回复数
}

// 检查评论信息是否无误
func (data *Comment) check() (code int) {
	if data.ArticleID <= 0 {
		return errmsg.ERROR_CANT_FIND_TAR_ART
	}

	// 对评论内容进行检查和基本纠错

	var tmpContent = data.Content
	tmpContent = strings.TrimSpace(tmpContent)

	// 处理字符串中的空白符号
	reg := regexp.MustCompile(`(\n{2,})`)
	tmpContent = reg.ReplaceAllString(tmpContent, "\n")
	reg = regexp.MustCompile(`( |\t|\r|\f|\v){2,}`)
	tmpContent = reg.ReplaceAllString(tmpContent, ` `)

	data.Content = tmpContent

	if strings.Count(data.Content, "") > 200 || strings.Count(data.Content, "") < 5 {
		return errmsg.ERROR_CMT_CONTENT_ILLEGAL
	}

	// 对其他信息进行检查
	if data.ToCommentID == 0 && data.ToUserID != 0 {
		return errmsg.ERROR_CANT_FIND_TAR_USER
	}

	err := db.Model(Article{}).Where("id = ?", data.ArticleID).Limit(1).Error
	if err != nil {
		return errmsg.ERROR_CANT_FIND_TAR_ART
	}
	// 判断目标评论情况（此时该comment为回复）
	if data.ToCommentID != 0 {
		var cmt Comment
		err := db.Where("id = ?", data.ToCommentID).First(&cmt).Error
		if err != nil || cmt.ArticleID != data.ArticleID {
			return errmsg.ERROR_CANT_FIND_TAR_COMMENT
		}
		// 判断目标用户情况
		if data.ToUserID != 0 {
			cmts, code := getReplysByCommentID(data.ToCommentID, cmt.Replys)
			if code != errmsg.SUCCES {
				return errmsg.ERROR_CANT_FIND_TAR_USER
			}

			var tmpBool = true
			for _, v := range cmts {
				if v.UserID == data.ToUserID {
					tmpBool = false
					break
				}
			}
			if tmpBool {
				return errmsg.ERROR_CANT_FIND_TAR_USER
			}
		}

		cmt.Replys++
		db.Model(&cmt).UpdateColumn("replys", cmt.Replys)
	}
	return errmsg.SUCCES
}

// AddComment 新增评论
func AddComment(data *Comment) int {
	code := data.check()
	if code != errmsg.SUCCES {
		return code
	}
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

// 删除评论
func DelComment(id int) int {
	err := db.Where("id = ?", id).Delete(&cate{}).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCES
}

// 获取评论下的回复不包括用户信息
func getReplysByCommentID(id int, limit int) ([]Comment, int) {
	var res []Comment

	// 使用limit进行优化
	err := db.Where("to_comment_id = ?", id).Limit(limit).Find(&res).Error
	if err != nil {
		return res, errmsg.ERROR
	}
	return res, errmsg.SUCCES
}

// 结构体用于从数据库获取联合数据，因此不对外开放
type commentWithUser struct {
	Comment
	Nickname  string `json:"nickname"`
	AvatarUrl string `json:"avatar_url"`
}

type CmtWithUserWithReplys struct {
	commentWithUser
	// for reply
	ToUserNickname string                  `json:"to_user_nickname"`
	Children       []CmtWithUserWithReplys `json:"children"`
}

// 获取评论下的回复包括用户信息
func GetReplysWithUserByCommentID(pageSize int, pageNum int, commentID int) (res []CmtWithUserWithReplys, code int) {
	var cmtWithUser []commentWithUser
	err := db.Model(&Comment{}).Select("user.nickname,user.avatar_url,comment.*").Joins("left join user on comment.user_id=user.id").Where("to_comment_id = ?",
		commentID).Order("likes DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Scan(&cmtWithUser).Error

	if err != nil {
		return nil, errmsg.ERROR
	}

	for _, v := range cmtWithUser {
		var tmpCmtWithUserWithReplys CmtWithUserWithReplys
		tmpCmtWithUserWithReplys.commentWithUser = v
		if v.ToUserID > 0 {
			var tmpUser User
			err := db.Where("id = ?", v.ToUserID).First(&tmpUser).Error
			if err != nil {
				tmpCmtWithUserWithReplys.ToUserNickname = "未知用户-" + strconv.Itoa(v.ToUserID)
			} else {
				tmpCmtWithUserWithReplys.ToUserNickname = tmpUser.Nickname
			}
		}
		res = append(res, tmpCmtWithUserWithReplys)
	}

	return res, errmsg.SUCCES
}

// 获取文章评论包括用户信息包括每个评论三个回复
func GetCommentsWithUserByArtID(pageSize int, pageNum int, articleID int) ([]CmtWithUserWithReplys, int64, int) {
	var cmts []commentWithUser
	var res []CmtWithUserWithReplys
	var total int64

	err := db.Model(&Comment{}).Select("user.nickname,user.avatar_url,comment.*").Joins("left join user on comment.user_id=user.id").Where("to_comment_id = ? AND article_id = ?",
		0, articleID).Order("likes DESC").Count(&total).Limit(pageSize).Offset((pageNum - 1) * pageSize).Scan(&cmts).Error

	if err != nil {
		return nil, total, errmsg.ERROR
	}

	for _, v := range cmts {
		var tmpCmtWithReplys CmtWithUserWithReplys
		tmpCmtWithReplys.commentWithUser = v
		if v.Replys > 0 {
			tmpReply, code := GetReplysWithUserByCommentID(3, 1, v.ID)
			if code == errmsg.SUCCES {
				tmpCmtWithReplys.Children = tmpReply
			}
		}
		res = append(res, tmpCmtWithReplys)
	}

	return res, total, errmsg.SUCCES

}

// 获取文章下的评论不包括用户信息包括三个回复
func GetCommentsByArtID(pageSize int, pageNum int, articleID int) ([]CmtWithUserWithReplys, int64, int) {
	var res []CmtWithUserWithReplys
	var cmts []Comment
	var total int64

	// 按点赞数排行
	err := db.Model(&cmts).Where("to_comment_id = ? AND article_id = ?", 0,
		articleID).Order("likes DESC").Count(&total).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&cmts).Error

	if err != nil {
		return nil, total, errmsg.ERROR
	}

	for _, v := range cmts {
		var tmpResCmt CmtWithUserWithReplys
		tmpResCmt.Nickname = "用户-" + strconv.Itoa(v.UserID)
		tmpResCmt.AvatarUrl = "https://avatars.githubusercontent.com/u/0" // GitHub的默认头像
		tmpResCmt.Comment = v

		if tmpResCmt.Replys > 0 {
			replys, code := getReplysByCommentID(v.ID, 3)
			if code == errmsg.SUCCES {
				var children []CmtWithUserWithReplys
				for _, cmt := range replys {
					var tmpReply CmtWithUserWithReplys
					tmpReply.AvatarUrl = "https://avatars.githubusercontent.com/u/0"
					tmpReply.Nickname = "用户-" + strconv.Itoa(cmt.UserID)
					tmpReply.Comment = cmt
					if cmt.ToUserID > 0 {
						tmpReply.ToUserNickname = "用户-" + strconv.Itoa(cmt.ToUserID)
					}
					children = append(children, tmpReply)
				}
				tmpResCmt.Children = children
			}
		}

		res = append(res, tmpResCmt)
	}
	return res, total, errmsg.SUCCES
}
