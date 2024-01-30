package main

import (
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"webook/m/internal/repository"
	"webook/m/internal/repository/dao"
	"webook/m/internal/service"
	"webook/m/internal/web"
	"webook/m/internal/web/middleware"
)

func main() {
	db := initDB()
	userDAO := dao.NewUserDAO(db)
	userRepository := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepository)
	svc := web.NewUserHandler(userService)
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "mumu.com")
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type"},
	}))
	store := sessions.NewCookieStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePath("/user/signup").IgnorePath("/user/login").Build())
	svc.RegisterRouter(server.Group("user"))
	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:dyt15918781398@tcp(111.229.199.247:3306)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
