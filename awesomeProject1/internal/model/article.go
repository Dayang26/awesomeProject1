package model

import "time"

const (
	STATUS_PUBLIC = iota + 1 // 公开
	STATUS_SECRET            // 私密
	STATUS_DRAFT             // 草稿
)

const (
	TYPE_ORIGINAL  = iota + 1 // 原创
	TYPE_REPRINT              // 转载
	TYPE_TRANSLATE            // 翻译
)

type Article struct {
	Model

	Title       string `gorm:"type:varchar(100);not null" json:"title"`
	Desc        string `json:"desc"`
	Content     string `json:"content"`
	Img         string `json:"img"`
	Type        int    `gorm:"type:tinyint;comment:类型(1-原创 2-转载 3-翻译)" json:"type"` // 1-原创 2-转载 3-翻译
	Status      int    `gorm:"type:tinyint;comment:状态(1-公开 2-私密)" json:"status"`    // 1-公开 2-私密
	IsTop       bool   `json:"is_top"`
	IsDelete    bool   `json:"is_delete"`
	OriginalUrl string `json:"original_url"`

	CategoryId int `json:"category_id"`
	UserId     int `json:"-"`

	Tags     []*Tag    `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`
	Category *Category `gorm:"foreignkey:CategoryId" json:"category"`
	User     *UserAuth `gorm:"foreignkey:UserId" json:"user"`
}

type ArticleTag struct {
	ArticleId int
	TagId     int
}

type BlogArticleVo struct {
	Article

	CommentCount int64 `json:"comment_count"`
	LikeCount    int64 `json:"like_count"`
	ViewCount    int64 `json:"view_count"`

	LastArticle       ArticlePaginationVo  `gorm:"-" json:"last_article"`
	NextArticle       ArticlePaginationVo  `gorm:"-" json:"next_article"`
	RecommendArticles []RecommendArticleVo `gorm:"-" json:"recommend_articles"`
	NewestArticles    []RecommendArticleVo `gorm:"-" json:"newest_articles"`
}

type ArticlePaginationVo struct {
	ID    int    `json:"id"`
	Img   string `json:"img"`
	Title string `json:"title"`
}

type RecommendArticleVo struct {
	ArticlePaginationVo
	CreatedAt time.Time `json:"created_at"`
}
