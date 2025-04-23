CREATE DATABASE IF NOT EXISTS gorder_v2;
USE gorder_v2;

DROP TABLE IF EXISTS `o_stock`;

CREATE TABLE `o_stock` (
                           id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                           product_id VARCHAR(255) NOT NULL,
                           quantity INT UNSIGNED NOT NULL DEFAULT 0,
                           version INT NOT NULL DEFAULT 0,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO o_stock (product_id, quantity, version)
VALUES ('prod_S8KbxrN4dAEE0y', 1000, 0) ('prod_S8Kcs4H65OdCxg', 1000, 0)