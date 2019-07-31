alter table sale_order add parent_id bigint not null default 0;

drop table payment;

CREATE TABLE IF NOT EXISTS `mymall`.`wechat_payment` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `sale_order_id` BIGINT NOT NULL DEFAULT 0,
  `sale_order_no` VARCHAR(32) NOT NULL,
  `amount` DECIMAL(12,2) NOT NULL DEFAULT 0.00,
  `status` VARCHAR(16) NOT NULL COMMENT '成功/失败等',
  `create_time` TIMESTAMP NOT NULL DEFAULT now(),
  `update_time` TIMESTAMP NULL,
  `delete_time` TIMESTAMP NULL,
  `description` VARCHAR(45) NULL COMMENT '付款备注',
  `transaction_id` VARCHAR(32) NULL DEFAULT 0,
  `transaction_type` VARCHAR(1) NOT NULL DEFAULT 'D',
  PRIMARY KEY (`id`))
ENGINE = InnoDB;