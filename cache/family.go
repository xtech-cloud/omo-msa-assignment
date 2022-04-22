package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy"
	"omo.msa.assignment/proxy/nosql"
	"time"
)

type FamilyInfo struct {
	Status uint8
	baseInfo
	Remark     string
	SN         string
	Cover      string
	Region     string
	Address    string
	Location   string
	Passwords  string
	Master     string
	Assistants []string
	Children   []string
	Tags       []string
	Agents     []string
	Members    []proxy.MemberInfo
}

func (mine *cacheContext) CreateFamily(info *pb.ReqFamilyAdd) (*FamilyInfo, error) {
	db := new(nosql.Family)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetFamilyNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = info.Operator
	db.Creator = info.Operator
	db.Name = info.Name
	db.Remark = info.Remark
	db.SN = info.Sn
	db.Region = info.Region
	db.Status = 0
	db.Location = info.Location
	db.Address = info.Address
	db.Passwords = info.Passwords
	db.Master = info.Master
	db.Assistants = make([]string, 0, 1)
	db.Children = make([]string, 0, 1)
	db.Tags = make([]string, 0, 1)
	db.Agents = make([]string, 0, 1)
	db.Members = make([]proxy.MemberInfo, 0, 2)
	for _, member := range info.Members {
		db.Members = append(db.Members, proxy.MemberInfo{User: member.User, Remark: member.Remark})
	}

	err := nosql.CreateFamily(db)
	if err == nil {
		tmp := new(FamilyInfo)
		tmp.initInfo(db)
		return tmp, nil
	}
	return nil, err
}

func (mine *cacheContext) GetFamily(uid string) (*FamilyInfo, error) {
	db, err := nosql.GetFamily(uid)
	if err != nil {
		return nil, err
	}
	info := new(FamilyInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetFamiliesByMember(uid string) ([]*FamilyInfo, error) {
	dbs, err := nosql.GetFamiliesByMember(uid)
	if err != nil {
		return make([]*FamilyInfo, 0, 1), err
	}
	list := make([]*FamilyInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(FamilyInfo)
		info.initInfo(db)
		list = append(list, info)
	}

	return list, nil
}

func (mine *cacheContext) GetFamiliesByAgent(uid string) ([]*FamilyInfo, error) {
	dbs, err := nosql.GetFamiliesByAgent(uid)
	if err != nil {
		return make([]*FamilyInfo, 0, 1), err
	}
	list := make([]*FamilyInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(FamilyInfo)
		info.initInfo(db)
		list = append(list, info)
	}

	return list, nil
}

func (mine *cacheContext) GetFamiliesByRegion(region string) ([]*FamilyInfo, error) {
	dbs, err := nosql.GetFamiliesByRegion(region)
	if err != nil {
		return make([]*FamilyInfo, 0, 1), err
	}
	list := make([]*FamilyInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(FamilyInfo)
		info.initInfo(db)
		list = append(list, info)
	}

	return list, nil
}

func (mine *cacheContext) GetFamiliesByRegions(regions []string) ([]*FamilyInfo, error) {
	list := make([]*FamilyInfo, 0, 100)
	for _, region := range regions {
		array,_ := mine.GetFamiliesByRegion(region)
		for _, info := range array {
			list = append(list, info)
		}
	}

	return list, nil
}

func (mine *cacheContext) GetFamilyByCreator(user string) (*FamilyInfo, error) {
	db, err := nosql.GetFamilyByCreator(user)
	if err != nil {
		return nil, err
	}
	info := new(FamilyInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetFamilyBySN(sn string) (*FamilyInfo, error) {
	db, err := nosql.GetFamilyBySN(sn)
	if err != nil {
		return nil, err
	}
	info := new(FamilyInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetFamilyByChild(child string) (*FamilyInfo, error) {
	db, err := nosql.GetFamilyByChild(child)
	if err != nil {
		return nil, err
	}
	info := new(FamilyInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) RemoveFamily(uid, operator string) error {
	return nosql.RemoveFamily(uid, operator)
}

func (mine *FamilyInfo) initInfo(db *nosql.Family) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Cover = db.Cover
	mine.Passwords = db.Passwords
	mine.Master = db.Master
	mine.Remark = db.Remark
	mine.SN = db.SN
	mine.Address = db.Address
	mine.Location = db.Location
	mine.Region = db.Region
	mine.Status = db.Status
	mine.Assistants = db.Assistants
	mine.Agents = db.Agents
	mine.Children = db.Children
	mine.Members = db.Members
	mine.Tags = db.Tags
}

func (mine *FamilyInfo) UpdateBase(name, remark, psw, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateFamilyBase(mine.UID, name, remark, psw, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *FamilyInfo) UpdateMaster(master, operator string) error {
	err := nosql.UpdateFamilyMaster(mine.UID, master, operator)
	if err == nil {
		mine.Master = master
		mine.Operator = operator
	}
	return err
}

func (mine *FamilyInfo) UpdateSN(sn, operator string) error {
	err := nosql.UpdateFamilySN(mine.UID, sn, operator)
	if err == nil {
		mine.SN = sn
		mine.Operator = operator
	}
	return err
}

func (mine *FamilyInfo) UpdateStatus(operator string, st uint8) error {
	err := nosql.UpdateFamilyStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *FamilyInfo) UpdatePasswords(psw, operator  string) error {
	err := nosql.UpdateFamilyPasswords(mine.UID, operator, psw)
	if err == nil {
		mine.Passwords = psw
		mine.Operator = operator
	}
	return err
}

func (mine *FamilyInfo) UpdateTags(operator string, tags []string) error {
	err := nosql.UpdateFamilyTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *FamilyInfo) UpdateAgents(operator string, list []string) error {
	err := nosql.UpdateFamilyAgents(mine.UID, operator, list)
	if err == nil {
		mine.Agents = list
		mine.Operator = operator
	}
	return err
}

func (mine *FamilyInfo) UpdateChildren(operator string, list []string) error {
	err := nosql.UpdateFamilyChildren(mine.UID, operator, list)
	if err == nil {
		mine.Children = list
		mine.Operator = operator
	}
	return err
}

func (mine *FamilyInfo) UpdateMemberIdentify(user, remark string) error {
	if !mine.HadMember(user) {
		return errors.New("not found the member in family")
	}
	er := mine.SubtractMember(user)
	if er != nil {
		return er
	}
	return mine.AppendMember(user, remark)
}

func (mine *FamilyInfo) UpdateAssistants(operator string, list []string) error {
	err := nosql.UpdateFamilyAssistants(mine.UID, operator, list)
	if err == nil {
		mine.Assistants = list
		mine.Operator = operator
	}
	return err
}

func (mine *FamilyInfo) HadMember(member string) bool {
	for i := 0; i < len(mine.Members); i += 1 {
		if mine.Members[i].User == member {
			return true
		}
	}
	return false
}

func (mine *FamilyInfo) AppendMember(user, remark string) error {
	if mine.HadMember(user) {
		return nil
	}
	t := proxy.MemberInfo{User: user, Remark: remark}
	err := nosql.AppendFamilyMember(mine.UID, t)
	if err == nil {
		mine.Members = append(mine.Members, t)
	}
	return err
}

func (mine *FamilyInfo) SubtractMember(member string) error {
	if !mine.HadMember(member) {
		return nil
	}
	err := nosql.SubtractFamilyMember(mine.UID, member)
	if err == nil {
		for i := 0; i < len(mine.Members); i += 1 {
			if mine.Members[i].User == member {
				if i == len(mine.Members)-1 {
					mine.Members = append(mine.Members[:i])
				} else {
					mine.Members = append(mine.Members[:i], mine.Members[i+1:]...)
				}

				break
			}
		}
	}
	return err
}
