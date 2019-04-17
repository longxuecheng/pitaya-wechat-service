-- table cart --
alter table cart drop sale_unit_price;
alter table cart drop market_price;
alter table cart add supplier_id bigint(20) not null default 0;

-- table stock --

alter table stock add supplier_id bigint(20) not null default 0;6