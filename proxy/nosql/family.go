package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy"
	"time"
)

type Family struct {
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
	SN         string   `json:"sn" bson:"sn"`
	Address    string   `json:"address" bson:"address"`
	Location   string   `json:"location" bson:"location"`
	Region     string   `json:"region" bson:"region"`
	Status     uint8    `json:"status" bson:"status"`
	Assistants []string `json:"assistants" bson:"assistants"`
	Tags       []string `json:"tags" bson:"tags"`
	Agents     []string `json:"agents" bson:"agents"`
	// 家庭小孩，没有User, 只有entity, 保存的是词条UID
	Children []string `json:"children" bson:"children"`
	// 家庭成员，保存User的uid
	Custodians []proxy.CustodianInfo `json:"custodians" bson:"custodians"`
	Members []proxy.MemberInfo `json:"members" bson:"members"`
}

func CreateFamily(info *Family) error {
	_, err := insertOne(TableFamily, info)
	if err != nil {
		return err
	}
	return nil
}

func GetFamilyNextID() uint64 {
	num, _ := getSequenceNext(TableFamily)
	return num
}

func GetInviterNextID() uint64 {
	num, _ := getSequenceNext("family_invitee")
	return num
}

func GetFamily(uid string) (*Family, error) {
	result, err := findOne(TableFamily, uid)
	if err != nil {
		return nil, err
	}
	model := new(Family)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetFamilyByCreator(creator string) (*Family, error) {
	msg := bson.M{"creator": creator}
	result, err := findOneBy(TableFamily, msg)
	if err != nil {
		return nil, err
	}
	model := new(Family)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetFamilyBySN(sn string) (*Family, error) {
	msg := bson.M{"sn": sn}
	result, err := findOneBy(TableFamily, msg)
	if err != nil {
		return nil, err
	}
	model := new(Family)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllFamilies() ([]*Family, error) {
	cursor, err1 := findAll(TableFamily, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Family, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Family)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetFamiliesByMember(user string) ([]*Family, error) {
	msg := bson.M{"members": bson.M{"$elemMatch": bson.M{"user": user}}}
	cursor, err1 := findMany(TableFamily, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Family, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(Family)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetFamilyByChild(entity string) (*Family, error) {
	msg := bson.M{"children": bson.M{"$elemMatch": bson.M{"$eq": entity}}}
	result, err1 := findOneBy(TableFamily, msg)
	if err1 != nil {
		return nil, err1
	}
	model := new(Family)
	err1 = result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateFamilyBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func UpdateFamilyMaster(uid, master, operator string) error {
	msg := bson.M{"master": master, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func UpdateFamilyStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func UpdateFamilyPasswords(uid, operator, psw string) error {
	msg := bson.M{"passwords": psw, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func UpdateFamilyTags(uid, operator string, list []string) error {
	msg := bson.M{"tags": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func UpdateFamilyAgents(uid, operator string, list []string) error {
	msg := bson.M{"agents": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func UpdateFamilyAssistants(uid, operator string, list []string) error {
	msg := bson.M{"assistants": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func UpdateFamilyChildren(uid, operator string, list []string) error {
	msg := bson.M{"children": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func UpdateFamilySN(uid, sn, operator string) error {
	msg := bson.M{"sn": sn, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func UpdateFamilyCover(uid string, icon, operator string) error {
	msg := bson.M{"cover": icon, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFamily, uid, msg)
	return err
}

func RemoveFamily(uid, operator string) error {
	_, err := removeOne(TableFamily, uid, operator)
	return err
}

func AppendFamilyChild(uid, child string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"children": child}
	_, err := appendElement(TableFamily, uid, msg)
	return err
}

func SubtractFamilyChild(uid, child string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"children": child}
	_, err := removeElement(TableFamily, uid, msg)
	return err
}

func AppendFamilyMember(uid string, invitee proxy.MemberInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": invitee}
	_, err := appendElement(TableFamily, uid, msg)
	return err
}

func SubtractFamilyMember(uid, user string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": bson.M{"user": user}}
	_, err := removeElement(TableFamily, uid, msg)
	return err
}
