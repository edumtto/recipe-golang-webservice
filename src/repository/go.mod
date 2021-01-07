module github.com/Edu15/recipe-golang-webservice/src/repository

go 1.15

replace github.com/Edu15/recipe-golang-webservice/src/domain => ../domain

require (
	github.com/Edu15/recipe-golang-webservice/src/domain v0.0.0-20210106130618-ba477f364819
	github.com/lib/pq v1.9.0
)
