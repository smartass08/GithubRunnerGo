package utils

import "C"
import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

var ALL []info

type DB struct {
	client *mongo.Client
}

type info struct {
	Repo    string `json:"repo"`
	Token   string `json:"token"`
	Running bool   `json:"running"`
}

func genjson(token string, repo string, running bool) bson.M {
	return bson.M{"repo": repo,"token": token, "running": running}
}

func (C *DB) Access(url string) {
	log.Println("Starting DB process")
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Println("Error while connecting DB", err)
		return
	}
	C.client = client

}

func (C *DB) GetAllConfigs() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := C.client.Database(GetDbName()).Collection(GetDbCollection())
	all, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
	}
	defer all.Close(ctx)
	for all.Next(ctx) {
		var a info
		var result bson.M
		err = all.Decode(&result)
		if err != nil {
			log.Println(err)
		} else {
			if result["repo"] != nil && result["running"] != nil && result["token"] != nil{
				_ = bson.Unmarshal(all.Current, a)
				ALL = append(ALL, a)
			}
		}
	}

}

func (C *DB) Insert(repo string, token string, running bool) {
	wg.Add(1)
	go func() {
		collection := C.client.Database(GetDbName()).Collection(GetDbCollection())
		_, err := collection.InsertOne(context.Background(), genjson(token, repo, running))
		if err != nil {
			log.Printf("Error inserting the id into db %s", err)
			return
		}
	}()
}

func CheckValid(repo string) bool {
	for _, i := range ALL {
		if i.Repo == repo {
			return true
		}
	}
	return false
}

func getIdIndex(repo string) (int, bool) {
	for i, j := range ALL {
		if j.Repo == repo {
			return i, true
		}
	}
	return 0, false
}

func (C *DB) Delete(repo string, token string, running bool) {
	defer wg.Done()
	index, found := getIdIndex(repo)
	if found {
		ALL[index] = ALL[len(ALL)-1]
		ALL[len(ALL)-1] = info{}
		ALL = ALL[:len(ALL)-1]
	}else {
		return
	}
	wg.Add(1)
	go func() {
		collection := C.client.Database(GetDbName()).Collection(GetDbCollection())
		_, err := collection.DeleteOne(context.Background(), genjson(token, repo, running))
		if err != nil {
			log.Println(err)
			return
		}
	}()
}

func GetValues(repo string)*info{
	for _,v := range ALL{
		if v.Repo == repo{
			return &v
		}
	}
	return &info{}
}