CREATE TABLE IF NOT EXISTS `settlement` (
  `id` BIGINT(16) NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(64) NOT NULL,
  `total_price` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  `create_time` TIMESTAMP NOT NULL DEFAULT now(),
  `update_time` TIMESTAMP NULL,
  `delete_time` TIMESTAMP NULL,
  `is_delete` TINYINT NOT NULL DEFAULT 0,
  `supplier_id` BIGINT(64) NOT NULL DEFAULT 0,
  `user_id` BIGINT(64) NOT NULL DEFAULT 0,
  `method` VARCHAR(8) NULL COMMENT '结算方式',
  PRIMARY KEY (`id`))
ENGINE = InnoDB;