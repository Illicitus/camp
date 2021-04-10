package blog

type ArticleForm struct {
	Title string `schema:"title"`
	Body  string `schema:"body"`
}
