module main

go 1.15

replace github.com/Edu15/recipe-golang-webservice/src/service => ../service

replace github.com/Edu15/recipe-golang-webservice/src/repository => ../repository

replace github.com/Edu15/recipe-golang-webservice/src/domain => ../domain

replace github.com/Edu15/recipe-golang-webservice/src/render => ../render

replace github.com/Edu15/recipe-golang-webservice/src/render/html => ../render/html

require (
	github.com/Edu15/recipe-golang-webservice/src/render v0.0.0-00010101000000-000000000000 // indirect
	github.com/Edu15/recipe-golang-webservice/src/service v0.0.0-00010101000000-000000000000
)
