package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Edu15/recipe-golang-webservice/src/domain"

	// Postgres driver
	_ "github.com/lib/pq"
)

// repository struct
type repository struct {
	db *sql.DB
}

// NewRepository method
func NewRepository(db *sql.DB) domain.Repository {
	return &repository{
		db: db,
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

// FetchRecipe returns the recipe with the inputed ID, if it exists.
func (repo repository) FetchRecipe(recipeID int) (*domain.Recipe, error) {
	sqlStatement := `SELECT id, title, description, author_id, category_id, difficulty_id, rating,
	preparation_time, serving, ingredients, steps, access_count, image, published_date
	FROM recipe WHERE id=$1;`

	row := repo.db.QueryRow(sqlStatement, recipeID)
	var recipe domain.Recipe
	var ingredients, steps string
	var publishedDateTime time.Time

	err := row.Scan(&recipe.ID, &recipe.Title, &recipe.Description, &recipe.Author.ID,
		&recipe.Category.ID, &recipe.Difficulty.ID, &recipe.Rating, &recipe.PreparationTime,
		&recipe.Serving, &ingredients, &steps, &recipe.AccessCount,
		&recipe.ImageURL, &publishedDateTime)

	recipe.Ingredients = strings.Split(ingredients, "|")
	recipe.Steps = strings.Split(steps, "|")
	recipe.PublishedDate = formatDate(publishedDateTime)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &recipe, nil
	default:
		return nil, err
	}
}

func formatDate(input time.Time) string {
	dateStr := input.Format(time.RFC3339)[:10]
	t, _ := time.Parse("2006-01-02", dateStr)
	return t.Format("02/Jan/2006")
}

// FetchAuthor returns the author with the inputed ID, if it exists.
func (repo repository) FetchAuthor(ID int) (*domain.RecipeAuthor, error) {
	var author domain.RecipeAuthor

	sqlStatement := `SELECT id, name FROM author WHERE id=$1;`
	row := repo.db.QueryRow(sqlStatement, ID)
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

// FetchCategory returns the category with the inputed ID, if it exists.
func (repo repository) FetchCategory(ID int) (*domain.RecipeCategory, error) {
	var category domain.RecipeCategory

	sqlStatement := `SELECT id, name FROM category WHERE id=$1;`
	row := repo.db.QueryRow(sqlStatement, ID)
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

// FetchDifficulty returns the difficulty with the inputed ID, if it exists.
func (repo repository) FetchDifficulty(ID int) (*domain.RecipeDifficulty, error) {
	var difficulty domain.RecipeDifficulty

	sqlStatement := `SELECT id, name FROM difficulty WHERE id=$1;`
	row := repo.db.QueryRow(sqlStatement, ID)
	err := row.Scan(&difficulty.ID, &difficulty.Name)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &difficulty, nil
	default:
		return nil, err
	}
}

// FetchRecipePreviews returns a list with short descriptions for each recipe.
func (repo repository) FetchRecipePreviews(w http.ResponseWriter, r *http.Request) (*[]domain.RecipePreview, error) {
	sqlStatement := `SELECT id, title, description, image FROM recipe LIMIT $1;`
	rows, err := repo.db.Query(sqlStatement, 10)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var previews []domain.RecipePreview

	for rows.Next() {
		var preview domain.RecipePreview
		if err := rows.Scan(&preview.ID, &preview.Title, &preview.Description, &preview.ImageURL); err != nil {
			log.Fatal(err)
		}
		previews = append(previews, preview)
	}

	return &previews, err
}

// UpdateRecipe updates the information related to a recipe with the inputed ID, if it exists.
func (repo repository) UpdateRecipe(w http.ResponseWriter, r *http.Request, id int) error {
	sqlStatement := `UPDATE recipe 
	SET title = $2, description = $3, preparation_time = $4, serving = $5, image =$6
	WHERE id = $1;`
	title := r.FormValue("title")
	description := r.FormValue("description")
	preparationTime, _ := strconv.Atoi(r.FormValue("preparation-time"))
	serving := r.FormValue("serving")
	imageURL := r.FormValue("imgURL")
	_, err := repo.db.Exec(sqlStatement, id, title, description, preparationTime, serving, imageURL)
	return err
}

// InsertRecipe adds a new recipe in the repository.
func (repo repository) InsertRecipe(w http.ResponseWriter, r *http.Request) (int, error) {
	sqlStatement := `
	INSERT INTO recipe (title, description, author_id, category_id, difficulty_id, preparation_time, serving, ingredients, steps, image)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id;`

	title := r.FormValue("title")
	description := r.FormValue("description")
	categoryID, _ := strconv.Atoi(r.FormValue("category"))
	difficultyID, _ := strconv.Atoi(r.FormValue("difficulty"))
	preparationTime, _ := strconv.Atoi(r.FormValue("preparation-time"))
	serving, _ := strconv.Atoi(r.FormValue("serving"))
	ingredients := r.FormValue("ingredients")
	steps := r.FormValue("steps")
	imageURL := r.FormValue("imgURL")

	var id int
	err := repo.db.QueryRow(sqlStatement, title, description, 2, categoryID, difficultyID, preparationTime, serving, ingredients, steps, imageURL).Scan(&id)
	return id, err
}

// RemoveRecipe removes a recipe from the repository.
func (repo repository) RemoveRecipe(w http.ResponseWriter, r *http.Request, id int) error {
	sqlStatement := `DELETE FROM recipe WHERE id = $1;`
	_, err := repo.db.Exec(sqlStatement, id)
	return err
}

// FetchCategories returns a list of category options.
func (repo repository) FetchCategories() (*[]domain.RecipeCategory, error) {
	sqlStatement := `SELECT id, name FROM category;`
	rows, err := repo.db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var categories []domain.RecipeCategory

	for rows.Next() {
		var category domain.RecipeCategory
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Fatal(err)
		}
		categories = append(categories, category)
	}

	return &categories, err
}

// FetchDifficulties returns a list of difficult options.
func (repo repository) FetchDifficulties() (*[]domain.RecipeDifficulty, error) {
	sqlStatement := `SELECT id, name FROM difficulty;`
	rows, err := repo.db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var dificulties []domain.RecipeDifficulty

	for rows.Next() {
		var difficulty domain.RecipeDifficulty
		if err := rows.Scan(&difficulty.ID, &difficulty.Name); err != nil {
			log.Fatal(err)
		}
		dificulties = append(dificulties, difficulty)
	}

	return &dificulties, err
}
