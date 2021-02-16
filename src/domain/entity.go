package domain

type ResponseFormat int

const (
	JSON ResponseFormat = iota
	HTML ResponseFormat = iota
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
	Difficulty      RecipeDifficulty
	Rating          int
	PreparationTime int
	Serving         int
	Ingredients     []string
	Steps           []string
	PublishedDate   string
	AccessCount     int
	ImageURL        string
}

// RecipePreview struct
type RecipePreview struct {
	ID          int
	Title       string
	Description string
}

// RecipeForm struct
type RecipeForm struct {
	Recipe       Recipe
	Categories   []RecipeCategory
	Difficulties []RecipeDifficulty
}
