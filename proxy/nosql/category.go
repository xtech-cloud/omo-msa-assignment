package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Category struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name   string `json:"name" bson:"name"`
	Parent string `json:"parent" bson:"parent"`
	Remark string `json:"remark" bson:"remark"`
	Weight uint32 `json:"weight" bson:"weight"`
	Quote  string `json:"quote" bson:"quote"`
	Owner  string `json:"owner" bson:"owner"`
}

func CreateCategory(info *Category) error {
	_, err := insertOne(TableCategory, info)
	if err != nil {
		return err
	}
	return nil
}
func UpdateCategoryStr(filter string, value, operator, uid string) error {
	msg := bson.M{filter: value, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCategory, uid, msg)
	return err
}

func UpdateCategoryOwner(uid, owner string) error {
	msg := bson.M{"owner": owner, "updatedAt": time.Now()}
	_, err := updateOne(TableCategory, uid, msg)
	return err
}

func UpdateCategoryBase(uid, name, remark, quote, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "quote": quote, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCategory, uid, msg)
	return err
}
func UpdateCategoryInt(filter, operator, uid string, value int64) error {
	msg := bson.M{filter: value, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCategory, uid, msg)
	return err
}
func GetCategoryListByWeight(weight int64) ([]*Category, error) {
	msg := bson.M{"weight": weight, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCategory, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Category, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Category)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
func GetAllCategories() ([]*Category, error) {
	msg := bson.M{"deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCategory, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Category, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Category)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
func GetCategoryListByTypeAndParent(parent string) ([]*Category, error) {
	msg := bson.M{"parent": parent, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCategory, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Category, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Category)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
func GetCategoryListByParent(parent string) ([]*Category, error) {
	msg := bson.M{"parent": parent, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCategory, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Category, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Category)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCategoryListByOwner(owner, parent string) ([]*Category, error) {
	msg := bson.M{"owner": owner, "parent": parent, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCategory, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Category, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Category)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

//带着大题库uid列表
func GetCategoryListByParentTwo(parent string, list []string) ([]*Category, error) {
	msg := bson.M{"parent": parent, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCategory, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Category, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Category)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
func GetCategory() ([]*Category, error) {
	var items = make([]*Category, 0, 10)
	cursor, err1 := findAll(TableCategory, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Category)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
func GetOneCategory(uid string) (*Category, error) {
	result, err := findOne(TableCategory, uid)
	if err != nil {
		return nil, err
	}
	model := new(Category)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}
func DeleteCategory(uid, operator string) error {
	_, err := removeOne(TableCategory, uid, operator)
	return err
}
