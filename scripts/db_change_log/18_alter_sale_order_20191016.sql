alter table sale_order add discount_amt DECIMAL(12,2) NOT NULL DEFAULT 0.00;
alter table sale_order add discount_type VARCHAR(16) NOT NULL DEFAULT 'None';
