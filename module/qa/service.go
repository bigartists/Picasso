package qa

import (
	"context"
	"fmt"
	"github.com/jianfengye/collection"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"picasso/client"
	"time"
)

var ServiceGetter Service

type ServiceGetterImpl struct {
}

func (s ServiceGetterImpl) GetQuestions(ctx context.Context, pager *Pager) ([]*Question, error) {
	questions := make([]*Question, 0, pager.Size)
	total := int64(0)
	if err := client.Orm.Count(&total).Error; err != nil {
		pager.Total = total
	}
	fmt.Println("service , GetQuestions pager=", pager.Total)
	if err := client.Orm.Debug().WithContext(ctx).Order("created_at desc").Offset(pager.Start).Limit(pager.Size).Find(&questions).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*Question{}, nil
		}
		return nil, err
	} else {
		fmt.Println("service2 , GetQuestions err=", err, questions)

	}
	return questions, nil
}

func (s ServiceGetterImpl) GetQuestion(ctx context.Context, questionID int64) (*Question, error) {
	question := &Question{ID: questionID}
	if err := client.Orm.Debug().WithContext(ctx).First(question).Error; err != nil {
		return nil, err
	}
	return question, nil
}

func (s ServiceGetterImpl) PostQuestion(ctx context.Context, question *Question) error {
	fmt.Println("service , PostQuestion=", question.AuthorID, question.Title, question.Context)
	if err := client.Orm.WithContext(ctx).Create(question).Error; err != nil {
		return err
	}
	return nil
}

func (s ServiceGetterImpl) QuestionLoadAuthor(ctx context.Context, question *Question) error {
	if err := client.Orm.Debug().WithContext(ctx).Preload("Author").First(question).Error; err != nil {
		return err
	}
	return nil
}

func (s ServiceGetterImpl) QuestionsLoadAuthor(ctx context.Context, questions *[]*Question) error {
	if err := client.Orm.WithContext(ctx).Preload("Author").Find(questions).Error; err != nil {
		return err
	}
	return nil
}

func (s ServiceGetterImpl) QuestionLoadAnswers(ctx context.Context, question *Question) error {
	if err := client.Orm.Debug().WithContext(ctx).Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("Answer.created_at desc")
	}).First(question).Error; err != nil {
		return err
	}
	return nil
}

func (s ServiceGetterImpl) QuestionsLoadAnswers(ctx context.Context, questions *[]*Question) error {
	//TODO implement me
	panic("implement me")
}

func (s ServiceGetterImpl) PostAnswer(ctx context.Context, answer *Answer) error {
	if answer.QuestionID == 0 {
		return errors.New("问题不存在")
	}
	question := &Question{ID: answer.QuestionID}
	if err := client.Orm.WithContext(ctx).First(question).Error; err != nil {
		return err
	}
	if err := client.Orm.WithContext(ctx).Create(answer).Error; err != nil {
		return err
	}
	question.AnswerNum = question.AnswerNum + 1
	if err := client.Orm.WithContext(ctx).Save(question).Error; err != nil {
		return err
	}
	return nil
}

func (s ServiceGetterImpl) GetAnswer(ctx context.Context, answerID int64) (*Answer, error) {
	answer := &Answer{ID: answerID}
	if err := client.Orm.WithContext(ctx).First(answer).Error; err != nil {
		return nil, err
	}
	return answer, nil
}

func (s ServiceGetterImpl) AnswerLoadAuthor(ctx context.Context, question *Answer) error {
	if err := client.Orm.WithContext(ctx).Preload("Author").First(question).Error; err != nil {
		return err
	}
	return nil
}

// todo 重点看看这个方法
func (s ServiceGetterImpl) AnswersLoadAuthor(ctx context.Context, answers *[]*Answer) error {
	if answers == nil {
		return nil
	}
	answerColl := collection.NewObjPointCollection(*answers)
	ids, err := answerColl.Pluck("ID").ToInt64s()
	fmt.Println("AnswersLoadAuthor ids=", ids)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	if err := client.Orm.Debug().WithContext(ctx).Preload("Author").Find(answers, ids).Error; err != nil {
		return err
	}
	return nil
}

// SELECT * FROM `Question` WHERE `Question`.`deleted_at` IS NULL AND `Question`.`id` = 3 ORDER BY `Question`.`id` LIMIT 1
func (s ServiceGetterImpl) DeleteQuestion(ctx context.Context, questionID int64) error {
	question := &Question{ID: questionID}
	if err := client.Orm.WithContext(ctx).Delete(question).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s ServiceGetterImpl) DeleteAnswer(ctx context.Context, answerID int64) error {
	answer := &Answer{ID: answerID}
	if err := client.Orm.WithContext(ctx).Delete(answer).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s ServiceGetterImpl) UpdateQuestion(ctx context.Context, question *Question) error {
	questionDB := &Question{ID: question.ID}
	if err := client.Orm.WithContext(ctx).First(questionDB).Error; err != nil {
		return errors.WithStack(err)
	}
	questionDB.UpdatedAt = time.Now()
	if question.Title != "" {
		questionDB.Title = question.Title
	}
	if question.Context != "" {
		questionDB.Context = question.Context
	}
	fmt.Println("questionDB=====", questionDB)
	// todo 注释部分有bug
	fmt.Println("questionDB.UpdatedAt=", questionDB.UpdatedAt, "questionDB.CreatedAt=", questionDB.CreatedAt, "questionDB.DeletedAt=", questionDB.DeletedAt)
	if err := client.Orm.Debug().WithContext(ctx).Save(questionDB).Error; err != nil {
		return errors.WithStack(err)
	}

	// 只更新指定字段
	//if err := client.Orm.Debug().WithContext(ctx).Select("Title", "Context", "UpdatedAt").Updates(questionDB).Error; err != nil {
	//	return errors.WithStack(err)
	//}
	return nil
}

func NewQaServiceGetterImpl() *ServiceGetterImpl {
	return &ServiceGetterImpl{}
}

func init() {
	ServiceGetter = NewQaServiceGetterImpl()
}
