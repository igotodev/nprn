package salemodel

type Sale struct {
	ID            string  `json:"id" bson:"_id,omitempty"`
	Article       string  `json:"article" bson:"article"`
	PriceForOne   float64 `json:"price_for_one" bson:"price_for_one"`
	NumberOfUnits int     `json:"number_of_units" bson:"number_of_units"`
	Amount        float64 `json:"amount" bson:"amount"`
	Date          string  `json:"date" bson:"date"`
	SellerID      string  `json:"seller_id" bson:"seller_id"`
}
