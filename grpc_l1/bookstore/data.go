package main

import (
	"context"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

const defaultShelfSize = 100

// 使用GORM

func NewDB(dsn string) (*gorm.DB, error) {
	// 使用sqlite
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Shelf{}, &Book{})
	return db, nil
}

// 定义Model

// Shelf 书架
type Shelf struct {
	ID       int64 `gorm:"primaryKey"`
	Theme    string
	Size     int64
	CreateAt time.Time
	UpdateAt time.Time
}

// Book 图书
type Book struct {
	ID       int64 `gorm:"primaryKey"`
	Author   string
	Title    string
	ShelfID  int64
	CreateAt time.Time
	UpdateAt time.Time
}

type bookstore struct {
	db *gorm.DB
}

// CreateShelf 创建书架
func (b *bookstore) CreateShelf(ctx context.Context, data Shelf) (*Shelf, error) {
	if len(data.Theme) <= 0 {
		return nil, errors.New("invalid theme")
	}
	size := data.Size
	if size <= 0 {
		size = defaultShelfSize
	}
	v := Shelf{Theme: data.Theme, Size: size, CreateAt: time.Now(), UpdateAt: time.Now()}
	err := b.db.WithContext(ctx).Create(&v).Error
	return &v, err
}

func (b *bookstore) GetShelf(ctx context.Context, id int64) (*Shelf, error) {
	v := Shelf{}
	err := b.db.WithContext(ctx).First(&v, id).Error
	return &v, err
}

func (b *bookstore) ListShelves(ctx context.Context) (*[]Shelf, error) {
	v := []Shelf{}
	err := b.db.WithContext(ctx).Find(&v).Error
	return &v, err
}

func (b *bookstore) DeleteShelf(ctx context.Context, id int64) error {
	return b.db.WithContext(ctx).Delete(&Shelf{}, id).Error
}
func (b *bookstore) GetBookListByShelfID(ctx context.Context, shelfID int64, cursor string, pageSize int) (*[]Book, error) {
	var vl *[]Book
	err := b.db.WithContext(ctx).Where("shelf_id  = ? and id > ?", shelfID, cursor).Order("id desc").Limit(pageSize).Find(&vl).Error
	return vl, err
}
