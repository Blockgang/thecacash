CREATE DATABASE `theca`;
USE theca;

-- Activate EVENT SCHEDULER
SET GLOBAL event_scheduler = on;

-- Disable Save Update Mode
SET SQL_SAFE_UPDATES = 0;

-- memoLikeCounter Events
DELIMITER |
CREATE EVENT memoLikeCounter
	ON SCHEDULE EVERY 10 second
	DO
		BEGIN
			SET @start = NOW();
		  -- count Memo Likes
        UPDATE theca.prefix_0xe901 AS e901 SET e901.likes = e901.likes + COALESCE((SELECT COUNT(*) AS cnt FROM theca.prefix_0x6d04 as pref WHERE pref.txhash = e901.txid AND new = 1 GROUP BY pref.txhash),0);
			  UPDATE theca.prefix_0x6d04 AS 6d04 SET 6d04.new = 0, 6d04.theca = IF(EXISTS(SELECT 1 FROM theca.prefix_0xe901 as pref WHERE pref.txid = 6d04.txhash), 1, 0);
		  -- cleanup Memo Likes
		    DELETE FROM theca.prefix_0x6d04 WHERE new = 0 AND theca = 0; -- AND timestamp < @start;
		    SET @end = NOW();
		    SET @duration = (SELECT timestampdiff(SECOND,@start,@end));
		    INSERT INTO theca.event_scheduler(message,created_at,duration) VALUES('memoLikeCounter',@start,@duration);
		END |


-- Scheduler Events Table
DROP TABLE IF EXISTS event_scheduler;
CREATE TABLE `event_scheduler` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `message` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  `duration` int(11) NOT NULL,
  PRIMARY KEY (`id`)
);


-------------- THECA --------------
-- Theca Prefix
DROP TABLE IF EXISTS prefix_0xe901;
CREATE TABLE `prefix_0xe901` (
  `txid` varchar(64) NOT NULL,
  `hash` varchar(255) NOT NULL,
  `type` varchar(5) NOT NULL,
  `title` varchar(255) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  `likes` int(11) DEFAULT 0,
  PRIMARY KEY (`txid`)
);

-- Theca User
DROP TABLE IF EXISTS theca.users;
CREATE TABLE `theca`.`users` (
  `username` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `encrypted_pk` TEXT NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`username`));

-------------- MEMO --------------
-- Memo Set Name
DROP TABLE IF EXISTS prefix_0x6d01;
CREATE TABLE `prefix_0x6d01` (
  `txid` varchar(64) NOT NULL,
  `name` varchar(217) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  `new` bit(1) NOT NULL DEFAULT b'1',
  `theca` bit(1) NOT NULL DEFAULT b'0',
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`txid`)
);

-- Memo Reply
DROP TABLE IF EXISTS prefix_0x6d03;
CREATE TABLE `prefix_0x6d03` (
  `txid` varchar(64) NOT NULL,
  `txhash` varchar(64) NOT NULL,
  `message` varchar(184) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  `new` bit(1) NOT NULL DEFAULT b'1',
  `theca` bit(1) NOT NULL DEFAULT b'0',
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`txid`)
);

-- Memo Like
DROP TABLE IF EXISTS prefix_0x6d04;
CREATE TABLE `prefix_0x6d04` (
  `txid` varchar(64) NOT NULL,
  `txhash` varchar(64) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  `new` bit(1) NOT NULL DEFAULT b'1',
  `theca` bit(1) NOT NULL DEFAULT b'0',
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`txid`),
  KEY `TXHASH` (`txhash`),
  KEY `NEW` (`new`)
);

-- Memo Follow User
DROP TABLE IF EXISTS prefix_0x6d06;
CREATE TABLE `prefix_0x6d06` (
  `txid` varchar(64) NOT NULL,
  `address` varchar(35) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  `new` bit(1) NOT NULL DEFAULT b'1',
  `theca` bit(1) NOT NULL DEFAULT b'0',
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`txid`)
);

-- Memo Unfollow user
DROP TABLE IF EXISTS prefix_0x6d07;
CREATE TABLE `prefix_0x6d07` (
  `txid` varchar(64) NOT NULL,
  `address` varchar(35) NOT NULL,
  `blocktimestamp` int(11) DEFAULT 0,
  `blockheight` int(11) DEFAULT 0,
  `sender` varchar(60) NOT NULL,
  `new` bit(1) NOT NULL DEFAULT b'1',
  `theca` bit(1) NOT NULL DEFAULT b'0',
  `timestamp` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`txid`)
);
