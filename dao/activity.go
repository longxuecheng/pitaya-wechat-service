package dao

import (
	"gotrue/model"

	"github.com/Masterminds/squirrel"
)

var ActivityDao *Activity

func initActivityDao() {
	m := &model.Activity{}
	ActivityDao = &Activity{
		table:   m.TableName(),
		columns: m.Columns(),
		db:      GetEasyDB(),
	}
}

type Activity struct {
	table   string
	columns []string
	db      *EasyDB
}

func (d *Activity) QueryOnlineActivities() ([]*model.Activity, error) {
	activities := []*model.Activity{}
	return activities, d.db.SelectDSL(&activities, d.columns, d.table, squirrel.Eq{"is_online": true})
}

func (d *Activity) QueryByID(id int64) (*model.Activity, error) {
	activity := &model.Activity{}
	return activity, d.db.SelectOneDSL(activity, d.columns, d.table, squirrel.Eq{"id": id, "is_online": true})
}
