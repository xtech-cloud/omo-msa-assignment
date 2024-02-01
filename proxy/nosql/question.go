package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy"
	"time"
)

type Question struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Title    string           `json:"title" bson:"title"`
	Remark   string           `json:"remark" bson:"remark"`
	Cd       uint16           `json:"cd" bson:"cd"`
	Category string           `json:"category" bson:"category"`
	Quote    string           `json:"quote" bson:"quote"`
	Answers  []uint32         `json:"answers" bson:"answers"`
	Options  []proxy.PairInfo `json:"options" bson:"options"`
}

func CreateQuestion(info *Question) error {
	_, err := insertOne(TableQuestion, info)
	if err != nil {
		return err
	}
	return nil
}
func GetQuestionNextID() uint64 {
	num, _ := getSequenceNext(TableQuestion)
	return num
}

func GetQuestion(uid string) (*Question, error) {
	result, err := findOne(TableQuestion, uid)
	if err != nil {
		return nil, err
	}
	model := new(Question)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetQuestionsByName(title string) ([]*Question, error) {
	msg := bson.M{"title": title, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableQuestion, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Question, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Question)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetQuestionsByQuote(quote string) ([]*Question, error) {
	msg := bson.M{"quote": quote, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableQuestion, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Question, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Question)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetQuestionsByTitle(title, category string) ([]*Question, error) {
	msg := bson.M{"title": title, "category": category, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableQuestion, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Question, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Question)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetQuestionsByCategory(category string) ([]*Question, error) {
	msg := bson.M{"category": category, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableQuestion, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Question, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Question)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

//func GetAllQuestions() ([]*Question, error) {
//	var items = make([]*Question, 0, 10)
//	cursor, err1 := findAll(TableQuestion, 0)
//	if err1 != nil {
//		return nil, err1
//	}
//	defer cursor.Close(context.Background())
//	for cursor.Next(context.Background()) {
//		var node = new(Question)
//		if err := cursor.Decode(node); err != nil {
//			return nil, err
//		} else {
//			items = append(items, node)
//		}
//	}
//	return items, nil
//}

func UpdateQuestionBase(uid, title, remark, operator, category string, cd uint16, answers []uint32, opts []proxy.PairInfo) error {
	msg := bson.M{"title": title, "remark": remark, "cd": cd, "category": category,
		"answers": answers, "options": opts, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableQuestion, uid, msg)
	return err
}

func UpdateQuestionAnswers(uid, operator string, answers []uint32) error {
	msg := bson.M{"answers": answers, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableQuestion, uid, msg)
	return err
}

func UpdateQuestionOptions(uid, operator string, list []proxy.PairInfo) error {
	msg := bson.M{"options": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableQuestion, uid, msg)
	return err
}

func RemoveQuestion(uid, operator string) error {
	_, err := removeOne(TableQuestion, uid, operator)
	return err
}

//func GetQuestionByKV(k, v string) ([]*Question, error) {
//	msg := bson.M{k: v, "deleteAt": new(time.Time)}
//	cursor, err1 := findMany(TableQuestion, msg, 0)
//	if err1 != nil {
//		return nil, err1
//	}
//	defer cursor.Close(context.Background())
//	var items = make([]*Question, 0, 100)
//	for cursor.Next(context.Background()) {
//		var node = new(Question)
//		if err := cursor.Decode(node); err != nil {
//			return nil, err
//		} else {
//			items = append(items, node)
//		}
//	}
//	return items, nil
//}
