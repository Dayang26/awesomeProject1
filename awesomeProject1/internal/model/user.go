package model

import (
	"gorm.io/gorm"
	"time"
)

type UserInfo struct {
	Model
	Email    string `json:"email" gorm:"type:varchar(30)"`
	Nickname string `json:"nickname" gorm:"unique;type:varchar(30);not null"`
	Avatar   string `json:"avatar" gorm:"type:varchar(1024);not null"`
	Intro    string `json:"intro" gorm:"type:varchar(255)"`
	Website  string `json:"website" gorm:"type:varchar(255)"`
}

type UserInfoVo struct {
	UserInfo
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
}

func GetUserInfoById(db *gorm.DB, id int) (*UserInfo, error) {
	var userInfo UserInfo
	result := db.Model(&userInfo).Where("id", id).First(&userInfo)
	return &userInfo, result.Error
}

func GetUserAuthInfoByName(db *gorm.DB, name string) (*UserAuth, error) {
	var userAuth UserAuth
	result := db.Where(&UserAuth{Username: name}).First(&userAuth)
	return &userAuth, result.Error
}

func GetUserList(db *gorm.DB, page, size int, loginType int8, nickname, username string) (list []UserAuth, total int64, err error) {
	if loginType != 0 {
		db = db.Where("login_type = ?", loginType)
	}

	if username != "" {
		db = db.Where("username Like ?", "%"+username+"%")
	}

	result := db.Model(&UserAuth{}).
		Joins("LEFT JOIN user_info ON user_info.id = user_auth.user_info_id").
		Where("user_info.nickname LIKE ?", "%"+nickname+"%").
		Preload("UserInfo").
		Count(&total).
		Scopes(Paginate(page, size)).
		Find(&list)

	return list, total, result.Error
}

// 更新用户昵称与角色信息
func UpdateUserNicknameAndRole(db *gorm.DB, authId int, nickname string, roleIds []int) error {
	userAuth, err := GetUserAuthInfoById(db, authId)
	if err != nil {
		return err
	}

	userinfo := UserInfo{
		Model:    Model{ID: userAuth.UserInfoId},
		Nickname: nickname,
	}
	result := db.Model(&userinfo).Updates(userinfo)
	if result.Error != nil {
		return result.Error
	}

	if len(roleIds) == 0 {
		return nil
	}

	result = db.Where(UserAuthRole{UserAuthId: userAuth.UserInfoId}).Delete(UserAuthRole{})
	if result.Error != nil {
		return result.Error
	}

	var userRoles []UserAuthRole
	for _, id := range roleIds {
		userRoles = append(userRoles, UserAuthRole{
			RoleId:     id,
			UserAuthId: userAuth.ID,
		})
	}

	result = db.Create(&userRoles)
	return result.Error
}

func UpdateUserPassword(db *gorm.DB, id int, password string) error {
	userAuth := UserAuth{
		Model:    Model{ID: id},
		Password: password,
	}

	result := db.Model(&userAuth).Updates(userAuth)
	return result.Error
}

func UpdateUserInfo(db *gorm.DB, id int, nickname, avatar, intro, website string) error {
	userInfo := UserInfo{
		Model:    Model{ID: id},
		Nickname: nickname,
		Avatar:   avatar,
		Intro:    intro,
		Website:  website,
	}

	result := db.
		Select("nickname", "avatar", "intro", "website").
		Updates(userInfo)
	return result.Error
}

func UpdateUserDisable(db *gorm.DB, id int, isDisable bool) error {
	userAuth := UserAuth{
		Model:     Model{ID: id},
		IsDisable: isDisable,
	}

	result := db.Model(&userAuth).Select("is_disable").Updates(&userAuth)
	return result.Error
}

func UpdateUserLoginInfo(db *gorm.DB, id int, ipAddress, ipSource string) error {
	now := time.Now()
	userAuth := UserAuth{
		IpAddress:     ipAddress,
		IpSource:      ipSource,
		LastLoginTime: &now,
	}

	result := db.Where("id", id).Updates(userAuth)
	return result.Error
}
