package blog

type ArticleForm struct {
	Title string `schema:"tile"`
	Body  string `schema:"body"`
}
