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

type MyLearning struct {
	Client *mongo.Client
	AppCfg models.Config
}

func (p *MyLearning) SetMyLearning(w http.ResponseWriter, r *http.Request) {
	//Creating instance of resp struct.
	resp := utiles.WBody{}
	var respData interface{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.UserCollection)

	// Accessing the request params from URL.
	queryValues := mux.Vars(r)
	emailId := queryValues["emailId"]
	// queryy_params==> r.URL.Query().Get()

	//Reading Json Data form Response Body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request body")
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}

	// Parsing the request body data in Card struct format.
	var obj models.Card
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("Error occured while reading the request body")
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}

	// Query to insert card into user's Mylearning Array of Database.
	filter := bson.M{"emailId": emailId}
	update := bson.M{"$push": bson.M{"myLearningList": bson.M{"category": obj.Category,
		"image": obj.Image, "courseName": obj.CourseName, "courseAuthor": obj.CourseAuthor,
		"ratingNumber": obj.RatingNumber, "numOfRatings": obj.NumOfRatings, "discountPrice": obj.DiscountPrice,
		"originalPrice": obj.OriginalPrice}}}

	res := collection.FindOneAndUpdate(context.Background(), filter, update)
	if res.Err() != nil {
		// ErrNoDocuments meresp that the filter did not match any documents in
		// the collection.
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			fmt.Println("no document found : ", res.Err().Error())
			resp.ResBodyBadRequest(respData, "Bad Request", false, w)
			return
		}
		fmt.Println("err", res.Err().Error())
		resp.ResBodyBadRequest(respData, "Bad Request", false, w)
		return
	}

	//Sending the response
	resp.ResBodyStatusOK(obj, "Successfully added to my learning", true, w)

}
func (p *MyLearning) DeleteMyLearning(w http.ResponseWriter, r *http.Request) {
	//Creating instance of resp struct.
	resp := utiles.WBody{}
	var respData interface{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.UserCollection)

	// Accessing the request params from URL.
	queryValues := mux.Vars(r)
	emailId := queryValues["emailId"]
	// queryy_params==> r.URL.Query().Get()

	//Reading the json data from Request Body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request body")
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}

	// Parsing the json data of request body into Card stuct form.
	var obj models.Card
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("Error occured while reading the request body")
		resp.ResBodyInternalServerError(respData, "Internal Server Error", false, w)
		return
	}

	// Query to remove card from user Mylearning array present inside DataBase.
	filter := bson.M{"emailId": emailId}
	update := bson.M{"$pull": bson.M{"myLearningList": bson.M{"category": obj.Category,
		"image": obj.Image, "courseName": obj.CourseName, "courseAuthor": obj.CourseAuthor,
		"ratingNumber": obj.RatingNumber, "numOfRatings": obj.NumOfRatings, "discountPrice": obj.DiscountPrice,
		"originalPrice": obj.OriginalPrice}}}

	result := collection.FindOneAndUpdate(context.Background(), filter, update)

	if result.Err() != nil {
		// ErrNoDocuments meresp that the filter did not match any documents in
		// the collection.
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			fmt.Println("no document found : ", result.Err().Error())
			resp.ResBodyBadRequest(respData, "Bad Request", false, w)
			return
		}
		fmt.Println("err", result.Err().Error())
		resp.ResBodyBadRequest(respData, "Bad Request", false, w)
		return
	}

	//Sending the object we added.
	resp.ResBodyStatusOK(obj, "Suceesfuly Removed the card ", true, w)

}
