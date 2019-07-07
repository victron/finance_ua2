package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"testing"
	"time"
)

type testParseTKData struct {
	file string
	out  string
	outE bool
}

var testsParseTKData = []testParseTKData{
	{file: "test_data/soybean-oil-price.html", out: "300002,24,0,333", outE: false},
	{file: "test_data/soybean-oil-price_no_TKData.html", out: "", outE: true},
	{file: "test_data/soybean-oil-price_empty_TKData.html", out: "", outE: true},
	{file: "test_data/soybean-oil-price_wrong_TKData.html", out: "", outE: true},
}

// help function to check that err != nil
func ErrTrue(err error) bool {
	if err != nil {
		return true
	}
	return false
}

type testData struct {
	file     string
	biPoints biDataPoints
	err      bool
	mDocs    mongoDocs
}

var testsData = []testData{
	{file: "test_data/biData.json",
		biPoints: biDataPoints{docs: &[]biDataPoint{
			{Close: 0.16, Open: 0, High: 0, Low: 0, Volume: 0,
				Date: biTime{time.Date(2017, 4, 25, 0, 0, 0, 0, time.UTC)}},
			{Close: 0.15, Open: 0, High: 0, Low: 0, Volume: 0,
				Date: biTime{time.Date(2017, 4, 26, 0, 0, 0, 0, time.UTC)}},
			{Close: 0.15, Open: 0, High: 0, Low: 0, Volume: 0,
				Date: biTime{time.Date(2017, 4, 27, 0, 0, 0, 0, time.UTC)}},
			{Close: 0.15, Open: 0, High: 0, Low: 0, Volume: 0,
				Date: biTime{time.Date(2017, 4, 28, 0, 0, 0, 0, time.UTC)}},
			{Close: 0.16, Open: 0, High: 0, Low: 0, Volume: 0,
				Date: biTime{time.Date(2017, 5, 1, 0, 0, 0, 0, time.UTC)}},
		},
			symbol: "iron-ore",},
		err: false,
		mDocs: mongoDocs{docs: []mongoDoc{
			{Close: 0.16, Open: 0, High: 0, Low: 0, Volume: 0, Source: "BI", Symbol: "iron-ore",
				Time: time.Date(2017, 4, 25, 0, 0, 0, 0, time.UTC),
				Id:   "2017-04-25T00:00:00Z_iron-ore"},
			{Close: 0.15, Open: 0, High: 0, Low: 0, Volume: 0, Source: "BI", Symbol: "iron-ore",
				Time: time.Date(2017, 4, 26, 0, 0, 0, 0, time.UTC),
				Id:   "2017-04-26T00:00:00Z_iron-ore"},
			{Close: 0.15, Open: 0, High: 0, Low: 0, Volume: 0, Source: "BI", Symbol: "iron-ore",
				Time: time.Date(2017, 4, 27, 0, 0, 0, 0, time.UTC),
				Id:   "2017-04-27T00:00:00Z_iron-ore"},
			{Close: 0.15, Open: 0, High: 0, Low: 0, Volume: 0, Source: "BI", Symbol: "iron-ore",
				Time: time.Date(2017, 4, 28, 0, 0, 0, 0, time.UTC),
				Id:   "2017-04-28T00:00:00Z_iron-ore"},
			{Close: 0.16, Open: 0, High: 0, Low: 0, Volume: 0, Source: "BI", Symbol: "iron-ore",
				Time: time.Date(2017, 5, 1, 0, 0, 0, 0, time.UTC),
				Id:   "2017-05-01T00:00:00Z_iron-ore"},
		}}},
}

func TestParseTKData(t *testing.T) {
	for n, testPair := range testsParseTKData {
		input, err := ioutil.ReadFile(testPair.file)
		check(err)
		result, err := ParseTKData(input)

		if result != testPair.out || ErrTrue(err) != testPair.outE {
			t.Fatal("wrong output in pair:", n,
				"\n exp:", testPair.out, "err:", testPair.outE,
				"\n got:", result, "err:", err)
		}
	}
}

func TestParseBIData(t *testing.T) {
	for n, testPair := range testsData {
		input, err := ioutil.ReadFile(testPair.file)
		check(err)
		result, err := ParseBIData(&input, "iron-ore")

		if ErrTrue(err) != testPair.err {
			t.Fatal("wrong error code in pair:", n,
				"\n exp:", testPair.err,
				"\n got:", err)
		}

		if len(*(*result).docs) != len(*testPair.biPoints.docs) {
			t.Fatal("wrong result size",
				"\n exp:", len(*testPair.biPoints.docs),
				"\n got:", len(*(*result).docs))
		}

		for docNum, doc := range *testPair.biPoints.docs {
			if !doc.Date.Equal((*(*result).docs)[docNum].Date.Time) || doc.Close != (*(*result).docs)[docNum].Close ||
				doc.Open != (*(*result).docs)[docNum].Open || doc.High != (*(*result).docs)[docNum].High ||
				doc.Low != (*(*result).docs)[docNum].Low {
				t.Fatal("wrong result in docNum:", docNum,
					"\n exp:", doc,
					"\n got:", (*(*result).docs)[docNum])
			}
		}
	}
}

// test mongo docs before inserting to DB
// this test success only if above test OK
func TestGenMongoDocs(t *testing.T) {
	for n, testPair := range testsData {
		input, err := ioutil.ReadFile(testPair.file)
		check(err)

		points, err := ParseBIData(&input, "iron-ore")
		check(err)
		result, err := points.genMongoDocs()
		check(err)

		if len((*result).docs) != len(testPair.mDocs.docs) {
			t.Fatal("wrong result size in testPair:", n,
				"\n exp:", len(testPair.mDocs.docs),
				"\n got:", len((*result).docs))
		}

		for docNum, doc := range (*result).docs {
			if doc.Id != testPair.mDocs.docs[docNum].Id || doc.Close != testPair.mDocs.docs[docNum].Close ||
				doc.Open != testPair.mDocs.docs[docNum].Open || doc.High != testPair.mDocs.docs[docNum].High ||
				doc.Low != testPair.mDocs.docs[docNum].Low || doc.Volume != testPair.mDocs.docs[docNum].Volume ||
				doc.Source != testPair.mDocs.docs[docNum].Source || doc.Symbol != testPair.mDocs.docs[docNum].Symbol ||
				!doc.Time.Equal(testPair.mDocs.docs[docNum].Time) {

				t.Fatal("wrong result in testPair:", n, "docNum:", docNum,
					"\n exp:", testPair.mDocs.docs[docNum],
					"\n got:", doc)
			}
		}
	}
}

// help function to Drop collection, before test
func DropCollection(db, collect string) error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	CheckF(err)
	defer client.Disconnect(context.TODO())
	err = client.Ping(context.TODO(), nil)
	CheckF(err)
	Debug.Println("Connected to MongoDB!")
	collection := client.Database(db).Collection(collect)
	err = (*collection).Drop(context.TODO())
	return err
}

// connect to DB should be established
func TestInsertDocs(t *testing.T) {
	for n, testPair := range testsData {
		input, err := ioutil.ReadFile(testPair.file)
		check(err)

		points, err := ParseBIData(&input, "iron-ore")
		check(err)
		mDocs, err := points.genMongoDocs()
		check(err)

		testDB := "test"
		testCollection := "commoditiesT"
		err = DropCollection(testDB, testCollection)
		CheckF(err)

		for i := 0; i < 2; i++ {
			// first insert
			// insert docs without duplicate
			if i == 0 {
				result, err := mDocs.InsertDocs(testDB, testCollection, true)
				if result != "<nil>" || err != nil {
					t.Fatal("Error not nill after insert, i=", i, "testPair:", n,
						"\n err:", err,
						"\n result:", result,
					)
				}
			}
			// second insert
			// insert docs with duplicate
			if i == 1 {
				result, err := mDocs.InsertDocs(testDB, testCollection, true)
				if result == "<nil>" || err != nil {
					t.Fatal("Error not nill after insert, i=", i, "testPair:", n,
						"\n err:", err,
						"\n result:", result,
					)
				}
			}
		}
	}
}
