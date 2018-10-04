CREATE DATABASE `theca`;

USE theca;

CREATE TABLE `opreturn` (
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
