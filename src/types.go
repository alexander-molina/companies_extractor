package main

// Company represents info about found company
type Company struct {
	Name     string
	URL      string
	Location string
	Emails   map[string]string
}
