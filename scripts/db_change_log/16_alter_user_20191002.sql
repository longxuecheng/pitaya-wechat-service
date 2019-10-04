ALTER TABLE user add channel_user_id BIGINT(16) NOT NULL DEFAULT 0;
ALTER TABLE user add bind_channel_time TIMESTAMP;
alter table user add channel_code varchar(45) not null default '';
