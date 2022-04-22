package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Agent struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Creator  string   `json:"creator" bson:"creator"`
	Operator string   `json:"operator" bson:"operator"`
	User     string   `json:"user" bson:"user"`
	Entity   string   `json:"entity" bson:"entity"`
	Remark   string   `json:"remark" bson:"remark"`
	Owner    string   `json:"owner" bson:"owner"`
	Way      string   `json:"way" bson:"way"`
	Type     uint8    `json:"type" bson:"type"`
	Status   uint8    `json:"status" bson:"status"`
	Regions  []string `json:"regions" bson:"regions"`
	Tags     []string `json:"tags" bson:"tags"`
}

func CreateAgent(info *Agent) error {
	_, err := insertOne(TableAgent, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAgentNextID() uint64 {
	num, _ := getSequenceNext(TableAgent)
	return num
}

func GetAgentCount() int64 {
	num, _ := getCount(TableAgent)
	return num
}

func GetAgent(uid string) (*Agent, error) {
	result, err := findOne(TableAgent, uid)
	if err != nil {
		return nil, err
	}
	model := new(Agent)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAgentByUser(user string) (*Agent, error) {
	msg := bson.M{"user": user}
	result, err := findOneBy(TableAgent, msg)
	if err != nil {
		return nil, err
	}
	model := new(Agent)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveAgent(uid, operator string) error {
	_, err := removeOne(TableAgent, uid, operator)
	return err
}

func GetAllAgents() ([]*Agent, error) {
	cursor, err1 := findAll(TableAgent, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Agent, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Agent)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAgentsByOwner(owner string) ([]*Agent, error) {
	cursor, err1 := findMany(TableAgent, bson.M{"owner": owner, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Agent, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Agent)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAgentsByRegion(region string) ([]*Agent, error) {
	msg := bson.M{"regions": bson.M{"$elemMatch": bson.M{"$eq": region}}, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableAgent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Agent, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Agent)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAgentsByWay(owner, way string) ([]*Agent, error) {
	cursor, err1 := findMany(TableAgent, bson.M{"owner": owner,"way": way, "deleteAt": new(time.Time)}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Agent, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Agent)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateAgentBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAgent, uid, msg)
	return err
}

func UpdateAgentEntity(uid, entity, operator string) error {
	msg := bson.M{"entity": entity, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAgent, uid, msg)
	return err
}

func UpdateAgentStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAgent, uid, msg)
	return err
}

func UpdateAgentTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAgent, uid, msg)
	return err
}

func UpdateAgentRegions(uid, operator string, list []string) error {
	msg := bson.M{"regions": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAgent, uid, msg)
	return err
}
