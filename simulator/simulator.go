// Simulator of toHuus
// This is called by the main application and do not need to get run
package simulator

import (
	"time"
	"toHuus/models"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"toHuus/db"
	"math/rand"
	"strconv"
	"strings"
	"fmt"
)

// Initialisation
var run = true
var changed = true
var tickerControl *time.Ticker
var tickerSimulation *time.Ticker
var simStep = time.Duration(1)
var realStep = time.Duration(1)
var realIntervall time.Duration = time.Second * realStep
var realZeit time.Duration = 0
var simIntervall time.Duration = time.Second * simStep // Multiplier
var simZeit time.Duration = 0
// DB
const dbCollRelEvents = "RelationEventsDevices"
const dbCollSim = "Simulator"
const dbCollEvents = "Events"

// Basic function for this application to start the Simulator
// This don't need to be run if other application do it
func Start() {
	tickerControl = time.NewTicker(time.Millisecond * time.Duration(1000 / simStep))
	tickerSimulation = time.NewTicker(time.Millisecond * time.Duration(1000 / simStep))
	simulatorControl()
}

// Function to get state from simulator
//  Return: SimState(type from model) -> State data
func readState() models.SimState{
	data := models.GetSimData()
	newData := models.SimState{}
	if len(data) > 0 {
		newData = data[0]
	}
	return newData
}

// Function to set state from simulator
//  Param: SimState(type from model) -> State data
func setState(data models.SimState) {
	models.SetSimData(data)
}

func simulatorControl()  {
	data := models.GetSimData()[0]
	for range tickerControl.C {
		run = data.State
		realStep = time.Duration(data.Multiplier)
		simZeit = time.Duration(data.CurrentTime)
		startSimulation(*db.OpenConnection())
	}
}

// Function convert time
func durationFromTime(val string) time.Duration{
	data := strings.Split(val, ":")
	hours , _ := strconv.Atoi(data[0])
	mins , _ := strconv.Atoi(data[1])
	return (time.Duration(hours)*time.Hour) + (time.Duration(mins)*time.Minute)
}

// Basic function to handle event data and practice simulation
//  Param: SimState(type from model) -> State data
func startSimulation(database mgo.Database){
	collE := database.C(dbCollEvents)
	collS := database.C(dbCollSim)
	// "Simulatorschleife":
	if run {
		// Update offset
		events := []models.Event{}
		collE.Find(nil).All(&events)
		for _, element := range events {
			oldTime,_ := strconv.Atoi(element.Time)
			offset,_ := strconv.Atoi(element.Offset)
			newTime := oldTime + (rand.Intn(offset))
			collE.Update(bson.M{"_id": element.Id}, bson.M{"$set": bson.M{"Time": newTime}})
		}
		gtEvents := []models.Event{}
		collE.Find(bson.M{"Time": bson.M{"$gte": simZeit}}).All(&gtEvents)
		for _, element := range gtEvents {
			if durationFromTime(element.Time) < simZeit+realStep {
				devices := models.GetRelationToEvent(element.Id.Hex())
				for i := range devices {
					if devices[i].Id != "" {
						devices[i].State = models.GetNewState(element.Id.Hex(), devices[i].Id.Hex())
					}
				}
			}
		}
		simZeit = simZeit + realStep
		sim := models.SimState{}
		collS.Find(nil).One(&sim)
		collS.Update(bson.M{"_id": sim.Id}, bson.M{"$set": bson.M{"time": simZeit}})
		fmt.Print(simZeit)
		changed = false
		/* //new ad
		time.Sleep(realIntervall)
		realZeit = realZeit + realIntervall
		simZeit = simZeit + simIntervall
		simulatorControl()

		// Handle data
		data := readState()
		simStep = time.Duration(data.Multiplier) //?
		run = data.State
		simZeit = time.Duration(data.CurrentTime)
		setState(data)
		if run {
			doSim(data)
		}*/
	}
}