package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.organization/cache"
	"omo.msa.organization/proxy"
	"strconv"
)

type TaskService struct {}

func switchTask(info *cache.TaskInfo) *pb.TaskInfo {
	tmp := new(pb.TaskInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Type = int32(info.Type)
	tmp.Status = int32(info.Status)
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Location = info.Location
	tmp.Master = info.Master
	tmp.Entity = info.Entity
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Supporter = info.Supporter
	tmp.Bucket = info.Bucket
	tmp.Short = info.ShortName
	tmp.Parents = info.Parents()
	tmp.Members = info.AllMembers()
	tmp.Domains = make([]*pb.ProductInfo, 0, len(info.Domains))
	for _, domain := range info.Domains {
		tmp.Domains = append(tmp.Domains, &pb.ProductInfo{Type: uint32(domain.Type), Uid: domain.UID,
			Remark: domain.Remark, Keywords: domain.Keywords, Name: domain.Name})
	}
	return tmp
}

func (mine *TaskService)AddOne(ctx context.Context, in *pb.ReqTaskAdd, out *pb.ReplyTaskOne) error {
	path := "Task.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := new(cache.TaskInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Master = in.Master
	info.Location = in.Location
	info.Type = cache.TaskType(in.Type)
	info.Cover = in.Cover
	info.Entity = in.Entity
	info.Creator = in.Operator
	info.ShortName = ""
	err := cache.Context().CreateTask(info)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchTask(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTaskOne) error {
	path := "Task.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchTask(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService)GetOneByMaster(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTaskOne) error {
	path := "Task.getByMaster"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTaskByMember(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchTask(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "Task.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.RemoveTask(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService)IsMasterUsed(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyMasterUsed) error {
	path := "Task.isMasterUsed"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	out.Master = in.Uid
	out.Used = cache.IsMasterUsed(in.Uid)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTaskList) error {
	path := "Task.getList"
	inLog(path, in)
	var total uint32 = 0
	var max uint32 = 0
	var list []*cache.TaskInfo
	if in.Parent == "" {
		total,max,list = cache.Context().GetTasks(in.Page, in.Number)
	}else{
		total,max,list = cache.Context().GetTasksByParent(in.Parent, in.Page, in.Number)
	}

	out.List = make([]*pb.TaskInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchTask(value))
	}
	out.PageNow = in.Page
	out.Total = total
	out.PageMax = max
	out.Status = &pb.ReplyStatus{
		Code: 0,
		Error: "",
	}
	return nil
}

func (mine *TaskService)GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyTaskList) error {
	path := "Task.getByFilter"
	inLog(path, in)
	var total uint32 = 0
	var max uint32 = 0
	var list []*cache.TaskInfo
	if in.Key == "shortname" {
		list = make([]*cache.TaskInfo, 0 ,1)
	}else if in.Key == "type" {
		tp,er := strconv.ParseUint(in.Parent, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		list = cache.Context().GetTasksByType(uint8(tp))
	}else if in.Key == "parent" {
		 total, max, list = cache.Context().GetTasksByParent(in.Value, in.Page, in.Number)
	}else if in.Key == "array" {
		list = cache.Context().GetTasksByArray(in.List)
	}else{
		list = make([]*cache.TaskInfo, 0 ,1)
	}

	out.List = make([]*pb.TaskInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchTask(value))
	}
	out.PageNow = in.Page
	out.Total = total
	out.PageMax = max
	out.Status = &pb.ReplyStatus{
		Code: 0,
		Error: "",
	}
	return nil
}

func (mine *TaskService)GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "Task.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) UpdateBase (ctx context.Context, in *pb.ReqTaskUpdate, out *pb.ReplyInfo) error {
	path := "Task.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if len(in.Cover) > 0 {
		err = info.UpdateCover(in.Cover, in.Operator)
	}
	if len(in.Master) > 0 {
		err = info.UpdateMaster(in.Master, in.Operator)
	}
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}
	if in.Type > 0 {
		err = info.UpdateType(in.Operator, uint8(in.Type))
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *TaskService) UpdateAddress (ctx context.Context, in *pb.RequestAddress, out *pb.ReplyTaskOne) error {
	path := "Task.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateAddress(in.Country, in.Province, in.City, in.Zone, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchTask(info)
	out.Status = outLog(path, out)
	return err
}

func (mine *TaskService) UpdateLocation (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "Task.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateLocation(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *TaskService) UpdateStatus (ctx context.Context, in *pb.ReqTaskStatus, out *pb.ReplyInfo) error {
	path := "Task.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateStatus(cache.TaskStatus(in.Status),in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) UpdateSupporter (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "Task.updateSupporter"
	inLog(path, in)
	if len(in.Flag) < 1 {
		out.Status = outError(path,"the supporter uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the Task uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateSupporter(in.Flag,in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) UpdateDomains (ctx context.Context, in *pb.ReqTaskDomains, out *pb.ReplyInfo) error {
	path := "Task.updateDomains"
	inLog(path, in)

	if len(in.Uid) < 1 {
		out.Status = outError(path,"the Task uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	arr := make([]proxy.DomainInfo, 0, len(in.List))
	for _, item := range in.List {
		arr = append(arr, proxy.DomainInfo{Type: uint8(item.Type), UID: item.Uid, Remark: item.Remark, Keywords: item.Keywords, Name: item.Name})
	}
	err := info.UpdateDomains(in.Operator, arr)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) UpdateBucket (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "Task.updateBucket"
	inLog(path, in)
	if len(in.Flag) < 1 {
		out.Status = outError(path,"the domain uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the Task uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateBucket(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) UpdateShortName (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "Task.updateShortName"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the short name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateShortName(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) AppendMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "Task.appendMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Parent)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.AllMembers()
	out.Status = outLog(path, out)
	return nil
}

func (mine *TaskService) SubtractMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "Task.subtractMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Parent)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractMember(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.AllMembers()
	out.Status = outLog(path, out)
	return nil
}

//func (mine *TaskService) UpdateDisplay (ctx context.Context, in *pb.ReqTaskDisplay, out *pb.ReplyTaskDisplays) error {
//	path := "Task.updateDisplay"
//	inLog(path, in)
//	if len(in.Task) < 1 {
//		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info := cache.Context().GetTask(in.Task)
//	if info == nil {
//		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//	if in.Slots == nil {
//		in.Slots = make([]string, 0, 1)
//	}
//	err := info.UpdateDisplay(in.Uid, in.Key, in.Skin, in.Operator, in.Slots)
//	if err != nil {
//		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Uid = in.Uid
//	out.List = switchExhibitions(info.Exhibitions)
//	out.Status = outLog(path, out)
//	return nil
//}

//func (mine *TaskService) PutOnDisplay (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTaskDisplays) error {
//	path := "Task.putOnDisplay"
//	inLog(path, in)
//	if len(in.Parent) < 1 {
//		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info := cache.Context().GetTask(in.Parent)
//	if info == nil {
//		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//
//	err := info.PutOnDisplay(in.Uid)
//	if err != nil {
//		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Uid = in.Uid
//	out.List = switchExhibitions(info.Exhibitions)
//	out.Status = outLog(path, out)
//	return nil
//}

//func (mine *TaskService) CancelDisplay (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTaskDisplays) error {
//	path := "Task.cancelDisplay"
//	inLog(path, in)
//	if len(in.Parent) < 1 {
//		out.Status = outError(path,"the parent is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info := cache.Context().GetTask(in.Parent)
//	if info == nil {
//		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//
//	err := info.CancelDisplay(in.Uid)
//	if err != nil {
//		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Uid = in.Uid
//	out.List = switchExhibitions(info.Exhibitions)
//	out.Status = outLog(path, out)
//	return nil
//}

func (mine *TaskService) UpdateParents (ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "Task.updateParents"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTask(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Task not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.UpdateParents(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Parents()
	out.Status = outLog(path, out)
	return nil
}


