CREATE DATABASE `theca`; 

CREATE TABLE `theca`.`opreturn` (
  `txid` VARCHAR(64) NOT NULL,
  `prefix` VARCHAR(4) NULL,
  `hash` VARCHAR(255) NULL,
  `type` VARCHAR(5) NULL,
  `title` VARCHAR(255) NULL,
  PRIMARY KEY (`txid`));
