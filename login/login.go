// Package for a login
// Call the CheckLogin function
// The construction of HTML forms and inputs is important for this (variables)
// Todo: Enable security again (hash)
package login

import (
	"net/http"
	"fmt"
	"sync"
	"strconv"
	"encoding/hex"
	"crypto/rand"
	"encoding/base64"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"toHuus/db"
)

// Initialisation
const dbCollName = "Userdata"
const cookieName = "session"
const buttonName = "authBtn"
const usernameName = "uname"
const passwordName = "passwd"
var formButtons = []string  { "Login", "Registration", "Logout" }
var storageMutex sync.RWMutex
var Message string  // Message returned by calling functions

// An User represents a user with basic data
type User struct {
	Id 	 bson.ObjectId 	`bson:"_id"`
	Username string 	`bson:"Username"`
	Title string 		`bson:"Title"`
	Password string 	`bson:"Password"`
	SessionId string 	`bson:"Session"`
	Avatar string 		`bson:"Avatar"`
}

// Basic function for this package
//  Params: ResponseWriter, Request -> For cookie handling
//  Return: Boolean -> Action was valid
func CheckLogin(w http.ResponseWriter, r *http.Request) bool{
	var valid bool
	database := db.GetConnection()
	coll := database.C(dbCollName)
	button := r.PostFormValue(buttonName) // r.Form[]
	if len(button) <= 2 {
		// Check for cookies to handle session, because no action found
		valid = CheckCookie(w, r)
	}else{
		uname := r.PostFormValue(usernameName)
		password := r.PostFormValue(passwordName)
		switch button {
		case formButtons[0]:
			// Call login
			valid = login(w, r, uname, password, coll)
		case formButtons[1]:
			// Call registration
			valid = register(uname, password, coll)
		case formButtons[2]:
			// Call logout
			logout(w, uname, coll)
			valid = false
		default:
			valid = false
			Message = "Error: Unknown"
		}
	}
	return valid
}

// Basic function for sessions
//  Params: ResponseWriter, Request -> For cookie handling
//  Return: Boolean -> Cookie is valid
func CheckCookie(w http.ResponseWriter, r *http.Request) bool{
	var valid bool
	cookie, err := r.Cookie(cookieName)
	database := db.GetConnection()
	coll := database.C(dbCollName)
	if err != nil {
		if err != http.ErrNoCookie {
			fmt.Fprint(w, err)
			valid = false
			Message = "Error: No Session"
		} else {
			err = nil
		}
	}
	if cookie != nil {
		// Check session is valid
		result := User{}
		storageMutex.RLock()
		coll.Find(bson.M{ "Session" : cookie.Value }).One(&result)
		storageMutex.RUnlock()
		if result.SessionId != "" {
			valid = true
		}
	}else{
		valid = false
	}
	return valid
}

// Function to handle registration
//  Params: Username(String), Password(String), Collection -> Get username and password
//  Return: Boolean -> Action was valid
func register(uname string, password string, coll *mgo.Collection) bool{
	var valid bool
	// Check if the input of username and password is valid
	if validatePassword(password) && validateUsername(uname) {
		// Insert to database
		coll.Insert(User{
			bson.NewObjectId(),
			uname,
			"",
			password, // hash(password) | disabled security
			"",
			"",
		})
		Message = "Rigistered"
		valid = false
	}else{
		valid = false
		if Message == "" {
			Message = "Successfully registered"
		}
	}
	return valid
}

// Function to handle login
//  Params: ResponseWriter, Request, Username(String), Password(String), Collection -> Get cookie, username and password
//  Return: Boolean -> Action was valid
func login(w http.ResponseWriter, r *http.Request, uname string, password string, coll *mgo.Collection) bool{
	var valid bool
	result := User{}
	coll.Find(bson.M{ "Username" : uname }).One(&result)
	// Check if username and password from input is valid with db
	if result.Username == uname && result.Password == password { // hash(password) | disabled security
		// Create session
		sessionId := generateSessionId()
		coll.Update(result, bson.M{"$set": bson.M{ "Session" : sessionId }})
		DeleteCookie(w)
		setCookie(w, r, sessionId, coll)
		valid = true
		Message = "Successfully logged in"
	}else{
		valid = false
		if len(Message) <= 2 {
			Message = "Error: Invalid username or password"
		}
	}
	return valid
}

// Function to handle logout
//  Params: ResponseWriter, Username(String), Collection -> Get cookie and username
func logout(w http.ResponseWriter, uname string, coll *mgo.Collection){
	// Delete the session and cookie
	result := User{}
	coll.Find(bson.M{ "Username" : uname }).One(&result)
	coll.Update(result, bson.M{ "$set": bson.M{ "Session" : "" }})
	DeleteCookie(w)
	Message = "Logged out"
}

// Helper-Function to delete an cookie of w
//  Params: ResponseWriter -> Get cookie
func DeleteCookie(w http.ResponseWriter)  {
	newCookie := http.Cookie{
		Name: cookieName,
		MaxAge: -1,
	}
	http.SetCookie(w, &newCookie)
}

// Function to set a new cookie to user
//  ResponseWriter, Request, session id(String), Collection -> Get cookie and session
func setCookie(w http.ResponseWriter, r *http.Request, sessionId string, coll *mgo.Collection) {
	// Check for present cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err != http.ErrNoCookie {
			fmt.Fprint(w, err)
			return
		} else {
			err = nil
		}
	}
	// Generate a new session
	if sessionId == "" {
		sessionId = generateSessionId()
	}
	// Set cookie to user and db
	cookie = &http.Cookie{
		Name: cookieName,
		Value: sessionId,
	}
	result := User{}
	storageMutex.Lock()
	coll.Find(bson.M{ "Session" : cookie.Value }).One(&result)
	coll.Update(result, bson.M{"$set": bson.M{ "Session" : sessionId }})
	storageMutex.Unlock()
	http.SetCookie(w, cookie)
}

// Helper-Function to generate an session
//  Return: String -> Session id
func generateSessionId() string{
	buffer := make([]byte, 32)
	// Random byte
	_, err := rand.Read(buffer)
	if err != nil {
		panic(err)
	}
	// Encode byte
	return hex.EncodeToString(buffer)
}

// Helper-Function to check password by rules
//  Param: String -> Password
//  Return: Boolean -> valid
func validatePassword(password string) bool{
	pLenght := 3
	var valid bool
	// Length
	if len(password) >= pLenght {
		valid = true
	}else{
		valid = false
		Message = "Error: Invalid password length (min. " + strconv.Itoa(pLenght) + ")"
	}
	return valid
}

// Helper-Function to check username by rules
//  Param: String -> Username
//  Return: Boolean -> valid
func validateUsername(uname string) bool{
	uLenght := 3
	var valid bool
	// Length
	if len(uname) >= uLenght {
		// Check username forgiven
		database := db.GetConnection()
		coll := database.C(dbCollName)
		result := User{}
		coll.Find(bson.M{ "Username" : uname }).One(&result)
		if result.Username != uname {
			valid = true
		}else{
			valid = false
			Message = "Error: User already exist"
		}
	}else{
		valid = false
		Message = "Error: Invalid username length (min. " + strconv.Itoa(uLenght) + ")"
	}
	return valid
}

// Helper-Function to hash an password for security
//  Param: String -> Clean password
//  Return: String -> Hashed string
func hash(data string) string{
	return base64.StdEncoding.EncodeToString([]byte(data))
}


