package main

import (
	"context"
	"flag"
	"fmt"
	"main/config"
	"main/controllers"
	"main/database"
	"net/http"

	"github.com/gorilla/mux"
)

// // closure with callback function
// func CheckToken(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println(" middleware ", r.URL)
// 		val := r.Header.Get("jwttoken")
// 		if len(val) == 0 {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		next(w, r)
// 	})
// }
// func CORSHandler(p http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println(" cors middleware ", r.URL)
// 		w.Header().Add("CustomHeader", "custoVal")
// 		w.Header().Add("Access-Control-Allow-Origin", "*")
// 		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, HEAD")
// 		w.Header().Add("Access-Control-Allow-Headers", " X-Requested-With")
// 		w.Header().Add("Access-Control-Max-Age", " 1728000")
// 		p.ServeHTTP(w, r)
// 	})
// }

func main() {
	PORT := flag.Int("port", 3000, "Enter the port number")
	ConfPath := flag.String("conf_path", "/tmp/blala.txt", "Application Configuration")
	flag.Parse()

	var cfg config.Config
	appCfg := cfg.Load(*ConfPath)

	//Connecting to Database.
	client, err := database.Connection(appCfg.DNS())
	if err != nil {
		fmt.Println("err", err.Error())

	}
	defer client.Disconnect(context.TODO())

	r := mux.NewRouter()

	//Mounting our Url.
	subRouter := r.PathPrefix("/api/v1").Subrouter()

	// // subRouter.Use(CheckToken)
	// subRouter.Use(CORSHandler)

	//Default Route.
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		obj := []byte("this is home page")
		w.Write(obj)
	}).Methods("GET", "OPTION")

	//Routes.
	card := controllers.Card{Client: client, AppCfg: appCfg}
	// subRouter.HandleFunc("/", CheckToken(card.GetAllCards)).Methods("GET", "OPTION")
	subRouter.HandleFunc("/", card.GetAllCards).Methods("GET", "OPTION")
	subRouter.HandleFunc("/", card.CreateCard).Methods("POST", "OPTION")

	auth := controllers.Auth{Client: client, AppCfg: appCfg}
	subRouter.HandleFunc("/signUp", auth.Signup).Methods("POST", "OPTION")
	subRouter.HandleFunc("/logIn", auth.Login).Methods("POST", "OPTION")

	wishList := controllers.WishList{Client: client, AppCfg: appCfg}
	subRouter.HandleFunc("/Wishlist/{emailId}", wishList.SetWishlist).Methods("PUT", "OPTION")
	subRouter.HandleFunc("/Wishlist/delete/{emailId}", wishList.DeleteWishlist).Methods("PUT", "OPTION")

	myLearning := controllers.MyLearning{Client: client, AppCfg: appCfg}
	subRouter.HandleFunc("/MyLearning/{emailId}", myLearning.SetMyLearning).Methods("PUT", "OPTION")
	subRouter.HandleFunc("/MyLearning/delete/{emailId}", myLearning.DeleteMyLearning).Methods("PUT", "OPTION")

	cartList := controllers.CartList{Client: client, AppCfg: appCfg}
	subRouter.HandleFunc("/Cartlist/{emailId}", cartList.SetCartlist).Methods("PUT", "OPTION")
	subRouter.HandleFunc("/Cartlist/delete/{emailId}", cartList.DeleteCartlist).Methods("PUT", "OPTION")
	subRouter.HandleFunc("/Cartlist/delete/all/{emailId}", cartList.DeleteAllCartlist).Methods("PUT", "OPTION")

	//Server listening at Port.
	Url := fmt.Sprintf(":%d", *PORT)
	fmt.Printf("Server Started at port %d\n", *PORT)
	http.ListenAndServe(Url, r)
}
