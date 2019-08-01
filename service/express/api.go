package express

import "gotrue/facility/errors"

type ExpressMethod string

const (
	ExpressMethodZTO  ExpressMethod = "ZTO"
	ExpressMethodSTO  ExpressMethod = "STO"
	ExpressMethodYTO  ExpressMethod = "YTO"
	ExpressMethodEMS  ExpressMethod = "EMS"
	ExpressMethodYDA  ExpressMethod = "YDA"
	ExpressMethodBSHT ExpressMethod = "BSHT"
)

var expressMethodMap = map[ExpressMethod]string{
	ExpressMethodZTO:  "中通快递",
	ExpressMethodSTO:  "申通快递",
	ExpressMethodYTO:  "圆通快递",
	ExpressMethodEMS:  "EMS",
	ExpressMethodYDA:  "韵达快递",
	ExpressMethodBSHT: "百世汇通快递",
}

func (e ExpressMethod) String() string {
	return string(e)
}

func (e ExpressMethod) Name() string {
	if method, ok := expressMethodMap[e]; ok {
		return method
	}
	return ""
}

// IsSupport check wether the given method was supported now
func IsSupport(method string) error {
	_, found := expressMethodMap[ExpressMethod(method)]
	if !found {
		return errors.NewWithCodef("UnsupportedExpress", "不支持的快递")
	}
	return nil
}
