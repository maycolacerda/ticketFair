
CREATE TABLE IF NOT EXISTS users (
    `user_id` VARCHAR(36) PRIMARY KEY,
    `email` VARCHAR(100) NOT NULL UNIQUE,
    `password` VARCHAR(64) NOT NULL,
    `username` VARCHAR(100) NOT NULL UNIQUE,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS profiles (
    `profile_id` VARCHAR(36) PRIMARY KEY,
    `user_id` VARCHAR(36) NOT NULL UNIQUE,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100) NOT NULL,
    `phone_number` VARCHAR(15) NOT NULL UNIQUE,
    `address` TEXT,
    `verified_email` BOOLEAN DEFAULT FALSE,
    `verified_phone` BOOLEAN DEFAULT FALSE,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS merchants (
    `merchant_id` VARCHAR(36) PRIMARY KEY,
    `name` VARCHAR(100) NOT NULL UNIQUE,
    `description` TEXT,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS merchant_reps (
    `merchant_rep_id` VARCHAR(36) PRIMARY KEY,
    `merchant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(100) NOT NULL UNIQUE,
    `password` VARCHAR(64) NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (merchant_id) REFERENCES merchants(merchant_id) ON DELETE CASCADE
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS events (
    `event_id` VARCHAR(36) PRIMARY KEY,
    `merchant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `description` TEXT NOT NULL,
    `date` DATETIME NOT NULL,
    `location` VARCHAR(255) NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (merchant_id) REFERENCES merchants(merchant_id) ON DELETE CASCADE
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;