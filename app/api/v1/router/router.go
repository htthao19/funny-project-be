package router

import (
	"github.com/beego/beego"

	"funny-project-be/app/api/v1/controller"
	"funny-project-be/domain/repo"
	"funny-project-be/infra/options"
)

// InitRouters initializes router of beego.
func InitRouters(
	uRepo repo.UserRepo,
	vRepo repo.VideoRepo,
	opts options.Options,
) {
	beego.AddNamespace(
		beego.NewNamespace("/funny-project",
			beego.NSNamespace("/v1",
				beego.NSNamespace("/rest",
					beego.NSNamespace("/users",
						beego.NSRouter("/me", &controller.UserController{BaseController: controller.BaseController{}, URepo: uRepo, Opts: opts}, "get:GetUser"),
					),
					beego.NSNamespace("/videos",
						beego.NSRouter("/:id", &controller.VideoController{BaseController: controller.BaseController{}, VRepo: vRepo, URepo: uRepo, Opts: opts}, "get:GetVideo"),
						beego.NSRouter("", &controller.VideoController{BaseController: controller.BaseController{}, VRepo: vRepo, URepo: uRepo, Opts: opts}, "post:CreateVideo"),
						beego.NSRouter("", &controller.VideoController{BaseController: controller.BaseController{}, VRepo: vRepo, URepo: uRepo, Opts: opts}, "get:ListVideos"),
					),
				),

				beego.NSNamespace("/rpc",
					beego.NSNamespace("/auth",
						beego.NSRouter("/login", &controller.UserController{BaseController: controller.BaseController{}, URepo: uRepo, Opts: opts}, "post:Login"),
					),
				),
			),
		),
	)
}
