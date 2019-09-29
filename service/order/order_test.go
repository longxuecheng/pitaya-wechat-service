package order

import (
	"fmt"
	"gotrue/dto/response"
	"gotrue/model"
	"log"
	"strconv"
	"testing"

	"github.com/bwmarrin/snowflake"
	"github.com/shopspring/decimal"

	"github.com/looplab/fsm"
)

type Door struct {
	To  string
	FSM *fsm.FSM
}

func NewDoor(to string) *Door {
	d := &Door{
		To: to,
	}

	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{
			"enter_closed": func(e *fsm.Event) { d.enterState(e) },
			"enter_open":   func(e *fsm.Event) { d.enterState(e) },
		},
	)

	return d
}

func (d *Door) enterState(e *fsm.Event) {
	fmt.Printf("The door to %s is %s\n", d.To, e.Dst)
}

func TestFSMExample(t *testing.T) {
	door := NewDoor("heaven")
	door.FSM.SetState("open")
	fmt.Println(door.FSM.Can("open"))
	err := door.FSM.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	err = door.FSM.Event("close")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("final state is ", door.FSM.Current())
}

func TestOrderFSM(t *testing.T) {
	saleOrder := &model.SaleOrder{
		Status: model.Paid,
	}
	sof := newOrderFSM(saleOrder)
	err := sof.can("pay")
	if err != nil {
		fmt.Println("pay ", err)
	}
	err = sof.can("cancel")
	if err != nil {
		fmt.Println("cancel ", err)
	}
	err = sof.can("send")
	if err != nil {
		fmt.Println("send1 ", err)
	}
	sof.fsm.Event("send")
	fmt.Println(sof.fsm.Current())

	err = sof.can("send")
	if err != nil {
		fmt.Println("send2 ", err)
	}
	sof.fsm.Event("send")
	fmt.Println(sof.fsm.Current())
}

func TestStockOrder(t *testing.T) {
	quantity, _ := decimal.NewFromString("2")
	unitExpressFee, _ := decimal.NewFromString("10")
	cost_unit_price, _ := decimal.NewFromString("20")
	sale_unit_price, _ := decimal.NewFromString("30")
	sb := StockOrderBuilder{
		UserID:         999,
		Quantity:       quantity,
		UnitExpressFee: unitExpressFee,
		Address: &response.UserAddress{
			ID:         10,
			ProvinceID: 1,
			CityID:     2,
			DistrictID: 3,
			Address:    "address detail",
			Mobile:     "189117923194",
		},
		Goods: &model.Goods{
			ID:         1,
			Name:       "test_goods_name",
			SupplierID: 90,
		},
		Stock: &model.Stock{
			ID:            10,
			SupplierID:    90,
			GoodsID:       1,
			CostUnitPrice: cost_unit_price,
			SaleUnitPrice: sale_unit_price,
			Splitable:     true,
		},
	}
	stockOrderList, err := sb.Build()
	if err != nil {
		t.Error(err)
	}
	for _, stockOrder := range stockOrderList {
		fmt.Printf("[sale order is %+v\n", stockOrder.SaleOrder)
		for _, detail := range stockOrder.SaleDetails {
			fmt.Printf("details is %+v\n", detail)
		}
	}
}

func TestOrderNo12(t *testing.T) {
	o := model.SaleOrder{
		OrderNo: "0123456789012",
	}
	fmt.Println(o.OrderNo12())
}

func TestSnowFlake(t *testing.T) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		// Generate a snowflake ID.
		id := node.Generate()
		fmt.Println(strconv.FormatInt(id.Int64(), 10))
	}
}
