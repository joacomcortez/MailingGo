package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MailOffer representa los datos para el envío de correos con ofertas.
type MailOffer struct {
	Articles []EnrichedArticle `bson:"articles" json:"articles"`
	Users    []EnrichedUser    `bson:"users" json:"users"`
}

// User representa un usuario con su estado de suscripción.
type User struct {
	ID        string `bson:"_id" json:"_id"` // Usamos `string` en lugar de `primitive.ObjectID`
	Subscribed bool  `bson:"subscribed" json:"subscribed"` // Renombrado a `subscribed` para mayor claridad
}


// Article representa un artículo con su estado de habilitación.
type Article struct {
	ID        string `bson:"_id" json:"_id"` // Usamos `string` en lugar de `primitive.ObjectID`
	Offer bool  `bson:"offer" json:"offer"` 
}


// EnrichedUser representa un usuario enriquecido con los detalles adicionales como el nombre y el correo.
type EnrichedUser struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Subscribed bool `json:"subscribed"`
}

type EnrichedArticle struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}


// Structs para recibir detalles de los artículos
type ArticleDetails struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	Description  string        `bson:"description" json:"description" validate:"required"`
	Price        float32            `bson:"price" json:"price"`
	Stock        int                `bson:"stock" json:"stock"`
	Offer        bool               `bson:"offer" json:"offer"`
	WasEmailSent bool               `bson:"wasEmailSent" json:"wasEmailSent"`
	Created      time.Time          `bson:"created" json:"created"`
	Updated      time.Time          `bson:"updated" json:"updated"`
	Enabled      bool               `bson:"enabled" json:"enabled"`
}

// Description contiene los detalles del artículo como nombre y descripción.
type Description struct {
	Name        string `bson:"name" json:"name" validate:"required,min=1,max=100"`
	Description string `bson:"description" json:"description" validate:"required,min=1,max=256"`
	Image       string `bson:"image" json:"image" validate:"max=100"` // Este campo puede eliminarse si no es necesario
}

// UserDetails contiene detalles más extensos sobre el usuario, como nombre, login, email y suscripción.
type UserDetails struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" validate:"required,min=1,max=100"`
	Login       string             `bson:"login" validate:"required,min=5,max=100"`
	Password    string             `bson:"password" validate:"required"`
	Email       string             `bson:"email" validate:"required,email"`
	Permissions []string           `bson:"permissions"`
	Enabled     bool               `bson:"enabled"`
	Created     time.Time          `bson:"created"`
	Updated     time.Time          `bson:"updated"`
	Subscribed  bool               `bson:"subscribed"`
}
