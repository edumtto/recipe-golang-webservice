Experimental project of webservice using Golang and Postgres created only for learning purposes.

### Running the service

- Configure a Postgres database with the specified format in `configuration/CreatingTablesCommands.txt`;
- Run the following command on terminal from the src directory: `go run main.go`
- To render responses in JSON format, add the json keyword in the end of the command: `go run main.go json`

### Calling the service

POST /recipes to add a new recipe
GET /recipes to fetch all existing recipes
GET /recipes/{recipeId} to fetch a single recipe using its ID
PUT /recipes/{recipeId} to update an existing recipe
DELETE /recipes/{itemId} to delete a recipe

**Only for HTML forms:**
GET /new to fech a form to create a new recipe
GET /create to fech a form to create a new recipe
GET /edit to fech a form to create a new recipe
GET /update/{recipeId} to update an existing recipe
GET /delete/{recipeId} to delete a recipe
