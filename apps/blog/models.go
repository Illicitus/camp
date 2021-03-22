package blog

import (
	"camp/apps/accounts"
	"camp/core/web"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var _ web.Model = &Article{}

type Article struct {
	gorm.Model
	AuthorID int           `gorm:"not null"`
	Author   accounts.User `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title    string        `gorm:"not null"`
	Body     string        `gorm:"not null"`
}

func (a Article) IsGormModel() {}

type ArticleDB interface{}

var _ web.Model = &Comment{}

type Comment struct {
	gorm.Model
	AuthorID  int           `gorm:"not null"`
	Author    accounts.User `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ArticleID int           `gorm:"not null"`
	Article   Article       `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Body      string        `gorm:"not null"`
}

func (c Comment) IsGormModel() {}

type CommentDB interface{}
