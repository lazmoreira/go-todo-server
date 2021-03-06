package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-http-utils/headers"
	"github.com/gorilla/mux"
	"github.com/lazmoreira/go-todo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://todoUser:todopass@zero-d4efp.gcp.mongodb.net/test?retryWrites=true&w=majority"

//const connectionString = "mongodb://0.0.0.0:27017/admin"

const dbName = "TodoDB"

const collectionName = "Todos"

var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database(dbName).Collection(collectionName)

	fmt.Println("Collection instance created!")
}

// GetAllTask get all task route
func GetAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headers.ContentType, "application/x-www-form-urlencoded")
	w.Header().Set(headers.AccessControlAllowOrigin, "*")

	payload := getAllTask()
	json.NewEncoder(w).Encode(payload)
}

// CreateTask create task route
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headers.ContentType, "application/x-www-form-urlencoded")
	w.Header().Set(headers.AccessControlAllowOrigin, "*")
	w.Header().Set(headers.AccessControlAllowMethods, "POST")
	w.Header().Set(headers.AccessControlAllowHeaders, "Content-Type")

	var task models.ToDoList

	_ = json.NewDecoder(r.Body).Decode(&task)

	insertOneResult := insertOneTask(task)

	if oid, ok := insertOneResult.InsertedID.(primitive.ObjectID); ok {
		task.ID = oid
	}

	json.NewEncoder(w).Encode(task)
}

//TaskComplete task complete route
func TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headers.ContentType, "application/x-www-form-urlencoded")
	w.Header().Set(headers.AccessControlAllowOrigin, "*")
	w.Header().Set(headers.AccessControlAllowMethods, "PUT")
	w.Header().Set(headers.AccessControlAllowHeaders, "Content-Type")

	if r.Method == "PUT" {
		params := mux.Vars(r)
		json.NewEncoder(w).Encode(taskComplete(params["id"], r))
	} else {
		return
	}
}

//UndoTask task complete route
func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headers.ContentType, "application/x-www-form-urlencoded")
	w.Header().Set(headers.AccessControlAllowOrigin, "*")
	w.Header().Set(headers.AccessControlAllowMethods, "PUT")
	w.Header().Set(headers.AccessControlAllowHeaders, "Content-Type")

	params := mux.Vars(r)

	undoTask(params["id"])

	json.NewEncoder(w).Encode(params["id"])
}

//DeleteTask task complete route
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headers.ContentType, "application/x-www-form-urlencoded")
	w.Header().Set(headers.AccessControlAllowOrigin, "*")
	w.Header().Set(headers.AccessControlAllowMethods, "DELETE")
	w.Header().Set(headers.AccessControlAllowHeaders, "Content-Type")

	if r.Method == "DELETE" {
		params := mux.Vars(r)

		deleteOneTask(params["id"])

		json.NewEncoder(w).Encode(params["id"])
	} else {
		return
	}
}

//DeleteAllTask task complete route
func DeleteAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(headers.ContentType, "application/x-www-form-urlencoded")
	w.Header().Set(headers.AccessControlAllowOrigin, "*")

	count := deleteAllTask()

	json.NewEncoder(w).Encode(count)
}

func getAllTask() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M

	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)

		if e != nil {
			log.Fatal(e)
		}

		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())

	return results
}

func insertOneTask(task models.ToDoList) *mongo.InsertOneResult {
	insertResult, err := collection.InsertOne(context.Background(), task)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New task added", insertResult.InsertedID)

	return insertResult
}

// task complete method, update task's status to true
func taskComplete(task string, r *http.Request) models.ToDoList {
	var initialTask models.ToDoList

	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}

	err := collection.FindOne(context.Background(), filter).Decode(&initialTask)

	if err != nil {
		log.Fatal(err)
	}

	update := bson.M{"$set": bson.M{"status": !initialTask.Status}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	var updatedTask models.ToDoList

	err = collection.FindOne(context.Background(), filter).Decode(&updatedTask)

	if err != nil {
		log.Fatal(err)
	}

	return updatedTask
}

// task undo method, update task's status to false
func undoTask(task string) {
	fmt.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": false}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

// delete one task from the DB, delete by ID
func deleteOneTask(task string) {
	fmt.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
}

// delete all the tasks from the DB
func deleteAllTask() int64 {
	d, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
	return d.DeletedCount
}
