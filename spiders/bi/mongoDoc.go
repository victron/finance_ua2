package main

import "time"

type mongoDoc struct {
	Id     string    `bson:"_id"` // "2019-04-11T00:00:00_iron-ore"
	Time   time.Time `bson:"time"`
	Source string    `bson:"source"` // "bi"
	Symbol string    `bson:"symbol"` // "iron-ore"
	Open   float64   `bson:"open"`
	High   float64   `bson:"high"`
	Low    float64   `bson:"low"`
	Close  float64   `bson:"close"`
	Volume float64   `bson:"volume"`
}

type mongoDocs struct {
	docs []mongoDoc
}
