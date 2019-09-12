ALTER TABLE sale_order add settlement_id BIGINT(64) NOT NULL DEFAULT 0;
ALTER TABLE sale_order add cost_amt DECIMAL(12,2) NOT NULL DEFAULT 0.00;