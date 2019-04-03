/*
Navicat MySQL Data Transfer

Source Server         : 本地
Source Server Version : 50725
Source Host           : localhost:3306
Source Database       : linkchain

Target Server Type    : MYSQL
Target Server Version : 50725
File Encoding         : 65001

Date: 2019-04-03 21:27:31
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for blocks
-- ----------------------------
DROP TABLE IF EXISTS `blocks`;
CREATE TABLE `blocks` (
  `height` int(10) unsigned DEFAULT NULL,
  `hash` varchar(255) DEFAULT NULL,
  `version` int(10) unsigned DEFAULT NULL,
  `time` timestamp NULL DEFAULT NULL,
  `nonce` int(10) unsigned DEFAULT NULL,
  `difficulty` int(10) unsigned DEFAULT NULL,
  `prev` varchar(255) DEFAULT NULL,
  `tx_root` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `sign` varbinary(255) DEFAULT NULL,
  `hex` int(64) DEFAULT NULL,
  UNIQUE KEY `uix_blocks_height` (`height`),
  UNIQUE KEY `uix_blocks_hash` (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
