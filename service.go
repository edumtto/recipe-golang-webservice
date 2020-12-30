package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "edu"
	password = "1234"
	dbname   = "recipes_db"
)

// RecipeAuthor struct
type RecipeAuthor struct {
	ID   int
	Name string
}

// RecipeCategory struct
type RecipeCategory struct {
	ID   int
	Name string
}

// RecipeDificulty struct
type RecipeDificulty struct {
	ID   int
	Name string
}

// Recipe struct
type Recipe struct {
	ID              int
	Title           string
	Description     string
	Author          RecipeAuthor
	Category        RecipeCategory
	Dificulty       RecipeDificulty
	Rating          int
	PreparationTime int
	Serving         int
	Ingredients     []string
	Steps           []string
	//PublishedDate Date
	AccessCount int
	ImageURL    string
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view-recipe.html"))
var validPath = regexp.MustCompile("^/receita/([a-zA-Z0-9]+)$")
var db *sql.DB

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1])
	}
}

func recipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := fetchFullRecipe(recipeID)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	if recipe == nil {
		fmt.Fprintf(w, "Receita não encontrada!")
		return
	}

	renderTemplate(w, "view-recipe", recipe)
}

func renderTemplate(w http.ResponseWriter, tmpl string, recipe *Recipe) {
	err := templates.ExecuteTemplate(w, tmpl+".html", recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func connectWithDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func fetchFullRecipe(recipeID int) (*Recipe, error) {
	recipe, err := fetchRecipe(recipeID)
	if err != nil {
		return nil, err
	}

	author, err := fetchAuthor(recipe.Author.ID)
	if err != nil {
		return nil, err
	}
	recipe.Author.Name = author.Name

	category, err := fetchCategory(recipe.Category.ID)
	if err != nil {
		return nil, err
	}
	recipe.Category.Name = category.Name

	dificulty, err := fetchDificulty(recipe.Dificulty.ID)
	if err != nil {
		return nil, err
	}
	recipe.Dificulty.Name = dificulty.Name

	return recipe, nil
}

func fetchRecipe(recipeID int) (*Recipe, error) {
	sqlStatement := `SELECT id, title, description, author_id, category_id, dificulty_id, rating,
	preparation_time, serving, ingredients, steps, access_count, image
	FROM recipe WHERE id=$1;`

	row := db.QueryRow(sqlStatement, recipeID)
	var recipe Recipe
	var ingredients, steps string

	err := row.Scan(&recipe.ID, &recipe.Title, &recipe.Description, &recipe.Author.ID,
		&recipe.Category.ID, &recipe.Dificulty.ID, &recipe.Rating, &recipe.PreparationTime,
		&recipe.Serving, &ingredients, &steps, &recipe.AccessCount,
		&recipe.ImageURL)

	recipe.Ingredients = strings.Split(ingredients, "|")
	recipe.Steps = strings.Split(steps, "|")

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &recipe, nil
	default:
		return nil, err
	}
}

func fetchAuthor(ID int) (*RecipeAuthor, error) {
	var author RecipeAuthor

	sqlStatement := `SELECT id, name FROM author WHERE id=$1;`
	row := db.QueryRow(sqlStatement, ID)
	err := row.Scan(&author.ID, &author.Name)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &author, nil
	default:
		return nil, err
	}
}

func fetchCategory(ID int) (*RecipeCategory, error) {
	var category RecipeCategory

	sqlStatement := `SELECT id, name FROM category WHERE id=$1;`
	row := db.QueryRow(sqlStatement, ID)
	err := row.Scan(&category.ID, &category.Name)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &category, nil
	default:
		return nil, err
	}
}

func fetchDificulty(ID int) (*RecipeDificulty, error) {
	var dificulty RecipeDificulty

	sqlStatement := `SELECT id, name FROM dificulty WHERE id=$1;`
	row := db.QueryRow(sqlStatement, ID)
	err := row.Scan(&dificulty.ID, &dificulty.Name)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &dificulty, nil
	default:
		return nil, err
	}
}

func insertRecipeTest() {
	sqlStatement := `
	INSERT INTO recipe (title, description, author_id, category_id, dificulty_id, preparation_time, serving, ingredients, steps, image)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id;`

	ingredients := "Massa:,1 lata de leite condensado, 1 xícara de leite de vaca, 4 ovos inteiros, Calda: 1 xícara (chá) de açúcar, 1/3 de xícara (chá) de água"
	steps := `Calda:, Em uma panela, misture a água e o açúcar até formar uma calda. Unte uma forma com a calda e reserve.
	Massa:, Bata todos os ingredientes no liquidificador e despeje na forma caramelizada., Leve para assar em banho-maria por 40 minutos.,
	Desenforme e sirva.`
	imgURL := "https://img.itdg.com.br/tdg/images/recipes/000/003/687/38788/38788_original.jpg?mode=crop&width=710&height=400"

	_, err := db.Exec(sqlStatement, "Pudim de doce de leite 2", "Receita de pudim de doce de leite.", 2, 10, 1, 40, 20, ingredients, steps, imgURL)
	if err != nil {
		panic(err)
	}
}

func main() {
	db = connectWithDatabase()
	defer db.Close()

	http.HandleFunc("/receita/", makeHandler(recipeHandler))
	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
