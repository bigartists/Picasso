package qa

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"picasso/pkg/middlewares"
	. "picasso/pkg/utils"

	"github.com/spf13/cast"
)

type QaController struct {
}

func NewQaController() *QaController {
	return &QaController{}
}

type questionCreateParam struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binging:"required"`
}

func (this *QaController) QuestionList(c *gin.Context) {
	start := c.DefaultQuery("start", "0")
	size := c.DefaultQuery("size", "20")
	page := &Pager{
		Start: cast.ToInt(start),
		Size:  cast.ToInt(size),
	}

	//questions, err := ServiceGetter.
	questions, err := ServiceGetter.GetQuestions(c, page)
	if err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}

	if len(questions) == 0 {
		c.JSON(200, ResultWrapper(c)([]*QuestionDTO{}, "")(OK))
		return
	}

	if err := ServiceGetter.QuestionsLoadAuthor(c, &questions); err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}

	questionsDTO := ConvertQuestionsToDTO(questions)

	c.JSON(200, ResultWrapper(c)(questionsDTO, "")(OK))
}

/*
SELECT * FROM `Question` WHERE `Question`.`deleted_at` IS NULL AND `Question`.`id` = 1 ORDER BY `Question`.`id` LIMIT 1
SELECT * FROM `user` WHERE `user`.`id` = 1
SELECT * FROM `Question` WHERE `Question`.`deleted_at` IS NULL AND `Question`.`id` = 1 ORDER BY `Question`.`id` LIMIT 1
SELECT * FROM `Answer` WHERE `Answer`.`question_id` = 1 AND `Answer`.`deleted_at` IS NULL ORDER BY Answer.created_at desc
SELECT * FROM `Question` WHERE `Question`.`deleted_at` IS NULL AND `Question`.`id` = 1 ORDER BY `Question`.`id` LIMIT 1
SELECT * FROM `user` WHERE `user`.`id` = 1
SELECT * FROM `Answer` WHERE `Answer`.`id` IN (3,1) AND `Answer`.`deleted_at` IS NULL
*/
func (this *QaController) QuestionDetail(c *gin.Context) {

	id := cast.ToInt64(c.DefaultQuery("id", "0"))
	fmt.Println("id:", id)
	if id == 0 {
		c.JSON(400, ResultWrapper(c)(nil, "参数错误")(Error))
		return
	}

	question, err := ServiceGetter.GetQuestion(c, id)
	fmt.Println("question:", question)
	if err != nil {

		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}

	if err := ServiceGetter.QuestionLoadAuthor(c, question); err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}

	if err := ServiceGetter.QuestionLoadAnswers(c, question); err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}

	if err := ServiceGetter.AnswersLoadAuthor(c, &(question.Answers)); err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}

	questionDTO := ConvertQuestionToDTO(question, nil)

	c.JSON(200, ResultWrapper(c)(questionDTO, "")(OK))
}

func (this *QaController) QuestionCreate(c *gin.Context) {
	params := &questionCreateParam{}

	if err := c.ShouldBind(params); err != nil {

		c.JSON(400, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}

	user := middlewares.GetAuthUser(c)

	if user == nil {
		c.JSON(500, ResultWrapper(c)(nil, "无权限操作")(Error))
		return
	}

	fmt.Println("user:", user.Id, user.Username, user.Nickname, user.Email)

	question := &Question{
		Title:    params.Title,
		Context:  params.Content,
		AuthorID: user.Id,
	}
	err := ServiceGetter.PostQuestion(c, question)
	if err != nil {
		ret := ResultWrapper(c)(nil, err.Error())(Error)
		c.JSON(400, ret)
		return
	}
	c.JSON(200, ResultWrapper(c)(question, "")(OK))
}

func (this *QaController) QuestionDelete(c *gin.Context) {
	id := cast.ToInt64(c.DefaultQuery("id", "0"))
	if id == 0 {
		c.JSON(400, ResultWrapper(c)(nil, "参数错误")(Error))
		return
	}
	question, err := ServiceGetter.GetQuestion(c, id)
	if err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}
	if question == nil {
		c.JSON(500, ResultWrapper(c)(nil, "问题不存在")(Error))
		return
	}

	user := middlewares.GetAuthUser(c)
	if user.Id != question.AuthorID {
		c.JSON(500, ResultWrapper(c)(nil, "无权限操作")(Error))
		return
	}

	if err := ServiceGetter.DeleteQuestion(c, id); err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}
	c.JSON(200, ResultWrapper(c)(true, "删除成功")(OK))
}

type questionEditParam struct {
	ID      int64  `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (this *QaController) QuestionEdit(c *gin.Context) {
	params := &questionEditParam{}
	if err := c.ShouldBind(params); err != nil {
		c.JSON(400, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}

	questionOld, err := ServiceGetter.GetQuestion(c, params.ID)
	if err != nil || questionOld == nil {
		c.JSON(500, ResultWrapper(c)(nil, "问题不存在")(Error))
		return
	}
	user := middlewares.GetAuthUser(c)
	if user == nil || user.Id != questionOld.AuthorID {
		c.JSON(500, ResultWrapper(c)(nil, "无权限操作")(Error))
		return
	}
	fmt.Println("question:===========", params)
	question := &Question{
		ID:      params.ID,
		Title:   params.Title,
		Context: params.Content,
	}
	fmt.Println("question:===========", params)

	if err := ServiceGetter.UpdateQuestion(c, question); err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return

	}
	c.JSON(200, ResultWrapper(c)(true, "操作成功")(OK))

}

type answerCreateParam struct {
	QuestionID int64  `json:"question_id" binding:"required"`
	Context    string `json:"context" binding:"required"`
}

func (this *QaController) AnswerCreate(c *gin.Context) {
	params := &answerCreateParam{}
	if err := c.ShouldBind(params); err != nil {
		c.JSON(400, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}
	user := middlewares.GetAuthUser(c)
	if user == nil {
		c.JSON(500, ResultWrapper(c)(nil, "请登录后再操作")(Error))
		return
	}
	answer := &Answer{
		QuestionID: params.QuestionID,
		Context:    params.Context,
		AuthorID:   user.Id,
	}

	if err := ServiceGetter.PostAnswer(c, answer); err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}

	c.JSON(200, ResultWrapper(c)(true, "操作成功")(OK))
}

func (this *QaController) AnswerDelete(c *gin.Context) {
	id := cast.ToInt64(c.DefaultQuery("id", "0"))
	if id == 0 {
		c.JSON(400, ResultWrapper(c)(nil, "参数错误")(Error))
		return
	}
	user := middlewares.GetAuthUser(c)
	if user == nil {
		c.JSON(500, ResultWrapper(c)(nil, "请登录后再操作")(Error))
		return
	}
	answer, err := ServiceGetter.GetAnswer(c, id)
	if err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}
	if answer.AuthorID != user.Id {
		c.JSON(500, ResultWrapper(c)(nil, "没有权限做此操作")(Error))
		return
	}
	if err := ServiceGetter.DeleteAnswer(c, id); err != nil {
		c.JSON(500, ResultWrapper(c)(nil, err.Error())(Error))
		return
	}
	c.JSON(200, ResultWrapper(c)(true, "删除成功")(OK))
}

func (this *QaController) Build(r *gin.RouterGroup) {
	r.GET("/question/list", this.QuestionList)
	r.GET("/question/detail", this.QuestionDetail)
	r.POST("/question/create", this.QuestionCreate)
	r.POST("/question/delete", this.QuestionDelete)
	r.POST("/question/edit", this.QuestionEdit)
	r.POST("/answer/create", this.AnswerCreate)
	r.POST("/answer/delete", this.AnswerDelete)
}
