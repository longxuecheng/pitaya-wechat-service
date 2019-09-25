package dao

import (
	"gotrue/model"
	"gotrue/sys"

	"github.com/Masterminds/squirrel"
)

var SupplierAdminDao *SupplierAdmin

type SupplierAdmin struct {
	columns []string
	table   string
	*sys.EasyDB
}

func initSupplierAdminDao() {
	if SupplierAdminDao != nil {
		return
	}
	admin := &model.SupplierAdmin{}
	SupplierAdminDao = &SupplierAdmin{
		admin.Columns(),
		admin.TableName(),
		sys.GetEasyDB(),
	}
}

func (s *SupplierAdmin) QueryByUserID(userID int64) (*model.SupplierAdminCollection, error) {
	supplierAdmins := []*model.SupplierAdmin{}
	err := s.SelectDSL(&supplierAdmins, s.columns, s.table, squirrel.Eq{"user_id": userID})
	return &model.SupplierAdminCollection{
		Items: supplierAdmins,
	}, err
}

func (s *SupplierAdmin) QueryBySupplierID(supplierID int64) (*model.SupplierAdminCollection, error) {
	supplierAdmins := []*model.SupplierAdmin{}
	err := s.SelectDSL(&supplierAdmins, s.columns, s.table, squirrel.Eq{"supplier_id": supplierID})
	return &model.SupplierAdminCollection{
		Items: supplierAdmins,
	}, err
}
