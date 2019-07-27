alter table goods_category modify COLUMN img_url varchar(255) not null default '';
alter table goods_category modify COLUMN banner_url varchar(255) not null default '';
alter table goods_category modify COLUMN wap_banner_url varchar(255) not null default '';
alter table goods_category modify COLUMN level varchar(16) not null default '';
alter table goods_category modify COLUMN front_name varchar(64) not null default '';