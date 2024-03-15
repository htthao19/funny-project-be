package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beego/beego"
	"github.com/beego/beego/plugins/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"funny-project-be/app/api/v1/router"
	"funny-project-be/domain/entity"
	"funny-project-be/infra/beego/plugin/authn"
	"funny-project-be/infra/options"
	"funny-project-be/infra/repo/repoimpl"
)

func main() {
	// Load config.
	if err := beego.LoadAppConfig("ini", "config/app.conf"); err != nil {
		log.Fatal(err)
	}
	opts, err := options.Load()
	if err != nil {
		log.Fatal(err)
	}

	var connectStr = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		opts.DBUser,
		opts.DBPass,
		opts.DBHost,
		opts.DBPort,
		opts.DBName,
	)

	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		NowFunc:        func() time.Time { return time.Now().Local() },
	})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema.
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Video{})

	uRepo := repoimpl.NewUserRepo(db)
	vRepo := repoimpl.NewVideoRepo(db)

	router.InitRouters(uRepo, vRepo, opts)

	// cors plugin
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type", "X-Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	beego.InsertFilter("*", beego.BeforeRouter, authn.VerifyToken(opts))
	beego.BConfig.WebConfig.AutoRender = false

	beego.Run()
}
