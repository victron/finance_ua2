package main

import "time"

type mongoDoc struct {
	Id         string    `bson:"_id"` // "2019-04-11T00:00:00_TIOc1:NYMEX"
	Time       time.Time `bson:"time"`
	Source     string    `bson:"source"`     // "ft"
	Symbol     string    `bson:"symbol"`     //TIOc1
	ExchangeID string    `bson:"exchangeID"` // NYMEX (NYM)
	Open       float64   `bson:"open"`
	High       float64   `bson:"high"`
	Low        float64   `bson:"low"`
	Close      float64   `bson:"close"`
	Volume     float64   `bson:"volume"`
}

type mongoDocs struct {
	docs []mongoDoc
}

type symbolInfo struct {
	symbol      string
	source      string
	exchangeID  string
	currency    string // "USc"
	description string //"Iron Ore 62% Fe, CFR China (TSI) Swa"
	//specific interface{}
}

type ftInfo struct {
	Label            string //"e1bb9889",
	CompanyName      string // "Iron Ore 62% Fe, CFR China (TSI) Swa",
	IssueType        string //"FU",
	Symbol           string // "27110161",
	UtcOffsetMinutes int    // -240,
	ExchangeId       string // "NYM",
	Currency         string // "USD",
}
