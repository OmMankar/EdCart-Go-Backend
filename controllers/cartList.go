package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/models"
	"main/utiles"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartList struct {
	Client *mongo.Client
	AppCfg models.Config
}

func (p *CartList) SetCartlist(w http.ResponseWriter, r *http.Request) {

	//Creating instance of resp struct.
	var respData interface{}
	resp := utiles.WBody{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.UserCollection)

	// Accessing the Request params present in the URL.
	queryValues := mux.Vars(r)
	emailId := queryValues["emailId"]
	// queryy_params==> r.URL.Query().Get()

	//Reading the json data from the Request Body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}

	//Parsing the json data of request body into Card struct format.
	var obj models.Card
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("Error occured while UnMarshalling the json request body : ", err.Error())
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}

	//Query to insert the card into user's CartList array in the Database.
	filter := bson.M{"emailId": emailId}
	update := bson.M{"$push": bson.M{"cartList": bson.M{"category": obj.Category,
		"image": obj.Image, "courseName": obj.CourseName, "courseAuthor": obj.CourseAuthor,
		"ratingNumber": obj.RatingNumber, "numOfRatings": obj.NumOfRatings, "discountPrice": obj.DiscountPrice,
		"originalPrice": obj.OriginalPrice}}}

	var updatedDoc models.User
	err = collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&updatedDoc)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection.
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no document found : ", err.Error())
			resp.ResBodyBadRequest(respData, "Bad request", false, w)
			return
		}
		fmt.Println("err", err.Error())
		resp.ResBodyBadRequest(respData, "Bad request", false, w)
		return
	}
	updatedDoc.Password = ""

	//Sending the Response.
	resp.ResBodyStatusOK(updatedDoc, "Adding item to the cartList of db", true, w)
}
func (p *CartList) DeleteCartlist(w http.ResponseWriter, r *http.Request) {
	//Creating instance of resp struct.
	resp := utiles.WBody{}
	var respData interface{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.UserCollection)

	// Accessing the Request params present in the URL.
	queryValues := mux.Vars(r)
	emailId := queryValues["emailId"]
	// queryy_params==> r.URL.Query().Get()

	//Reading the json data from the Request Body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return

	}

	//Parsing the json data of request body into Card struct format.
	var obj models.Card
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}

	//Query to remove the card from user's CartList.
	filter := bson.M{"emailId": emailId}
	update := bson.M{"$pull": bson.M{"cartList": bson.M{"category": obj.Category,
		"image": obj.Image, "courseName": obj.CourseName, "courseAuthor": obj.CourseAuthor,
		"ratingNumber": obj.RatingNumber, "numOfRatings": obj.NumOfRatings, "discountPrice": obj.DiscountPrice,
		"originalPrice": obj.OriginalPrice}}}

	result := collection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection.
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			fmt.Println("no document found : ", result.Err().Error())
			resp.ResBodyBadRequest(respData, "Bad Request", false, w)
			return
		}
		fmt.Println("err", result.Err().Error())
		resp.ResBodyBadRequest(respData, "Bad Request", false, w)
		return
	}

	//Sending in response.
	resp.ResBodyStatusOK(obj, "Removed item to the cartlist of db", true, w)
}
func (p *CartList) DeleteAllCartlist(w http.ResponseWriter, r *http.Request) {

	//Creating instance of resp struct.
	resp := utiles.WBody{}
	var respData interface{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.UserCollection)

	// Extract emailId from URL parameters.
	queryValues := mux.Vars(r)
	emailId := queryValues["emailId"]

	// Define the filter and update operation.
	filter := bson.M{"emailId": emailId}
	update := bson.M{"$unset": bson.M{"cartList": ""}}

	// Perform the update operation.
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("Error occurred during the update:", err.Error())
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}
	//Sending response.
	resp.ResBodyStatusOK("", "Successfully cleared the cartList for emailId", true, w)

}
