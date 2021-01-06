module main

go 1.15

replace github.com/Edu15/recipe-golang-webservice/service => ./service

replace github.com/Edu15/recipe-golang-webservice/repository => ./repository

replace github.com/Edu15/recipe-golang-webservice/domain => ./domain

require github.com/Edu15/recipe-golang-webservice/service v0.0.0-00010101000000-000000000000
