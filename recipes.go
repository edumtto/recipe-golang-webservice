package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "recipes_db"
)

// Recipe struct
type Recipe struct {
	ID              int
	Title           string
	Description     string
	AuthorID        int
	CategoryID      int
	DificultyID     int
	Rating          int
	PreparationTime int
	Serving         int
	Ingredients     string
	Steps           string
	//PublishedDate Date
	AccessCount int
	ImageURL    *string
}

//var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
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

	recipe, err := fetchRecipe(recipeID)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	if recipe == nil {
		fmt.Fprintf(w, "Receita não encontrada!")
		return
	}

	//fmt.Println(recipe)
	var imageHtml string
	if recipe.ImageURL == nil {
		imageHtml = ""
	} else {
		imageHtml = "<img src='" + *recipe.ImageURL + "' height='200px'></img>"
	}
	fmt.Fprintf(w, "<h1>%s</h1>%s<p>%s</p>", recipe.Title, imageHtml, recipe.Description)
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

func fetchRecipe(recipeID int) (*Recipe, error) {
	sqlStatement := `SELECT id, title, description, author_id, category_id, dificulty_id, rating,
	preparation_time, serving, ingredients, steps, access_count, image
	FROM recipe WHERE id=$1;`

	row := db.QueryRow(sqlStatement, recipeID)
	var recipe Recipe

	err := row.Scan(&recipe.ID, &recipe.Title, &recipe.Description, &recipe.AuthorID,
		&recipe.CategoryID, &recipe.DificultyID, &recipe.Rating, &recipe.PreparationTime,
		&recipe.Serving, &recipe.Ingredients, &recipe.Steps, &recipe.AccessCount,
		&recipe.ImageURL)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &recipe, nil
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
