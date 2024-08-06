package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type UserAuth struct {
	Model
	Username      string     `gorm:"unique;type:varchar(50)" json:"username"`
	Password      string     `gorm:"type:varchar(100)" json:"-"`
	LoginType     int        `gorm:"type:tinyint(1);comment:登录类型" json:"login_type"`
	IpAddress     string     `gorm:"type:varchar(20);comment:登录IP地址" json:"ip_address"`
	IpSource      string     `gorm:"type:varchar(50);comment:IP来源" json:"ip_source"`
	LastLoginTime *time.Time `json:"last_login_time"`
	IsDisable     bool       `json:"is_disable"`
	IsSuper       bool       `json:"is_super"` // 超级管理员只能后台设置

	UserInfoId int       `json:"user_info_id"`
	UserInfo   *UserInfo `json:"info"`
	Roles      []*Role   `json:"roles" gorm:"many2many:user_auth_role"`
}

func (u *UserAuth) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

type Role struct {
	Model
	Name      string `gorm:"unique" json:"name"`
	Label     string `gorm:"unique" json:"label"`
	IsDisable bool   `json:"is_disable"`

	Resources []Resource `json:"resources" gorm:"many2many:role_resource"`
	Menus     []Menu     `json:"menus" gorm:"many2many:role_menu"`
	Users     []UserAuth `json:"users" gorm:"many2many:user_auth_role"`
}

type Resource struct {
	Model
	Name      string `gorm:"unique;type:varchar(50)" json:"name"`
	ParentId  int    `json:"parent_id"`
	Url       string `gorm:"type:varchar(255)" json:"url"`
	Method    string `gorm:"type:varchar(10)" json:"request_method"`
	Anonymous bool   `json:"is_anonymous"`

	Roles []*Role `json:"roles" gorm:"many2many:role_resource"`
}

type Menu struct {
	Model
	ParentId     int    `json:"parent_id"`
	Name         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(20);unique" json:"name"`
	Path         string `gorm:"uniqueIndex"`
	Component    string `gorm:"type:varchar(50)" json:"component"`
	Icon         string `gorm:"type:varchar(50)" json:"icon"`
	OrderNum     int8   `json:"order_num"`
	Redirect     string `gorm:"type:varchar(50)" json:"redirect"`
	Catalogue    bool   `json:"is_catalogue"`
	Hidden       bool   `json:"is_hidden"`
	KeepAlive    bool   `json:"keep_alive"`
	External     bool   `json:"is_external"`
	ExternalLink string `gorm:"type:varchar(255)" json:"external_link"`

	Roles []*Role `json:"roles" gorm:"many2many:role_menu"`
}

type UserAuthRole struct {
	UserAuthId int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
	RoleId     int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
}

type RoleResource struct {
	RoleId     int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
	ResourceId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
}

type RoleMenu struct {
	RoleId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
	MenuId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
}

type RoleVO struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	IsDisable   bool      `json:"is_disable"`
	ResourceIds []int     `json:"resource_ids" gorm:"-"`
	MenuIds     []int     `json:"menu_ids" gorm:"-"`
}

// Menu
func SaveOrUpdateMenu(db *gorm.DB, menu *Menu) error {
	var result *gorm.DB

	if menu.ID > 0 {
		result = db.Model(menu).
			Select("name", "path", "component", "icon", "redirect", "parent_id", "order_num", "catalogue", "hidden", "keep_alive", "external").
			Updates(menu)
	} else {
		result = db.Create(menu)
	}

	return result.Error
}

// GetMenuIdsByRoleId 根据角色ID获取菜单ID列表
func GetMenuIdsByRoleId(db *gorm.DB, roleId int) (ids []int, err error) {
	// 执行查询，将结果存储到ids变量中
	result := db.Model(&RoleMenu{}).Where("role_id = ?", roleId).Pluck("menu_id", &ids)
	// 返回菜单ID列表和可能的错误信息
	return ids, result.Error
}

func GetMenuById(db *gorm.DB, id int) (menu *Menu, err error) {
	result := db.First(&menu, id)
	return menu, result.Error
}

func CheckMenuInUse(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&RoleMenu{}).Where("menu_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func CheckMenuHasChild(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&Menu{}).Where("parent_id = ?", id).Count(&count)
	return count > 0, result.Error
}

func GetAllMenuList(db *gorm.DB) (menu []Menu, err error) {
	result := db.Find(&menu)
	return menu, result.Error
}

func GetMenuListByUserId(db *gorm.DB, id int) (menus []Menu, err error) {
	var userAuth UserAuth
	result := db.Where(&UserAuth{Model: Model{ID: id}}).
		Preload("Roles").Preload("Roles.Menus").
		First(&userAuth)

	if result.Error != nil {
		return nil, result.Error
	}

	set := make(map[int]Menu)
	for _, role := range userAuth.Roles {
		for _, menu := range role.Menus {
			set[menu.ID] = menu
		}
	}

	for _, menu := range set {
		menus = append(menus, menu)
	}

	return menus, nil
}

func GetMenuList(db *gorm.DB, keyword string) (list []Menu, total int64, err error) {
	db = db.Model(&Menu{})
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}

	result := db.Count(&total).Find(&list)

	return list, total, result.Error
}

func DeleteMenu(db *gorm.DB, id int) error {
	result := db.Delete(&Menu{}, id)
	return result.Error
}

//Resource

func SaveOrUpdateResource(db *gorm.DB, id, pid int, name, url, method string) error {
	resource := Resource{
		Model:    Model{ID: id},
		Name:     name,
		Url:      url,
		Method:   method,
		ParentId: pid,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&resource)
	} else {
		result = db.Create(&resource)
	}

	return result.Error
}

func GetResourceIdsByRoleId(db *gorm.DB, roleId int) (ids []int, err error) {
	result := db.Model(&RoleResource{}).
		Where("role_id = ?", roleId).
		Pluck("resource_id", &ids)
	return ids, result.Error
}

func GetResourceList(db *gorm.DB, keyword string) (list []Resource, err error) {
	if keyword != "" {
		db = db.Where("name like ?", "%"+"keyword"+"%")
	}
	result := db.Find(&list)
	return list, result.Error
}

func GetResourceListByIds(db *gorm.DB, ids []int) (list []Resource, err error) {
	result := db.Where("id in ?", ids).Find(&list)
	return list, result.Error
}

// role
func SaveOrUpdateRole(db *gorm.DB, id int, name, label string, isDisable bool) error {
	role := Role{
		Model:     Model{ID: id},
		Name:      name,
		Label:     label,
		IsDisable: isDisable,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&role)
	} else {
		result = db.Create(&role)
	}
	return result.Error
}

func GetRoleOption(db *gorm.DB) (list []OptionVo, err error) {
	result := db.Model(&Role{}).Select("id", "name").Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func GetRoleList(db *gorm.DB, num, size int, keyword string) (list []RoleVO, total int64, err error) {
	db = db.Model(&Role{})
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	db.Count(&total)

	result := db.Select("id", "name", "label", "created_at", "is_disable").
		Scopes(Paginate(num, size)).
		Find(&list)
	return list, total, result.Error
}

func GetRoleIdsByUserId(db *gorm.DB, userAuthId int) (ids []int, err error) {
	result := db.Model(&UserAuthRole{UserAuthId: userAuthId}).
		Pluck("role_id", &ids)
	return ids, result.Error
}

func SaveRole(db *gorm.DB, name, label string) error {
	role := Role{
		Name:  name,
		Label: label,
	}
	result := db.Create(&role)
	return result.Error
}

func UpdateRole(db *gorm.DB, id int, name, label string, isDisable bool, resourceIds, menuIds []int) error {
	role := Role{
		Model:     Model{ID: id},
		Name:      name,
		Label:     label,
		IsDisable: isDisable,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := db.Model(&role).Select("name", "label", "is_disable").Updates(&role).Error; err != nil {
			return err
		}

		//role_resource
		if err := db.Delete(&RoleResource{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		for _, rid := range resourceIds {
			if err := db.Create(&RoleResource{RoleId: role.ID, ResourceId: rid}).Error; err != nil {
				return err
			}
		}

		// role_menu
		if err := db.Delete(&RoleMenu{}, "role_id = ?", id).Error; err != nil {
			return err
		}
		for _, mid := range menuIds {
			if err := db.Create(&RoleMenu{RoleId: role.ID, MenuId: mid}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func DeleteRoles(db *gorm.DB, ids []int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := db.Delete(&Role{}, "id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		result = db.Delete(&RoleResource{}, "role_id in ?", ids)
		if result.Error != nil {
			return result.Error
		}

		result = db.Delete(&RoleMenu{}, "role_id in ?", ids)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func GetUserAuthInfoById(db *gorm.DB, id int) (*UserAuth, error) {
	var userAuth = UserAuth{Model: Model{ID: id}}
	result := db.Model(&userAuth).
		Preload("Roles").Preload("UserInfo").
		First(&userAuth)
	return &userAuth, result.Error
}
