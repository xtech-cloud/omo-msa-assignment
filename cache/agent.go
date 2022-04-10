package cache

import (
	"errors"
	"omo.msa.organization/proxy/nosql"
)

type AgentInfo struct {
	baseInfo
	Remark string
	Contact string
	Cover string
	Master string
	Assistant string
	Address nosql.AddressInfo
	Location string
	Scene string
	members []string
}

func (mine *cacheContext)GetAgent(uid string) *AgentInfo {
	for _, scene := range mine.scenes {
		Agent := scene.GetAgent(uid)
		if Agent != nil {
			return Agent
		}
	}
	return nil
}

func (mine *cacheContext)GetAgentByMember(uid string) []*AgentInfo {
	list := make([]*AgentInfo, 0, 5)
	for _, scene := range mine.scenes {
		for _, Agent := range scene.Agents {
			if Agent.HadMember(uid) {
				list = append(list, Agent)
			}
		}
	}
	return list
}

func (mine *cacheContext)GetAgentByContact(phone string) []*AgentInfo {
	list := make([]*AgentInfo, 0, 5)
	for _, scene := range mine.scenes {
		for _, Agent := range scene.Agents {
			if Agent.Contact == phone {
				list = append(list, Agent)
			}
		}
	}
	return list
}

func (mine *cacheContext)RemoveAgent(uid, operator string) error {
	for _, scene := range mine.scenes {
		if scene.HadAgent(uid) {
			return scene.RemoveAgent(uid, operator)
		}
	}
	return nil
}

func (mine *AgentInfo)initInfo(db *nosql.Agent)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Cover = db.Cover
	mine.Remark = db.Remark
	mine.Master = db.Master
	mine.Location = db.Location
	mine.Assistant = db.Assistant
	mine.Contact = db.Contact
	mine.members = db.Members
	mine.Address = db.Address
	mine.Scene = db.Scene
}

func (mine *AgentInfo)UpdateBase(name, remark, operator string) error {
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

func (mine *AgentInfo)UpdateContact(phone, operator string) error {
	err := nosql.UpdateAgentContact(mine.UID, phone, operator)
	if err == nil {
		mine.Contact = phone
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo)UpdateMaster(master, operator string) error {
	err := nosql.UpdateAgentMaster(mine.UID, master, operator)
	if err == nil {
		mine.Master = master
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo)UpdateAssistant(uid, operator string) error {
	err := nosql.UpdateAgentAssistant(mine.UID, uid, operator)
	if err == nil {
		mine.Assistant = uid
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo)UpdateCover(cover, operator string) error {
	err := nosql.UpdateAgentCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo)UpdateLocation(local, operator string) error {
	err := nosql.UpdateAgentLocation(mine.UID, local, operator)
	if err == nil {
		mine.Location = local
		mine.Operator = operator
	}
	return err
}

func (mine *AgentInfo)UpdateAddress(country, province, city, zone, operator string) error {
	addr := nosql.AddressInfo{Country: country, Province: province, City: city, Zone: zone}
	err := nosql.UpdateAgentAddress(mine.UID, operator, addr)
	if err == nil {
		mine.Address = addr
		mine.Operator = operator
	}
	return err
}


func (mine *AgentInfo)HadMember(member string) bool {
	if mine.Master == member || mine.Assistant == member {
		return true
	}
	for i := 0;i < len(mine.members);i += 1 {
		if mine.members[i] == member {
			return true
		}
	}
	return false
}

func (mine *AgentInfo)AllMembers() []string {
	return mine.members
}

func (mine *AgentInfo)AppendMember(member string) error {
	if mine.HadMember(member){
		return errors.New("the member had existed")
	}
	err := nosql.AppendAgentMember(mine.UID, member)
	if err == nil {
		mine.members = append(mine.members, member)
	}
	return err
}

func (mine *AgentInfo)SubtractMember(member string) error {
	if !mine.HadMember(member){
		return errors.New("the member not existed")
	}
	err := nosql.SubtractAgentMember(mine.UID, member)
	if err == nil {
		for i := 0;i < len(mine.members);i += 1 {
			if mine.members[i] == member {
				mine.members = append(mine.members[:i], mine.members[i+1:]...)
				break
			}
		}
	}
	return err
}

