// Package model has domain models used throughout the application.
package model

import "recipeze/parsing"

type Recipe struct {
	ID          int
	Name        string
	Url         string
	Description string
	ImageURL    string
	GroupID     int
	Data        *parsing.RecipeCollection
}

type User struct {
	ID    int
	Name  string
	Email string
}

type Group struct {
	ID      int
	Name    string
	Members []GroupMember
}

type GroupMember struct {
	ID      int
	Name    string
	Email   string
	IsAdmin bool
}
