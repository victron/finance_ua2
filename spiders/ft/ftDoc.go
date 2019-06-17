package main

type ftData struct {
	Dates    []string    `json:"Dates"`
	Elements []ftElement `json:"Elements"`
}

type ftElement struct {
	Type            string        `json:"Type"`
	CompanyName     string        `json:"CompanyName"`
	Symbol          string        `json:"Symbol"`
	ExchangeId      string        `json:"ExchangeId"`
	Currency        string        `json:"Currency"`
	ComponentSeries []ftComponent `json:"ComponentSeries"`
}

type ftComponent struct {
	Type         string    `json:"Type"`
	MaxValue     float64   `json:"MaxValue"`
	MinValue     float64   `json:"MinValue"`
	MaxValueDate string    `json:"MaxValueDate"`
	MinValueDate string    `json:"MinValueDate"`
	Values       []float64 `json:"Values"`
}

var ftSymbols = map[string]symbolInfo{
	"27110161": {"TIOc1", "FT", "NYM", "USD", "Iron Ore 62% Fe, CFR China (TSI) Swa"},
	"1037373":  {"Wc1", "FT", "CBT", "USc", "Wheat Front Month Futures"},
	"1061323": {"QW.1","FT","IEU", "USD", "SUGAR NO5 DEC6"},
	"1046731": {"SB.1", "FT", "IUS", "USc", "No.11 Sugar Front Month Futures"},
	"1038557": {"1SMc1", "FT", "CBT", "USD", "c1 Soybean Meal (Electronic trades)"},
	"1037017": {"Sc1", "FT", "CBT", "USc", "Soybeans Front Month Futures"},
	"1044733": {"LB.1", "FT","CME", "USD", "Random Length Lumber Front Month Futures"},
	"1044843":{"LH.1", "FT","CME", "USc", "Lean Hogs Front Month Futures"},
	"1045268":{"FC.1", "FT","CME", "USc", "Feeder Cattle Front Month Futures"},
	"1039187": {"C.1", "FT", "CBT","USc", "Corn Front Month Futures"},
	"1044934": {"LC.1", "FT","CME", "USc","Live Cattle Front Month Futures"},
	"1046328": {"US@HG.1", "FT","CMX", "USD", "Copper High Grade Front Month Futures"},
	"1070572": {"CL.1", "FT","NYM", "USD", "Crude Oil Front Month Futures"},
	"3334534": {"US@RB.1", "FT","NYM", "USD", "New York Harbor RBOB Gasoline Front Month Futures"},
	"1069936": {"US@NG.1", "FT", "NYM", "USD", "Henry Hub Natural Gas Front Month Futures"},
	"1072386": {"US@HO.1", "FT", "NYM", "USD", "Heating Oil Front Month Futures"},
	"1518990": {"1ZEC1", "FT", "CBT", "USD", "CBT ETHANOL"},
	"1054972": {"IB.1", "FT","IEU", "USD", "ICE Brent Crude Energy Future c1"},

}
