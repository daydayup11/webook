package web

import (
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/m/internal/domain"
	"webook/m/internal/service"
)

var UserIdKey = "userId"

type UserHandler struct {
	EmailRegexp    *regexp2.Regexp
	PasswordRegexp *regexp2.Regexp
	svc            *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	emailRegex := regexp2.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, regexp2.None)
	passwordRegex := regexp2.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`, regexp2.None)
	return &UserHandler{
		EmailRegexp:    emailRegex,
		PasswordRegexp: passwordRegex,
		svc:            userService,
	}
}

func (u *UserHandler) RegisterRouter(ug *gin.RouterGroup) {
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}
func (u *UserHandler) RegisterRouter1(server *gin.Engine) {
	ug := server.Group("/user")
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
}
func (u *UserHandler) Profile(ctx *gin.Context) {
	type Profile struct {
		Email string
	}
	session := sessions.Default(ctx)
	id := session.Get(UserIdKey)
	user, err := u.svc.Profile(ctx, id.(int64))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误，查询失败")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: user.Email,
	})
}
func (u *UserHandler) Signup(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	// 当我们调用 Bind 方法的时候，如果有问题，Bind 方法已经直接写响应回去了
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//传指针才能修改值
	isEmail, err := u.EmailRegexp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱错误")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不相同")
		return
	}

	isPassword, err := u.PasswordRegexp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK,
			"密码必须包含数字、特殊字符，并且长度不能小于 8 位")
		return
	}
	err = u.svc.Signup(ctx.Request.Context(),
		domain.User{Email: req.Email, Password: req.ConfirmPassword})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "重复邮箱，请换一个邮箱")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "服务器异常，注册失败")
		return
	}
	ctx.String(http.StatusOK, "hello, 注册成功")
}
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req LoginReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不正确")
		return
	}
	//设置session
	session := sessions.Default(ctx)
	session.Set(UserIdKey, user.Id)
	if err := session.Save(); err != nil {
		ctx.String(http.StatusOK, "服务器异常")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req EditReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	isEmail, err := u.EmailRegexp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱错误")
		return
	}

	id := sessions.Default(ctx).Get(UserIdKey)
	user, err := u.svc.Edit(ctx, domain.User{
		Id:       id.(int64),
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		ctx.Error(err)
	}
	ctx.JSON(http.StatusOK, user)
}
