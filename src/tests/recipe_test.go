package tests

/*
import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
	"github.com/Edu15/recipe-golang-webservice/src/recipe"
	"github.com/Edu15/recipe-golang-webservice/src/tests/mock_domain"
	"github.com/golang/mock/gomock"
)

// func TestNewService(t *testing.T) {
// 	recipeService := recipe.NewService(domain.JSON)
// 	if recipeService == nil {
// 		t.Errorf("Service is nil")
// 	}
// }

func TestList(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repository := mock_domain.NewMockRepository(controller)
	//repo.EXPECT().FetchRecipePreviews()
	renderer := mock_domain.NewMockRender(controller)

	sut := recipe.NewService(repository, renderer, domain.JSON)

	req, err := http.NewRequest("GET", "/recipes/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	preview := domain.RecipePreview{
		ID:          6,
		Title:       "Brigadeiro",
		Description: "Uma delícia",
		ImageURL:    "",
	}
	previews := &[]domain.RecipePreview{preview}

	repository.EXPECT().FetchRecipePreviews().Return(previews, nil)
	renderer.EXPECT().RenderRecipeList(gomock.Eq(rr), gomock.Eq(previews))

	sut.List(rr, req)

	// sut := recipe.NewService(domain.JSON)

	// req, err := http.NewRequest("GET", "/recipes/", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// rr := httptest.NewRecorder()

	// sut.List(rr, req, "")

	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }

	// expected := `[{"id":6,"title":"Pudim de doce de leite (uma delícia)","description":"Receita de pudim padrão. Com leite condensado e coco...."},{"id":36,"title":"Bolo de fubá","description":"Um bolo delicioso"},{"id":8,"title":"Abóbora com carne seca","description":"Uma receita deliciosa para seu almoço."},{"id":9,"title":"Bolo de pote","description":"Bolo muito bom!"},{"id":23,"title":"rwer","description":"werwer"},{"id":4,"title":"Bolo comum","description":"Receita de bolo padrão."},{"id":5,"title":"Brigadeiro","description":"O brigadeiro é um doce genuinamente brasileiro. Um orgulho só! Essa delícia de chocolate faz a alegria da criançada e de muita gente grande em qualquer circunstância."},{"id":39,"title":"título123","description":"descrição"},{"id":40,"title":"título123","description":"descrição"},{"id":53,"title":"título1234","description":"descrição"}]`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
*/
