package handle

import (
	g "awesomeProject1/internal/global"
	"awesomeProject1/internal/model"
	"awesomeProject1/internal/utils"
	"awesomeProject1/internal/utils/jwt"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type UserAuth struct{}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
	Code     string `json:"code" binding:"required"`
}

type LoginVo struct {
	model.UserInfo
	ArticleLikeSet []string `json:"article_like_set"`
	CommentListSet []string `json:"comment_list_set"`
	Token          string   `json:"token"`
}

func (*UserAuth) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	userAuth, err := model.GetUserAuthInfoByName(db, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	if !utils.BcryptCheck(req.Password, userAuth.Password) {
		ReturnError(c, g.ErrPassword, nil)
		return
	}

	// TODO 获取ip相关信息
	ipAddress := "none ipAddress"
	ipSource := "none ipSource"

	//获取个人信息
	userInfo, err := model.GetUserInfoById(db, userAuth.UserInfoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, nil)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	roleIds, err := model.GetRoleIdsByUserId(db, userAuth.ID)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	articleLikeSet, err := rdb.SMembers(rctx, g.ARTICLE_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	commentLikeSet, err := rdb.SMembers(rctx, g.COMMENT_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	conf := g.Conf.JWT
	token, err := jwt.GenToken(conf.Secret, conf.Issuer, int(conf.Expire), userAuth.ID, roleIds)
	if err != nil {
		ReturnError(c, g.ErrTokenCreate, err)
		return
	}

	err = model.UpdateUserLoginInfo(db, userAuth.ID, ipAddress, ipSource)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	slog.Info("User Login successful! User:", userAuth.Username)

	session := sessions.Default(c)
	session.Set(g.CTX_USER_AUTH, userAuth.ID)
	session.Save()

	// Delete offline status information in redis cache
	offlineKey := g.OFFLINE_USER + strconv.Itoa(userAuth.ID)
	rdb.Del(rctx, offlineKey).Result()

	ReturnSuccess(c, LoginVo{
		UserInfo:       *userInfo,
		ArticleLikeSet: articleLikeSet,
		CommentListSet: commentLikeSet,
		Token:          token,
	})
}

func (*UserAuth) Register(c *gin.Context) {
	ReturnSuccess(c, "Register")
}

func (*UserAuth) Logout(c *gin.Context) {
	c.Set(g.CTX_USER_AUTH, nil)

	auth, _ := CurrentUserAuth(c)
	if auth == nil {
		ReturnSuccess(c, nil)
		return
	}

	session := sessions.Default(c)
	session.Delete(g.CTX_USER_AUTH)
	session.Save()

	rdb := GetRDB(c)
	onlineKey := g.ONLINE_USER + strconv.Itoa(auth.ID)
	rdb.Del(rctx, onlineKey)

	ReturnSuccess(c, nil)
}

func (*UserAuth) SendCode(c *gin.Context) {
	ReturnSuccess(c, "send verity code By Email!")
}

// Todo :refresh token
