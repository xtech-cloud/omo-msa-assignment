package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Team struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Creator    string   `json:"creator" bson:"creator"`
	Operator   string   `json:"operator" bson:"operator"`
	Owner      string   `json:"owner" bson:"owner"`
	Remark     string   `json:"remark" bson:"remark"`
	Status     uint8    `json:"status" bson:"status"`
	Region     string   `json:"region" bson:"region"`
	Master     string   `json:"master" bson:"master"`
	Assistants []string `json:"assistants" bson:"assistants"`
	Tags       []string `json:"tags" bson:"tags"`
	Members    []string `json:"members" bson:"members"`
}

func CreateTeam(info *Team) error {
	_, err := insertOne(TableTeam, info)
	if err != nil {
		return err
	}
	return nil
}

func GetTeamNextID() uint64 {
	num, _ := getSequenceNext(TableTeam)
	return num
}

func GetTeamCount() int64 {
	num, _ := getCount(TableTeam)
	return num
}

func GetTeam(uid string) (*Team, error) {
	result, err := findOne(TableTeam, uid)
	if err != nil {
		return nil, err
	}
	model := new(Team)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTeamByID(id uint64) (*Team, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableTeam, msg)
	if err != nil {
		return nil, err
	}
	model := new(Team)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTeamByName(owner, name string) (*Team, error) {
	msg := bson.M{"owner": owner, "name":name, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableTeam, msg)
	if err != nil {
		return nil, err
	}
	model := new(Team)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveTeam(uid, operator string) error {
	_, err := removeOne(TableTeam, uid, operator)
	return err
}

func GetAllTeams() ([]*Team, error) {
	cursor, err1 := findAll(TableTeam, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Team, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Team)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTeamsByOwner(owner string) ([]*Team, error) {
	cursor, err1 := findMany(TableTeam, bson.M{"owner": owner, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Team, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Team)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTeamsByMember(user string) ([]*Team, error) {
	msg := bson.M{"members": bson.M{"$elemMatch": bson.M{"$eq": user}}}
	cursor, err1 := findMany(TableTeam, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Team, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(Team)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateTeamBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func UpdateTeamAssistants(uid, operator string, list []string) error {
	msg := bson.M{"operator": operator, "assistants": list, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func UpdateTeamTags(uid, operator string, list []string) error {
	msg := bson.M{"operator": operator, "tags": list, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func UpdateTeamStatus(uid, operator string, st uint8) error {
	msg := bson.M{"operator": operator, "status": st, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func UpdateTeamRegion(uid, region, operator string) error {
	msg := bson.M{"operator": operator, "region": region, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func UpdateTeamMembers(uid, operator string, members []string) error {
	msg := bson.M{"operator": operator, "members": members, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func UpdateTeamMaster(uid, member, operator string) error {
	msg := bson.M{"master": member, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func AppendTeamMember(uid, member string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := appendElement(TableTeam, uid, msg)
	return err
}

func SubtractTeamMember(uid string, member string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := removeElement(TableTeam, uid, msg)
	return err
}
