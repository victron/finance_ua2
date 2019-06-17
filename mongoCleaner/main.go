package main

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func CheckF(e error){
	if e != nil {
		Error.Fatal(e)
	}
}

func main() {

	Debug.Println("cleaning all before=", *hours, "in db=", *db, "collection=", *collectionArg)

	zone, err := time.LoadLocation("UTC")
	//zone, err:= time.LoadLocation("Local")
	CheckF(err)
	currentTime := time.Now().In(zone)

	olderThen := currentTime.Add(time.Hour * -time.Duration(*hours))
	Info.Printf("clean OlderThen= %v \n", olderThen)

	filterOlderThen := bson.M{"time": bson.M{"$lte": olderThen}}

	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	CheckF(err)

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	CheckF(err)
	Debug.Println("Connected to MongoDB!")

	collection := client.Database(*db).Collection(*collectionArg)
	delRes, err := collection.DeleteMany(context.TODO(), filterOlderThen)
	CheckF(err)
	Info.Printf("Deleted %v documents in db= %v collection= %v \n", delRes.DeletedCount, *db, *collectionArg)

	// close when connection not need any more
	err = client.Disconnect(context.TODO())
	CheckF(err)
	Debug.Println("Connection to MongoDB closed.")

}