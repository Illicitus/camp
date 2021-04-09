package blog

import (
	"camp/apps/accounts"
	"camp/core/web"
	"log"
	"net/http"
)

var (
	LayoutDir   = "apps/core/views/layouts/"
	TemplateDir = "apps/blog/views/"
)

type ArticleController struct {
	ListView          *web.View
	ArticleCreateView *web.View
	as                ArticleService
}

func NewArticleController(db *web.DB, cfg *web.AppConfig) *ArticleController {
	return &ArticleController{
		ListView:          web.NewView(TemplateDir, LayoutDir, "bootstrap", "list"),
		ArticleCreateView: web.NewView(TemplateDir, LayoutDir, "bootstrap", "createArticle"),
		as:                NewArticleService(db.Conn),
	}
}

func (ac *ArticleController) List(w http.ResponseWriter, r *http.Request) {
	articles, err := ac.as.All()
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	var vd web.Data
	vd.Yield = articles
	hub.ErrorHandler(ac.ListView.Render(w, r, vd))
}

func (ac *ArticleController) Create(w http.ResponseWriter, r *http.Request) {
	var vd web.Data
	var form ArticleForm
	vd.Yield = &form
	if err := web.ParseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		hub.ErrorHandler(ac.ArticleCreateView.Render(w, r, vd))
		return
	}

	user := accounts.UserContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/accounts/login", http.StatusFound)
		return
	}

	article := ArticleModel{
		Title:    form.Title,
		Body:     form.Body,
		AuthorID: user.ID,
	}
	if err := ac.as.Create(&article); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		hub.ErrorHandler(ac.ArticleCreateView.Render(w, r, vd))
		return
	}
	//url, err := g.r.Get("show_gallery").URL("id", fmt.Sprintf("%v", gallery.ID))
	//if err != nil {
	//	http.Redirect(w, r, "/", http.StatusFound)
	//	return
	//}
	//http.Redirect(w, r, url.Path, http.StatusFound)
}
