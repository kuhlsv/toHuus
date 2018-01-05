// Class for struct handling
// This class just provide public struct for controllers and work with private helper structs
// With this help, other classes do not need imports for use of structs
// Package models is needed
package controllers

import (
	"toHuus/models"
)

// A Load represents a type needed to build templates with additional data like Nav(Struct)
type Load struct{
	Nav 		Nav
	Message 	string
	User 		models.UserData
}

// A Nav represents an array with Nav elements
type Nav struct{
	Elements	[]NavElement
}

// A NavElement represents a type with data for the nav
type NavElement struct{
	Name 	string
	Ref 	string
	Icon 	string
}
