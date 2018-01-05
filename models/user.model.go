// Class for user manipulation
// This class processing data of users (Set data/Delete/Get)
// Package db, login is needed
package models

import (
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"toHuus/db"
	"toHuus/login"
)

// Initialisation
const CookieName = "session"
const DbCollUserdata = "Userdata"
var UserMessage = ""

// Function to get all data from User by Session
//  Params: Session(String) -> Get data from specific user
//  Return: UserData(type from model) -> Struct with data of the user
func GetUserBySession(session string) UserData {
	result := UserData{}
	if session != "" {
		database := db.OpenConnection()
		coll := database.C(DbCollUserdata)
		coll.Find(bson.M{"Session": session}).One(&result)
		db.CloseConnection()
	}
	return result
}

// Function to get all data from User by username
//  Params: Username(String) -> Get data from specific user
//  Return: UserData(type from model) -> Struct with data of the user
func GetUserByUname(uname string) UserData {
	result := UserData{}
	if uname != "" {
		database := db.OpenConnection()
		coll := database.C(DbCollUserdata)
		coll.Find(bson.M{"Username": uname}).One(&result)
		db.CloseConnection()
	}
	return result
}

// Function to a title(like admin/guest) to a user
//  Params: Request, Title(String) -> Set title to user by session
func SetTitle(r *http.Request, title string){
	user := GetUserData(r)
	database := db.OpenConnection()
	coll := database.C(DbCollUserdata)
	coll.Update(user, bson.M{"$set": bson.M{ "Title" : title }})
	db.CloseConnection()
}

// Function to a avatar image to a user
//  Params: Request, Path(String) -> Set avatar to user by session into db
func SetAvatar(r *http.Request, path string){
	user := GetUserData(r)
	database := db.OpenConnection()
	coll := database.C(DbCollUserdata)
	coll.Update(user, bson.M{"$set": bson.M{ "Avatar" : path }})
	db.CloseConnection()
}

// Function to delete the user
// Deleting Avatar, Events/Relations, Session, User
//  Params: ResponseWriter, Request -> For execute
func DeleteUser(w http.ResponseWriter, r *http.Request){
	database := db.OpenConnection()
	coll := database.C(DbCollUserdata)
	user := GetUserData(r)
	coll.Remove(bson.M{ "Session" : user.SessionId })
	// Logout
	db.CloseConnection()
	login.DeleteCookie(w)
	// Delete avatar and events/relations
	DeleteAvatar(user.Avatar)
	DelEventsById(user.Id)
}

// Function to get all data from current User (Session)
//  Params: Request -> Get data from user by session
//  Return: UserData(type from model) -> Struct with data of the user
func GetUserData(r *http.Request) UserData{
	cookie, _ := r.Cookie(CookieName)
	data := UserData{}
	if cookie != nil {
		data = GetUserBySession(cookie.Value)
	}
	return data
}

// Function to get all data from User
//  Return: []UserData(type from model) -> Struct arrays with data of the users
func GetAllUserData() []UserData{
	data := []UserData{}
	database := db.OpenConnection()
	coll := database.C(DbCollUserdata)
	coll.Find(nil).All(&data)
	db.CloseConnection()
	return data
}



