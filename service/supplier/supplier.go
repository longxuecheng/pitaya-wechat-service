package supplier

import (
	"database/sql"
	"gotrue/dao"
	"gotrue/dto/response"
	"gotrue/facility/errors"
	"gotrue/facility/strings"
)

var SupplierService *supplier

type supplier struct {
	dao *dao.Supplier
}

func initSupplierService() {
	if SupplierService != nil {
		return
	}
	SupplierService = &supplier{
		dao.SupplierDao,
	}
}

func (s *supplier) QueryByAdmin(adminID int64) (*response.Supplier, error) {
	supplier, err := s.dao.SelectByAdmin(adminID)
	if err == sql.ErrNoRows {
		return nil, errors.NewWithCodef("SupplierNotFound", "您不是供应商的管理员")
	}
	if err != nil {
		return nil, err
	}
	return &response.Supplier{
		ID:   supplier.ID,
		Name: strings.PtrValue(supplier.Name),
	}, nil
}
