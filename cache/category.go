package cache

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.assignment/proxy/nosql"
	"time"
)

const DefaultParent = "0"

type CategoryInfo struct {
	baseInfo
	Remark string
	Parent string
	Weight uint32
	Quote  string
}

func (mine *cacheContext) NewCategory(name, parent, source, operator string, weight uint32) (*CategoryInfo, error) {
	db := new(nosql.Category)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetQuestionNextID()
	db.CreatedTime = time.Now()
	db.Operator = operator
	db.Name = name
	db.Parent = DefaultParent
	if parent != "" {
		db.Parent = parent
	}
	db.Owner = DefaultOwner
	db.Quote = source
	db.Weight = weight
	if db.Weight > 0 {
		categoryList, err := nosql.GetCategoryListByParent(db.Parent)
		if err != nil {
			return nil, err
		}
		if (len(categoryList) + 1) < int(weight) {
			return nil, err
		}
		for _, v := range categoryList {
			if v.Weight >= db.Weight {
				err1 := nosql.UpdateCategoryInt("weight", "", v.UID.Hex(), int64(v.Weight+1))
				if err1 != nil {
					return nil, err1
				}
			}
		}
	}
	err := nosql.CreateCategory(db)
	if err != nil {
		return nil, err
	}
	info := new(CategoryInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetOneCategory(uid string) (*CategoryInfo, error) {
	category, err := nosql.GetOneCategory(uid)
	if err != nil {
		return nil, err
	}
	//判断数据是否是被标记删除
	if category.DeleteTime.Format("2006") != "0001" {
		return nil, nil
	}
	info := new(CategoryInfo)
	info.initInfo(category)
	return info, nil
}

func (mine *cacheContext) GetCategoriesByParent(parent string) ([]*CategoryInfo, error) {
	if parent == "" {
		return nil, errors.New("parent is null")
	}
	array, err := nosql.GetCategoryListByParent(parent)
	if err != nil {
		return nil, err
	}
	list := make([]*CategoryInfo, 0, len(array))
	for _, v := range array {
		info := new(CategoryInfo)
		info.initInfo(v)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetTopCategoriesByScene(owner string) ([]*CategoryInfo, error) {
	var array []*nosql.Category
	var err error
	if len(owner) > 2 {
		array, err = nosql.GetCategoryListByOwner(owner, DefaultParent)
	} else {
		array, err = nosql.GetCategoryListByOwner(DefaultOwner, DefaultParent)
	}

	if err != nil {
		return nil, err
	}
	list := make([]*CategoryInfo, 0, len(array))
	for _, v := range array {
		info := new(CategoryInfo)
		info.initInfo(v)
		list = append(list, info)
	}
	return list, nil
}

func (mine *CategoryInfo) initInfo(db *nosql.Category) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Name = db.Name
	mine.Parent = db.Parent
	mine.Quote = db.Quote
	mine.Weight = db.Weight
	mine.Remark = db.Remark
}
func (mine *CategoryInfo) Export(dir string) {

}
func (mine *CategoryInfo) Update(name, remark, quote, operator string, weight uint32) error {
	err := nosql.UpdateCategoryBase(mine.UID, name, remark, quote, operator)
	if err != nil {
		return err
	}
	mine.Remark = remark
	mine.Quote = remark
	mine.Name = name
	mine.Operator = operator
	if weight != 0 {
		arry, err := nosql.GetCategoryListByParent(mine.Parent)
		if err != nil {
			return err
		}
		if uint32(len(arry)) < weight {
			return errors.New("the newWeight is exceed the limit")
		}
		//小变大
		if weight > mine.Weight {
			for _, v := range arry {
				if v.Weight > mine.Weight && v.Weight <= weight {
					err := nosql.UpdateCategoryInt("weight", "", v.UID.Hex(), int64(v.Weight-1))
					if err != nil {
						return errors.New("the newWeight is exceed the limit")
					}
				}
			}
			err = nosql.UpdateCategoryInt("weight", "", mine.UID, int64(weight))
			if err != nil {
				return errors.New("the newWeight update is err")
			}
			//mine.Weight = weight
			//return nil
		} else if weight < mine.Weight {
			//大变小
			for _, v := range arry {
				if v.Weight >= weight && v.Weight < mine.Weight {
					err := nosql.UpdateCategoryInt("weight", "", v.UID.Hex(), int64(v.Weight+1))
					if err != nil {
						return errors.New("the newWeight is exceed the limit")
					}
				}
			}
			err = nosql.UpdateCategoryInt("weight", "", mine.UID, int64(weight))
			if err != nil {
				return errors.New("the newWeight update is err")
			}
			//mine.Weight = weight
		}
		//if len(arry) > 0 {
		//	for _, v := range arry {
		//		if weight < oldWeight {
		//			if weight <= v.Weight && v.Weight < oldWeight {
		//				err := nosql.UpdateCategoryInt("weight", "", v.UID.Hex(), int64(v.Weight+1))
		//				if err != nil {
		//					return errors.New("the newWeight is exceed the limit")
		//				}
		//			}
		//		} else if weight > oldWeight {
		//			if weight < v.Weight && v.Weight <= weight {
		//				err := nosql.UpdateCategoryInt("weight", "", v.UID.Hex(), int64(v.Weight-1))
		//				if err != nil {
		//					return errors.New("the newWeight is exceed the limit")
		//				}
		//			}
		//		} else {
		//			return errors.New("the newWeight is exceed the limit")
		//		}
		//	}
		//	err = nosql.UpdateCategoryInt("weight", "", mine.UID, int64(weight))
		//	if err != nil {
		//		return errors.New("the newWeight update is err")
		//	}
		//}
		mine.Weight = weight
	}
	return nil
}
func (mine *CategoryInfo) Delete(operator string) error {
	list, err := nosql.GetCategoryListByParent(mine.UID)
	if err != nil {
		return err
	}
	if len(list) > 0 {
		return errors.New(fmt.Sprintf("the children count = %d", len(list)))
	}
	if mine.Parent == DefaultParent {
		err = nosql.DeleteCategory(mine.UID, operator)
		if err != nil {
			return err
		}
		return nil
	}
	num := nosql.GetQuestionCount(mine.UID)
	if num > 0 {
		return errors.New(fmt.Sprintf("the children question count = %d", num))
	}
	err = nosql.DeleteCategory(mine.UID, operator)
	if err != nil {
		return nil
	}
	infos, err := nosql.GetCategoryListByParent(mine.Parent)
	if err != nil {
		return nil
	}
	for _, v := range infos {
		if v.Weight > mine.Weight {
			err1 := nosql.UpdateCategoryInt("weight", "", v.UID.Hex(), int64(v.Weight-1))
			if err1 != nil {
				return err1
			}
		}
	}
	return nil
}
