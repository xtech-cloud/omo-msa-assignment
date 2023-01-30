package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Meeting struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	StopTime  time.Time `json:"stopAt" bson:"stopAt"`
	StartTime time.Time `json:"startAt" bson:"startAt"`
	Creator   string    `json:"creator" bson:"creator"`
	Operator  string    `json:"operator" bson:"operator"`

	Status uint8  `json:"status" bson:"status"`
	Type   uint8  `json:"type" bson:"type"`
	Owner  string `json:"owner" bson:"owner"`
	/**
	所属组织或者部门
	*/
	Group  string `json:"group" bson:"group"`
	Remark string `json:"remark" bson:"remark"`

	/**
	预约时间
	*/
	Appointed string   `json:"appointed" bson:"appointed"`
	Location  string   `json:"location" bson:"location"`
	Signs     []string `json:"signs" bson:"signs"`
	Submits   []string `json:"submits" bson:"submits"`
	Notifies  []string `json:"notifies" bson:"notifies"`
}

func CreateMeeting(info *Meeting) error {
	_, err := insertOne(TableMeeting, info)
	if err != nil {
		return err
	}
	return nil
}

func GetMeetingNextID() uint64 {
	num, _ := getSequenceNext(TableMeeting)
	return num
}

func GetAllMeetings() ([]*Meeting, error) {
	cursor, err1 := findAll(TableMeeting, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Meeting, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(Meeting)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMeetingsByGroup(group string) ([]*Meeting, error) {
	msg := bson.M{"group": group, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableMeeting, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Meeting, 0, 30)
	for cursor.Next(context.Background()) {
		var node = new(Meeting)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMeetingsByScene(owner string) ([]*Meeting, error) {
	msg := bson.M{"owner": owner, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableMeeting, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Meeting, 0, 30)
	for cursor.Next(context.Background()) {
		var node = new(Meeting)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateMeetingBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableMeeting, uid, msg)
	return err
}

func UpdateMeetingLocation(uid, location, operator string, kind uint8) error {
	msg := bson.M{"type": kind, "location": location, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableMeeting, uid, msg)
	return err
}

func UpdateMeetingStatus(uid string, status uint16) error {
	msg := bson.M{"status": status, "updatedAt": time.Now()}
	_, err := updateOne(TableMeeting, uid, msg)
	return err
}

func StopMeeting(uid, operator string) error {
	msg := bson.M{"status": 3, "operator": operator, "stopAt": time.Now()}
	_, err := updateOne(TableMeeting, uid, msg)
	return err
}

func GetMeeting(uid string) (*Meeting, error) {
	result, err := findOne(TableMeeting, uid)
	if err != nil {
		return nil, err
	}
	model := new(Meeting)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveMeeting(uid, operator string) error {
	_, err := removeOne(TableMeeting, uid, operator)
	return err
}

func AppendMeetingSign(uid, member, operator string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"signs": member}
	_, err := appendElement(TableMeeting, uid, msg)
	return err
}

func AppendMeetingNotify(uid, member, operator string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"notifies": member}
	_, err := appendElement(TableMeeting, uid, msg)
	return err
}

func AppendMeetingSubmit(uid, member, operator string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"submits": member}
	_, err := appendElement(TableMeeting, uid, msg)
	return err
}
