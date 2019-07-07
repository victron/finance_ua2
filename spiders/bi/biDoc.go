package main

import (
	"strings"
	"time"
)

// ---------- custom time in JSON -----------
//https://www.evergreeninnovations.co/blog-working-with-package-time-in-go/
const biLayout string = "2006-01-02 15:04"

type biTime struct {
	time.Time
}

func (bt *biTime) UnmarshalJSON(b []byte) error {
	// Convert to string and remove quotes
	s := strings.Trim(string(b), "\"")

	// Parse the time using the biLayout
	t, err := time.Parse(biLayout, s)
	if err != nil {
		return err
	}

	// Assign the parsed time to our variable
	(*bt).Time = t
	return nil
}

//=======================================================

type biDataPoint struct {
	Close  float64 `json:"Close"`
	Open   float64 `json:"Open"`
	High   float64 `json:"High"`
	Low    float64 `json:"Low"`
	Volume float64 `json:"Volume"`
	Date   biTime  `json:"Date"`
}

type biDataPoints struct {
	docs   *[]biDataPoint
	symbol string
}

type symbolInfo struct {
	symbol      string // "iron-ore"
	source      string // "BI"
	exchangeID  string
	currency    string // "USc" or "USD"
	description string //"Iron Ore 62% Fe, CFR China (TSI) Swa"
	url         string
	//specific interface{}
}

var biSymbols = map[string]symbolInfo{
	"iron-ore": {"iron-ore", "BI", "", "USD", "Iron Ore",
		"https://markets.businessinsider.com/commodities/iron-ore-price"},
	"gold": {"gold", "BI", "", "USD", "Gold",
		"https://markets.businessinsider.com/commodities/gold-price"},
	"palladium": {"palladium", "BI", "", "USD", "Palladium",
		"https://markets.businessinsider.com/commodities/palladium-price"},
	"rhodium": {"rhodium", "BI", "", "USD", "Rhodium",
		"https://markets.businessinsider.com/commodities/rhodium-price"},
	"silver": {"silver", "BI", "", "USD", "Silver",
		"https://markets.businessinsider.com/commodities/silver-price"},
	"natural-gas": {"natural-gas", "BI", "", "USD", "natural gas",
		"https://markets.businessinsider.com/commodities/natural-gas-price"},
	"ethanol": {"ethanol", "BI", "", "USD", "Ethanol",
		"https://markets.businessinsider.com/commodities/ethanol-price"},
	"heating-oil": {"heating-oil", "BI", "", "USD", "heating oil",
		"https://markets.businessinsider.com/commodities/heating-oil-price"},
	"coal": {"coal", "BI", "", "USD", "Coal",
		"https://markets.businessinsider.com/commodities/coal-price"},
	"uranium": {"uranium", "BI", "", "USD", "uranium",
		"https://markets.businessinsider.com/commodities/uranium-price"},
	"oil-Brent": {"oil-Brent", "BI", "", "USD", "oil Brent",
		"https://markets.businessinsider.com/commodities/oil-price?type=Brent"},
	"oil-WTI": {"oil-WTI", "BI", "", "USD", "oil WTI",
		"https://markets.businessinsider.com/commodities/oil-price?type=WTI"},
	"aluminum": {"aluminum", "BI", "", "USD", "aluminum",
		"https://markets.businessinsider.com/commodities/aluminum-price"},
	"lead": {"lead", "BI", "", "USD", "lead",
		"https://markets.businessinsider.com/commodities/lead-price"},
	"copper": {"copper", "BI", "", "USD", "copper",
		"https://markets.businessinsider.com/commodities/copper-price"},
	"zinc": {"zinc", "BI", "", "USD", "zinc",
		"https://markets.businessinsider.com/commodities/zinc-price"},
	"tin": {"tin", "BI", "", "USD", "tin",
		"https://markets.businessinsider.com/commodities/tin-price"},
	"cotton": {"cotton", "BI", "", "USD", "cotton",
		"https://markets.businessinsider.com/commodities/cotton-price"},
	"oats": {"oats", "BI", "", "USD", "oats",
		"https://markets.businessinsider.com/commodities/oats-price"},
	"lumber": {"lumber", "BI", "", "USD", "lumber",
		"https://markets.businessinsider.com/commodities/lumber-price"},
	"coffee": {"coffee", "BI", "", "USD", "coffee",
		"https://markets.businessinsider.com/commodities/coffee-price"},
	"cocoa": {"cocoa", "BI", "", "USD", "cocoa",
		"https://markets.businessinsider.com/commodities/cocoa-price"},
	"live-cattle": {"live-cattle", "BI", "", "USD", "live cattle",
		"https://markets.businessinsider.com/commodities/live-cattle-price"},
	"lean-hog": {"lean-hog", "BI", "", "USD", "lean hog",
		"https://markets.businessinsider.com/commodities/lean-hog-price"},
	"corn": {"corn", "BI", "", "USD", "corn",
		"https://markets.businessinsider.com/commodities/corn-price"},
	"feeder-cattle": {"feeder-cattle", "BI", "", "USD", "feeder cattle",
		"https://markets.businessinsider.com/commodities/feeder-cattle-price"},
	"milk": {"milk", "BI", "", "USD", "milk",
		"https://markets.businessinsider.com/commodities/milk-price"},
	"orange-juice": {"orange-juice", "BI", "", "USD", "orange juice",
		"https://markets.businessinsider.com/commodities/orange-juice-price"},
	"palm-oil": {"palm-oil", "BI", "", "USD", "palm oil",
		"https://markets.businessinsider.com/commodities/palm-oil-price"},
	"rapeseed": {"rapeseed", "BI", "", "USD", "rapeseed",
		"https://markets.businessinsider.com/commodities/rapeseed-price"},
	"rice": {"rice", "BI", "", "USD", "rice",
		"https://markets.businessinsider.com/commodities/rice-price"},
	"soybean-meal": {"soybean-meal", "BI", "", "USD", "soybean meal",
		"https://markets.businessinsider.com/commodities/soybean-meal-price"},
	"soybeans": {"soybeans", "BI", "", "USD", "soybeans",
		"https://markets.businessinsider.com/commodities/soybeans-price"},
	"soybean-oil": {"soybean-oil", "BI", "", "USD", "soybean oil",
		"https://markets.businessinsider.com/commodities/soybean-oil-price"},
	"wheat": {"wheat", "BI", "", "USD", "wheat",
		"https://markets.businessinsider.com/commodities/wheat-price"},
	"sugar": {"sugar", "BI", "", "USD", "sugar",
		"https://markets.businessinsider.com/commodities/sugar-price"},
}

