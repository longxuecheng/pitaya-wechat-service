alter table user add avatar_url VARCHAR(256);
alter table user add city VARCHAR(16);
alter table user add country VARCHAR(16);
alter table user add gender TINYINT(1);
alter table user add nick_name VARCHAR(32);
alter table user add province VARCHAR(16);
alter table user modify name varchar(45) not null default '';
