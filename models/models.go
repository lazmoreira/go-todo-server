package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ToDoList struct
type ToDoList struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Status bool               `json:"status"`
	Task   string             `json:"task"`
}
