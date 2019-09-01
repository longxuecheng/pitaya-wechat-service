CREATE TABLE IF NOT EXISTS `mymall`.`goods_express_constraint` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `goods_id` BIGINT NOT NULL DEFAULT 0,
  `is_free` TINYINT NOT NULL DEFAULT 0 COMMENT '是否包邮',
  `province_id` INT NOT NULL DEFAULT 0,
  `express_fee` DECIMAL(10,2) NOT NULL DEFAULT 0.00 COMMENT '运费',
  `stock_id` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;