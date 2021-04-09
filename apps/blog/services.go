package blog

import (
	"github.com/jinzhu/gorm"
)

type ArticleService interface {
	//All() ([]*ArticleModel, error)
	ArticleDB
}

func NewArticleService(db *gorm.DB) ArticleService {
	ag := &articleGorm{db}
	av := newArticleValidator(ag)

	return &articleService{
		ArticleDB: av,
	}
}

var _ ArticleService = &articleService{}

type articleService struct {
	ArticleDB
}

//func (as *articleService) All() ([]*ArticleModel, error) {
//	articles, err := as.db
//	foundUser, err := us.ByEmail(email)
//	if err != nil {
//		return nil, err
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+us.pepper))
//	if err != nil {
//		switch err {
//		case bcrypt.ErrMismatchedHashAndPassword:
//			return nil, utils.GormErr.InvalidPassword
//		default:
//			return nil, err
//		}
//	}
//	return foundUser, nil
//}
