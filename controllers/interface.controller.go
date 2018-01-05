// Class for handling interface actions
// This class check und serialize actions for the interface
// Package db, login, models is needed
package controllers

import (
	"net/http"
	"html/template"
	"encoding/json"
	"toHuus/login"
	"toHuus/models"
	"toHuus/db"
	"os"
	"fmt"
)

// An Config represents the data for database config
type config struct{
	Url 		string `json:"url"`
	Db		 	string `json:"database"`
	Port		string `json:"port"`
}

// Initialisation
var url = UrlConfig()

// Function to show the user interface
//  Params: ResponseWriter, Request -> For execute
func ShowInterface(w http.ResponseWriter, r *http.Request){
	// Load templates
	t := template.Must(template.ParseFiles("./toHuus/views/header.html",
		"./toHuus/views/interface.html", "./toHuus/views/footer.html"))
	// Executes templates with data (Nav, Message, Userdata, AllItemData)
	t.ExecuteTemplate(w, "header",
		Load{getInterfaceNav(), models.UserMessage, models.GetUserData(r)})
	t.ExecuteTemplate(w, "content", models.GetAllDTE(models.GetUserData(r).Id))
	t.ExecuteTemplate(w, "footer", nil)
}

// Basic function to handle the login check and carry out another action to interface
//  Params: ResponseWriter, Request -> For execute
func InterfaceHandler(w http.ResponseWriter, r *http.Request){
	// DB
	db.OpenConnection()
	// Check Cookies
	if(login.CheckCookie(w, r)){
		// Show Interface
		if(login.Message != ""){
			models.UserMessage = login.Message
		}
		ShowInterface(w, r)
	}else{
		// Redirect to Login
		models.UserMessage = login.Message
		http.Redirect(w, r, url, 301)
	}
	db.CloseConnection()
}

// Helper function to show an specific Nav/Menu in header template
// Specific for every site-handler
//  Return: Nav(type from controller) -> Navigation data (Name, Anchor, Icon)
func getInterfaceNav() Nav{
	elements := []NavElement{
		{"Overview","home", "home"},
		{"Devices","devices", "th"},
		{"Events","events", "calendar"},
		{"Types","types", "th-list"},
		{"User","user", "user"},
		{"About","about", "info"},
	}
	return Nav{elements}
}

// Function to handle action at user data got by /ui/user
// After that return to ui
//  Params: ResponseWriter, Request -> For execute
func UserHandler(w http.ResponseWriter, r *http.Request){
	// Get action from button
	set := r.FormValue("set")
	del := r.FormValue("del")
	if set != "" {
		if set == "avatar" {
			models.UploadAvatar(w, r)
		}else if set == "title" {
			models.SetTitle(r, r.FormValue("title"))
		}
	}else if del == "user" {
		models.DeleteUser(w, r)
	}
	http.Redirect(w, r, url + "ui#user", 301)
}

// Function to handle adding/update new items got by /ui/add
// After that return to ui
//  Params: ResponseWriter, Request -> For execute
func AddHandler(w http.ResponseWriter, r *http.Request){
	// Get action from button
	back := "home"
	device := r.FormValue("addDevice")
	group := r.FormValue("addType")
	event := r.FormValue("addEvent")
	// Handle Add or Update(Edit)
	if device == "Add" {
		back = "devices"
		models.AddData(r, back)
	} else if device == "Update" {
		back = "devices"
		models.UpdateData(r, back)
	}
	if group == "Add" {
		back = "types"
		models.AddData(r, back)
	} else if group == "Update" {
		back = "types"
		models.UpdateData(r, back)
	}
	if event == "Add" {
		back = "events"
		models.AddData(r, back)
	} else if event == "Update" {
		back = "events"
		models.UpdateData(r, back)
	}
	// Return to the ui with an anchor
	http.Redirect(w, r, url + "ui#"+back, 301)
}

// Function to handle update state from overview action
//  Params: ResponseWriter, Request -> For execute
func StateHandler(w http.ResponseWriter, r *http.Request) {
	// Get data
	state := r.FormValue("State")
	name := r.FormValue("Name")
	// Update
	models.UpdateState(name, state)
}

// Function to handle the get requests of data from item got by /ui/get
// Executes json data
//  Params: ResponseWriter, Request -> For execute
func GetHandler(w http.ResponseWriter, r *http.Request){
	var result []byte
	get := r.FormValue("Get")
	// Return devices
	if get == "AllDevicesByType" {
		result, _ = json.Marshal(models.GetAllDevicesByType(r.FormValue("Type")))
	}else if get == "AllDevices" {
		// Change type to art for overview
		data := models.GetAllDevices()
		for i := 0; i<len(data); i++ {
			data[i].Type = models.GetKindByType(data[i].Type)
		}
		result, _ = json.Marshal(data)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

// Function to handle deleting requests got by /ui/del
// After that return to ui
//  Params: ResponseWriter, Request -> For execute
func DelHandler(w http.ResponseWriter, r *http.Request){
	back := "home"
	item := r.FormValue("Item")
	name := r.FormValue("Name")
	if item == "D" {
		back = "devices"
		models.DelData(r, name, back)
	} else if item == "T" {
		back = "types"
		models.DelData(r, name, back)
	} else if item == "E" {
		back = "events"
		models.DelData(r, name, back)
	}
	http.Redirect(w, r, url + "ui#"+back, 301)
}

// Helper function for main to add basic types
func SetDefaultTypes(){
	if len(models.GetAllTypes()) < 1 {
		models.AddType("Light", "Switch", "0", "1")
		models.AddType("Roll Shutter", "Range", "0", "100")
		models.AddType("Dimmer", "Range", "0", "100")
		models.AddType("Heater", "Number", "0", "6")
		models.AddType("Coffee Maschine", "Switch", "0", "1")
	}
}

// Loading the configuration file for db
func DbConfig() (string, string){
	f, err := os.Open("./toHuus/conf/dbConf.json")
	if err != nil {
		fmt.Println("No file")
	}
	d := json.NewDecoder(f)
	conf := config{}
	err = d.Decode(&conf)
	if err != nil {
		fmt.Println("Bad file")
	}
	return conf.Db, conf.Url + ":" + conf.Port
}

// Loading the configuration file for url to client
func UrlConfig() string{
	f, err := os.Open("./toHuus/conf/urlConf.json")
	if err != nil {
		fmt.Println("No file")
	}
	d := json.NewDecoder(f)
	conf := config{}
	err = d.Decode(&conf)
	if err != nil {
		fmt.Println("Bad file")
	}
	return conf.Url
}

