package cache

import (
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy/nosql"
	"omo.msa.assignment/tool"
	"time"
)

const (
	/**

	 */
	Pending MeetingStatus = 0
	/**
	已经开始
	*/
	Idle MeetingStatus = 1
	/**
	自动停止
	*/
	AutoStop MeetingStatus = 2
	/**
	关闭
	*/
	Close MeetingStatus = 3
)

const (
	InRoom  LocationType = 0
	Outside LocationType = 1
)

type MeetingStatus uint8

type LocationType uint8

type MeetingInfo struct {
	baseInfo
	Status   MeetingStatus
	Type     LocationType
	StopTime time.Time
	/**
	会议持续多少分钟
	*/
	Duration  uint16
	Owner     string
	Group     string
	Remark    string
	StartTime time.Time
	// 预约时间
	Appointed string
	Location  string
	Signs     []string
	Submits   []string
	Notifies  []string
}

func (mine *cacheContext) CreateMeeting(in *pb.ReqMeetingAdd) (*MeetingInfo, error) {
	db := new(nosql.Meeting)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetMeetingNextID()
	db.CreatedTime = time.Now()
	db.Creator = in.Operator
	db.Name = in.Name
	db.Remark = in.Remark
	db.Group = in.Group
	db.Owner = in.Owner
	db.Status = uint8(Pending)
	db.Signs = make([]string, 0, 1)
	db.Submits = make([]string, 0, 1)
	db.Notifies = make([]string, 0, 1)
	db.Appointed = in.Appointed
	db.Location = in.Location
	db.StartTime, _ = Context().formatTime(in.Appointed)
	db.Type = uint8(in.Type)

	err := nosql.CreateMeeting(db)
	if err != nil {
		return nil, err
	}
	info := new(MeetingInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetMeeting(uid string) (*MeetingInfo, error) {
	if uid == "" {
		return nil, nil
	}
	db, err := nosql.GetMeeting(uid)
	if err != nil {
		return nil, err
	}
	var info = new(MeetingInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) RemoveMeeting(uid, operator string) error {
	if uid == "" {
		return nil
	}
	return nosql.RemoveMeeting(uid, operator)
}

func (mine *cacheContext) GetMeetingsByGroup(uid string) []*MeetingInfo {
	list := make([]*MeetingInfo, 0, 5)
	array, err := nosql.GetMeetingsByGroup(uid)
	if err == nil {
		for _, item := range array {
			info := new(MeetingInfo)
			info.initInfo(item)

			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetMeetingsByTime(group, from, to string) ([]*MeetingInfo, error) {
	list := make([]*MeetingInfo, 0, 5)
	begin, err := mine.formatDate(from)
	if err != nil {
		return list, err
	}
	end, err := mine.formatDate(to)
	if err != nil {
		return list, err
	}
	utcB := begin.Unix()
	utcE := end.Unix()

	array, err := nosql.GetMeetingsByGroup(group)
	if err == nil {
		for _, item := range array {
			if item.StartTime.Unix() > utcB && item.StopTime.Unix() < utcE {
				info := new(MeetingInfo)
				info.initInfo(item)
				list = append(list, info)
			}
		}
	}
	return list, nil
}

func (mine *cacheContext) GetMeetingsByOwner(uid string) []*MeetingInfo {
	list := make([]*MeetingInfo, 0, 5)
	array, err := nosql.GetMeetingsByScene(uid)
	if err == nil {
		for _, item := range array {
			info := new(MeetingInfo)
			info.initInfo(item)

			list = append(list, info)
		}
	}
	return list
}

func (mine *MeetingInfo) initInfo(db *nosql.Meeting) bool {
	if db == nil {
		return false
	}
	mine.UID = db.UID.Hex()
	mine.Name = db.Name
	mine.ID = db.ID
	mine.Creator = db.Creator
	mine.Group = db.Group
	mine.Owner = db.Owner
	mine.CreateTime = db.CreatedTime
	mine.StopTime = db.StopTime
	mine.StartTime = db.StartTime
	mine.Remark = db.Remark
	mine.Status = MeetingStatus(db.Status)
	mine.Type = LocationType(db.Type)
	mine.Location = db.Location
	mine.Appointed = db.Appointed
	mine.Duration = uint16(mine.StopTime.Unix() - mine.StartTime.Unix())
	mine.Signs = db.Signs
	if mine.Signs == nil {
		mine.Signs = make([]string, 0, 1)
	}
	mine.Submits = db.Submits
	if mine.Submits == nil {
		mine.Submits = make([]string, 0, 1)
	}
	mine.Notifies = db.Notifies
	if mine.Notifies == nil {
		mine.Notifies = make([]string, 0, 1)
	}
	return true
}

func (mine *MeetingInfo) CheckStatus() MeetingStatus {
	if mine.Status == AutoStop || mine.Status == Close {
		return mine.Status
	}
	now := time.Now().Unix()
	if now > mine.StopTime.Unix() {
		mine.Status = AutoStop
	} else if now > mine.StartTime.Unix() {
		mine.Status = Idle
	}
	return mine.Status
}

func (mine *MeetingInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateMeetingBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *MeetingInfo) UpdateLocation(location, operator string, kind LocationType) error {
	err := nosql.UpdateMeetingLocation(mine.UID, location, operator, uint8(kind))
	if err == nil {
		mine.Type = kind
		mine.Location = location
		mine.Operator = operator
	}
	return err
}

func (mine *MeetingInfo) UpdateStartEnd(begin, end int64, operator string) error {
	from := time.Unix(begin, 0).UTC()
	to := time.Unix(end, 0).UTC()
	err := nosql.UpdateMeetingDate(mine.UID, operator, from, to)
	if err == nil {
		mine.StartTime = from
		mine.Operator = operator
		mine.StopTime = to
	}
	return err
}

func (mine *MeetingInfo) UpdateGroup(group, operator string) error {
	err := nosql.UpdateMeetingGroup(mine.UID, operator, group)
	if err == nil {
		mine.Group = group
		mine.Operator = operator
	}
	return err
}

func (mine *MeetingInfo) UpdateStop(date, operator string) error {
	t, err := Context().formatTime(date)
	if err != nil {
		return err
	}
	err = nosql.UpdateMeetingStop(mine.UID, operator, t)
	if err == nil {
		mine.StopTime = t
		mine.UpdateTime = time.Now()
		mine.Operator = operator
	}
	return err
}

func (mine *MeetingInfo) checkRestTime() uint16 {
	if mine.Status != Pending {
		return 0
	}
	return 0
}

func (mine *MeetingInfo) HadSigned(user string) bool {
	for i := 0; i < len(mine.Signs); i += 1 {
		if mine.Signs[i] == user {
			return true
		}
	}
	return false
}

func (mine *MeetingInfo) Sign(member, operator, location string) error {
	if mine.HadSigned(member) {
		return nil
	}
	//if len(mine.Signs) == 0 {
	//	mine.StartTime = time.Now()
	//}
	//if mine.Type == Outside && !Context().checkDistance(mine.Location, location) {
	//	return errors.New("the user location incorrect")
	//}
	err := nosql.AppendMeetingSign(mine.UID, member, operator)
	if err == nil {
		mine.Signs = append(mine.Signs, member)
		mine.Operator = operator
	}
	return err
}

func (mine *MeetingInfo) Submit(member, operator string) error {
	if tool.HasItem(mine.Submits, member) {
		return nil
	}
	err := nosql.AppendMeetingSubmit(mine.UID, member, operator)
	if err == nil {
		mine.Submits = append(mine.Submits, member)
		mine.Operator = operator
	}
	return err
}

func (mine *MeetingInfo) Close(operator string) error {
	err := nosql.StopMeeting(mine.UID, operator)
	if err == nil {
		mine.Status = Close
		mine.StopTime = time.Now()
	}
	return err
}
