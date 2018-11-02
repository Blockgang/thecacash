USE theca;
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
DELIMITER ;
