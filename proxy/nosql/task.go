package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy"
	"time"
)

type Task struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name   string `json:"name" bson:"name"`
	Type   uint8  `json:"type" bson:"type"`
	Status uint8  `json:"status" bson:"status"`
	Remark string `json:"remark" bson:"remark"`
	Target string `json:"target" bson:"target"`

	Owner     string             `json:"owner" bson:"owner"`
	Way       string             `json:"way" bson:"way"`
	Duration  proxy.DateInfo     `json:"duration" bson:"duration"`
	Executors []string           `json:"executors" bson:"executors"`
	PreTasks  []string           `json:"preTasks" bson:"preTasks"`
	Regions   []string           `json:"regions" bson:"regions"`
	Tags      []string           `json:"tags" bson:"tags"`
	Assets    []string           `json:"assets" bson:"assets"`
	Records   []proxy.RecordInfo `json:"records" bson:"records"`
}

func CreateTask(info *Task) error {
	_, err := insertOne(TableTask, info)
	if err != nil {
		return err
	}
	return nil
}

func GetTaskNextID() uint64 {
	num, _ := getSequenceNext(TableTask)
	return num
}

func GetTask(uid string) (*Task, error) {
	result, err := findOne(TableTask, uid)
	if err != nil {
		return nil, err
	}
	model := new(Task)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTasksByOwner(uid string, st uint8) ([]*Task, error) {
	msg := bson.M{"owner": uid, "status":st, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTask, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Task, 0, 100)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTasksByOwner2(uid string) ([]*Task, error) {
	msg := bson.M{"owner": uid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTask, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Task, 0, 100)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTasksByType(owner string, tp uint8) ([]*Task, error) {
	msg := bson.M{"owner": owner, "type": tp, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTask, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Task, 0, 100)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTasksByRegion(region string, st uint8) ([]*Task, error) {
	msg := bson.M{"regions": bson.M{"$elemMatch": bson.M{"$eq": region}}, "status":st, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTask, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Task, 0, 100)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTasksByRegion2(region string) ([]*Task, error) {
	msg := bson.M{"regions": bson.M{"$elemMatch": bson.M{"$eq": region}}, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTask, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Task, 0, 100)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTasksByAgent(agent string, st uint8) ([]*Task, error) {
	msg := bson.M{"executors": bson.M{"$elemMatch": bson.M{"$eq": agent}}, "status":st, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTask, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Task, 0, 100)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTasksByAgent2(agent string) ([]*Task, error) {
	msg := bson.M{"executors": bson.M{"$elemMatch": bson.M{"$eq": agent}}, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTask, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Task, 0, 100)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTasksByTarget(client string, st uint8) ([]*Task, error) {
	msg := bson.M{"target": client, "status":st, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTask, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Task, 0, 100)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTasksByTarget2(client string) ([]*Task, error) {
	msg := bson.M{"target": client, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTask, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Task, 0, 100)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllTasks() ([]*Task, error) {
	var items = make([]*Task, 0, 100)
	cursor, err1 := findAll(TableTask, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Task)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateTaskBase(uid, name, remark, operator string, assets []string) error {
	msg := bson.M{"name": name, "remark": remark, "assets":assets, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskExecutors(uid, operator string, list []string) error {
	msg := bson.M{"executors": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskType(uid, operator string, tp uint8) error {
	msg := bson.M{"type": tp, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskStatus(uid string, status uint8, operator string) error {
	msg := bson.M{"status": status, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func RemoveTask(uid, operator string) error {
	_, err := removeOne(TableTask, uid, operator)
	return err
}

func AppendTaskRecord(uid string, data proxy.RecordInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"records": data}
	_, err := appendElement(TableTask, uid, msg)
	return err
}

func SubtractTaskRecord(uid, record string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"records.uid": record}
	_, err := removeElement(TableTask, uid, msg)
	return err
}

func AppendTaskExecutor(uid, user string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"executors": user}
	_, err := appendElement(TableTask, uid, msg)
	return err
}

func SubtractTaskExecutor(uid, user string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"executors": user}
	_, err := removeElement(TableTask, uid, msg)
	return err
}
