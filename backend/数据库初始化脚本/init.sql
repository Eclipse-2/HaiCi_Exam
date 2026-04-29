SET NAMES utf8mb4;
CREATE DATABASE IF NOT EXISTS hospital_registry DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE hospital_registry;

-- 1. 用户表
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `phone` varchar(20) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `role` varchar(20) DEFAULT 'USER' COMMENT 'USER, ADMIN',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 就诊人表
CREATE TABLE `patients` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `name` varchar(50) NOT NULL,
  `id_card` varchar(20) NOT NULL,
  `phone` varchar(20) NOT NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_id_card` (`id_card`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 科室表
CREATE TABLE `departments` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. 医生表
CREATE TABLE `doctors` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `dept_id` bigint(20) NOT NULL,
  `name` varchar(50) NOT NULL,
  `title` varchar(50) NOT NULL COMMENT '职称',
  `specialty` varchar(255) DEFAULT '' COMMENT '专长',
  `description` text COMMENT '简介',
  PRIMARY KEY (`id`),
  KEY `idx_dept_id` (`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5. 排班表 (号源)
CREATE TABLE `schedules` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `doctor_id` bigint(20) NOT NULL,
  `date` date NOT NULL,
  `session` varchar(20) NOT NULL COMMENT 'MORNING, AFTERNOON',
  `total_slots` int(11) NOT NULL DEFAULT 0,
  `available_slots` int(11) NOT NULL DEFAULT 0,
  `status` varchar(20) DEFAULT 'AVAILABLE' COMMENT 'AVAILABLE, FULL, SUSPENDED',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_doc_date_session` (`doctor_id`,`date`,`session`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 6. 预约记录表
CREATE TABLE `appointments` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `appointment_no` varchar(50) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `patient_id` bigint(20) NOT NULL,
  `schedule_id` bigint(20) NOT NULL,
  `status` varchar(20) DEFAULT 'PENDING' COMMENT 'PENDING(待就诊), FINISHED(已完成), CANCELLED(已取消)',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_appointment_no` (`appointment_no`),
  KEY `idx_user_idx` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入一些演示数据
INSERT INTO `departments` (name, description) VALUES ('内科', '内科系统疾病'), ('外科', '外科系统疾病'), ('儿科', '儿童常见病');
INSERT INTO `doctors` (dept_id, name, title, specialty, description) VALUES (1, '张医生', '主任医师', '心内科', '擅长心脑血管疾病');
INSERT INTO `schedules` (doctor_id, date, session, total_slots, available_slots, status) VALUES (1, CURDATE() + INTERVAL 1 DAY, 'MORNING', 10, 10, 'AVAILABLE');
