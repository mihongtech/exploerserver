/*
Navicat MySQL Data Transfer

Source Server         : 本地
Source Server Version : 50725
Source Host           : localhost:3306
Source Database       : linkchain

Target Server Type    : MYSQL
Target Server Version : 50725
File Encoding         : 65001

Date: 2019-04-03 21:27:47
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for transactions
-- ----------------------------
DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions` (
  `height` int(10) unsigned DEFAULT NULL,
  `tx_id` varchar(255) DEFAULT NULL,
  `version` int(10) unsigned DEFAULT NULL,
  `type` int(10) unsigned DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
