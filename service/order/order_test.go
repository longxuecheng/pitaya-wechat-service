package order

import (
	"fmt"
	"gotrue/model"
	"testing"

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
