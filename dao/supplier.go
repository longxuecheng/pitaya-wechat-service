package dao

import (
	"gotrue/facility/utils"
	"gotrue/model"
	"gotrue/sys"

	"github.com/Masterminds/squirrel"
)

var SupplierDao *Supplier

type Supplier struct {
	columns []string
	table   string
	*sys.EasyDB
}

func initSupplierDao() {
	if SupplierDao != nil {
		return
	}
	splr := &model.Supplier{}
	SupplierDao = &Supplier{
		utils.TagValues(splr, "db"),
		splr.TableName(),
		sys.GetEasyDB(),
	}
}

func (s *Supplier) SelectByAdmin(adminID int64) (*model.Supplier, error) {
	supplier := &model.Supplier{}
	err := s.SelectOneDSL(supplier, s.columns, s.table, squirrel.Eq{"admin_id": adminID})
	return supplier, err
}
