package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.assignment/cache"
	"omo.msa.assignment/proxy"
	"strconv"
)

type TaskService struct{}

func switchTask(info *cache.TaskInfo) *pb.TaskInfo {
	tmp := new(pb.TaskInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Target = info.Target
	tmp.Type =uint32(info.Type)
	tmp.Status = uint32(info.Status)
	tmp.Remark = info.Remark
	tmp.Way = info.Way
	tmp.Duration = &pb.DateInfo{Begin: info.Duration.Begin, End: info.Duration.End}
	tmp.Executors = info.Executors
	tmp.Pretasks = info.PreTasks
	tmp.Tags = info.Tags
	tmp.Regions = info.Regions
	tmp.Assets = info.Assets
	tmp.Records = switchRecords(info.Records)
	return tmp
}

func switchRecords(array []proxy.RecordInfo) []*pb.RecordInfo {
	list := make([]*pb.RecordInfo, 0, len(array))
	for _, info := range array {
		tmp := new(pb.RecordInfo)
		tmp.Uid = info.UID
		tmp.Creator = info.Creator
		tmp.Created = info.CreatedTime.Unix()
		tmp.Name = info.Name
		tmp.Status = uint32(info.Status)
		tmp.Remark = info.Remark
		tmp.Executor = info.Executor
		tmp.Tags = info.Tags
		tmp.Assets = info.Assets
		list = append(list, tmp)
	}
	return list
}

func (mine *TaskService) AddOne(ctx context.Context, in *pb.ReqTaskAdd, out *pb.ReplyTaskOne) error {
	path := "task.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateTask(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchTask(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTaskOne) error {
	path := "task.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTask(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchTask(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "task.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.RemoveTask(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTaskList) error {
	path := "task.search"
	inLog(path, in)
	var total uint32 = 0
	var max uint32 = 0

	out.Total = total
	out.PageMax = max
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TaskService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyTaskList) error {
	path := "task.getListByFilter"
	inLog(path, in)
	var total uint32 = 0
	var max uint32 = 0
	var list []*cache.TaskInfo
	var err error
	if in.Key == "status" {
		st := parseStringToInt(in.Value)
		total, max, list = cache.Context().GetTasksByOwner(in.Owner, st, in.Page, in.Number)
	}else if in.Key == "target;status" {
		uid,st := parseString(in.Value, ";")
		list,err = cache.Context().GetTasksByTarget(uid, st)
	} else if in.Key == "agent;status" {
		uid,st := parseString(in.Value, ";")
		list,err = cache.Context().GetTasksByAgent(uid, st)
	} else if in.Key == "regions" {
		list = cache.Context().GetTasksByRegions(in.Values, -1)
	} else if in.Key == "regions;status" {
		st := parseStringToInt(in.Value)
		list = cache.Context().GetTasksByRegions(in.Values, st)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.List = make([]*pb.TaskInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchTask(value))
	}
	out.PageNow = in.Page
	out.Total = total
	out.PageMax = max
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TaskService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "task.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) UpdateBase(ctx context.Context, in *pb.ReqTaskUpdate, out *pb.ReplyInfo) error {
	path := "task.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTask(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.Remark, in.Operator, in.Assets)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *TaskService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "task.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTask(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "" {
		val,_ := strconv.ParseUint(in.Value, 10, 32)
		err = info.UpdateType(in.Operator, uint8(val))
	}

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *TaskService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "task.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTask(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateStatus(cache.TaskStatus(in.Flag), in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) AppendAgent(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "task.appendAgent"
	inLog(path, in)
	if len(in.Flag) < 1 {
		out.Status = outError(path, "the supporter uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the Task uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTask(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendExecutor(in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Executors
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) SubtractAgent(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "task.subtractAgent"
	inLog(path, in)
	if len(in.Flag) < 1 {
		out.Status = outError(path, "the supporter uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the Task uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTask(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractExecutor(in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Executors
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) AppendRecord(ctx context.Context, in *pb.ReqTaskRecord, out *pb.ReplyTaskRecords) error {
	path := "task.appendRecord"
	inLog(path, in)
	if len(in.Task) < 1 {
		out.Status = outError(path, "the parent is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTask(in.Task)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AddRecord(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = switchRecords(info.Records)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) SubtractRecord(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTaskRecords) error {
	path := "task.subtractRecord"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the task uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetTask(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.RemoveRecord(in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = switchRecords(info.Records)
	out.Status = outLog(path, out)
	return nil
}
