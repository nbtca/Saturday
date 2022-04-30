/*
 Navicat Premium Data Transfer

 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 80022
 Source Host           : localhost:3306
 Source Schema         : saturday_dev

 Target Server Type    : MySQL
 Target Server Version : 80022
 File Encoding         : 65001

 Date: 30/04/2022 23:10:14
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for client
-- ----------------------------
DROP TABLE IF EXISTS `client`;
CREATE TABLE `client` (
  `client_id` bigint NOT NULL AUTO_INCREMENT,
  `openid` char(28) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '微信',
  `gmt_create` datetime NOT NULL,
  `gmt_modified` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`client_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of client
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for event
-- ----------------------------
DROP TABLE IF EXISTS `event`;
CREATE TABLE `event` (
  `client_id` bigint NOT NULL,
  `event_id` bigint NOT NULL AUTO_INCREMENT,
  `model` varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '型号',
  `phone` varchar(11) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `qq` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '',
  `contact_preference` varchar(20) NOT NULL DEFAULT 'qq' COMMENT '联系偏好',
  `event_description` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '事件（用户）描述',
  `repair_description` varchar(1000) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '维修描述',
  `member_id` char(10) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '最后由谁维修',
  `closed_by` char(10) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL COMMENT '由谁关闭',
  `gmt_create` datetime NOT NULL,
  `gmt_modified` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`event_id`) USING BTREE,
  KEY `fk_Event_Admin_2` (`closed_by`) USING BTREE,
  KEY `fk_Event_User_1` (`client_id`) USING BTREE,
  KEY `fk_Event_repairElements_1` (`member_id`) USING BTREE,
  CONSTRAINT `event_ibfk_1` FOREIGN KEY (`closed_by`) REFERENCES `member` (`member_id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `event_ibfk_2` FOREIGN KEY (`client_id`) REFERENCES `client` (`client_id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_Event_repairElements_1` FOREIGN KEY (`member_id`) REFERENCES `member` (`member_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of event
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for event_action
-- ----------------------------
DROP TABLE IF EXISTS `event_action`;
CREATE TABLE `event_action` (
  `event_action_id` tinyint NOT NULL,
  `action` varchar(30) NOT NULL DEFAULT '',
  PRIMARY KEY (`event_action_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of event_action
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for event_action_relation
-- ----------------------------
DROP TABLE IF EXISTS `event_action_relation`;
CREATE TABLE `event_action_relation` (
  `event_log_id` bigint NOT NULL AUTO_INCREMENT,
  `event_action_id` tinyint NOT NULL,
  PRIMARY KEY (`event_log_id`,`event_action_id`),
  KEY `fk_event_action_relation_event_action_1` (`event_action_id`),
  CONSTRAINT `fk_event_action_relation_event_action_1` FOREIGN KEY (`event_action_id`) REFERENCES `event_action` (`event_action_id`),
  CONSTRAINT `fk_event_action_relation_event_log_1` FOREIGN KEY (`event_log_id`) REFERENCES `event_log` (`event_log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of event_action_relation
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for event_event_status_relation
-- ----------------------------
DROP TABLE IF EXISTS `event_event_status_relation`;
CREATE TABLE `event_event_status_relation` (
  `event_status_id` tinyint NOT NULL,
  `event_id` bigint NOT NULL,
  PRIMARY KEY (`event_status_id`,`event_id`),
  KEY `fk_event_event_status_relation_event_1` (`event_id`),
  CONSTRAINT `fk_event_event_status_relation_event_1` FOREIGN KEY (`event_id`) REFERENCES `event` (`event_id`),
  CONSTRAINT `fk_event_event_status_relation_event_status_1` FOREIGN KEY (`event_status_id`) REFERENCES `event_status` (`event_status_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of event_event_status_relation
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for event_log
-- ----------------------------
DROP TABLE IF EXISTS `event_log`;
CREATE TABLE `event_log` (
  `event_log_id` bigint NOT NULL AUTO_INCREMENT,
  `event_id` bigint NOT NULL,
  `description` varchar(255) DEFAULT '',
  `member_id` char(10) DEFAULT NULL,
  `gmt_create` datetime NOT NULL,
  PRIMARY KEY (`event_log_id`),
  KEY `fk_event_log_element_1` (`member_id`),
  KEY `fk_event_log_event_1` (`event_id`),
  CONSTRAINT `fk_event_log_element_1` FOREIGN KEY (`member_id`) REFERENCES `member` (`member_id`),
  CONSTRAINT `fk_event_log_event_1` FOREIGN KEY (`event_id`) REFERENCES `event` (`event_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of event_log
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for event_status
-- ----------------------------
DROP TABLE IF EXISTS `event_status`;
CREATE TABLE `event_status` (
  `event_status_id` tinyint NOT NULL,
  `status` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`event_status_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of event_status
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for member
-- ----------------------------
DROP TABLE IF EXISTS `member`;
CREATE TABLE `member` (
  `member_id` char(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `alias` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '昵称',
  `password` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '',
  `name` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '',
  `section` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '班级（计算机196）',
  `profile` varchar(1000) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '个人简介',
  `phone` varchar(11) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '',
  `qq` varchar(20) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '',
  `avatar` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '头像地址',
  `created_by` char(10) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '由谁添加',
  `gmt_create` datetime NOT NULL,
  `gmt_modified` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`member_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of member
-- ----------------------------
BEGIN;
INSERT INTO `member` VALUES ('0000000000', '管理', '000000', '管理', '计算机000', '', '', '', '', '', '2022-04-30 17:28:42', '2022-04-30 17:28:44');
INSERT INTO `member` VALUES ('2333333333', '滑稽', '123456', '滑稽', '计算机233', 'relaxing', '12356839487', '123456', '', '0000000000', '2022-04-23 15:49:59', '2022-04-30 17:29:46');
INSERT INTO `member` VALUES ('3000000001', '小稽', '', '滑小稽', '计算机233', '。。。', '', '123456', '', '2333333333', '2022-04-30 23:06:44', '2022-04-30 23:06:44');
COMMIT;

-- ----------------------------
-- Table structure for member_role_relation
-- ----------------------------
DROP TABLE IF EXISTS `member_role_relation`;
CREATE TABLE `member_role_relation` (
  `member_id` varchar(10) NOT NULL,
  `role_id` tinyint NOT NULL,
  PRIMARY KEY (`member_id`,`role_id`),
  KEY `fk_member_role_role_1` (`role_id`),
  CONSTRAINT `fk_member_role_member_1` FOREIGN KEY (`member_id`) REFERENCES `member` (`member_id`),
  CONSTRAINT `fk_member_role_role_1` FOREIGN KEY (`role_id`) REFERENCES `role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of member_role_relation
-- ----------------------------
BEGIN;
INSERT INTO `member_role_relation` VALUES ('2333333333', 2);
INSERT INTO `member_role_relation` VALUES ('0000000000', 3);
COMMIT;

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `role_id` tinyint NOT NULL,
  `role` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of role
-- ----------------------------
BEGIN;
INSERT INTO `role` VALUES (0, 'member_inactive');
INSERT INTO `role` VALUES (1, 'admin_inavtive');
INSERT INTO `role` VALUES (2, 'member');
INSERT INTO `role` VALUES (3, 'admin');
COMMIT;

-- ----------------------------
-- Table structure for setting
-- ----------------------------
DROP TABLE IF EXISTS `setting`;
CREATE TABLE `setting` (
  `setting` varchar(10000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of setting
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for status
-- ----------------------------
DROP TABLE IF EXISTS `status`;
CREATE TABLE `status` (
  `status_id` tinyint NOT NULL,
  `status` varchar(255) NOT NULL,
  PRIMARY KEY (`status_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of status
-- ----------------------------
BEGIN;
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
