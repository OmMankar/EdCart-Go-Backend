package models

// Card represents the schema for a course card
type Card struct {
	Category      string `bson:"category" json:"category" validate:"required,max=100"`         // Category with a max length of 100
	Image         string `bson:"image" json:"image" validate:"required,max=100"`               // Image field with a max length of 100
	CourseName    string `bson:"courseName" json:"courseName" validate:"required,max=240"`     // Course name with a max length of 240
	CourseAuthor  string `bson:"courseAuthor" json:"courseAuthor" validate:"required,max=100"` // Course author with a max length of 100
	RatingNumber  string `bson:"ratingNumber" json:"ratingNumber" validate:"required,max=100"` // Rating number with a max length of 100
	NumOfRatings  string `bson:"numOfRatings" json:"numOfRatings" validate:"required"`         // Number of ratings field
	DiscountPrice string `bson:"discountPrice" json:"discountPrice" validate:"required"`       // Discounted price field
	OriginalPrice string `bson:"originalPrice" json:"originalPrice" validate:"required"`       // Original price field
}
