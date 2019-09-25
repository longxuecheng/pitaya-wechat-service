package test

import (
	"fmt"
	"gotrue/service"
	"testing"
)

func TestGeo(t *testing.T) {
	g := service.GaodeMap{}
	geo, err := g.AddressLocation("", "安慧里一区")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", geo)
}

func TestInputTip(t *testing.T) {
	s := service.GaodeMap{}
	resp, err := s.AddressTips("安徽里一区", "", "")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)
}
