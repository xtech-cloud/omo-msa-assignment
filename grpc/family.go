package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.assignment/cache"
)

type FamilyService struct{}

func switchFamily(info *cache.FamilyInfo) *pb.FamilyInfo {
	tmp := new(pb.FamilyInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Location = info.Location
	tmp.Master = info.Master
	tmp.Sn = info.SN
	tmp.Address = info.Address
	tmp.Region = info.Region
	tmp.Status = uint32(info.Status)
	tmp.Location = info.Location
	tmp.Passwords = info.Passwords
	tmp.Assistants = info.Assistants
	tmp.Children = info.Children
	tmp.Tags = info.Tags
	tmp.Agents = info.Agents
	tmp.Members = make([]*pb.IdentifyInfo, 0, len(info.Members))
	for _, member := range info.Members {
		tmp.Members = append(tmp.Members, &pb.IdentifyInfo{User: member.User, Remark: member.Remark})
	}
	return tmp
}

func (mine *FamilyService) AddOne(ctx context.Context, in *pb.ReqFamilyAdd, out *pb.ReplyFamilyInfo) error {
	path := "family.addOne"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateFamily(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFamily(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FamilyService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFamilyInfo) error {
	path := "family.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetFamily(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchFamily(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FamilyService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFamilyList) error {
	path := "family.search"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *FamilyService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "family.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the key is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *FamilyService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "family.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveFamily(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *FamilyService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyFamilyList) error {
	path := "family.getList"
	inLog(path, in)


	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *FamilyService) UpdateBase(ctx context.Context, in *pb.ReqFamilyUpdate, out *pb.ReplyInfo) error {
	path := "family.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetFamily(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return err
}

func (mine *FamilyService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "family.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetFamily(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "passwords" {
		err = info.UpdatePasswords(in.Operator, in.Value)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return err
}

func (mine *FamilyService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "family.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetFamily(in.Uid)
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

func (mine *FamilyService) AppendMember(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFamilyMembers) error {
	path := "family.appendMember"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetFamily(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendMember(in.Operator, in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.IdentifyInfo, 0, len(info.Members))
	for _, member := range info.Members {
		out.List = append(out.List, &pb.IdentifyInfo{User: member.User, Remark: member.Remark})
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *FamilyService) SubtractMember(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFamilyMembers) error {
	path := "family.subtractMember"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info,er := cache.Context().GetFamily(in.Uid)
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
		out.List = append(out.List, &pb.IdentifyInfo{User: member.User, Remark: member.Remark})
	}
	out.Status = outLog(path, out)
	return nil
}
