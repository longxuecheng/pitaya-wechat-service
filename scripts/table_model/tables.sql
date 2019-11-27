CREATE TABLE IF NOT EXISTS `mymall`.`activity` (
  `id` BIGINT(16) NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(45) NOT NULL,
  `banner_url` VARCHAR(128) NOT NULL,
  `page_route` VARCHAR(45) NOT NULL COMMENT '类似这种/pages/activities/coupon/coupon',
  `type` VARCHAR(16) NOT NULL,
  `bg_url` VARCHAR(128) NULL,
  `start_time` TIMESTAMP NOT NULL DEFAULT now(),
  `expire_time` TIMESTAMP NOT NULL DEFAULT now(),
  `is_delete` TINYINT NOT NULL DEFAULT 0,
  `is_online` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `mymall`.`activity_coupon` (
  `id` BIGINT(16) NOT NULL AUTO_INCREMENT,
  `activity_id` BIGINT(16) NOT NULL DEFAULT 0,
  `coupon_type` VARCHAR(45) NOT NULL,
  `coupon_price` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `composable_type` VARCHAR(16) NOT NULL,
  `total_quantity` INT NOT NULL DEFAULT 0,
  `available_quantity` INT NOT NULL DEFAULT 0,
  `category_id` BIGINT(16) NOT NULL DEFAULT 0,
  `goods_id` BIGINT(16) NOT NULL DEFAULT 0,
  `create_time` TIMESTAMP NOT NULL DEFAULT now(),
  PRIMARY KEY (`id`))
ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `mymall`.`coupon` (
  `id` BIGINT(16) NOT NULL AUTO_INCREMENT,
  `coupon_no` VARCHAR(45) NOT NULL COMMENT '优惠券码,唯一',
  `price` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `create_time` TIMESTAMP NOT NULL DEFAULT now(),
  `consumed` TINYINT NOT NULL DEFAULT 0 COMMENT '是否已经使用',
  `expire_time` TIMESTAMP NOT NULL DEFAULT now(),
  `type` VARCHAR(16) NOT NULL COMMENT '优惠券类型(全品类/特定分类/特定商品)',
  `category_id` BIGINT(16) NOT NULL DEFAULT 0 COMMENT '优惠券类型为特定分类时使用',
  `goods_id` BIGINT(16) NULL COMMENT '优惠券为特定商品时可用',
  `composable_type` VARCHAR(16) NOT NULL DEFAULT 'None' COMMENT '可叠加使用优惠券类型(默认不可以组合使用)',
  `user_id` BIGINT(64) NULL,
  `consumed_time` TIMESTAMP NULL COMMENT '使用时间',
  `sale_order_id` BIGINT(16) NOT NULL DEFAULT 0 COMMENT '使用该优惠券的订单ID',
  `activity_id` BIGINT(16) NOT NULL DEFAULT 0 COMMENT '优惠券活动来源',
  `received` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;