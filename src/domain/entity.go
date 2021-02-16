package domain

type ResponseFormat int

const (
	JSON ResponseFormat = iota
	HTML ResponseFormat = iota
)

// RecipeAuthor struct
type RecipeAuthor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RecipeCategory struct
type RecipeCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RecipeDifficulty struct
type RecipeDifficulty struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Recipe struct
type Recipe struct {
	ID              int              `json:"id"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	Author          RecipeAuthor     `json:"author"`
	Category        RecipeCategory   `json:"category"`
	Difficulty      RecipeDifficulty `json:"difficulty"`
	Rating          int              `json:"rating"`
	PreparationTime int              `json:"preparation_time"`
	Serving         int              `json:"serving"`
	Ingredients     []string         `json:"ingredients"`
	Steps           []string         `json:"steps"`
	PublishedDate   string           `json:"published_date"`
	AccessCount     int              `json:"access_count"`
	ImageURL        string           `json:"image_url"`
}

// RecipePreview struct
type RecipePreview struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// RecipeForm struct
type RecipeForm struct {
	Recipe       Recipe             `json:"recipe"`
	Categories   []RecipeCategory   `json:"categories"`
	Difficulties []RecipeDifficulty `json:"difficulties"`
}
