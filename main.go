// Main for toHuus
// Call this to start the application
// Just this need to be called, the simulation will start by it self
// Package db, controllers, simulator is needed
package main

import (
	"fmt"
	"net/http"
	"toHuus/controllers"
	"toHuus/db"
	"toHuus/simulator"
)

// Main function for this application
// This need to be run
func main() {
	// Initialise DB from config
	db.Database(controllers.DbConfig())
	// This delete the db to start clean
	//db.OpenConnection().DropDatabase()
	// Preset types
	controllers.SetDefaultTypes()
	// Preset simulator states
	controllers.SetSimStates()
	// Start simulator as thread
	go simulator.Start() // Interoperability
	// Handler
	http.Handle("/images/", http.FileServer(http.Dir("./toHuus/views/")))
	http.Handle("/assets/", http.FileServer(http.Dir("./toHuus/views/")))
	http.Handle("/avatar/", http.FileServer(http.Dir("./toHuus/conf/")))
	http.HandleFunc("/", controllers.CheckLogin)
	http.HandleFunc("/ui", controllers.InterfaceHandler)
	http.HandleFunc("/ui/user", controllers.UserHandler)
	http.HandleFunc("/ui/add", controllers.AddHandler)
	http.HandleFunc("/ui/get", controllers.GetHandler)
	http.HandleFunc("/ui/del", controllers.DelHandler)
	http.HandleFunc("/ui/set", controllers.StateHandler)
	http.HandleFunc("/sim", controllers.SimulatorHandler)
	http.HandleFunc("/sim/data", controllers.DataHandler)
	http.HandleFunc("/sim/get", controllers.SimGetHandler)
	http.HandleFunc("/sim/set", controllers.SimSetHandler)
	err := http.ListenAndServe(":4242", nil)
	if err != nil {
		fmt.Println(err)
	}
}