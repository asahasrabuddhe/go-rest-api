package expenseDb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go-rest-api/ers"
	"go-rest-api/requests"
	"go-rest-api/response"
	"go-rest-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
	"time"
)

const DefaultDatabase = "expenseCollection"

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

func (mh *MongoHandler) ExpenseCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		expenseID := chi.URLParam(r, "id")
		id, _ := strconv.Atoi(expenseID)
		expense := &types.Expense{}
		err := mh.GetOne_DB(expense, bson.M{"id": id})
		if err != nil {
			log.Println(err)
		}
		ctx := context.WithValue(r.Context(), "expense", expense)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mh *MongoHandler) GetOne(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(*types.Expense)
	err := render.Render(writer, request, response.Listexpense(expense))
	if err != nil {
		log.Println(err)
		_ = render.Render(writer, request, ers.ErrRender(err))
		return
	}
}

func (mh *MongoHandler) GetOne_DB(c *types.Expense, filter interface{}) error {
	//Will automatically create a collection if not available
	collection := mh.client.Database(mh.database).Collection("expenseCollection")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(c)
	return err
}

func (mh *MongoHandler) Get(filter interface{}) []*types.Expense {
	collection := mh.client.Database(mh.database).Collection("expenseCollection")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var result []*types.Expense
	for cur.Next(ctx) {
		contact := &types.Expense{}
		er := cur.Decode(contact)
		if er != nil {
			log.Fatal(er)
		}
		result = append(result, contact)
	}
	return result
}

func (mh *MongoHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	expenses := mh.Get(bson.M{})
	_ = json.NewEncoder(writer).Encode(expenses)
}

func (mh *MongoHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var req requests.CreateExpenseRequest
	err := render.Bind(request, &req)
	if err != nil {
		log.Println(err)
		return
	}
	req.Expense.CreatedOn = time.Now()
	_, err = mh.AddOne(req.Expense)
	if err != nil {
		log.Println(err)
		return
	}
	j, err := json.Marshal(req.Expense)
	if err != nil {
		_ = render.Render(writer, request, ers.ErrRender(err))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(writer, `{"success": true, "data": %v}`, string(j))
}

func (mh *MongoHandler) AddOne(c *types.Expense) (*mongo.InsertOneResult, error) {
	collection := mh.client.Database(mh.database).Collection("expenseCollection")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, c)
	return result, err
}

func (mh *MongoHandler) Update(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(*types.Expense)
	var req requests.UpdateExpenseRequest
	err := render.Bind(request, &req)
	if err != nil {
		log.Println(err)
		return
	}
	req.Expense.CreatedOn = expense.CreatedOn
	req.Expense.UpdatedOn = time.Now()
	_, err = mh.Update_DB(bson.M{"id": expense.Id}, *req.Expense)
	if err != nil {
		log.Println(err)
	}
	_ = mh.GetOne_DB(expense, bson.M{"id": expense.Id})
	err = render.Render(writer, request, response.Listexpense(expense))
	if err != nil {
		_ = render.Render(writer, request, ers.ErrRender(err))
		return
	}
}

func (mh *MongoHandler) Update_DB(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mh.client.Database(mh.database).Collection("expenseCollection")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	expense := bson.D{{"$set", update}}
	result, err := collection.UpdateOne(ctx, filter, expense)
	return result, err
}

func (mh *MongoHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(*types.Expense)
	_, err := mh.RemoveOne(bson.M{"id": expense.Id})
	if err != nil {
		log.Println(err)
		return
	}
}

func (mh *MongoHandler) RemoveOne(filter interface{}) (*mongo.DeleteResult, error) {
	collection := mh.client.Database(mh.database).Collection("expenseCollection")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}
