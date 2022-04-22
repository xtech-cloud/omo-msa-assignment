package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy"
	"omo.msa.assignment/proxy/nosql"
	"time"
)

const (
	TaskStatusIdle  TaskStatus = 0
	TaskStatusBusy  TaskStatus = 1
	TaskStatusEnd   TaskStatus = 2
	TaskStatusFroze TaskStatus = 99
)

type TaskStatus uint8

type TaskInfo struct {
	Type      uint8
	Status    TaskStatus
	baseInfo
	Remark    string
	Target    string
	Owner     string
	Way       string
	Duration  proxy.DateInfo
	Regions   []string
	Executors []string
	PreTasks  []string
	Tags      []string
	Assets    []string
	Records   []proxy.RecordInfo
}

func (mine *cacheContext) CreateTask(info *pb.ReqTaskAdd) (*TaskInfo, error) {
	db := new(nosql.Task)
	db.UID = primitive.NewObjectID()
	db.Type = uint8(info.Type)
	db.ID = nosql.GetTaskNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = info.Operator
	db.Creator = info.Operator
	db.Name = info.Name
	db.Remark = info.Remark
	db.Target = info.Target
	db.Status = uint8(TaskStatusIdle)
	db.Owner = info.Owner
	db.Duration = proxy.DateInfo{Begin: info.Duration.Begin, End: info.Duration.End}
	db.Regions = info.Regions
	db.PreTasks = info.Pretasks
	db.Tags = info.Tags
	db.Assets = info.Assets
	db.Records = make([]proxy.RecordInfo, 0, 1)
	if db.Regions == nil {
		db.Regions = make([]string, 0, 1)
	}
	if db.PreTasks == nil {
		db.PreTasks = make([]string, 0, 1)
	}
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	if db.Assets == nil {
		db.Assets = make([]string, 0, 1)
	}
	err := nosql.CreateTask(db)
	if err == nil {
		tmp := new(TaskInfo)
		tmp.initInfo(db)
		return tmp, nil
	}
	return nil, err
}

func (mine *cacheContext) GetTask(uid string) (*TaskInfo,error) {
	if len(uid) < 2 {
		return nil,errors.New("the task uid is empty")
	}
	db, err := nosql.GetTask(uid)
	if err == nil {
		info := new(TaskInfo)
		info.initInfo(db)
		return info,nil
	}
	return nil,err
}

func (mine *cacheContext) GetTasksByOwner(parent string, st int, page, number uint32) (uint32, uint32, []*TaskInfo) {
	if number < 1 {
		number = 10
	}
	var dbs []*nosql.Task
	var err error
	if st < 0 {
		dbs, err = nosql.GetTasksByOwner2(parent)
	}else{
		dbs, err = nosql.GetTasksByOwner(parent, uint8(st))
	}
	if err != nil {
		return 0, 0, make([]*TaskInfo, 0, 1)
	}
	all := make([]*TaskInfo, 0, len(dbs))
	for _, item := range dbs {
		info := new(TaskInfo)
		info.initInfo(item)
		all = append(all, info)
	}
	total, maxPage, list := checkPage(page, number, all)
	return total, maxPage, list.([]*TaskInfo)
}

func (mine *cacheContext) GetTasksByType(owner string, tp uint8) []*TaskInfo {
	dbs, err := nosql.GetTasksByType(owner, tp)
	if err != nil {
		return make([]*TaskInfo, 0, 1)
	}
	all := make([]*TaskInfo, 0, len(dbs))
	for _, item := range dbs {
		info := new(TaskInfo)
		info.initInfo(item)
		all = append(all, info)
	}
	return all
}

func (mine *cacheContext) GetTasksByRegions(regions []string, st int) []*TaskInfo {
	all := make([]*TaskInfo, 0, len(regions)*100)
	for _, region := range regions {
		list, err := mine.GetTasksByRegion(region, st)
		if err == nil {
			all = append(all, list...)
		}
	}
	return all
}

func (mine *cacheContext) GetTasksByRegion(region string, st int) ([]*TaskInfo,error) {
	if region == "" {
		return nil, errors.New("the region is empty")
	}
	var dbs []*nosql.Task
	var err error
	if st < int(TaskStatusIdle) {
		dbs, err = nosql.GetTasksByRegion2(region)
	}else{
		dbs, err = nosql.GetTasksByRegion(region, uint8(st))
	}

	if err != nil {
		return nil, err
	}
	all := make([]*TaskInfo, 0, len(dbs))
	for _, item := range dbs {
		info := new(TaskInfo)
		info.initInfo(item)
		all = append(all, info)
	}
	return all,nil
}

func (mine *cacheContext) GetTasksByAgent(agent string, st int) ([]*TaskInfo,error) {
	if agent == "" {
		return nil, errors.New("the agent is empty")
	}
	var dbs []*nosql.Task
	var err error
	if st < int(TaskStatusIdle) {
		dbs, err = nosql.GetTasksByAgent2(agent)
	}else{
		dbs, err = nosql.GetTasksByAgent(agent, uint8(st))
	}

	if err != nil {
		return nil,err
	}
	all := make([]*TaskInfo, 0, len(dbs))
	for _, item := range dbs {
		info := new(TaskInfo)
		info.initInfo(item)
		all = append(all, info)
	}
	return all,nil
}

func (mine *cacheContext) GetTasksByTarget(target string, st int) ([]*TaskInfo,error) {
	if target == "" {
		return nil, errors.New("the target is empty")
	}
	var dbs []*nosql.Task
	var err error
	if st < int(TaskStatusIdle) {
		dbs, err = nosql.GetTasksByTarget2(target)
	}else{
		dbs, err = nosql.GetTasksByTarget(target, uint8(st))
	}

	if err != nil {
		return nil,err
	}
	all := make([]*TaskInfo, 0, len(dbs))
	for _, item := range dbs {
		info := new(TaskInfo)
		info.initInfo(item)
		all = append(all, info)
	}
	return all,nil
}

func RemoveTask(uid, operator string) error {
	if len(uid) < 1 {
		return errors.New("the team uid is empty")
	}
	err := nosql.RemoveTask(uid, operator)
	return err
}

func (mine *TaskInfo) initInfo(db *nosql.Task) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Status = TaskStatus(db.Status)
	mine.Type = db.Type
	mine.Target = db.Target
	mine.Owner = db.Owner
	mine.Way = db.Way
	mine.Duration = db.Duration

	mine.Regions = db.Regions
	mine.Executors = db.Executors
	mine.PreTasks = db.PreTasks
	mine.Tags = db.Tags
	mine.Assets = db.Assets
	mine.Records = db.Records
}

func (mine *TaskInfo) UpdateBase(name, remark, operator string, assets []string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateTaskBase(mine.UID, name, remark, operator, assets)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Assets = assets
	}
	return err
}

func (mine *TaskInfo) UpdateType(operator string, tp uint8) error {
	if uint8(mine.Type) == tp {
		return nil
	}
	err := nosql.UpdateTaskType(mine.UID, operator, tp)
	if err == nil {
		mine.Type = tp
		mine.Operator = operator
	}
	return err
}

func (mine *TaskInfo) UpdateExecutors(operator string, agents []string) error {
	err := nosql.UpdateTaskExecutors(mine.UID, operator, agents)
	if err == nil {
		mine.Executors = agents
		mine.Operator = operator
	}
	return err
}

func (mine *TaskInfo) UpdateTags(operator string, list []string) error {
	err := nosql.UpdateTaskTags(mine.UID, operator, list)
	if err == nil {
		mine.Tags = list
		mine.Operator = operator
	}
	return err
}

func (mine *TaskInfo) UpdateStatus(st TaskStatus, operator string) error {
	err := nosql.UpdateTaskStatus(mine.UID, uint8(st), operator)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *TaskInfo) HadExecutor(member string) bool {
	for i := 0;i < len(mine.Executors);i += 1 {
		if mine.Executors[i] == member {
			return true
		}
	}
	return false
}

func (mine *TaskInfo)HadRecord(uid string) bool {
	for i := 0;i < len(mine.Records);i += 1 {
		if mine.Records[i].UID == uid {
			return true
		}
	}
	return false
}

func (mine *TaskInfo) AppendExecutor(member string) error {
	if mine.HadExecutor(member){
		return nil
	}
	err := nosql.AppendTaskExecutor(mine.UID, member)
	if err == nil {
		mine.Executors = append(mine.Executors, member)
	}
	return err
}

func (mine *TaskInfo) SubtractExecutor(member string) error {
	if !mine.HadExecutor(member){
		return nil
	}
	err := nosql.SubtractTaskExecutor(mine.UID, member)
	if err == nil {
		for i := 0;i < len(mine.Executors);i += 1 {
			if mine.Executors[i] == member {
				if i == len(mine.Executors) - 1 {
					mine.Executors = append(mine.Executors[:i])
				}else{
					mine.Executors = append(mine.Executors[:i], mine.Executors[i+1:]...)
				}

				break
			}
		}
	}
	return err
}

func (mine *TaskInfo)AddRecord(tmp *pb.ReqTaskRecord) error {
	st := TaskStatus(tmp.Status)
	var err error
	if mine.Status != st {
		err = mine.UpdateStatus(st, tmp.Creator)
		if err != nil {
			return err
		}
	}
	info := proxy.RecordInfo{
		Creator: tmp.Creator,
		CreatedTime: time.Now(),
		Name: tmp.Name,
		Remark: tmp.Remark,
		Executor: tmp.Executor,
		Status: uint8(tmp.Status),
		Tags: tmp.Tags,
		Assets: tmp.Assets,
	}
	err = nosql.AppendTaskRecord(mine.UID, info)
	if err == nil {
		mine.Records = append(mine.Records, info)
		arr := make([]string, 0, 1)
		arr = append(arr, tmp.Executor)
		_ = mine.UpdateExecutors(tmp.Creator, arr)
	}
	return err
}

func (mine *TaskInfo)RemoveRecord(uid string) error {
	if !mine.HadRecord(uid){
		return nil
	}
	err := nosql.SubtractTaskRecord(mine.UID, uid)
	if err == nil {
		for i := 0;i < len(mine.Records);i += 1 {
			if mine.Records[i].UID == uid {
				if i == len(mine.Records) - 1 {
					mine.Records = append(mine.Records[:i])
				}else{
					mine.Records = append(mine.Records[:i], mine.Records[i+1:]...)
				}
				break
			}
		}
	}
	return err
}


