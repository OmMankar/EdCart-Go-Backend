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
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Client *mongo.Client
	AppCfg models.Config
}

func (p *Auth) Signup(w http.ResponseWriter, r *http.Request) {
	resp := utiles.WBody{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.UserCollection)

	//Reading the data from request body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	}
	obj := models.User{}
	//To set the default value as empty array.
	obj.CartList = []models.Card{}
	obj.WishList = []models.Card{}
	obj.MyLearningList = []models.Card{}

	//Parsing the json data from r.body into User struct.
	if err := json.Unmarshal(data, &obj); err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	}

	//Validating the input information.
	if obj.Name == "" || obj.EmailId == "" || obj.Password == "" {
		fmt.Println("Bad Request try again")
		resp.ResBodyBadRequest("", "Bad Request", false, w)
		return
	}

	// Check if emailId already exists in the database and Retrieves a document that matches the filter.
	filter := bson.M{"emailId": obj.EmailId}
	opts := options.FindOne()
	result := models.User{}
	if err := collection.FindOne(context.TODO(), filter, opts).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			//Do nothing.
		} else {
			fmt.Println("Error occured while checking wheter email id  already present : ", err.Error())
			resp.ResBodyBadRequest("", "Bad Request", false, w)
			return
		}
	} else {
		fmt.Println("User Already exist")
		resp.ResBodyBadRequest("", "Bad Request", false, w)
		return
	}

	//Hashing Password.
	hash, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error occured while hashing the password : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	}
	//Copying the Hash Password to our obj.Password.
	var res models.User = obj
	res.Password = string(hash)

	//Inserting in New User into Database.
	if result, err := collection.InsertOne(context.TODO(), res); err != nil {
		fmt.Println("Error occured while inserting data in database : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	} else {
		fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	}
	obj.Password = ""

	//Sending Response
	resp.ResBodyStatusOK("", "Create Your account successfully.", true, w)

}
func (p *Auth) Login(w http.ResponseWriter, r *http.Request) {
	resp := utiles.WBody{}

	// Configuring the DB collection.
	collection := p.Client.Database(p.AppCfg.DbName).Collection(p.AppCfg.UserCollection)

	//Reading json data from request body.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request body : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	}

	var obj models.User

	//Parsing the json data of  r.body into User struct format.
	if err := json.Unmarshal(data, &obj); err != nil {
		fmt.Println("Error occured Unmarshalling json type request body : ", err.Error())
		resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
		return
	}

	//validating the input information.
	if obj.EmailId == "" || obj.Password == "" {
		fmt.Println("Bad Request try again")
		resp.ResBodyBadRequest("", "Bad Request", false, w)
		return
	}

	// Check if emailId already exists in the database and Retrieves a document that matches the filter.
	filter := bson.M{"emailId": obj.EmailId}
	opts := options.FindOne()
	var result models.User
	if err := collection.FindOne(context.TODO(), filter, opts).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Invalid Email Id / Please sign Up : ", err.Error())
			resp.ResBodyBadRequest("", "Bad Request", false, w)
			return
		} else {
			fmt.Println("Error occured while checking wheter email id  already present")
			resp.ResBodyInternalServerError("", "Internal Server Error", false, w)
			return
		}
	}

	// Verifying if the password from the client side matches the one retrieved from the database
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(obj.Password))
	if err != nil {
		fmt.Println("Incorrect password")
		resp.ResBodyBadRequest("", "Bad Request", false, w)
		return
	}
	result.Password = ""

	//Sending Response.
	resp.ResBodyStatusOK(result, "Successfully Logged In", true, w)

}
