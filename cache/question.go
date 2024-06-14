package cache

import (
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy"
	"omo.msa.assignment/proxy/nosql"
	"time"
)

type QuestionInfo struct {
	baseInfo
	Quote    string
	Remark   string
	Cd       int
	Category string
	Answers  []uint32
	Assets   []string
	Options  []proxy.PairInfo
}

func (mine *cacheContext) NewQuestion(title, remark, category, entity, operator string, cd int, answers []uint32, options []*pb.QuestionOption) (*QuestionInfo, error) {
	db := new(nosql.Question)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetQuestionNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Title = title
	db.Remark = remark
	db.Quote = entity
	db.Cd = uint16(cd)
	db.Answers = answers
	db.Assets = make([]string, 0, 1)
	db.Category = category
	db.Options = make([]proxy.PairInfo, 0, len(options))
	for _, v := range options {
		db.Options = append(db.Options, proxy.PairInfo{
			Key:   fmt.Sprintf("%d", v.Id),
			Value: v.Desc,
		})
	}
	err := nosql.CreateQuestion(db)
	if err != nil {
		return nil, err
	}
	info := new(QuestionInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetQuestionsByNameAndKind(title, category string) ([]*QuestionInfo, error) {
	if title == "" {
		return nil, errors.New("the parent is null")
	}
	array, err := nosql.GetQuestionsByTitle(title, category)
	if err != nil {
		return nil, err
	}
	questions := make([]*QuestionInfo, 0, len(array))
	for _, v := range array {
		question := new(QuestionInfo)
		question.initInfo(v)
		questions = append(questions, question)
	}
	return questions, nil
}

func (mine *cacheContext) GetQuestion(uid string) (*QuestionInfo, error) {
	db, err := nosql.GetQuestion(uid)
	if err == nil {
		info := new(QuestionInfo)
		info.initInfo(db)
		return info, err
	}
	return nil, err
}

func (mine *cacheContext) GetQuestionsByName(title string) ([]*QuestionInfo, error) {
	dbs, err := nosql.GetQuestionsByName(title)
	list := make([]*QuestionInfo, 0, 100)
	if err != nil {
		return nil, err
	}
	for _, question := range dbs {
		info := new(QuestionInfo)
		info.initInfo(question)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetQuestionsByCategory(kind string) ([]*QuestionInfo, error) {
	list := make([]*QuestionInfo, 0, 100)
	array, err := nosql.GetQuestionsByCategory(kind)
	if err != nil {
		return nil, err
	}
	for _, question := range array {
		info := new(QuestionInfo)
		info.initInfo(question)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetQuestionsByEntity(entity string) ([]*QuestionInfo, error) {
	list := make([]*QuestionInfo, 0, 100)
	array, err := nosql.GetQuestionsByQuote(entity)
	if err != nil {
		return nil, err
	}
	for _, question := range array {
		info := new(QuestionInfo)
		info.initInfo(question)
		list = append(list, info)
	}
	return list, nil
}

func (mine *QuestionInfo) initInfo(db *nosql.Question) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Name = db.Title
	mine.Remark = db.Remark
	mine.Cd = int(db.Cd)
	mine.Quote = db.Quote
	mine.Answers = db.Answers
	mine.Assets = db.Assets
	mine.Options = db.Options
	mine.Category = db.Category
}

func (mine *QuestionInfo) UpdateAnswers(operator string, answers []uint32) error {
	err := nosql.UpdateQuestionAnswers(mine.UID, operator, answers)
	if err == nil {
		mine.Answers = answers
		mine.Operator = operator
	}
	return err
}

func (mine *QuestionInfo) UpdateAssets(operator string, arr []string) error {
	if arr == nil {
		arr = make([]string, 0, 1)
	}
	err := nosql.UpdateQuestionAssets(mine.UID, operator, arr)
	if err == nil {
		mine.Assets = arr
		mine.Operator = operator
	}
	return err
}

func (mine *QuestionInfo) UpdateOptions(operator string, lis []proxy.PairInfo) error {
	err := nosql.UpdateQuestionOptions(mine.UID, operator, lis)
	if err == nil {
		mine.Options = lis
		mine.Operator = operator
	}
	return err
}
func (mine *QuestionInfo) Delete(uid string) error {
	err := nosql.RemoveQuestion(mine.UID, uid)
	if err == nil {
		return err
	}
	return err
}
func (mine *QuestionInfo) UpdateBase(title, remark, operator, category string, cd uint16, answers []uint32, opts []*pb.QuestionOption) error {
	arr := make([]proxy.PairInfo, 0, len(opts))
	for _, v := range opts {
		arr = append(arr, proxy.PairInfo{
			Key:   fmt.Sprintf("%d", v.Id),
			Value: v.Desc,
		})
	}
	err := nosql.UpdateQuestionBase(mine.UID, title, remark, operator, category, cd, answers, arr)
	if err == nil {
		mine.Name = title
		mine.Remark = remark
		mine.Cd = int(cd)
		mine.Category = category
		mine.Answers = answers
		mine.Options = arr
	}
	return err
}
