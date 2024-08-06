package handle

import (
	g "awesomeProject1/internal/global"
	"awesomeProject1/internal/model"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

// 响应体结构
type Response[T any] struct {
	Code    int    `json:"code"`    //业务状态码
	Message string `json:"message"` // 响应消息
	Data    T      `json:"data"`    // response data
}

func ReturnHttpResponse(c *gin.Context, httpCode, code int, msg string, data any) {
	c.JSON(httpCode, Response[any]{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func ReturnResponse(c *gin.Context, r g.Result, data any) {
	ReturnHttpResponse(c, http.StatusOK, r.Code(), r.Msg(), data)
}

func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c, g.OkResult, data)
}

// 所有可预料的错误 = 业务错误 + 系统错误 ,在业务层面处理，并且返回http200状态码
// 对于不可预料的错误，会出发 panic ,有gin中间件捕获，返回500
// err 是业务错误, data 是错误数据 (可以是 error 或 string)
func ReturnError(c *gin.Context, r g.Result, data any) {
	slog.Info("[Func-ReturnError] ", r.Msg())

	var val string = r.Msg()

	if data != nil {
		switch v := data.(type) {
		case error:
			val = v.Error()
		case string:
			val = v
		}
		slog.Error(val) //错误日志
	}

	c.AbortWithStatusJSON(
		http.StatusOK,
		Response[any]{
			Code:    r.Code(),
			Message: r.Msg(),
			Data:    val,
		},
	)
}

type PageQuery struct {
	Page    int    `form:"page_num"`  //当前页
	Size    int    `form:"page_size"` //每页条数
	Keyword string `form:"keyword"`   // 搜索关键字
}

// 分页响应数据
type PageResult[T any] struct {
	Page  int   `json:"page_num"`  //每页条数
	Size  int   `json:"page_size"` // 上次页数
	Total int64 `json:"total"`     //总条数
	List  []T   `json:"page_data"` //分页数据
}

func GetDB(c *gin.Context) *gorm.DB { return c.MustGet(g.CTX_DB).(*gorm.DB) }

func GetRDB(c *gin.Context) *redis.Client { return c.MustGet(g.CTX_RDB).(*redis.Client) }

func CurrentUserAuth(c *gin.Context) (*model.UserAuth, error) {
	key := g.CTX_USER_AUTH

	if cache, exist := c.Get(key); exist && cache != nil {
		slog.Debug("[Func-CurrentUserAuth]get from cache " + cache.(*model.UserAuth).Username)
		return cache.(*model.UserAuth), nil
	}

	session := sessions.Default(c)
	id := session.Get(key)
	if id == nil {
		return nil, errors.New("session中没有 user_auth_id")
	}

	db := GetDB(c)
	user, err := model.GetUserAuthInfoById(db, id.(int))
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}
