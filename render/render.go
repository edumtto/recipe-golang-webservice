package render

import (
	"net/http"
	"text/template"

	"github.com/Edu15/recipe-golang-webservice/domain"
)

type Interface interface {
	RenderRecipeList(w http.ResponseWriter, recipePreviews *[]domain.RecipePreview)
	RenderRecipe(w http.ResponseWriter, recipe *domain.Recipe)
	RenderRecipeEditor(w http.ResponseWriter, recipe *domain.Recipe)
	RenderNewRecipeForm(w http.ResponseWriter)
}

type HTMLRenderer struct{}

var templates = template.Must(template.ParseFiles("./tmpl/edit-recipe.html", "./tmpl/view-recipe.html", "./tmpl/recipe-list.html", "./tmpl/new-recipe.html"))

func (HTMLRenderer) RenderRecipeList(w http.ResponseWriter, recipePreviews *[]domain.RecipePreview) {
	err := templates.ExecuteTemplate(w, "recipe-list.html", recipePreviews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (HTMLRenderer) RenderRecipe(w http.ResponseWriter, recipe *domain.Recipe) {
	err := templates.ExecuteTemplate(w, "view-recipe.html", recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (HTMLRenderer) RenderRecipeEditor(w http.ResponseWriter, recipe *domain.Recipe) {
	err := templates.ExecuteTemplate(w, "edit-recipe.html", recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (HTMLRenderer) RenderNewRecipeForm(w http.ResponseWriter) {
	err := templates.ExecuteTemplate(w, "new-recipe.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
