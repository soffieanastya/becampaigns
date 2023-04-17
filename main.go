package main

import (
	// "fmt"
	"log"
	// "net/http"

	// "github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"campaigns/handler"
	"campaigns/user"
)

func main() {
	// koneksi db, pake mysql. phpmyadmin
	dsn := "root:@tcp(127.0.0.1:3306)/webcampaign?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	
	// userInput := user.RegisterUserInput{}
	// userInput.Name = "soffie"
	// userInput.Email = "soffie@gmail.com"
	// userInput.Occupation = "fullstack dev"
	// userInput.Password = "soffie"
	
	// userService.RegisterUser(userInput)
	
	// input sesuai isidari FE (postman)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default() 
	api := router.Group("/api/v1")

	api.POST("/users",userHandler.RegisterUser)

	router.Run(":5000")
	// user := user.User{
	// 	Name: "Test simoan",
	// }
	// userRepository.Save(user)





	// fmt.Println("connection to database is good")

	// var users []user.User	// sediain variabel yang udah di deskripsi di struct user
	// db.Find(&users)		// data yang ada di tabel users ini otomatis bakal disimpen di variabel users di atas
	// // jadi kalau di db plural, di package nya ga plural jd otomatis bisa kebaca tuh isi di db tabel users

	// // length := len(users)
	// // fmt.Println(users)

	// for _, user := range users {
	// 	fmt.Println(user.Name)
	// 	fmt.Println(user.Email)
	// }

	// router := gin.Default()
	// router.GET("/",userHandler)
	// router.Run(":5000")
}

// func userHandler(c *gin.Context){
// 	// koneksi db, pake mysql. phpmyadmin
// 	dsn := "root:@tcp(127.0.0.1:3306)/webcampaign?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
	
// 	var users []user.User
// 	db.Find(&users)

// 	// tampilkan ke web isi dari tabel users di db campaigns
// 	c.JSON(http.StatusOK, users)
// }

// input dai user
// handler, mapping inputan ke struct nput
// service, mapping dari struct input ke struct users
// repository