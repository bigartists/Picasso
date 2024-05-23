package user

import (
	"github.com/gin-gonic/gin"
	"picasso/config"
	"picasso/pkg/result"
	. "picasso/pkg/utils"
	"time"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (this *UserController) Login(c *gin.Context) {
	// 校验输入参数是否合法
	params := &struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}
	// 校验参数
	result.Result(c.ShouldBindJSON(params)).Unwrap()

	user, err := ServiceGetter.SignIn(params.Username, params.Password)
	if err != nil {
		ret := ResultWrapper(c)(nil, err.Error())(Error)
		c.JSON(400, ret)
		return
	}

	//// 生成 token
	prikey := []byte(config.SysYamlconfig.Jwt.PrivateKey)
	curTime := time.Now().Add(time.Second * 60 * 60 * 24)
	token, _ := GenerateToken(user.ID, prikey, curTime)

	c.Set("token", token)
	ret := ResultWrapper(c)(user, "")(OK)
	c.JSON(200, ret)
}

//func (this *UserController) SignUp(c *gin.Context) {
//	// 校验输入参数是否合法
//	params := &dto.SignupRequest{}
//	// 校验参数
//	result.Result(c.ShouldBindJSON(params)).Unwrap()
//
//	err := ServiceGetter.SignUp(params.Email, params.Username, params.Password)
//	if err != nil {
//		ret := ResultWrapper(c)(nil, err.Error())(Error)
//		c.JSON(400, ret)
//	}
//	ret := ResultWrapper(c)(true, "")(Created)
//	c.JSON(201, ret)
//}

func (this *UserController) UserList(c *gin.Context) {
	ret := ResultWrapper(c)(ServiceGetter.GetUserList(), "")(OK)
	c.JSON(200, ret)
}

func (this *UserController) UserDetail(c *gin.Context) {
	id := &struct {
		Id int64 `uri:"id" binding:"required"`
	}{}
	result.Result(c.ShouldBindUri(id)).Unwrap()
	ret := ResultWrapper(c)(ServiceGetter.GetUserDetail(id.Id).Unwrap(), "")(OK)
	c.JSON(200, ret)
}

func (this *UserController) Build(r *gin.RouterGroup) {
	r.POST("/login", this.Login)
	//r.POST("/register", this.SignUp)
	r.GET("/users", this.UserList)
	r.GET("/user/:id", this.UserDetail)
}
