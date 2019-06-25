package expenseDB

import (
	"context"
	"go-rest-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DefaultDatabase = "expensedb"

type MongoHandler struct {
	client   *mongo.Client
	database string
}

//MongoHandler Constructor
func NewHandler(address string) *MongoHandler {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cl, _ := mongo.Connect(ctx, options.Client().ApplyURI(address))
	mh := &MongoHandler{
		client:   cl,
		database: DefaultDatabase,
	}
	return mh
}

func (mh *MongoHandler) GetOne(expense *types.Expense, filter interface{}) error {
	//Will automatically create a collection if not available
	collection := mh.client.Database(mh.database).Collection("expensecoll")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(expense)
	return err
}

func (mh *MongoHandler) Get(filter interface{}) []*types.Expense {
	collection := mh.client.Database(mh.database).Collection("expensecoll")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var result []*types.Expense
	for cur.Next(ctx) {
		expense := &types.Expense{}
		er := cur.Decode(expense)
		if er != nil {
			log.Fatal(er)
		}
		result = append(result, expense)
	}
	return result
}

func (mh *MongoHandler) AddOne(exp *types.Expense) (*mongo.InsertOneResult, error) {
	collection := mh.client.Database(mh.database).Collection("expensecoll")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, exp)
	return result, err
}

func (mh *MongoHandler) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mh.client.Database(mh.database).Collection("expensecoll")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	exp:= bson.D{{"$set",update}}
	result, err := collection.UpdateOne(ctx, filter,exp )
	return result, err
}

func (mh *MongoHandler) RemoveOne(filter interface{}) (*mongo.DeleteResult, error) {
	collection := mh.client.Database(mh.database).Collection("expensecoll")
	ctx, _ := context.WithTimeout(context.Background(), 50*time.Second)

	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}
