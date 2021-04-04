package blog

import (
	"camp/apps/accounts"
	"camp/core/utils"
	"camp/core/web"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"reflect"
)

var _ web.Model = &ArticleModel{}

type ArticleModel struct {
	gorm.Model
	AuthorID int                `gorm:"not null"`
	Author   accounts.UserModel `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title    string             `gorm:"not null"`
	Body     string             `gorm:"not null"`
}

func (_ ArticleModel) IsGormModel() {}
func (m ArticleModel) TableName() string {
	return utils.NormalizeModelName(SubAppName, reflect.TypeOf(m).Name())
}

type ArticleDB interface{}

var _ web.Model = &CommentModel{}

type CommentModel struct {
	gorm.Model
	AuthorID  int                `gorm:"not null"`
	Author    accounts.UserModel `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ArticleID int                `gorm:"not null"`
	Article   ArticleModel       `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Body      string             `gorm:"not null"`
}

func (_ CommentModel) IsGormModel() {}
func (m CommentModel) TableName() string {
	return utils.NormalizeModelName(SubAppName, reflect.TypeOf(m).Name())
}

type CommentDB interface{}
