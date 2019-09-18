CREATE TABLE IF NOT EXISTS `mymall`.`banner` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(16) NOT NULL DEFAULT '',
  `src` VARCHAR(128) NOT NULL,
  `type` VARCHAR(16) NOT NULL DEFAULT 'banner',
  `link` VARCHAR(64) NOT NULL DEFAULT '',
  `create_time` TIMESTAMP NOT NULL DEFAULT now(),
  `update_time` TIMESTAMP NULL,
  `delete_time` TIMESTAMP NULL,
  `is_delete` TINYINT NOT NULL DEFAULT 0,
  `online_time` TIMESTAMP NOT NULL DEFAULT now() COMMENT '上线时间',
  `online_duration` BIGINT NOT NULL DEFAULT -1 COMMENT '上线持续时间单位/秒 -1代表无限制',
  `is_online` TINYINT NOT NULL DEFAULT 1,
  `priority` INT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;

INSERT INTO banner (src,priority) VALUES ('https://glxy-goods-1258625730.cos.ap-chengdu.myqcloud.com/banner00.jpeg',1),
('https://glxy-goods-1258625730.cos.ap-chengdu.myqcloud.com/b02.png',2),
('https://glxy-goods-1258625730.cos.ap-chengdu.myqcloud.com/b01.png',3);
