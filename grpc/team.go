package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.organization/cache"
)

type TeamService struct {}

func switchTeam(info *cache.TeamInfo) *pb.TeamInfo {
	tmp := new(pb.TeamInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Location = info.Location
	tmp.Master = info.Master
	tmp.Assistant = info.Assistant
	tmp.Contact = info.Contact
	tmp.Members = info.AllMembers()
	tmp.Scene = info.Scene
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Address = new(pb.AddressInfo)
	tmp.Address.Country = info.Address.Country
	tmp.Address.Province = info.Address.Province
	tmp.Address.City = info.Address.City
	tmp.Address.Zone = info.Address.Zone
	return tmp
}

func (mine *TeamService)AddOne(ctx context.Context, in *pb.ReqTeamAdd, out *pb.ReplyTeamOne) error {
	path := "Team.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path,"the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	scene := cache.Context().GetScene(in.Scene)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	if scene.HadTeamByName(in.Name) {
		out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_Repeated)
		return nil
	}

	Team, err := scene.CreateTeam(in)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchTeam(Team)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTeamOne) error {
	path := "Team.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.TeamInfo
	if len(in.Parent) > 0 {
		scene := cache.Context().GetScene(in.Parent)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		info = scene.GetTeam(in.Uid)
	}else{
		info = cache.Context().GetTeam(in.Uid)
	}

	if info == nil {
		out.Status = outError(path,"the Team not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchTeam(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService)GetByContact(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTeamList) error {
	path := "Team.getByContact"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the phone is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	list := cache.Context().GetTeamByContact(in.Uid)
	out.List = make([]*pb.TeamInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchTeam(info))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TeamService)GetByUser(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTeamList) error {
	path := "Team.getByUser"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	list := cache.Context().GetTeamByMember(in.Uid)
	out.List = make([]*pb.TeamInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchTeam(info))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TeamService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "Team.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path,"the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "Team.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveTeam(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTeamList) error {
	path := "Team.getList"
	inLog(path, in)
	scene := cache.Context().GetScene(in.Parent)
	if scene == nil {
		out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	total,max,list := scene.GetTeams(in.Number, in.Page)
	out.List = make([]*pb.TeamInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchTeam(value))
	}
	out.PageNow = in.Page
	out.Total = total
	out.PageMax = max
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TeamService) UpdateBase (ctx context.Context, in *pb.ReqTeamUpdate, out *pb.ReplyInfo) error {
	path := "Team.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTeam(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Team not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if len(in.Cover) > 0 {
		err = info.UpdateCover(in.Cover, in.Operator)
	}

	if len(in.Name) > 0 || len(in.Remark) > 0 {
		scene := cache.Context().GetScene(info.Scene)
		if scene == nil {
			out.Status = outError(path,"not found the scene ", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Name != info.Name && scene.HadTeamByName(in.Name) {
			out.Status = outError(path,"the department name repeated ", pbstatus.ResultStatus_Repeated)
			return nil
		}
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *TeamService) UpdateAddress (ctx context.Context, in *pb.RequestAddress, out *pb.ReplyTeamOne) error {
	path := "Team.updateAddress"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTeam(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Team not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Location != info.Location {
		err = info.UpdateLocation(in.Location, in.Operator)
		if err != nil {
			out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}
	err = info.UpdateAddress(in.Country, in.Province, in.City, in.Zone, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchTeam(info)
	out.Status = outLog(path, out)
	return err
}

func (mine *TeamService) UpdateLocation (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "Team.updateLocation"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTeam(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Team not found ", pbstatus.ResultStatus_NotExisted)
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

func (mine *TeamService) UpdateContact (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "Team.updateContact"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTeam(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Team not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateContact(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *TeamService) UpdateMaster (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "Team.updateMaster"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTeam(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Team not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateMaster(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) UpdateAssistant (ctx context.Context, in *pb.RequestFlag, out *pb.ReplyInfo) error {
	path := "Team.updateAssistant"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTeam(in.Uid)
	if info == nil {
		out.Status = outError(path,"the Team not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateMaster(in.Flag, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeamService) AppendMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "Team.appendMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTeam(in.Parent)
	if info == nil {
		out.Status = outError(path,"the Team not found ", pbstatus.ResultStatus_NotExisted)
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

func (mine *TeamService) SubtractMember (ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "Team.subtractMember"
	inLog(path, in)
	if len(in.Parent) < 1 {
		out.Status = outError(path,"the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetTeam(in.Parent)
	if info == nil {
		out.Status = outError(path,"the Team not found ", pbstatus.ResultStatus_NotExisted)
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


