package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	// "net/http"

	// "github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"campaigns/auth"
	"campaigns/campaign"
	"campaigns/handler"
	"campaigns/helper"
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
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)

	// cek pake service
	campaigns, _ := campaignService.FindCampaigns(0)

	fmt.Println(campaigns)

	// // mbil semua data
	// // campaigns, err := campaignRepository.FindAll()

	// // berdasar user

	// campaigns, err := campaignRepository.FindByUserid(1)
	// fmt.Println("debug")
	// fmt.Println(len(campaigns))

	// // print semua isi campaigns
	// for _, campaign := range campaigns {
	// 	fmt.Println(campaign.Name)
	// 	fmt.Println(campaign.CampaignImages)
	// }

	// token, validasi
	// pake middleware aja jgn manual gini
	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0fQ.XIQmRQxTNSxJo4fFWKdn1C5N3xxHl8reZoeBrNifXvI")

	// if err != err {
	// 	fmt.Println("ERRORRRRR")
	// }

	// if token.Valid {
	// 	fmt.Println("VALIDDD")
	// }else{
	// 	fmt.Println("INVALIDD")
	// }

	// fmt.Println(authService.GenerateToken(1001))

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "soffie"
	// userInput.Email = "soffie@gmail.com"
	// userInput.Occupation = "fullstack dev"
	// userInput.Password = "soffie"

	// userService.RegisterUser(userInput)

	// login
	// userByEmail, err := userRepository.FindByEmail("soffie@gmail.com")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// if userByEmail.ID == 0 {
	// 	fmt.Println("User tidak ditemukan")
	// } else {
	// 	fmt.Println(userByEmail.Name)
	// }

	// // login pake service
	// input := user.LoginInput{
	// 	Email: "soffie@gmail.com",
	// 	Password: "soffie",
	// }
	// user, err := userService.Login(input)
	// if err != nil {
	// 	fmt.Println("Terjadi kesalahan")
	// 	fmt.Println(err.Error())
	// }

	// fmt.Println(user.Email)
	// fmt.Println(user.Name)

	// input sesuai isidari FE (postman)
	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	// // update avatar
	// userService.SaveAvatar(1, "imags/1-profile.png")

	// register
	api.POST("/users", userHandler.RegisterUser)

	// login
	api.POST("login", userHandler.Login)

	// cek email sudah teraftar atau belum
	api.POST("/email_checkers", userHandler.CheckEmailAVailability)

	// update path avatar per-user
	api.POST("/avatar", authMiddleware(authService, userService), userHandler.UploadAvatar)

	// jwt

	// cek em
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

// middleware
// ambil nilai header authorization: Bearer tokennnn
// dari header authorization, ambil nilai tokenny saja
// validasi token, pake service validatetoken
// ambil userID
// ambil user dari db berdasar user_id lewat service
// set context isinya user

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ambil nilai header authorization: Bearer tokennnn
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			// abort = hentikan proses
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// dariheader authorization, ambil nilai tokenny saja
		// jadi isinya 2 buaharray
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		// validasi token, pake service validatetoken
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			// abort = hentikan proses
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			// abort = hentikan proses
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// kalau semua ok, ambil userid
		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			// abort = hentikan proses
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// kalau usernyna ada dan semuanya lancar
		// set context isinya user

		// lempar data user yang login ke handler
		c.Set("currentUser", user)
	}
}
