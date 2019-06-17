package main

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

type testParseResp struct {
	input  []byte
	result *ftData
}

var testsParseResp = new(testParseResp)
var testFTdata = &ftData{
	[]string{"2019-04-11T00:00:00", "2019-04-12T00:00:00", "2019-04-15T00:00:00", "2019-04-16T00:00:00",
		"2019-04-17T00:00:00", "2019-04-18T00:00:00", "2019-04-22T00:00:00", "2019-04-23T00:00:00",
		"2019-04-24T00:00:00", "2019-04-25T00:00:00", "2019-04-26T00:00:00", "2019-04-29T00:00:00",
		"2019-04-30T00:00:00", "2019-05-01T00:00:00", "2019-05-02T00:00:00", "2019-05-03T00:00:00",
		"2019-05-06T00:00:00", "2019-05-07T00:00:00", "2019-05-08T00:00:00", "2019-05-09T00:00:00",
		"2019-05-10T00:00:00"},
	[]ftElement{
		{"price", "Iron Ore 62% Fe, CFR China (TSI) Swa", "27110161", "NYM",
			"USD", []ftComponent{
			{"Open", 95.08, 92.91,
				"2019-05-10T00:00:00", "2019-04-18T00:00:00",
				[]float64{93.6, 93.79, 93.97, 93.44, 92.96, 92.91, 93.2, 93.29, 92.95, 93.14, 93.16, 93.17,
					93.24, 93.65, 93.6, 93.81, 93.51, 94.75, 94.19, 94.24, 95.08}},
			{"High", 95.08, 92.91,
				"2019-05-10T00:00:00", "2019-04-18T00:00:00",
				[]float64{93.6, 93.79, 93.97, 93.44, 92.96, 92.91, 93.2, 93.29, 92.95, 93.14, 93.16, 93.17,
					93.24, 93.65, 93.6, 93.81, 93.51, 94.75, 94.19, 94.24, 95.08}},
			{"Low", 95.08, 92.91,
				"2019-05-10T00:00:00", "2019-04-18T00:00:00",
				[]float64{93.6, 93.79, 93.97, 93.44, 92.96, 92.91, 93.2, 93.29, 92.95, 93.14, 93.16, 93.17,
					93.24, 93.65, 93.6, 93.81, 93.51, 94.75, 94.19, 94.24, 95.08}},
			{"Close", 95.08, 92.91,
				"2019-05-10T00:00:00", "2019-04-18T00:00:00",
				[]float64{93.6, 93.79, 93.97, 93.44, 92.96, 92.91, 93.2, 93.29, 92.95, 93.14, 93.16, 93.17,
					93.24, 93.65, 93.6, 93.81, 93.51, 94.75, 94.19, 94.24, 95.08}},
		},},
		{"volume", "Iron Ore 62% Fe, CFR China (TSI) Swa", "27110161", "NYM",
			"USD", []ftComponent{
			{"Volume", 0, 0,
				"2019-04-11T00:00:00", "2019-04-11T00:00:00",
				[]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		},},
	},
}
var etalonMongoDocs = []mongoDoc{
	{Id: "2019-04-11T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 11, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.6, High: 93.6, Low: 93.6, Close: 93.6, Volume: 0},
	{Id: "2019-04-12T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 12, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.79, High: 93.79, Low: 93.79, Close: 93.79, Volume: 0},
	{Id: "2019-04-15T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 15, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.97, High: 93.97, Low: 93.97, Close: 93.97, Volume: 0},
	{Id: "2019-04-16T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 16, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.44, High: 93.44, Low: 93.44, Close: 93.44, Volume: 0},
	{Id: "2019-04-17T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 17, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 92.96, High: 92.96, Low: 92.96, Close: 92.96, Volume: 0},
	{Id: "2019-04-18T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 18, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 92.91, High: 92.91, Low: 92.91, Close: 92.91, Volume: 0},
	{Id: "2019-04-22T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 22, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.2, High: 93.2, Low: 93.2, Close: 93.2, Volume: 0},
	{Id: "2019-04-23T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 23, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.29, High: 93.29, Low: 93.29, Close: 93.29, Volume: 0},
	{Id: "2019-04-24T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 24, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 92.95, High: 92.95, Low: 92.95, Close: 92.95, Volume: 0},
	{Id: "2019-04-25T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 25, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.14, High: 93.14, Low: 93.14, Close: 93.14, Volume: 0},
	{Id: "2019-04-26T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 26, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.16, High: 93.16, Low: 93.16, Close: 93.16, Volume: 0},
	{Id: "2019-04-29T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 29, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.17, High: 93.17, Low: 93.17, Close: 93.17, Volume: 0},
	{Id: "2019-04-30T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 4, 30, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.24, High: 93.24, Low: 93.24, Close: 93.24, Volume: 0},
	{Id: "2019-05-01T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.65, High: 93.65, Low: 93.65, Close: 93.65, Volume: 0},
	{Id: "2019-05-02T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 5, 2, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.6, High: 93.6, Low: 93.6, Close: 93.6, Volume: 0},
	{Id: "2019-05-03T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 5, 3, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.81, High: 93.81, Low: 93.81, Close: 93.81, Volume: 0},
	{Id: "2019-05-06T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 5, 6, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 93.51, High: 93.51, Low: 93.51, Close: 93.51, Volume: 0},
	{Id: "2019-05-07T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 5, 7, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 94.75, High: 94.75, Low: 94.75, Close: 94.75, Volume: 0},
	{Id: "2019-05-08T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 5, 8, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 94.19, High: 94.19, Low: 94.19, Close: 94.19, Volume: 0},
	{Id: "2019-05-09T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 5, 9, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 94.24, High: 94.24, Low: 94.24, Close: 94.24, Volume: 0},
	{Id: "2019-05-10T00:00:00_TIOc1:NYM",
		Time:   time.Date(2019, 5, 10, 0, 0, 0, 0, time.UTC),
		Source: "FT", Symbol: "TIOc1", ExchangeID: "NYM",
		Open: 95.08, High: 95.08, Low: 95.08, Close: 95.08, Volume: 0},
}

func TestMain(m *testing.M) {
	jsonFile := "ft-resp-30day.json"
	input, e := ioutil.ReadFile(jsonFile)
	check(e)
	(*testsParseResp).input = input
	(*testsParseResp).result = testFTdata

	exitCode := m.Run()

	os.Exit(exitCode)

}

// help function to Drop collection, before test
func DropCollection(db, collect string) error {
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	CheckF(err)
	defer client.Disconnect(context.TODO())
	err = client.Ping(context.TODO(), nil)
	CheckF(err)
	Debug.Println("Connected to MongoDB!")
	collection := client.Database(db).Collection(collect)
	err = (*collection).Drop(context.TODO())
	return err
}

// --------------- tests ----------------------

func TestPrepareReq(t *testing.T){
	prepareReq(7, "kkkkkk")
}

func TestParseResp(t *testing.T) {
	urlRes, err := parseResp(testsParseResp.input)

	if err != nil {
		t.Fatal("Error returned by Unmarshal")
	}

	///////////////// check root keys //////////////////////
	if len((*urlRes).Dates) != len((*testsParseResp).result.Dates) {
		t.Fatal(
			"Number of Dates incorrect:",
			"\n exp:", len((*testsParseResp).result.Dates),
			"\n got:", len((*urlRes).Dates),
		)
	}
	for i := 0; i < len((*urlRes).Dates); i++ {
		if (*urlRes).Dates[i] != (*testsParseResp).result.Dates[i] {
			t.Fatal("Error in Dates, index=", i,
				"\n epx:", (*testsParseResp).result.Dates[i],
				"\n got:", (*urlRes).Dates[i])
		}
	}

	////////////////// check Elements keys ///////////////////
	if len((*urlRes).Elements) != len((*testsParseResp).result.Elements) {
		t.Fatal(
			"Number of Elements != 2",
			"\n exp:", len((*testsParseResp).result.Elements),
			"\n got:", len((*urlRes).Elements),
			"Element1=", (*urlRes).Elements[0],
		)
	}

	elementFields := []string{"Type", "CompanyName", "Symbol", "ExchangeId", "Currency"}

	for i := 0; i < len((*urlRes).Elements); i++ {
		refl_urlElement := reflect.ValueOf(urlRes.Elements[i])
		refl_resEment := reflect.ValueOf(testsParseResp.result.Elements[i])
		for _, field := range elementFields {
			if reflect.Indirect(refl_urlElement).FieldByName(field).String() !=
				reflect.Indirect(refl_resEment).FieldByName(field).String() {
				t.Fatal("Error in Element=", i, "field=", field,
					"\n exp:", reflect.Indirect(refl_resEment).FieldByName(field),
					"\n got:", reflect.Indirect(refl_urlElement).FieldByName(field))
			}
		}
	}

	//////////////// check ComponentSeries //////////////////////
	for i := 0; i < len((*urlRes).Elements); i++ {
		if len((*urlRes).Elements[i].ComponentSeries) != len((*testsParseResp).result.Elements[i].ComponentSeries) {
			t.Fatal(
				"in Elements=", i, "number of ComponentSeries not eq",
				"\n exp:", len((*testsParseResp).result.Elements[i].ComponentSeries),
				"\n got:", len((*urlRes).Elements[i].ComponentSeries),
			)
		}
	}
	componentFieldsStr := []string{"Type", "MaxValueDate", "MinValueDate"}
	componentFieldsFloat := []string{"MaxValue", "MinValue"}

	for i := 0; i < len((*urlRes).Elements); i++ {
		for j := 0; j < len((*urlRes).Elements[i].ComponentSeries); j++ {
			refl_urlComp := reflect.ValueOf(urlRes.Elements[i].ComponentSeries[j])
			refl_resComp := reflect.ValueOf(testsParseResp.result.Elements[i].ComponentSeries[j])
			for _, field := range componentFieldsStr {
				if reflect.Indirect(refl_urlComp).FieldByName(field).String() !=
					reflect.Indirect(refl_resComp).FieldByName(field).String() {
					t.Fatal("Error in Component=", i, "field=", field,
						"\n exp:", reflect.Indirect(refl_resComp).FieldByName(field),
						"\n got:", reflect.Indirect(refl_urlComp).FieldByName(field))
				}
			}
			for _, field := range componentFieldsFloat {
				if reflect.Indirect(refl_urlComp).FieldByName(field).Float() !=
					reflect.Indirect(refl_resComp).FieldByName(field).Float() {
					t.Fatal("Error in Component=", i, "field=", field,
						"\n exp:", reflect.Indirect(refl_resComp).FieldByName(field),
						"\n got:", reflect.Indirect(refl_urlComp).FieldByName(field))
				}

			}

			//	//////// Values ///////////////////////////
			if len((*urlRes).Elements[i].ComponentSeries[j].Values) !=
				len((*testsParseResp).result.Elements[i].ComponentSeries[j].Values) {
				t.Fatal(
					"Number of Values not eq",
					"\n exp:", len((*testsParseResp).result.Elements[i].ComponentSeries[j].Values),
					"\n got:", len((*urlRes).Elements[i].ComponentSeries[j].Values),
				)
			}
			for ii := 0; ii < len((*urlRes).Elements[i].ComponentSeries[j].Values); ii++ {
				if (*urlRes).Elements[i].ComponentSeries[j].Values[ii] !=
					(*testsParseResp).result.Elements[i].ComponentSeries[j].Values[ii] {
					t.Fatal(
						"Error in Values, index=", ii,
						"\n exp:", (*testsParseResp).result.Elements[i].ComponentSeries[j].Values[ii],
						"\n got:", (*urlRes).Elements[i].ComponentSeries[j].Values[ii],
					)
				}
			}
		}
	}

	//////////////// final check struct /////////////
	if !reflect.DeepEqual(*urlRes, *(testsParseResp.result)) {
		t.Fatal("Result struct not as expected")
	}
}

func TestCheckDataFT(t *testing.T) {
	err := testFTdata.checkData()
	if err != nil {
		t.Fatal("Error in checking data",
			"\n exp:", "nil",
			"\n got:", err)
	}
}

func TestGenMongoDocs(t *testing.T) {
	mondoDocs, err := testFTdata.genMongoDocs()
	if err != nil {
		t.Fatal("Error in parsing ft data")
	}

	for i := 0; i < len(mondoDocs.docs); i++ {
		//gotDoc := mondoDocs[i]
		//expDoc := etalonMongoDocs[i]
		//refl_doc := reflect.ValueOf(gotDoc)
		//refl_etalon := reflect.ValueOf(expDoc)
		//if refl_doc.NumField() != refl_etalon.NumField() {
		//	t.Fatal("Number of field different in slice idx=", i,
		//		"\n exp:", etalonMongoDocs[i],
		//		"\n got:", mondoDocs[i])
		//}
		//
		//for f := 0; f < refl_doc.NumField(); f++ {
		//	if !refl_etalon.Field(f).CanInterface(){
		//
		//	fmt.Println(refl_etalon.Field(f).String())
		//	}
		//
		//	//if refl_etalon.Field(f).Interface() != reflect.Indirect(refl_doc).Field(f).Interface() {
		//	//	t.Fatal("values in fields not eq, docNum=", i, "field=", refl_etalon.Type().Field(f).Name,
		//	//		"\n exp:", refl_etalon.Field(f).Interface(),
		//	//		"\n got:", refl_doc.Field(f).Interface())
		//	//}
		//}

		if ! reflect.DeepEqual(mondoDocs.docs[i], etalonMongoDocs[i]) {
			t.Fatalf("doc not as expected in doc=%v \n exp: %+v \n got: %+v",
				i, etalonMongoDocs[i], mondoDocs.docs[i])
		}
	}
}

// connect to DB should be established
func TestInsertDocs(t *testing.T) {
	mDocs := mongoDocs{etalonMongoDocs}
	testDB := "test"
	testCollection := "commoditiesT"
	err := DropCollection(testDB, testCollection)
	CheckF(err)

	for i := 0; i < 2; i++ {
		// insert docs without duplicate
		if i == 0 {
			result, err := mDocs.InsertDocs(testDB, testCollection)
			if result != "<nil>" || err != nil {
				t.Fatal("Error not nill after insert, i=", i,
					"\n err:", err,
					"\n result:", result,
				)
			}
		}
		// insert docs with duplicate
		if i == 1 {
			result, err := mDocs.InsertDocs(testDB, testCollection)
			if result == "<nil>" || err != nil {
				t.Fatal("Error not nill after insert, i=", i,
					"\n err:", err,
					"\n result:", result,
				)
			}
		}
	}

}
