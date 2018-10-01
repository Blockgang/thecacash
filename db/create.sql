CREATE DATABASE `theca`;


  CREATE TABLE `opreturn` (
	    `txid` varchar(64) NOT NULL,
	    `prefix` varchar(4) DEFAULT NULL,
	    `hash` varchar(255) DEFAULT NULL,
	    `type` varchar(5) DEFAULT NULL,
	    `title` varchar(255) DEFAULT NULL,
	    `blocktimestamp` int(11) DEFAULT NULL,
	    `blockheight` int(11) DEFAULT NULL,
	    PRIMARY KEY (`txid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
