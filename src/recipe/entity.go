package recipe

// Author struct
type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Category struct
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Difficulty struct
type Difficulty struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Recipe struct
type Entity struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Author          Author     `json:"author"`
	Category        Category   `json:"category"`
	Difficulty      Difficulty `json:"difficulty"`
	Rating          int        `json:"rating"`
	PreparationTime int        `json:"preparation_time"`
	Serving         int        `json:"serving"`
	Ingredients     []string   `json:"ingredients"`
	Steps           []string   `json:"steps"`
	PublishedDate   string     `json:"published_date"`
	AccessCount     int        `json:"access_count"`
	ImageURL        string     `json:"image_url"`
}

// Preview struct
type Preview struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// Form struct
type Form struct {
	Recipe       Entity       `json:"recipe"`
	Categories   []Category   `json:"categories"`
	Difficulties []Difficulty `json:"difficulties"`
}
