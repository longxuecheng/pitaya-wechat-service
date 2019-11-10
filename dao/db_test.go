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

func TestStockColumns(t *testing.T) {
	c1 := []string{"x1", "x2"}
	c2 := c1
	c2 = append(c2, "x3")
	fmt.Printf("c1 is %+v\n", c1)
	fmt.Printf("c2 is %+v\n", c2)
}
