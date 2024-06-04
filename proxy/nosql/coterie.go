package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"omo.msa.assignment/proxy"
	"time"
)

type Coterie struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name       string   `json:"name" bson:"name"`
	Cover      string   `json:"cover" bson:"cover"`
	Passwords  string   `json:"passwords" bson:"passwords"`
	Master     string   `json:"master" bson:"master"`
	Remark     string   `json:"remark" bson:"remark"`
	Type       uint8    `json:"type" bson:"type"`
	Status     uint8    `json:"status" bson:"status"`
	Centre     string   `json:"centre" bson:"centre"`
	Meta       string   `json:"meta" bson:"meta"`
	Assistants []string `json:"assistants" bson:"assistants"`
	Tags       []string `json:"tags" bson:"tags"`
	// 家庭成员，保存User的uid
	Members []proxy.MemberInfo `json:"members" bson:"members"`
}

func CreateCoterie(info *Coterie) error {
	_, err := insertOne(TableCoterie, info)
	if err != nil {
		return err
	}
	return nil
}

func GetCoterieNextID() uint64 {
	num, _ := getSequenceNext(TableCoterie)
	return num
}

func GetCoterie(uid string) (*Coterie, error) {
	result, err := findOne(TableCoterie, uid)
	if err != nil {
		return nil, err
	}
	model := new(Coterie)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCoterieByCreator(creator string) (*Coterie, error) {
	msg := bson.M{"creator": creator, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableCoterie, msg)
	if err != nil {
		return nil, err
	}
	model := new(Coterie)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCoterieByCentre(centre string) (*Coterie, error) {
	msg := bson.M{"centre": centre, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableCoterie, msg)
	if err != nil {
		return nil, err
	}
	model := new(Coterie)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetActivitiesCount() int64 {
	def := new(time.Time)
	filter := bson.M{"deleteAt": def}
	num, err1 := getCountBy(TableCoterie, filter)
	if err1 != nil {
		return num
	}

	return num
}

func GetAllCoteries(page, num int64) ([]*Coterie, error) {
	msg := bson.M{"deleteAt": new(time.Time)}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(num).SetSkip(page)
	cursor, err1 := findManyByOpts(TableCoterie, msg, opts)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Coterie, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Coterie)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCoteriesByMember(user string) ([]*Coterie, error) {
	msg := bson.M{"members": bson.M{"$elemMatch": bson.M{"user": user}}, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCoterie, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Coterie, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(Coterie)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCoteriesByCreator(user string) ([]*Coterie, error) {
	msg := bson.M{"creator": user, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCoterie, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Coterie, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(Coterie)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCoteriesByMaster(user string) ([]*Coterie, error) {
	msg := bson.M{"master": user, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCoterie, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Coterie, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(Coterie)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateCoterieBase(uid, name, remark, psw, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "passwords": psw, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCoterie, uid, msg)
	return err
}

func UpdateCoterieMaster(uid, master, operator string) error {
	msg := bson.M{"master": master, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCoterie, uid, msg)
	return err
}

func UpdateCoterieStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCoterie, uid, msg)
	return err
}

func UpdateCoteriePasswords(uid, operator, psw string) error {
	msg := bson.M{"passwords": psw, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCoterie, uid, msg)
	return err
}

func UpdateCoterieTags(uid, operator string, list []string) error {
	msg := bson.M{"tags": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCoterie, uid, msg)
	return err
}

func UpdateCoterieAssistants(uid, operator string, list []string) error {
	msg := bson.M{"assistants": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCoterie, uid, msg)
	return err
}

func UpdateCoterieCover(uid string, icon, operator string) error {
	msg := bson.M{"cover": icon, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableCoterie, uid, msg)
	return err
}

func RemoveCoterie(uid, operator string) error {
	_, err := removeOne(TableCoterie, uid, operator)
	return err
}

func AppendCoterieMember(uid string, invitee proxy.MemberInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": invitee}
	_, err := appendElement(TableCoterie, uid, msg)
	return err
}

func SubtractCoterieMember(uid, user string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": bson.M{"user": user}}
	_, err := removeElement(TableCoterie, uid, msg)
	return err
}
