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

 Date: 05/06/2022 13:39:42
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of client
-- ----------------------------
BEGIN;
INSERT INTO `client` (`client_id`, `openid`, `gmt_create`, `gmt_modified`) VALUES (1, '', '2022-05-10 10:23:19', '2022-05-10 10:23:21');
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
  `problem` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '事件（用户）描述',
  `member_id` char(10) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '最后由谁维修',
  `closed_by` char(10) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' COMMENT '由谁关闭',
  `gmt_create` datetime NOT NULL,
  `gmt_modified` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`event_id`) USING BTREE,
  KEY `fk_Event_Admin_2` (`closed_by`) USING BTREE,
  KEY `fk_Event_User_1` (`client_id`) USING BTREE,
  KEY `fk_Event_repairElements_1` (`member_id`) USING BTREE,
  CONSTRAINT `event_ibfk_2` FOREIGN KEY (`client_id`) REFERENCES `client` (`client_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of event
-- ----------------------------
BEGIN;
INSERT INTO `event` (`client_id`, `event_id`, `model`, `phone`, `qq`, `contact_preference`, `problem`, `member_id`, `closed_by`, `gmt_create`, `gmt_modified`) VALUES (1, 1, '7590', '17557209007', '709196390', 'qq', 'hackintosh', '2333333333', '0000000000', '2022-05-10 10:23:54', '2022-06-02 16:18:52');
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
INSERT INTO `event_action` (`event_action_id`, `action`) VALUES (1, 'create');
INSERT INTO `event_action` (`event_action_id`, `action`) VALUES (2, 'accept');
INSERT INTO `event_action` (`event_action_id`, `action`) VALUES (3, 'cancel');
INSERT INTO `event_action` (`event_action_id`, `action`) VALUES (4, 'commit');
INSERT INTO `event_action` (`event_action_id`, `action`) VALUES (5, 'alterCommit');
INSERT INTO `event_action` (`event_action_id`, `action`) VALUES (6, 'drop');
INSERT INTO `event_action` (`event_action_id`, `action`) VALUES (7, 'close');
INSERT INTO `event_action` (`event_action_id`, `action`) VALUES (8, 'reject');
INSERT INTO `event_action` (`event_action_id`, `action`) VALUES (9, 'update');
COMMIT;

-- ----------------------------
-- Table structure for event_event_action_relation
-- ----------------------------
DROP TABLE IF EXISTS `event_event_action_relation`;
CREATE TABLE `event_event_action_relation` (
  `event_log_id` bigint NOT NULL AUTO_INCREMENT,
  `event_action_id` tinyint NOT NULL,
  PRIMARY KEY (`event_log_id`) USING BTREE,
  KEY `fk_event_action_relation_event_action_1` (`event_action_id`),
  CONSTRAINT `event_event_action_relation_ibfk_1` FOREIGN KEY (`event_log_id`) REFERENCES `event_log` (`event_log_id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `event_event_action_relation_ibfk_2` FOREIGN KEY (`event_action_id`) REFERENCES `event_action` (`event_action_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of event_event_action_relation
-- ----------------------------
BEGIN;
INSERT INTO `event_event_action_relation` (`event_log_id`, `event_action_id`) VALUES (8, 1);
INSERT INTO `event_event_action_relation` (`event_log_id`, `event_action_id`) VALUES (30, 2);
INSERT INTO `event_event_action_relation` (`event_log_id`, `event_action_id`) VALUES (32, 2);
INSERT INTO `event_event_action_relation` (`event_log_id`, `event_action_id`) VALUES (34, 2);
INSERT INTO `event_event_action_relation` (`event_log_id`, `event_action_id`) VALUES (36, 2);
INSERT INTO `event_event_action_relation` (`event_log_id`, `event_action_id`) VALUES (31, 5);
INSERT INTO `event_event_action_relation` (`event_log_id`, `event_action_id`) VALUES (33, 5);
INSERT INTO `event_event_action_relation` (`event_log_id`, `event_action_id`) VALUES (35, 5);
COMMIT;

-- ----------------------------
-- Table structure for event_event_status_relation
-- ----------------------------
DROP TABLE IF EXISTS `event_event_status_relation`;
CREATE TABLE `event_event_status_relation` (
  `event_id` bigint NOT NULL,
  `event_status_id` tinyint NOT NULL,
  PRIMARY KEY (`event_id`) USING BTREE,
  KEY `fk_event_event_status_relation_event_1` (`event_id`),
  KEY `event_status_id` (`event_status_id`),
  CONSTRAINT `event_event_status_relation_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `event` (`event_id`),
  CONSTRAINT `event_event_status_relation_ibfk_2` FOREIGN KEY (`event_status_id`) REFERENCES `event_status` (`event_status_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of event_event_status_relation
-- ----------------------------
BEGIN;
INSERT INTO `event_event_status_relation` (`event_id`, `event_status_id`) VALUES (1, 2);
COMMIT;

-- ----------------------------
-- Table structure for event_log
-- ----------------------------
DROP TABLE IF EXISTS `event_log`;
CREATE TABLE `event_log` (
  `event_log_id` bigint NOT NULL AUTO_INCREMENT,
  `event_id` bigint NOT NULL,
  `description` varchar(255) DEFAULT '',
  `member_id` char(10) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '',
  `gmt_create` datetime NOT NULL,
  PRIMARY KEY (`event_log_id`),
  KEY `fk_event_log_element_1` (`member_id`),
  KEY `fk_event_log_event_1` (`event_id`),
  CONSTRAINT `fk_event_log_event_1` FOREIGN KEY (`event_id`) REFERENCES `event` (`event_id`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of event_log
-- ----------------------------
BEGIN;
INSERT INTO `event_log` (`event_log_id`, `event_id`, `description`, `member_id`, `gmt_create`) VALUES (8, 1, '', '', '2022-05-13 16:50:55');
INSERT INTO `event_log` (`event_log_id`, `event_id`, `description`, `member_id`, `gmt_create`) VALUES (30, 1, '', '2333333333', '2022-05-16 14:49:55');
INSERT INTO `event_log` (`event_log_id`, `event_id`, `description`, `member_id`, `gmt_create`) VALUES (31, 1, '', '2333333333', '2022-05-16 14:49:55');
INSERT INTO `event_log` (`event_log_id`, `event_id`, `description`, `member_id`, `gmt_create`) VALUES (32, 1, '', '2333333333', '2022-05-17 10:39:55');
INSERT INTO `event_log` (`event_log_id`, `event_id`, `description`, `member_id`, `gmt_create`) VALUES (33, 1, '', '2333333333', '2022-05-17 10:58:55');
INSERT INTO `event_log` (`event_log_id`, `event_id`, `description`, `member_id`, `gmt_create`) VALUES (34, 1, '', '2333333333', '2022-05-17 10:58:55');
INSERT INTO `event_log` (`event_log_id`, `event_id`, `description`, `member_id`, `gmt_create`) VALUES (35, 1, '', '2333333333', '2022-05-17 11:08:55');
INSERT INTO `event_log` (`event_log_id`, `event_id`, `description`, `member_id`, `gmt_create`) VALUES (36, 1, '', '2333333333', '2022-05-17 11:08:55');
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
INSERT INTO `event_status` (`event_status_id`, `status`) VALUES (1, 'open');
INSERT INTO `event_status` (`event_status_id`, `status`) VALUES (2, 'accepted');
INSERT INTO `event_status` (`event_status_id`, `status`) VALUES (3, 'cancelled');
INSERT INTO `event_status` (`event_status_id`, `status`) VALUES (4, 'committed');
INSERT INTO `event_status` (`event_status_id`, `status`) VALUES (5, 'closed');
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
INSERT INTO `member` (`member_id`, `alias`, `password`, `name`, `section`, `profile`, `phone`, `qq`, `avatar`, `created_by`, `gmt_create`, `gmt_modified`) VALUES ('0000000000', '管理', '000000', '管理', '计算机000', '', '', '', '', '', '2022-04-30 17:28:42', '2022-04-30 17:28:44');
INSERT INTO `member` (`member_id`, `alias`, `password`, `name`, `section`, `profile`, `phone`, `qq`, `avatar`, `created_by`, `gmt_create`, `gmt_modified`) VALUES ('2333333333', '滑稽', '123456', '滑稽', '计算机233', 'relaxing', '12356839487', '123456', '', '0000000000', '2022-04-23 15:49:59', '2022-04-30 17:29:46');
COMMIT;

-- ----------------------------
-- Table structure for member_role_relation
-- ----------------------------
DROP TABLE IF EXISTS `member_role_relation`;
CREATE TABLE `member_role_relation` (
  `member_id` varchar(10) NOT NULL,
  `role_id` tinyint NOT NULL,
  PRIMARY KEY (`member_id`) USING BTREE,
  KEY `fk_member_role_role_1` (`role_id`),
  CONSTRAINT `fk_member_role_member_1` FOREIGN KEY (`member_id`) REFERENCES `member` (`member_id`),
  CONSTRAINT `fk_member_role_role_1` FOREIGN KEY (`role_id`) REFERENCES `role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of member_role_relation
-- ----------------------------
BEGIN;
INSERT INTO `member_role_relation` (`member_id`, `role_id`) VALUES ('2333333333', 2);
INSERT INTO `member_role_relation` (`member_id`, `role_id`) VALUES ('0000000000', 3);
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
INSERT INTO `role` (`role_id`, `role`) VALUES (0, 'member_inactive');
INSERT INTO `role` (`role_id`, `role`) VALUES (1, 'admin_inavtive');
INSERT INTO `role` (`role_id`, `role`) VALUES (2, 'member');
INSERT INTO `role` (`role_id`, `role`) VALUES (3, 'admin');
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
-- View structure for event_log_view
-- ----------------------------
DROP VIEW IF EXISTS `event_log_view`;
CREATE ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `event_log_view` AS select `event_log`.`event_log_id` AS `event_log_id`,`event_log`.`event_id` AS `event_id`,`event_log`.`description` AS `description`,`event_log`.`member_id` AS `member_id`,`event_log`.`gmt_create` AS `gmt_create`,`event_action`.`action` AS `action` from ((`event_log` left join `event_event_action_relation` on((`event_log`.`event_log_id` = `event_event_action_relation`.`event_log_id`))) left join `event_action` on((`event_event_action_relation`.`event_action_id` = `event_action`.`event_action_id`)));

-- ----------------------------
-- View structure for event_view
-- ----------------------------
DROP VIEW IF EXISTS `event_view`;
CREATE ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `event_view` AS select `event`.`event_id` AS `event_id`,`event`.`client_id` AS `client_id`,`event`.`model` AS `model`,`event`.`phone` AS `phone`,`event`.`qq` AS `qq`,`event`.`contact_preference` AS `contact_preference`,`event`.`problem` AS `problem`,`event`.`member_id` AS `member_id`,`event`.`closed_by` AS `closed_by`,`event`.`gmt_create` AS `gmt_create`,`event`.`gmt_modified` AS `gmt_modified`,`event_status`.`status` AS `status` from ((((`event` left join `event_event_status_relation` on((`event`.`event_id` = `event_event_status_relation`.`event_id`))) left join `event_status` on((`event_event_status_relation`.`event_status_id` = `event_status`.`event_status_id`))) left join `member` on((`event`.`member_id` = `member`.`member_id`))) left join `member` `admin` on((`event`.`closed_by` = `member`.`member_id`)));

-- ----------------------------
-- View structure for member_view
-- ----------------------------
DROP VIEW IF EXISTS `member_view`;
CREATE ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `member_view` AS select `member`.`member_id` AS `member_id`,`member`.`alias` AS `alias`,`member`.`password` AS `password`,`member`.`name` AS `name`,`member`.`section` AS `section`,`member`.`profile` AS `profile`,`member`.`phone` AS `phone`,`member`.`qq` AS `qq`,`member`.`avatar` AS `avatar`,`member`.`created_by` AS `created_by`,`member`.`gmt_create` AS `gmt_create`,`member`.`gmt_modified` AS `gmt_modified`,`role`.`role` AS `role` from ((`member` left join `member_role_relation` on((`member`.`member_id` = `member_role_relation`.`member_id`))) left join `role` on((`member_role_relation`.`role_id` = `role`.`role_id`)));

SET FOREIGN_KEY_CHECKS = 1;
