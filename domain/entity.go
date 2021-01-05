package domain

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
