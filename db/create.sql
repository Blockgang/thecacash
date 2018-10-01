CREATE DATABASE `theca`;

USE theca;

CREATE TABLE `opreturn` (
	    `txid` varchar(64) NOT NULL,
	    `prefix` varchar(4) DEFAULT NULL,
	    `hash` varchar(255) DEFAULT NULL,
	    `type` varchar(5) DEFAULT NULL,
	    `title` varchar(255) DEFAULT NULL,
	    `blocktimestamp` int(11) DEFAULT NULL,
	    `blockheight` int(11) DEFAULT NULL,
	    PRIMARY KEY (`txid`)
)
