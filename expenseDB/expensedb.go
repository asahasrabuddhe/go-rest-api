package expenseDB

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go-rest-api/errors"
	"go-rest-api/requests"
	"go-rest-api/responses"
	"go-rest-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
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

func (mh *MongoHandler) ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		expenseID := chi.URLParam(r, "id")
		id, _:=strconv.Atoi(expenseID)

		expense := &types.Expense{}
		err := mh.GetOne_DB(expense, bson.M{"id":id})
		if err !=nil{
			log.Println(err)
		}
		ctx := context.WithValue(r.Context(), "expense", expense )
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}



func (mh *MongoHandler) GetOne_DB(expense *types.Expense, filter interface{}) error {
	//Will automatically create a collection if not available
	collection := mh.client.Database(mh.database).Collection("expensecoll")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(expense)
	return err
}
func  (mh *MongoHandler) GetId(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(*types.Expense)
	err := render.Render(writer, request, responses.Listexpense(expense))
	if err != nil{
		log.Println(err)
		render.Render(writer,request,errors.ErrRender(err))
		return
	}
}

func (mh *MongoHandler) GetAll_DB(filter interface{}) []*types.Expense {
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
func (mh *MongoHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	expenses := mh.GetAll_DB(bson.M{})
	if err := render.Render(writer, request, responses.ExpensesResponse(expenses)); err != nil{
		render.Render(writer,request,errors.ErrRender(err))
		return
	}
}

func (mh *MongoHandler) AddOne_DB(exp *types.Expense) (*mongo.InsertOneResult, error) {
	collection := mh.client.Database(mh.database).Collection("expensecoll")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, exp)
	return result, err
}
func (mh *MongoHandler)Create(writer http.ResponseWriter, request *http.Request) {
	var req requests.CreateExpenseRequest

	err := render.Bind(request, &req)
	if err != nil {
		render.Render(writer, request, errors.ErrInvalidRequest(err))
		return
	}


	//expenses = append(expenses, *req.Expense)

	req.Expense.CreatedOn=time.Now()
	_,err =mh.AddOne_DB(req.Expense)
	if err!= nil{
		log.Println(err)
	}

	j, _ := json.Marshal(req.Expense)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	_, _ = fmt.Fprintf(writer, `{"success": true, "data": %v}`, string(j))
}

func (mh *MongoHandler) Update_DB(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mh.client.Database(mh.database).Collection("expensecoll")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	exp:= bson.D{{"$set",update}}
	result, err := collection.UpdateOne(ctx, filter,exp )
	return result, err
}
func (mh *MongoHandler) Update(writer http.ResponseWriter, request *http.Request) {

	expense := request.Context().Value("expense").(*types.Expense)

	var req requests.UpdateExpenseRequest

	err := render.Bind(request, &req)
	if err != nil {
		render.Render(writer,request,errors.ErrRender(err))
		return
	}

	req.Expense.CreatedOn = expense.CreatedOn
	req.Expense.UpdatedOn=time.Now()
	_, err = mh.Update_DB(bson.D{{"id",expense.Id}},req.Expense)
	if err!=nil{
		log.Println(err)
	}

	//expenses[expense.Id-1] = *req.Expense
	_=mh.GetOne_DB(expense,bson.M{"id":expense.Id})

	if err = render.Render(writer, request, responses.Listexpense(expense)) ; err != nil{
		render.Render(writer,request,errors.ErrRender(err))
		return

	}
}


func (mh *MongoHandler) RemoveOne_DB(filter interface{}) (*mongo.DeleteResult, error) {
	collection := mh.client.Database(mh.database).Collection("expensecoll")
	ctx, _ := context.WithTimeout(context.Background(), 50*time.Second)

	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}

func (mh *MongoHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	expense := request.Context().Value("expense").(*types.Expense)
	_, err := mh.RemoveOne_DB(bson.D{{"id", expense.Id}})
	if err!=nil{
		log.Println(err)
		return
	}

}
