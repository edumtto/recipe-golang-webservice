module main

go 1.15

replace github.com/Edu15/recipe-golang-webservice/domain => ./domain

replace github.com/Edu15/recipe-golang-webservice/repository => ./repository

require (
	github.com/Edu15/recipe-golang-webservice/domain v0.0.0-00010101000000-000000000000
	github.com/Edu15/recipe-golang-webservice/repository v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.9.0
)
