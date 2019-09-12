ALTER TABLE `user` CHARACTER SET = utf8mb4 , COLLATE = utf8mb4_unicode_ci;

ALTER TABLE `user` CHANGE COLUMN `nick_name` `nick_name` VARCHAR(32) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NULL DEFAULT NULL ;

