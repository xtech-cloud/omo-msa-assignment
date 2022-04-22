package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.assignment/cache"
)

type CoterieService struct{}

func switchCoterie(info *cache.CoterieInfo) *pb.CoterieInfo {
	tmp := new(pb.CoterieInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Master = info.Master
	tmp.Centre = info.Centre
	tmp.Meta = info.Meta
	tmp.Type = uint32(info.Type)
	tmp.Status = uint32(info.Status)
	tmp.Passwords = info.Passwords
	tmp.Assistants = info.Assistants
	tmp.Tags = info.Tags
	tmp.Members = make([]*pb.IdentifyInfo, 0, len(info.Members))
	for _, member := range info.Members {
		tmp.Members = append(tmp.Members, &pb.IdentifyInfo{User: member.User, Name: member.Name, Remark: member.Remark})
	}
	return tmp
}

func (mine *CoterieService) AddOne(ctx context.Context, in *pb.ReqCoterieAdd, out *pb.ReplyCoterieInfo) error {
	path := "family.addOne"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateCoterie(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCoterie(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CoterieService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCoterieInfo) error {
	path := "family.getOne"
	inLog(path, in)
	var info *cache.CoterieInfo
	var er error
	if len(in.Uid) > 1 {
		info,er = cache.Context().GetCoterie(in.Uid)
	}else if in.Flag == "centre" {
		info,er = cache.Context().GetCoterieByCentre(in.User)
	}else if in.Flag == "creator" {
		info,er = cache.Context().GetCoterieByCreator(in.User)
	}else {
		er = errors.New("")
	}

	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchCoterie(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CoterieService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCoterieList) error {
	path := "family.search"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CoterieService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "family.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the key is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CoterieService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "family.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveCoterie(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *CoterieService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyCoterieList) error {
	path := "family.getListByFilter"
	inLog(path, in)
	var list []*cache.CoterieInfo
	var err error
	if in.Key == "user" {
		list,err = cache.Context().GetCoteriesByMember(in.Value)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.CoterieInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchCoterie(value))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CoterieService) UpdateBase(ctx context.Context, in *pb.ReqCoterieUpdate, out *pb.ReplyInfo) error {
	path := "family.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetCoterie(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.Remark, in.Passwords, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *CoterieService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "family.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetCoterie(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "passwords" {
		err = info.UpdatePasswords(in.Value, in.Operator)
	}else if in.Key == "master" {
		err = info.UpdateMaster(in.Value, in.Operator)
	}else if in.Key == "sn" {
		err = info.UpdateMaster(in.Value, in.Operator)
	} else if in.Key == "identify" {
		err = info.UpdateMemberIdentify(in.Operator, in.Value, in.Values[0])
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return err
}

func (mine *CoterieService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "family.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetCoterie(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateStatus(in.Operator, uint8(in.Flag))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CoterieService) AppendMember(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCoterieMembers) error {
	path := "family.appendMember"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetCoterie(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendMember(in.User, in.Name, in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.IdentifyInfo, 0, len(info.Members))
	for _, member := range info.Members {
		out.List = append(out.List, &pb.IdentifyInfo{User: member.User, Name: member.Name, Remark: member.Remark})
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *CoterieService) SubtractMember(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCoterieMembers) error {
	path := "family.subtractMember"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetCoterie(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractMember(in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.IdentifyInfo, 0, len(info.Members))
	for _, member := range info.Members {
		out.List = append(out.List, &pb.IdentifyInfo{User: member.User, Name: member.Name, Remark: member.Remark})
	}
	out.Status = outLog(path, out)
	return nil
}
