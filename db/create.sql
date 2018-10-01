CREATE DATABASE `theca` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */;

CREATE TABLE `theca`.`opreturn` (
  `txid` VARCHAR(64) NOT NULL,
  `prefix` VARCHAR(4) NULL,
  `hash` VARCHAR(255) NULL,
  `type` VARCHAR(5) NULL,
  `title` VARCHAR(255) NULL,
  PRIMARY KEY (`txid`));
