package dao

func Init() {
	initAttributeDao()
	initCartDao()
	initCategoryDao()
	initGoodsAttributeDao()
	initGoodsDao()
	initGoodsImageDao()
	initGoodsSpecDao()
	initRegionDao()
	initSaleDetailDao()
	initSaleOrderDao()
	initSpecificationDao()
	initStockDao()
	initSupplierDao()
	initUserAddressDao()
	initUserDao()
	initWechatPaymentDao()
	InitGoodsExpressConstraintDao()
}
