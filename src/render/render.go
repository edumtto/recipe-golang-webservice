package render

import (
	"net/http"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
)

type Interface interface {
	RenderRecipeList(w http.ResponseWriter, recipePreviews *[]domain.RecipePreview)
	RenderRecipe(w http.ResponseWriter, recipe *domain.Recipe)
	RenderRecipeEditor(w http.ResponseWriter, recipe *domain.Recipe)
	RenderNewRecipeForm(w http.ResponseWriter)
}
