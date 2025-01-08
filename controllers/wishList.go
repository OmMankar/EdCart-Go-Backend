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

type WishList struct {
	Client *mongo.Client
	AppCfg models.Config
}

func (p *WishList) SetWishlist(w http.ResponseWriter, r *http.Request) {
	//Creating instance of resp struct.
	resp := utiles.WBody{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.UserCollection)

	// Accessing the request params from URL.
	queryValues := mux.Vars(r)
	emailId := queryValues["emailId"]
	// queryy_params==> r.URL.Query().Get()

	//Reading the json data from the request body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	}
	//Parsing the request body data into the Card struct format.
	var obj models.Card
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	}
	//Query to insert card into user WiushList Array.
	filter := bson.M{"emailId": emailId}
	update := bson.M{"$push": bson.M{"wishList": bson.M{"category": obj.Category,
		"image": obj.Image, "courseName": obj.CourseName, "courseAuthor": obj.CourseAuthor,
		"ratingNumber": obj.RatingNumber, "numOfRatings": obj.NumOfRatings, "discountPrice": obj.DiscountPrice,
		"originalPrice": obj.OriginalPrice}}}
	var updatedDoc models.User
	err = collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&updatedDoc)

	if err != nil {
		// ErrNoDocuments meresp that the filter did not match any documents in
		// the collection.
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no document found : ", err.Error())
			resp.ResBodyBadRequest("", "Bad Request", false, w)
			return
		}
		fmt.Println("err", err.Error())
		resp.ResBodyBadRequest("", "Bad Request", false, w)
		return
	}
	//sending response.
	resp.ResBodyStatusOK(obj, "Successfully entired Cards", true, w)

}
func (p *WishList) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	//Creating instance of resp struct.
	resp := utiles.WBody{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.UserCollection)

	// Accessing the request params from URL.
	queryValues := mux.Vars(r)
	emailId := queryValues["emailId"]
	// queryy_params==> r.URL.Query().Get()

	//Reading the json data from the request body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	}

	//Parsing the request body data into the Card struct format.
	var obj models.Card
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	}

	//Query to remove card into user WiushList Array.
	filter := bson.M{"emailId": emailId}
	update := bson.M{"$pull": bson.M{"wishList": bson.M{"category": obj.Category,
		"image": obj.Image, "courseName": obj.CourseName, "courseAuthor": obj.CourseAuthor,
		"ratingNumber": obj.RatingNumber, "numOfRatings": obj.NumOfRatings, "discountPrice": obj.DiscountPrice,
		"originalPrice": obj.OriginalPrice}}}

	var updatedDoc models.User

	err = collection.FindOneAndUpdate(context.Background(), filter, update).Decode(&updatedDoc)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no document found : ", err.Error())
			resp.ResBodyBadRequest("", "Bad Request", false, w)
			return
		}
		fmt.Println("err", err.Error())
		resp.ResBodyBadRequest("", "Bad Request", false, w)
		return
	}

	//sending Response.
	resp.ResBodyStatusOK(obj, "Successfully removed entire Cards", true, w)
}
