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

// RecipeDifficulty struct
type RecipeDifficulty struct {
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
	Dificulty       RecipeDifficulty
	Rating          int
	PreparationTime int
	Serving         int
	Ingredients     []string
	Steps           []string
	//PublishedDate Date
	AccessCount int
	ImageURL    string
}

// RecipePreview struct
type RecipePreview struct {
	ID          int
	Title       string
	Description string
}

var templates = template.Must(template.ParseFiles("tmpl/edit-recipe.html", "tmpl/view-recipe.html", "tmpl/recipe-list.html", "tmpl/new-recipe.html"))
var validPath = regexp.MustCompile("^/(view|edit|create|update|new|delete|home)/?([a-zA-Z0-9]+)?$")
var db *sql.DB

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			fmt.Println("Invalid path!")
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func listRecipesHandler(w http.ResponseWriter, r *http.Request, id string) {
	recipePreviews, err := fetchRecipePreviews(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = templates.ExecuteTemplate(w, "recipe-list.html", recipePreviews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
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

func editRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
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

	renderTemplate(w, "edit-recipe", recipe)
}

func newRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	renderTemplate(w, "new-recipe", nil)
}

func createRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	insertedID, err := insertRecipe(w, r)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("/view/%d", insertedID)
	http.Redirect(w, r, url, http.StatusFound)
}

func updateRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := fetchRecipe(recipeID)
	if recipe == nil {
		panic(err)
	}

	err = updateRecipe(w, r, recipeID)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/view/"+id, http.StatusFound)
}

func deleteRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = removeRecipe(w, r, recipeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusFound)
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

func fetchDificulty(ID int) (*RecipeDifficulty, error) {
	var dificulty RecipeDifficulty

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

func fetchRecipePreviews(w http.ResponseWriter, r *http.Request) (*[]RecipePreview, error) {
	sqlStatement := `SELECT id, title, description FROM recipe LIMIT $1;`
	rows, err := db.Query(sqlStatement, 10)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var previews []RecipePreview

	for rows.Next() {
		var preview RecipePreview
		if err := rows.Scan(&preview.ID, &preview.Title, &preview.Description); err != nil {
			log.Fatal(err)
		}
		previews = append(previews, preview)
	}

	return &previews, err
}

func updateRecipe(w http.ResponseWriter, r *http.Request, id int) error {
	sqlStatement := `UPDATE recipe 
	SET title = $2, description = $3, preparation_time = $4, serving = $5, image =$6
	WHERE id = $1;`
	title := r.FormValue("title")
	description := r.FormValue("description")
	//difficultyID, _ := strconv.Atoi(r.FormValue("difficulty"))
	preparationTime, _ := strconv.Atoi(r.FormValue("preparation-time"))
	// ingredients := strings.Join(recipe.Ingredients, "|")
	// steps := strings.Join(recipe.Steps, "|")
	serving := r.FormValue("serving")
	imageURL := r.FormValue("imgURL")
	_, err := db.Exec(sqlStatement, id, title, description, preparationTime, serving, imageURL)
	return err
}

func insertRecipe(w http.ResponseWriter, r *http.Request) (int, error) {
	sqlStatement := `
	INSERT INTO recipe (title, description, author_id, category_id, dificulty_id, preparation_time, serving, ingredients, steps, image)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id;`

	title := r.FormValue("title")
	description := r.FormValue("description")
	categoryID, _ := strconv.Atoi(r.FormValue("category"))
	difficultyID, _ := strconv.Atoi(r.FormValue("difficulty"))
	preparationTime, _ := strconv.Atoi(r.FormValue("preparation-time"))
	serving := r.FormValue("serving")
	ingredients := r.FormValue("ingredients")
	steps := r.FormValue("steps")
	imageURL := r.FormValue("imgURL")
	var id int
	err := db.QueryRow(sqlStatement, title, description, 2, categoryID, difficultyID, preparationTime, serving, ingredients, steps, imageURL).Scan(&id)
	fmt.Println(id)
	return id, err
}

func removeRecipe(w http.ResponseWriter, r *http.Request, id int) error {
	sqlStatement := `DELETE FROM recipe WHERE id = $1;`
	_, err := db.Exec(sqlStatement, id)
	return err
}

func main() {
	db = connectWithDatabase()
	defer db.Close()

	http.HandleFunc("/home/", makeHandler(listRecipesHandler))
	http.HandleFunc("/view/", makeHandler(viewRecipeHandler))
	http.HandleFunc("/edit/", makeHandler(editRecipeHandler))
	http.HandleFunc("/new/", makeHandler(newRecipeHandler))
	http.HandleFunc("/create/", makeHandler(createRecipeHandler))
	http.HandleFunc("/update/", makeHandler(updateRecipeHandler))
	http.HandleFunc("/delete/", makeHandler(deleteRecipeHandler))
	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
