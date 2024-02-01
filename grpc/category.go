package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-assignment/proto/assignment"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.assignment/cache"
)

type CategoryService struct{}

func switchCategory(info *cache.CategoryInfo) *pb.CategoryInfo {
	tmp := new(pb.CategoryInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Parent = info.Parent
	tmp.Source = info.Quote
	tmp.Weight = info.Weight
	children, _ := cache.Context().GetCategoriesByParent(info.UID)
	tmp.Children = make([]*pb.CategoryInfo, 0, len(children))
	for _, child := range children {
		v := switchCategory(child)
		tmp.Children = append(tmp.Children, v)
	}
	return tmp
}

func (mine *CategoryService) AddOne(ctx context.Context, in *pb.ReqCategoryAdd, out *pb.ReplyCategoryOne) error {
	path := "category.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().NewCategory(in.Name, in.Parent, in.Source, in.Operator, in.Weight)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCategory(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CategoryService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCategoryOne) error {
	path := "category.getOne"
	inLog(path, in)

	var info *cache.CategoryInfo
	var err error
	info, err = cache.Context().GetOneCategory(in.Uid)

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchCategory(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CategoryService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "category.getStatistic"
	inLog(path, in)
	if len(in.Key) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CategoryService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "category.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetOneCategory(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	err = info.Delete(in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *CategoryService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyCategoryList) error {
	path := "category.getListByFilter"
	inLog(path, in)
	var list []*cache.CategoryInfo
	var err error

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	if in.Key == "scene" {
		list, err = cache.Context().GetCategoriesByScene(in.Value)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.List = make([]*pb.CategoryInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchCategory(value))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CategoryService) UpdateBase(ctx context.Context, in *pb.ReqCategoryUpdate, out *pb.ReplyCategoryOne) error {
	path := "category.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetOneCategory(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.Update(in.Name, in.Remark, in.Source, in.Operator, in.Weight)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCategory(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CategoryService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "category.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	//info, er := cache.Context().GetOneCategory(in.Uid)
	//if er != nil {
	//	out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
	//	return nil
	//}
	var err error

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}
