package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"main/utiles"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Card struct {
	Client *mongo.Client
	AppCfg models.Config
}

func (p *Card) GetAllCards(w http.ResponseWriter, r *http.Request) {
	//Creating instance of resp struct.
	resp := utiles.WBody{}
	var respData interface{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.CardCollection)

	//Running Query to obtain all the Cards Present in Database.
	filter := bson.D{{}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("when running find query")
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}
	var result []models.Card
	if err = cursor.All(context.TODO(), &result); err != nil {
		fmt.Println("when cursor is converted to the struct card form")
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}

	//Sending all the cards in response.
	resp.ResBodyStatusOK(result, "All the cards are fetched", true, w)
}

func (p *Card) CreateCard(w http.ResponseWriter, r *http.Request) {
	//Creating instance of resp struct.
	resp := utiles.WBody{}
	var respData interface{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.CardCollection)

	// Reading the json data present in Request body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}
	var obj models.Card
	//Parsing the json Data in to ard stuct format.
	if err := json.Unmarshal(data, &obj); err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}

	//Inserting the new card into Database.
	if _, err := collection.InsertOne(context.TODO(), obj); err != nil {
		fmt.Println("Error occured while inserting data in database : ", err.Error())
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
	}

	//Sending Response
	resp.ResBodyStatusOK(obj, "Create the teacher in db successfully", true, w)

}
