package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func CheckF(e error){
	if e != nil {
		log.Fatal(e)
	}
}

var hours *int

func main() {
	db := flag.String("db", "", "mandatory arg, DB to use")
	collectionArg := flag.String("collection", "", "mandatory arg, collection to use")
	hours = flag.Int("hours", 8, "clean docs older then num hurs")
	flag.Parse()
	if *db == "" || *collectionArg == ""{
		log.Fatal("--db and --collection arguments is mandatory")
	}
	log.Println("DEBUG:", "cleaning all before=", *hours, "in db=", *db, "collection=", *collectionArg)

	zone, err := time.LoadLocation("UTC")
	//zone, err:= time.LoadLocation("Local")
	CheckF(err)
	currentTime := time.Now().In(zone)

	olderThen := currentTime.Add(time.Hour * -time.Duration(*hours))
	fmt.Printf("time= %v \n", olderThen)

	filterOlderThen := bson.M{"time": bson.M{"$lte": olderThen}}

	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	CheckF(err)

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	CheckF(err)
	fmt.Println("Connected to MongoDB!")

	collection := client.Database(*db).Collection(*collectionArg)
	delRes, err := collection.DeleteMany(context.TODO(), filterOlderThen)
	CheckF(err)
	fmt.Printf("Deleted %v documents in db= %v collection= %v \n", delRes.DeletedCount, *db, *collectionArg)

	// close when connection not need any more
	err = client.Disconnect(context.TODO())
	CheckF(err)
	fmt.Println("Connection to MongoDB closed.")

}