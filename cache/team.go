package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy/nosql"
	"time"
)

type TeamInfo struct {
	Status uint8
	baseInfo
	Remark string
	Owner string
	Master string
	Region string
	Tags []string
	Assistants []string
	Members []string
}

func (mine *cacheContext) CreateTeam(info *pb.ReqTeamAdd) (*TeamInfo, error) {
	db := new(nosql.Team)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetTeamNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = info.Operator
	db.Creator = info.Operator
	db.Name = info.Name
	db.Remark = info.Remark
	db.Status = uint8(TaskStatusIdle)
	db.Owner = info.Owner
	db.Master = ""
	db.Tags = make([]string, 0, 1)
	db.Members = make([]string, 0, 1)
	db.Assistants = make([]string, 0, 1)
	err := nosql.CreateTeam(db)
	if err == nil {
		tmp := new(TeamInfo)
		tmp.initInfo(db)
		return tmp, nil
	}
	return nil, err
}

func (mine *cacheContext) GetTeam(uid string) (*TeamInfo,error) {
	db,err := nosql.GetTeam(uid)
	if err == nil {
		info := new(TeamInfo)
		info.initInfo(db)
		return info,nil
	}
	return nil,err
}

func (mine *cacheContext) GetTeamByOwner(scene string) []*TeamInfo {
	list := make([]*TeamInfo, 0, 10)
	dbs,err := nosql.GetTeamsByOwner(scene)
	if err == nil {
		for _, db := range dbs {
			info := new(TeamInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) HadTeamByName(scene, name string) bool {
	db,_ := nosql.GetTeamByName(scene, name)
	if db == nil {
		return false
	}
	return true
}

func (mine *cacheContext) RemoveTeam(uid, operator string) error {
	if len(uid) < 1 {
		return errors.New("the team uid is empty")
	}
	err := nosql.RemoveTeam(uid, operator)
	return err
}

func (mine *TeamInfo) initInfo(db *nosql.Team) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Owner = db.Owner
	mine.Master = db.Master
	mine.Region = db.Region
	mine.Status = db.Status
	mine.Tags = db.Tags
	mine.Assistants = db.Assistants
	mine.Members = db.Members
}

func (mine *TeamInfo) UpdateBase(name, remark, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateTeamBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *TeamInfo) UpdateMaster(master, operator string) error {
	err := nosql.UpdateTeamMaster(mine.UID, master, operator)
	if err == nil {
		mine.Master = master
		mine.Operator = operator
	}
	return err
}

func (mine *TeamInfo) UpdateStatus(operator string, st uint8) error {
	err := nosql.UpdateTeamStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *TeamInfo) UpdateRegion(region, operator string) error {
	err := nosql.UpdateTeamRegion(mine.UID, region, operator)
	if err == nil {
		mine.Region = region
		mine.Operator = operator
	}
	return err
}

func (mine *TeamInfo) UpdateTags(operator string, tags []string) error {
	err := nosql.UpdateTeamTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *TeamInfo) UpdateAssistants(operator string, list []string) error {
	err := nosql.UpdateTeamAssistants(mine.UID, operator, list)
	if err == nil {
		mine.Assistants = list
		mine.Operator = operator
	}
	return err
}

func (mine *TeamInfo) UpdateMembers(operator string, list []string) error {
	err := nosql.UpdateTeamMembers(mine.UID, operator, list)
	if err == nil {
		mine.Members = list
		mine.Operator = operator
	}
	return err
}

func (mine *TeamInfo) HadMember(member string) bool {
	for i := 0;i < len(mine.Members);i += 1 {
		if mine.Members[i] == member {
			return true
		}
	}
	return false
}

func (mine *TeamInfo)AppendMember(member string) error {
	if mine.HadMember(member){
		return nil
	}
	err := nosql.AppendTeamMember(mine.UID, member)
	if err == nil {
		mine.Members = append(mine.Members, member)
	}
	return err
}

func (mine *TeamInfo)SubtractMember(member string) error {
	if !mine.HadMember(member){
		return nil
	}
	err := nosql.SubtractTeamMember(mine.UID, member)
	if err == nil {
		for i := 0;i < len(mine.Members);i += 1 {
			if mine.Members[i] == member {
				if i == len(mine.Members) - 1 {
					mine.Members = append(mine.Members[:i])
				}else{
					mine.Members = append(mine.Members[:i], mine.Members[i+1:]...)
				}
				break
			}
		}
	}
	return err
}
