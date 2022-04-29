package cache

import (
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy/nosql"
	"time"
)

const (
	AgentStatusIdle  uint8 = 0
	AgentStatusCheck  uint8 = 1
	AgentStatusFroze uint8 = 99
)

const (
	AgentTypeUnknown uint8 = 0
	AgentTypeSpec  uint8 = 1
	AgentTypeFree  uint8 = 2
)

type AgentInfo struct {
	Type   uint8
	Status uint8
	baseInfo
	Remark  string
	User    string
	Entity  string
	Owner   string
	Way     string
	Attaches []string
	Regions []string
	Tags    []string
}

func (mine *cacheContext) CreateAgent(info *pb.ReqAgentAdd) (*AgentInfo, error) {
	db := new(nosql.Agent)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAgentNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = info.Operator
	db.Creator = info.Operator
	db.Name = info.Name
	db.Remark = info.Remark

	if uint8(info.Type) == AgentTypeSpec {
		db.Type = AgentTypeSpec
		db.Status = AgentStatusIdle
	}else if uint8(info.Type) == AgentTypeFree {
		db.Type = AgentTypeFree
		db.Status = AgentStatusCheck
	}else{
		db.Type = AgentTypeUnknown
		db.Status = AgentStatusCheck
	}
	db.Owner = info.Owner
	db.User = info.User
	db.Entity = info.Entity
	db.Way = info.Way
	db.Regions = info.Regions
	db.Attaches = make([]string, 0, 1)
	db.Attaches = append(db.Attaches, info.Owner)
	db.Tags = make([]string, 0, 1)
	err := nosql.CreateAgent(db)
	if err == nil {
		tmp := new(AgentInfo)
		tmp.initInfo(db)
		return tmp, nil
	}
	return nil, err
}

func (mine *cacheContext) GetAgent(uid string) (*AgentInfo, error) {
	db, err := nosql.GetAgent(uid)
	if err != nil {
		return nil, err
	}
	info := new(AgentInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAgentByUser(uid string) (*AgentInfo, error) {
	db, err := nosql.GetAgentByUser(uid)
	if err != nil {
		return nil, err
	}
	info := new(AgentInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAgentsByOwner(uid string) []*AgentInfo {
	list := make([]*AgentInfo, 0, 5)
	dbs, err := nosql.GetAgentsByOwner(uid)
	if err != nil {
		return list
	}
	for _, db := range dbs {
		info := new(AgentInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAgentsByRegion(region string) []*AgentInfo {
	list := make([]*AgentInfo, 0, 5)
	dbs, err := nosql.GetAgentsByRegion(region)
	if err != nil {
		return list
	}
	for _, db := range dbs {
		info := new(AgentInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAgentsByAttach(scene string) []*AgentInfo {
	list := make([]*AgentInfo, 0, 5)
	dbs, err := nosql.GetAgentsByAttach(scene)
	if err != nil {
		return list
	}
	for _, db := range dbs {
		info := new(AgentInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAgentsByWay(owner, way string) []*AgentInfo {
	list := make([]*AgentInfo, 0, 5)
	dbs, err := nosql.GetAgentsByWay(owner, way)
	if err != nil {
		return list
	}
	for _, db := range dbs {
		info := new(AgentInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAgentsByArray(array []string) []*AgentInfo {
	list := make([]*AgentInfo, 0, len(array))
	for _, uid := range array {
		db, err := nosql.GetAgentByUser(uid)
		if err == nil {
			info := new(AgentInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) RemoveAgent(uid, operator string) error {
	return nosql.RemoveAgent(uid, operator)
}

func (mine *AgentInfo) initInfo(db *nosql.Agent) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Type = db.Type
	mine.Status = db.Status
	mine.User = db.User
	mine.Entity = db.Entity
	mine.Owner = db.Owner
	mine.Way = db.Way
	mine.Attaches = db.Attaches
	mine.Tags = db.Tags
	mine.Regions = db.Regions
}

func (mine *AgentInfo) UpdateBase(name, remark, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateAgentBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo) UpdateStatus(operator string, st uint32) error {
	err := nosql.UpdateAgentStatus(mine.UID, operator, uint8(st))
	if err == nil {
		mine.Status = uint8(st)
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo) UpdateEntity(entity, operator string) error {
	err := nosql.UpdateAgentEntity(mine.UID, entity, operator)
	if err == nil {
		mine.Entity = entity
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo) UpdateTags(operator string, tags []string) error {
	err := nosql.UpdateAgentTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo) UpdateRegions(operator string, list []string) error {
	err := nosql.UpdateAgentRegions(mine.UID, operator, list)
	if err == nil {
		mine.Regions = list
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo) UpdateAttaches(operator string, list []string) error {
	err := nosql.UpdateAgentAttaches(mine.UID, operator, list)
	if err == nil {
		mine.Regions = list
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo) hadAttach(uid string) bool {
	for _, attach := range mine.Attaches {
		if attach == uid {
			return true
		}
	}
	return false
}

func (mine *AgentInfo)AddAttach(uid string) error {
	if mine.hadAttach(uid){
		return nil
	}
	err := nosql.AppendAgentAttach(mine.UID, uid)
	if err == nil {
		mine.Attaches = append(mine.Attaches, uid)
	}
	return err
}

func (mine *AgentInfo)RemoveAttach(uid string) error {
	if !mine.hadAttach(uid){
		return nil
	}
	err := nosql.SubtractAgentAttach(mine.UID, uid)
	if err == nil {
		for i := 0;i < len(mine.Attaches);i += 1 {
			if mine.Attaches[i] == uid {
				if i == len(mine.Attaches) - 1 {
					mine.Attaches = append(mine.Attaches[:i])
				}else{
					mine.Attaches = append(mine.Attaches[:i], mine.Attaches[i+1:]...)
				}
				break
			}
		}
	}
	return err
}
