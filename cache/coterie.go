package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy"
	"omo.msa.assignment/proxy/nosql"
	"time"
)

type CoterieInfo struct {
	Status uint8
	Type   uint8
	baseInfo
	Remark     string
	Cover      string
	Centre     string
	Passwords  string
	Master     string
	Meta       string
	Assistants []string
	Tags       []string
	Members    []proxy.MemberInfo
}

func (mine *cacheContext) CreateCoterie(info *pb.ReqCoterieAdd) (*CoterieInfo, error) {
	db := new(nosql.Coterie)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetCoterieNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = info.Operator
	db.Creator = info.Operator
	db.Name = info.Name
	db.Remark = info.Remark
	db.Status = 0
	db.Cover = info.Cover
	db.Centre = info.Centre
	db.Type = uint8(info.Type)
	db.Passwords = info.Passwords
	db.Master = info.Master
	db.Assistants = make([]string, 0, 1)
	db.Tags = make([]string, 0, 1)
	db.Members = make([]proxy.MemberInfo, 0, 2)
	for _, member := range info.Members {
		db.Members = append(db.Members, proxy.MemberInfo{User: member.User, Name: member.Name, Remark: member.Remark})
	}

	err := nosql.CreateCoterie(db)
	if err == nil {
		tmp := new(CoterieInfo)
		tmp.initInfo(db)
		return tmp, nil
	}
	return nil, err
}

func (mine *cacheContext) GetCoterie(uid string) (*CoterieInfo, error) {
	db, err := nosql.GetCoterie(uid)
	if err != nil {
		return nil, err
	}
	info := new(CoterieInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetCoteriesByMember(uid string) ([]*CoterieInfo, error) {
	dbs, err := nosql.GetCoteriesByMember(uid)
	if err != nil {
		return make([]*CoterieInfo, 0, 1), err
	}
	list := make([]*CoterieInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(CoterieInfo)
		info.initInfo(db)
		list = append(list, info)
	}

	return list, nil
}

func (mine *cacheContext) GetCoterieByCreator(user string) (*CoterieInfo, error) {
	db, err := nosql.GetCoterieByCreator(user)
	if err != nil {
		return nil, err
	}
	info := new(CoterieInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetCoterieByCentre(centre string) (*CoterieInfo, error) {
	db, err := nosql.GetCoterieByCentre(centre)
	if err != nil {
		return nil, err
	}
	info := new(CoterieInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) RemoveCoterie(uid, operator string) error {
	return nosql.RemoveCoterie(uid, operator)
}

func (mine *CoterieInfo) initInfo(db *nosql.Coterie) {
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
	mine.Meta = db.Meta
	mine.Centre = db.Centre
	mine.Type = db.Type
	mine.Status = db.Status
	mine.Assistants = db.Assistants
	mine.Tags = db.Tags
	mine.Members = db.Members
}

func (mine *CoterieInfo) UpdateBase(name, remark, psw, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateCoterieBase(mine.UID, name, remark, psw, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *CoterieInfo) UpdateMaster(master, operator string) error {
	err := nosql.UpdateCoterieMaster(mine.UID, master, operator)
	if err == nil {
		mine.Master = master
		mine.Operator = operator
	}
	return err
}

func (mine *CoterieInfo) UpdateStatus(operator string, st uint8) error {
	err := nosql.UpdateCoterieStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *CoterieInfo) UpdatePasswords(psw, operator  string) error {
	err := nosql.UpdateCoteriePasswords(mine.UID, operator, psw)
	if err == nil {
		mine.Passwords = psw
		mine.Operator = operator
	}
	return err
}

func (mine *CoterieInfo) UpdateTags(operator string, tags []string) error {
	err := nosql.UpdateCoterieTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *CoterieInfo) UpdateMemberIdentify(user, name, remark string) error {
	if !mine.HadMember(user) {
		return errors.New("not found the member in family")
	}
	er := mine.SubtractMember(user)
	if er != nil {
		return er
	}
	return mine.AppendMember(user, name, remark)
}

func (mine *CoterieInfo) UpdateAssistants(operator string, list []string) error {
	err := nosql.UpdateCoterieAssistants(mine.UID, operator, list)
	if err == nil {
		mine.Assistants = list
		mine.Operator = operator
	}
	return err
}

func (mine *CoterieInfo) HadMember(member string) bool {
	for i := 0; i < len(mine.Members); i += 1 {
		if mine.Members[i].User == member {
			return true
		}
	}
	return false
}

func (mine *CoterieInfo) AppendMember(user, name, remark string) error {
	if mine.HadMember(user) {
		return nil
	}
	t := proxy.MemberInfo{User: user, Name: name, Remark: remark}
	err := nosql.AppendCoterieMember(mine.UID, t)
	if err == nil {
		mine.Members = append(mine.Members, t)
	}
	return err
}

func (mine *CoterieInfo) SubtractMember(member string) error {
	if !mine.HadMember(member) {
		return nil
	}
	err := nosql.SubtractCoterieMember(mine.UID, member)
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
