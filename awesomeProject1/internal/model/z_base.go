package model

import (
	"gorm.io/gorm"
	"time"
)

// 通用模型
type Model struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OptionVo struct {
	ID   int    `json:"value"`
	Name string `json:"label"`
}

// 分页
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// 通用CRUD
// 创建数据（可以创建[单挑]数据，也可以批量创建）
func Create[T any](db *gorm.DB, data *T) (*T, error) {
	result := db.Create(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func Get[T any](db *gorm.DB, data *T, query string, args ...any) (*T, error) {
	result := db.Where(query, args...).First(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func Update[T any](db *gorm.DB, data *T, slt ...string) error {
	db = db.Model(&data)
	if len(slt) > 0 {
		db = db.Select(slt)
	}
	result := db.Updates(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdatesMap[T any](db *gorm.DB, data *T, maps map[string]any, query string, args ...any) error {
	result := db.Model(data).Where(query, args...).Updates(maps)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func Updates[T any](db *gorm.DB, data T, query string, args ...any) error {
	result := db.Model(&data).Where(query, args...).Updates(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func List[T any](db *gorm.DB, data T, slt, order, query string, args ...any) (T, error) {
	db = db.Model(data).Select(slt).Order(order)
	if query != "" {
		db = db.Where(query, args...)
	}
	result := db.Find(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

func Delete[T any](db *gorm.DB, data T, query string, args ...any) error {
	result := db.Where(query, args...).Delete(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func Count[T any](db *gorm.DB, data *T, where ...any) (int, error) {
	var total int64
	db = db.Model(data)
	if len(where) > 0 {
		db = db.Where(where[0], where[1:]...)
	}
	result := db.Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(total), nil
}
