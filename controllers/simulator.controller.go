// Class for handling simulator actions
// This class check und serialize actions for the simulator
// Package db, login, models, simulator is needed
package controllers

import (
	"net/http"
	"text/template"
	"encoding/json"
	"time"
	"strconv"
	"gopkg.in/mgo.v2/bson"
	"toHuus/login"
	"toHuus/models"
	"toHuus/db"
	"strings"
)

// Function to show the simulator interface
//  Params: ResponseWriter, Request -> For execute
func ShowSimulator(w http.ResponseWriter, r *http.Request){
	// Load templates
	t := template.Must(template.ParseFiles("./toHuus/views/header.html",
		"./toHuus/views/simulator.html", "./toHuus/views/footer.html"))
	// Executes templates with data (Nav, Message, Userdata)
	t.ExecuteTemplate(w, "header",
		Load{getSimulatorNav(), models.UserMessage, models.GetUserData(r)})
	t.ExecuteTemplate(w, "content", nil)
	t.ExecuteTemplate(w, "footer", nil)
}

// Basic function to handle the login check and carry out another action to simulator
//  Params: ResponseWriter, Request -> For execute
func SimulatorHandler(w http.ResponseWriter, r *http.Request){
	// DB
	db.OpenConnection()
	// Check Cookies
	if login.CheckCookie(w,r) {
		// Show simulator
		if login.Message != "" {
			models.UserMessage = login.Message
		}
		ShowSimulator(w,r)
	}else{
		// Redirect to the interface/login
		http.Redirect(w, r, url, 301)
	}
	db.CloseConnection()
}

// Function to handle data like import and export for simulator
//  Params: ResponseWriter, Request -> For execute
func DataHandler(w http.ResponseWriter, r *http.Request){
	ntype := r.FormValue("type")
	// Handle Add or Update(Edit)
	if ntype == "import" {
		models.Import(w, r)
		// Return to the sim with an anchor
		http.Redirect(w, r, url + "sim#data", 301)
	} else if ntype == "export" {
		models.Export(w, r)
	}
}

// Function to set data to db for the simulator
//  Params: ResponseWriter, Request -> For execute
func SimSetHandler(w http.ResponseWriter, r *http.Request){
	set := r.FormValue("Set")
	val := r.FormValue("Value")
	data := models.GetSimData()[0]
	if(set == "State"){
		newVal, err := strconv.ParseBool(val)
		if err == nil{
			data.State = newVal
		}
		models.SetSimData(data)
	}else if(set == "Time"){
		// Convert to duration
		hours , _ := strconv.Atoi(strings.Split(val, ":")[0])
		mins , _ := strconv.Atoi(strings.Split(val, ":")[1])
		newVal := (time.Duration(hours)*time.Hour) +
			(time.Duration(mins)*time.Minute)
		data.CurrentTime = newVal.Nanoseconds()
		models.SetSimData(data)
	}else if(set == "Multiplier"){
		newVal, err := strconv.Atoi(val)
		if err == nil{
			data.Multiplier = newVal
		}
		models.SetSimData(data)
	}
}

// Function to handle the get requests of data from simulator got by /ui/sim
// Executes json data
//  Params: ResponseWriter, Request -> For execute
func SimGetHandler(w http.ResponseWriter, r *http.Request){
	var result []byte
	get := r.FormValue("Get")
	// Return devices
	data := models.GetSimData()
	if len(data)>0 {
		if get == "States" {
			result, _ = json.Marshal(data[0])
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

// Helper function to show an specific Nav/Menu in header template
// Specific for every site-handler
//  Return: Nav(type from controller) -> Navigation data (Name, Anchor, Icon)
func getSimulatorNav() Nav{
	elements := []NavElement{
		{"Simulator","simUi", "home"},
		{"Data","data", "cogs"},
		{"About","about", "info"},
	}
	return Nav{elements}
}

// Helper function for main to add default simulator states
func SetSimStates(){
	// Initialise time
	time.LoadLocation("Europe/Berlin")
	// Set default states
	models.SetSimData(models.SimState{bson.NewObjectId(), time.Now().Unix(),
	"", "", false, 1})
}