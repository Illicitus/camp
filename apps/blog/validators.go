package blog

import "camp/core/utils"

type articleValFunc func(model *ArticleModel) error

func runArticleValFuncs(user *ArticleModel, fns ...articleValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ ArticleDB = &articleValidator{}

type articleValidator struct {
	ArticleDB
}

func newArticleValidator(adb ArticleDB) *articleValidator {
	return &articleValidator{
		ArticleDB: adb,
	}
}

func (av *articleValidator) Create(article *ArticleModel) error {
	valFuncs := []articleValFunc{
		av.titleRequired,
		av.bodyRequired,
		av.userIDRequired,
	}
	if err := runArticleValFuncs(article, valFuncs...); err != nil {
		return err
	}
	return av.ArticleDB.Create(article)
}

func (av *articleValidator) userIDRequired(article *ArticleModel) error {
	if article.AuthorID <= 0 {
		return utils.ValErr.UserIDRequired
	}
	return nil
}

func (av *articleValidator) titleRequired(article *ArticleModel) error {
	if article.Title == "" {
		return utils.ValErr.TitleRequired
	}
	return nil
}

func (av *articleValidator) bodyRequired(article *ArticleModel) error {
	if article.Body == "" {
		return utils.ValErr.BodyRequired
	}
	return nil
}
