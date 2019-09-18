ALTER TABLE goods_express_constraint add reachable TINYINT not NULL DEFAULT 1;

ALTER TABLE supplier DROP COLUMN admin_id;

CREATE TABLE IF NOT EXISTS `supplier_admin` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL DEFAULT 0,
  `priority` INT NOT NULL DEFAULT 1,
  `supplier_id` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;