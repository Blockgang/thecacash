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

CREATE TABLE `prefix_0x6d01` (
  `txid` varchar(64) NOT NULL,
  `name` varchar(217) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  PRIMARY KEY (`txid`)
);

CREATE TABLE `prefix_0x6d03` (
  `txid` varchar(64) NOT NULL,
  `txhash` varchar(30) NOT NULL,
  `message` varchar(184) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  PRIMARY KEY (`txid`)
);

CREATE TABLE `prefix_0x6d04` (
  `txid` varchar(64) NOT NULL,
  `txhash` varchar(30) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  PRIMARY KEY (`txid`)
);

CREATE TABLE `prefix_0x6d06` (
  `txid` varchar(64) NOT NULL,
  `address` varchar(35) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  PRIMARY KEY (`txid`)
);

CREATE TABLE `prefix_0x6d07` (
  `txid` varchar(64) NOT NULL,
  `address` varchar(35) NOT NULL,
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
