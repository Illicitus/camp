package blog

import "github.com/jinzhu/gorm"

var _ ArticleDB = &articleGorm{}

type articleGorm struct {
	db *gorm.DB
}

func (ag *articleGorm) All() ([]ArticleModel, error) {
	var articles []ArticleModel

	result := ag.db.Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}
	return articles, nil
}

func (ag *articleGorm) Create(article *ArticleModel) error {
	return ag.db.Create(article).Error
}
