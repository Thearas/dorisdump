CREATE TABLE `my_db`.`complex_table` (
  `id` INT UNSIGNED NOT NULL,
  `product_name` VARCHAR(255) DEFAULT NULL,
  `price` DECIMAL(10,2) DEFAULT 0.00,
  `description` TEXT,
  `is_active` BOOLEAN DEFAULT TRUE,
  `code` CHAR(10),
  `stock` BIGINT DEFAULT 0,
  `last_updated` TIMESTAMP
);
