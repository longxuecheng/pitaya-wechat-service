package region

type regionCategory string

const (
	East3 regionCategory = "East3"
	West6 regionCategory = "West6"
)

// East3ProviceList Northeast 3 provinces
var East3ProviceList = []int64{7, 8, 9}

// West6ProvinceList west 6 provinces
var West6ProvinceList = []int64{6, 29, 30, 31, 32, 27}

// CommonProvinceList province
var CommonProvinceList = []int64{2, 3, 4, 5, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 28, 33, 34, 35, 36}
