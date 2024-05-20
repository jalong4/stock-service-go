package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the database
type User struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    FirstName       string    `bson:"firstName" json:"firstName"`
    LastName        string    `bson:"lastName" json:"lastName"`
    Email           string    `bson:"email" json:"email"`
    Password        string    `bson:"password" json:"-"` // '-' in JSON tag to prevent it from being sent to the client
    Timezone        string    `bson:"timezone" json:"timezone"`
	ProfileImageURL string    `json:"profileImageUrl"`
    Date            time.Time `bson:"date" json:"date"`
}

// RegistrationRequest struct represents the registration request
type RegistrationRequest struct {
    FirstName       string `json:"firstName"`
    LastName        string `json:"lastName"`
    Email           string `json:"email"`
    Password        string `json:"password"`
    Password2       string `json:"password2"`
    Timezone        string `json:"timezone"`
    ProfileImageURL string `json:"profileImageUrl"`
}



// Holding represents the structure of a holding record in the database
type Holding struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    Ticker     string             `bson:"ticker" json:"ticker"`
    Quantity   float64            `bson:"quantity" json:"quantity"`
    TotalCost  float64            `bson:"totalCost" json:"totalCost"`
    Account    string             `bson:"account" json:"account"`
}