package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

type Student struct {
	Name  string `bson:"name"`
	Grade int    `bson:"grade"`
}

func connection() (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client.Database("belajar_mongodb"), nil
}

func insert(db *mongo.Database, data Student) {
	_, err := db.Collection("student").InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Insert success!")
}

func find(db *mongo.Database) {
	csr, err := db.Collection("student").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(ctx)

	result := make([]Student, 0)

	for csr.Next(ctx) {
		var row Student
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}
	if len(result) > 0 {
		for i := 0; i < len(result); i++ {
			fmt.Print("Name :", result[i].Name + "  |  ")
			fmt.Println("Grade :", result[i].Grade)
		}
	}
}

func update(db *mongo.Database, name string, data Student) {
	var selector = bson.M{"name": name}
	var changes = data

	_, err := db.Collection("student").UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	fmt.Println("Update success")
}

func main() {
	conn, err := connection()
	if err != nil {
		log.Println("Error connection : ", err.Error())
	}

	var number int
	fmt.Println("1. Insert data")
	fmt.Println("2. Get data")
	fmt.Println("3. Update data")
	fmt.Print("Select number : ")
	fmt.Scanln(&number)

	switch number{
	case 1:
		var name string
		var grade int

		fmt.Print("Name : ")
		fmt.Scanln(&name)
		fmt.Print("Grade : ")
		fmt.Scanln(&grade)

		data := Student{Name: name, Grade: grade}
		insert(conn, data)
	case 2:
		find(conn)
	case 3:
		var name, newData string
		var grade int

		fmt.Print("Siapa yang ingin diubah? : ")
		fmt.Scanln(&newData)
		fmt.Print("Name : ")
		fmt.Scanln(&name)
		fmt.Print("Grade : ")
		fmt.Scanln(&grade)

		data := Student{Name: name, Grade: grade}
		update(conn, newData, data)
	}
}