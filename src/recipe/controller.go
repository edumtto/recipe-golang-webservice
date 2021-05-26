package recipe

import (
	"fmt"
	"net/http"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
)

// Controller interface
type Controller interface {
	List(w http.ResponseWriter, r *http.Request)
	View(w http.ResponseWriter, r *http.Request, id string)
	Edit(w http.ResponseWriter, r *http.Request, id string)
	New(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request, id string)
	Delete(w http.ResponseWriter, r *http.Request, id string)
}

// Controller
type controller struct {
	service   Service
	presenter domain.Render
	format    domain.ResponseFormat
}

// NewController creates a new instance o Controller injecting a service.
func NewController(service Service, presenter domain.Render, format domain.ResponseFormat) Controller {
	return &controller{
		service:   service,
		presenter: presenter,
		format:    format,
	}
}

// List fetches a list of all recipes and present a formated result.
func (c controller) List(w http.ResponseWriter, r *http.Request) {
	recipePreviews, err := c.service.List(w, r)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.presenter.RenderRecipeList(w, recipePreviews)
}

func (c controller) View(w http.ResponseWriter, r *http.Request, id string) {
	recipe, err := c.service.View(w, r, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	c.presenter.RenderRecipe(w, recipe)
}

func (c controller) Edit(w http.ResponseWriter, r *http.Request, id string) {
	recipeForm, err := c.service.Edit(w, r, id)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	c.presenter.RenderRecipeEditor(w, recipeForm)
}

func (c controller) New(w http.ResponseWriter, r *http.Request) {
	recipeForm, err := c.service.New(w, r)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.presenter.RenderNewRecipeForm(w, recipeForm)
}

func (c controller) Create(w http.ResponseWriter, r *http.Request) {
	insertedID, err := c.service.Create(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	if c.format == domain.HTML {
		url := fmt.Sprintf("/recipes/%d", insertedID)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func (c controller) Update(w http.ResponseWriter, r *http.Request, id string) {
	err := c.service.Update(w, r, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	if c.format == domain.HTML {
		http.Redirect(w, r, "/recipes/"+id, http.StatusFound)
	}
}

func (c controller) Delete(w http.ResponseWriter, r *http.Request, id string) {
	err := c.service.Delete(w, r, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if c.format == domain.HTML {
		http.Redirect(w, r, "/recipes/all/", http.StatusFound)
	}
}
