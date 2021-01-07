module github.com/Edu15/recipe-golang-webservice/service

go 1.15

replace github.com/Edu15/recipe-golang-webservice/domain => ../domain

replace github.com/Edu15/recipe-golang-webservice/repository => ../repository

replace github.com/Edu15/recipe-golang-webservice/render => ../render

require (
	github.com/Edu15/recipe-golang-webservice/domain v0.0.0-20210106130618-ba477f364819
	github.com/Edu15/recipe-golang-webservice/render v0.0.0-00010101000000-000000000000
	github.com/Edu15/recipe-golang-webservice/repository v0.0.0-00010101000000-000000000000
)
