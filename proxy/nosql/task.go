package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy"
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

	Name        string              `json:"name" bson:"name"`
	Type        uint8               `json:"type" bson:"type"`
	Status      uint8               `json:"status" bson:"status"`
	Short       string 				`json:"short" bson:"short"`
	Cover       string              `json:"cover" bson:"cover"`
	Master      string              `json:"master" bson:"master"`
	Remark      string              `json:"remark" bson:"remark"`
	Entity      string              `json:"entity" bson:"entity"`
	Location    string              `json:"location" bson:"location"`
	Supporter   string              `json:"supporter" bson:"supporter"`
	Bucket      string `json:"bucket" bson:"bucket"`
	Address     AddressInfo         `json:"address" bson:"address"`
	//Exhibitions []string            `json:"exhibitions" bson:"exhibitions"`
	//Displays    []proxy.ShowingInfo `json:"displays" bson:"displays"`
	Members     []string            `json:"members" bson:"members"`
	Parents     []string            `json:"parents" bson:"parents"`
	Domains     []proxy.DomainInfo 	`json:"domains" bson:"domains"`
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

func GetTaskByMaster(user string) (*Task, error) {
	msg := bson.M{"master": user}
	result, err := findOneBy(TableTask, msg)
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

func UpdateTaskBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskMaster(uid, master, operator string) error {
	msg := bson.M{"master": master, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskCover(uid, icon, operator string) error {
	msg := bson.M{"cover": icon, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskType(uid, operator string, tp uint8) error {
	msg := bson.M{"type": tp, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskLocal(uid, local, operator string) error {
	msg := bson.M{"location": local, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskAddress(uid, operator string, address AddressInfo) error {
	msg := bson.M{"address": address, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskStatus(uid string, status uint8, operator string) error {
	msg := bson.M{"status": status, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskDomains(uid, operator string, domains []proxy.DomainInfo) error {
	msg := bson.M{"domains": domains, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskShort(uid, operator, name string) error {
	msg := bson.M{"short": name, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskSupporter(uid, supporter, operator string) error {
	msg := bson.M{"supporter": supporter, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskBucket(uid, bucket, operator string) error {
	msg := bson.M{"bucket": bucket, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func UpdateTaskParents(uid, operator string, list []string) error {
	msg := bson.M{"parents": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTask, uid, msg)
	return err
}

func RemoveTask(uid, operator string) error {
	_, err := removeOne(TableTask, uid, operator)
	return err
}

func AppendTaskMember(uid string, member string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := appendElement(TableTask, uid, msg)
	return err
}

func SubtractTaskMember(uid, member string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := removeElement(TableTask, uid, msg)
	return err
}

//func UpdateTaskDisplay(uid, operator string, list []proxy.ShowingInfo) error {
//	msg := bson.M{"displays": list, "operator": operator, "updatedAt": time.Now()}
//	_, err := updateOne(TableTask, uid, msg)
//	return err
//}
//
//func AppendTaskDisplay(uid string, display *proxy.ShowingInfo) error {
//	if len(uid) < 1 {
//		return errors.New("the uid is empty")
//	}
//	msg := bson.M{"displays": display}
//	_, err := appendElement(TableTask, uid, msg)
//	return err
//}
//
//func SubtractTaskDisplay(uid, display string) error {
//	if len(uid) < 1 {
//		return errors.New("the uid is empty")
//	}
//	msg := bson.M{"displays": bson.M{"uid": display}}
//	_, err := removeElement(TableTask, uid, msg)
//	return err
//}

