package models

type User struct {
	Name           string `bson:"name" json:"name" `
	EmailId        string `bson:"emailId" json:"emailId" validation:"required"`
	Password       string `bson:"password" json:"password" validation:"required"`
	WishList       []Card `bson:"wishList" json:"wishList" `
	CartList       []Card `bson:"cartList" json:"cartList" `
	MyLearningList []Card `bson:"myLearningList" json:"myLearningList" `
}
