package express

import (
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
	"strings"
	"testing"
)

func TestExpressInfo(t *testing.T) {
	expressInfo, err := ExpressService.ExpressInfo("YZ", "9896590990776")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("express detail is %+v\n", expressInfo)
	for i, e := range expressInfo.Traces {
		fmt.Printf("express trace %d is % v \n", i, e)
	}
}

func TestConfirmBlock(t *testing.T) {
	pngReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader("post"))
	img, err := png.Decode(pngReader)
	if err != nil {
		log.Fatal(err)
	}
	leftTop := img.Bounds().Min
	rightBotom := img.Bounds().Max
	found := false
	for x := leftTop.X; x < rightBotom.X; x++ {
		if !found {
			for y := leftTop.Y; y < rightBotom.Y; y++ {
				f1 := false
				f2 := false

				r, g, b, a := img.At(x, y).RGBA()

				if r == 65535 && g == 65535 && b == 65535 && a == 65535 {
					f1 = true
				}
				r1, g1, b1, a1 := img.At(x, y+10).RGBA()
				if r1 == 65535 && g1 == 65535 && b1 == 65535 && a1 == 65535 {
					f2 = true
				}
				if f1 && f2 {
					fmt.Printf("Found start x position %d y position %d\n", x, y)
					found = true
					break
				}
			}
		}

	}
}

func TestSlideDecoder(t *testing.T) {
	slideDecode := NewSlideDecoder()
	err := slideDecode.LoadVerifyCode()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	err = slideDecode.CheckStartPosition()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	traces, err := slideDecode.QueryExpress("9896590990776")
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	for _, trace := range traces {
		fmt.Printf("%+v\n", trace)
	}

}
