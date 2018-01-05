// Class for handling login actions
// This class check und serialize actions for the login
// Package db, login, models is needed
package controllers

import (
	"html/template"
	"net/http"
	"toHuus/db"
	"toHuus/login"
	"toHuus/models"
)

// Function to show the login
//  Params: ResponseWriter, Request -> For execute
func ShowLogin(w http.ResponseWriter, r *http.Request){
	// Load templates
	t := template.Must(template.ParseFiles("./toHuus/views/header.html",
		"./toHuus/views/login.html", "./toHuus/views/footer.html"))
	// Executes templates with data (Nav, Message)
	t.ExecuteTemplate(w, "header",
		Load{getLoginNav(), models.UserMessage, models.UserData{}})
	t.ExecuteTemplate(w, "content", nil)
	t.ExecuteTemplate(w, "footer", nil)
}

// Basic function to handle the login check and carry out another action to login
//  Params: ResponseWriter, Request -> For execute
func CheckLogin(w http.ResponseWriter, r *http.Request) {
	// DB
	db.OpenConnection()
	// Check Cookies
	if login.CheckLogin(w,r) {
		// Redirect to UI
		models.UserMessage = login.Message
		http.Redirect(w, r, url + "ui", 301)
	} else {
		// Show Login
		models.UserMessage = login.Message
		ShowLogin(w,r)
	}
	db.CloseConnection()
}

// Helper function to show an specific Nav/Menu in header template
// Specific for every site-handler
//  Return: Nav(type from controller) -> Navigation data (Name, Anchor, Icon)
func getLoginNav() Nav{
	elements := []NavElement{
		{"Login","login", "home"},
		{"Data","data", "cogs"},
		{"About","about", "info"},
	}
	return Nav{elements}
}