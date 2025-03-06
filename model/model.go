// Package model has domain models used throughout the application.
package model

type Recipe struct {
	ID          int
	Name        string
	Url         string
	Description string
	ImageURL    string
}
