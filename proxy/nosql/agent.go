package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy"
	"time"
)

type Agent struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Creator  string `json:"creator" bson:"creator"`
	Operator string `json:"operator" bson:"operator"`
	Scene    string `json:"scene" bson:"scene"`
	Remark   string `json:"remark" bson:"remark"`
	Quotes   []string `json:"quotes" bson::"quotes"`
	//Displays  []proxy.DisplayInfo `json:"displays" bson:"displays"`
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

func GetAgentByID(id uint64) (*Agent, error) {
	msg := bson.M{"id": id}
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

func RemoveAgent(uid,operator string) error {
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

func GetAgentsByScene(scene string) ([]*Agent, error) {
	cursor, err1 := findMany(TableAgent, bson.M{"scene": scene, "deleteAt": new(time.Time)}, 0)
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

func UpdateAgentDisplays(uid, operator string, list []*proxy.DisplayInfo) error {
	msg := bson.M{"displays": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableAgent, uid, msg)
	return err
}

func UpdateAgentQuotes(uid, operator string, arr []string) error {
	msg := bson.M{"operator": operator, "quotes": arr, "updatedAt": time.Now()}
	_, err := updateOne(TableAgent, uid, msg)
	return err
}


