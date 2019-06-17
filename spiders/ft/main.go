package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"math/rand"
	"sort"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckF(e error) {
	if e != nil {
		Error.Fatal(e)
	}
}

func prepareReq(days int, symbol string) []byte {
	var dat map[string]interface{}
	reqBody, err := ioutil.ReadFile(*config)
	//file, err := os.Open("ft-iron.json")
	CheckF(err)
	err = json.Unmarshal(reqBody, &dat)
	CheckF(err)
	Debug.Printf("Before modif dat= %v", dat)
	dat["days"] = days
	elements := dat["elements"].([]interface{})
	for _, e := range elements {
		element := e.(map[string]interface{})
		element["Symbol"] = symbol
	}
	Debug.Printf("After modif dat= %v", dat)
	reqBody, err = json.Marshal(dat)
	CheckF(err)
	return reqBody
}

// parsing http replay into FT structure
func parseResp(resp []byte) (*ftData, error) {
	data := new(ftData)
	err := json.Unmarshal(resp, data)
	return data, err
}

// check FT structure
// mainly check does number of dates is == to number of values
func (d *ftData) checkData() error {
	datesNum := len((*d).Dates)
	symbol := (*d).Elements[0].Symbol
	//TODO: can be returned metadata
	elementNum := len((*d).Elements)
	if elementNum < 1 {
		return errors.New("Number of Elements < 1")
	}
	var openNum, closeNum, lowNum, highNum, volumNum int
	for e := 0; e < elementNum; e++ {
		for _, component := range (*d).Elements[e].ComponentSeries {
			switch component.Type {
			case "Open":
				openNum = len(component.Values)
			case "High":
				highNum = len(component.Values)
			case "Low":
				lowNum = len(component.Values)
			case "Close":
				closeNum = len(component.Values)
			case "Volume":
				volumNum = len(component.Values)
			}
		}
		switch {
		case (*d).Elements[e].ExchangeId != ftSymbols[symbol].exchangeID:
			e := fmt.Sprintln("recived ExchangeId=", (*d).Elements[e].ExchangeId, "according to ftSymbols should =",
				ftSymbols[symbol].exchangeID)
			Error.Println(e)
			return errors.New(e)
		case (*d).Elements[e].CompanyName != ftSymbols[symbol].description:
			e := fmt.Sprintln("recived CompanyName=", (*d).Elements[e].CompanyName, "according to ftSymbols should =",
				ftSymbols[symbol].description)
			Error.Println(e)
			return errors.New(e)
		case (*d).Elements[e].Currency != ftSymbols[symbol].currency:
			e := fmt.Sprintln("recived Currency=", (*d).Elements[e].Currency, "according to ftSymbols should =",
				ftSymbols[symbol].currency)
			Error.Println(e)
			return errors.New(e)

		}

	}
	if datesNum != openNum || datesNum != closeNum || datesNum != lowNum || datesNum != highNum ||
		datesNum != volumNum {
		err := fmt.Sprintf("Number of dates != valuse; datesNum=%v, closeNum=%v, openNum=%v, lowNum=%v, "+
			"highNum=%v, volumNum=%v", datesNum, closeNum, openNum, lowNum, highNum, volumNum)
		return errors.New(err)
	}
	return nil
}

// generate docs ready to insert in mongo from FT structure
func (d *ftData) genMongoDocs() (mongoDocs, error) {
	datesNum := len((*d).Dates)
	elementNum := len((*d).Elements)
	symbol := (*d).Elements[0].Symbol
	result := make([]mongoDoc, datesNum)
	var err error
	for i := 0; i < datesNum; i++ {
		// put values
		for e := 0; e < elementNum; e++ {
			for _, component := range (*d).Elements[e].ComponentSeries {
				switch component.Type {
				case "Open":
					result[i].Open = component.Values[i]
				case "High":
					result[i].High = component.Values[i]
				case "Low":
					result[i].Low = component.Values[i]
				case "Close":
					result[i].Close = component.Values[i]
				case "Volume":
					result[i].Volume = component.Values[i]
				}
			}
		}
		//	put dates
		result[i].Time, err = time.Parse("2006-01-02T15:04:05", (*d).Dates[i])
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in parsing Date= %s", (*d).Dates[i]))
			return mongoDocs{result}, err
		}
		// TODO: check recived info for symbol with data in ftSymbols
		result[i].Id = (*d).Dates[i] + "_" + ftSymbols[symbol].symbol + ":" + ftSymbols[symbol].exchangeID
		result[i].Source = ftSymbols[symbol].source
		result[i].Symbol = ftSymbols[symbol].symbol
		result[i].ExchangeID = ftSymbols[symbol].exchangeID
	}
	Info.Println("preapared =", datesNum, "docs for symbol=", ftSymbols[symbol].symbol, "from =",
		ftSymbols[symbol].source, "exchangeID=", ftSymbols[symbol].exchangeID)
	return mongoDocs{result}, nil
}

// inserting docs to DB
// order is strict in InsertMany
// return: (result of insert, nil if error as expected)
func (docs mongoDocs) InsertDocs(db, collect string) (string, error) {
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	CheckF(err)
	defer client.Disconnect(context.TODO())
	err = client.Ping(context.TODO(), nil)
	CheckF(err)
	Debug.Println("Connected to MongoDB!")
	collection := client.Database(db).Collection(collect)
	var insert []interface{}
	for _, doc := range docs.docs {
		insert = append(insert, doc)
	}
	t := true
	opts := &options.InsertManyOptions{Ordered: &t}
	Debug.Println("using db=", db, "collection=", collect)
	inserResult, err := collection.InsertMany(context.TODO(), insert, opts)
	if len((*inserResult).InsertedIDs) == len(docs.docs) {
		if err != nil {
			Warning.Println(err)
			return fmt.Sprint(err), nil
		}
		return fmt.Sprint(err), nil
	}
	return fmt.Sprint(err), err
}

// help methods for sorting
func (docs mongoDocs) Len() int           { return len(docs.docs) }
func (docs mongoDocs) Swap(i, j int)      { docs.docs[i], docs.docs[j] = docs.docs[j], docs.docs[i] }
func (docs mongoDocs) Less(i, j int) bool { return docs.docs[i].Time.Before(docs.docs[j].Time) }

// wrapper around prepareReq, checkData, genMongoDocs, InsertDocs
// colecting data from provider and putting them in db
// NOTE: connection is closed at the end
// TODO: reuse connection, for multiple work (not needed now)
func CollectData(symbol string) error {
	reqBody := prepareReq(*days, symbol)
	resp, err := httpClient.Post(*url, "application/json", bytes.NewBuffer(reqBody))
	defer resp.Body.Close()
	CheckF(err)
	respBody, err := ioutil.ReadAll(resp.Body)
	CheckF(err)
	ftData, err := parseResp(respBody)
	err = ftData.checkData()
	CheckF(err)
	mDocs, err := ftData.genMongoDocs()
	if *days == daysDefault {
		sort.Sort(sort.Reverse(mDocs))
	}
	_, err = mDocs.InsertDocs(*db, *collectionArg)
	CheckF(err)
	Info.Println("docs inserted for days=", *days, "in db=", *db, "collection=", *collectionArg)
	return nil
}

func main() {
	for key, _ := range ftSymbols {
		err := CollectData(key)
		CheckF(err)
		sleepTime := time.Duration(rand.Intn(10)) * time.Minute + time.Second * time.Duration(rand.Intn(60))
		Info.Println("sleep before next request =", sleepTime)
		time.Sleep(sleepTime)
	}

}
