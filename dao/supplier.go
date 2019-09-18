package dao

import (
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
		splr.Columns(),
		splr.TableName(),
		sys.GetEasyDB(),
	}
}

func (s *Supplier) SelectByAdmin(adminID int64) (*model.Supplier, error) {
	supplier := &model.Supplier{}
	err := s.SelectOneDSL(supplier, s.columns, s.table, squirrel.Eq{"admin_id": adminID})
	return supplier, err
}

func (s *Supplier) SelectByID(id int64) (*model.Supplier, error) {
	supplier := &model.Supplier{}
	err := s.SelectOneDSL(supplier, s.columns, s.table, squirrel.Eq{"id": id})
	return supplier, err
}
