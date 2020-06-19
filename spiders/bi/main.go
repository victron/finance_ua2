package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"time"
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

// Get page where TKdata present
func GetRootPage(url string) ([]byte, error){
	resp, err := httpClient.Get(url)
	if err != nil{
		return []byte{}, errors.New("error getting page: " + url)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	CheckF(err)
	return data, nil

}


// Receiving TKdata string
func ParseTKData(page []byte) (string, error){
	expr := regexp.MustCompile(`"TKData":"(?P<TKData>[0-9,]+)"`)
	match := expr.FindSubmatch(page)

	if len(match) < 2 {
		return "", errors.New("TKData not present")
	}
	return string(match[1]), nil

}


// Receiving data (based on TKdata)
func GetData(tkdata string, start, stop time.Time) (*[]byte, error) {
	rootUrl := "https://markets.businessinsider.com/Ajax/Chart_GetChartData"

	// prepare url
	u, err := url.Parse(rootUrl)
	CheckF(err)
	q := u.Query()
	q.Set("instrumentType", "Commodity")
	q.Set("tkData", tkdata)

	from := start.Format("20060102")
	to := stop.Format("20060102")
	q.Set("from", from)
	q.Set("to", to)

	u.RawQuery = q.Encode()

	resp, err := httpClient.Get(u.String())
	if err != nil{
		return &[]byte{}, errors.New("error getting page: " + u.String())
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	CheckF(err)
	return &data, nil

}


// parse data from json into GO structure
func ParseBIData(rawData *[]byte, symbol string) (*biDataPoints, error){
	var points []biDataPoint
	data := new(biDataPoints)
	(*data).symbol = symbol
	err := json.Unmarshal(*rawData, &points)
	if err != nil {
		return data, errors.New("error during parsing:" + symbol)
	}
	data.docs = &points
	return data, nil
}

//create mondo docs
func (points *biDataPoints) genMongoDocs() (*mongoDocs, error) {
	docs := new(mongoDocs)
	for _, point := range *(*points).docs {
		var doc mongoDoc
		doc.Id = point.Date.Format(time.RFC3339) + "_" + (*points).symbol
		doc.Time = point.Date.Time
		doc.Source = "BI"
		doc.Symbol = (*points).symbol
		doc.Open = point.Open
		doc.Close = point.Close
		doc.Low = point.Low
		doc.High = point.High
		doc.Volume = point.Volume

		(*docs).docs = append((*docs).docs, doc)
	}
	return docs, nil
}

// TODO: ------------------------ commons --------------------------

// help methods for sorting
func (docs mongoDocs) Len() int           { return len(docs.docs) }
func (docs mongoDocs) Swap(i, j int)      { docs.docs[i], docs.docs[j] = docs.docs[j], docs.docs[i] }
func (docs mongoDocs) Less(i, j int) bool { return docs.docs[i].Time.Before(docs.docs[j].Time) }

// inserting docs to DB
// order is strict in InsertMany
// return: (result of insert, nil if error as expected)
func (docs mongoDocs) InsertDocs(db, collect string, rSort bool) (string, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
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

	if rSort{
	sort.Sort(sort.Reverse(docs))
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



// wrapper around prepareReq, checkData, genMongoDocs, InsertDocs
// colecting data from provider and putting them in db
// NOTE: connection is closed at the end
// TODO: reuse connection, for multiple work (not needed now)
func CollectData(symbol string) error {
	page, err := GetRootPage(biSymbols[symbol].url)
	if err != nil{
		return err
	}

	tkData, err := ParseTKData(page)
	if err != nil {
		return err
	}

	// calculate start and stop dates
	stop := time.Now()
	start := stop.Add(- time.Duration(*days) * time.Hour * 24)

	rawData, err := GetData(tkData, start, stop)
	if err != nil {
		return err
	}

	points, err := ParseBIData(rawData, biSymbols[symbol].symbol)
	if err != nil {
		return err
	}

	mDocs, err := points.genMongoDocs()
	if err != nil{
		return err
	}

	// thinking that if days > 100 it's initial data collection - sorting not needed
	rSort := true
	if *days > 100 {
		rSort = false
	}

	_, err = mDocs.InsertDocs(*db, *collectionArg, rSort)
	CheckF(err)
	Info.Println("docs inserted for days=", *days, "in db=", *db, "collection=", *collectionArg)
	return nil
}

func main() {
	for symbol, _ := range biSymbols {
		Info.Println("collecting stat for map symbol:", symbol)
		err := CollectData(symbol)
		CheckF(err)
		sleepTime := time.Duration(rand.Intn(10)) * time.Minute + time.Second * time.Duration(rand.Intn(60))
		Info.Println("sleep before next request =", sleepTime)
		time.Sleep(sleepTime)
	}
	Info.Println("Normal exit...")
}