package dao

import (
	"fmt"
	"gotrue/sys"
	"testing"
)

func TestMain(m *testing.M) {
	Init()
	m.Run()
}

func TestDbCount(t *testing.T) {
	for i := 0; i < 100; i++ {
		total, err := CartDao.SelectCountByUserID(0)
		fmt.Println(sys.GetEasyDB().Stats())
		if err != nil {
			t.Error(err)
		}
		fmt.Println(total)
	}
}
