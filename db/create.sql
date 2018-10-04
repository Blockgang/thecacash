CREATE DATABASE `theca`;

USE theca;

CREATE TABLE `prefix_0xe901` (
  `txid` varchar(64) NOT NULL,
  `prefix` varchar(4) NOT NULL,
  `hash` varchar(255) NOT NULL,
  `type` varchar(5) NOT NULL,
  `title` varchar(255) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  PRIMARY KEY (`txid`)
);


CREATE TABLE `theca`.`users` (
  `username` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `encrypted_pk` TEXT NOT NULL,
  PRIMARY KEY (`username`));
