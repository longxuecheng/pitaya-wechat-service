CREATE TABLE IF NOT EXISTS `mymall`.`cut_order` (
  `id` BIGINT(16) NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT(16) NOT NULL DEFAULT 0,
  `cut_no` VARCHAR(45) NOT NULL COMMENT '砍价码',
  `goods_id` BIGINT(16) NOT NULL DEFAULT 0,
  `stock_id` BIGINT(16) NOT NULL DEFAULT 0,
  `create_time` TIMESTAMP NOT NULL DEFAULT now(),
  `expire_time` TIMESTAMP NULL,
  `consume_time` TIMESTAMP NULL COMMENT '使用时间',
  `consumed` TINYINT NOT NULL DEFAULT 0,
  `sale_order_id` BIGINT(16) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `mymall`.`cut_detail` (
  `id` BIGINT(16) NOT NULL AUTO_INCREMENT,
  `cut_order_id` BIGINT(16) NOT NULL DEFAULT 0,
  `user_id` BIGINT(16) NOT NULL DEFAULT 0,
  `cut_price` DECIMAL(12,2) NOT NULL DEFAULT 0.0,
  `is_delete` TINYINT NOT NULL DEFAULT 0,
  `create_time` TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`))
ENGINE = InnoDB;