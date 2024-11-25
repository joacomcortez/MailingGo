package mailoffer

// import "go.mongodb.org/mongo-driver/bson/primitive"

// "go.mongodb.org/mongo-driver/bson/primitive"

//LAS ESTRUCTURAS ESTAN EN MAILINGGO -> MAILER -> MODELS.GOS

type Cart struct{
	UserID      string             `bson:"userId" json:"userId"`
}

type CartOpen struct {
	Message string  `bson:"message" json:"message"`
	Name string     `bson:"name" json:"name"`
	Email string  `bson:"email" json:"email"`
}

