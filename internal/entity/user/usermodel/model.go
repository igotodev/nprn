package usermodel

// UserInternal only internal use!!!
type UserInternal struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"password" bson:"password"`
	Email        string `json:"email" bson:"email"`
}

// UserTransfer for sharing
type UserTransfer struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}
