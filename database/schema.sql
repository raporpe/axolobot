CREATE SCHEMA IF NOT EXISTS axolobot;
USE axolobot;

DROP TABLE IF EXISTS mention;
CREATE TABLE `mention` (
  `mention_id` varchar(64) NOT NULL,
  PRIMARY KEY (`mention_id`));

(SELECT count(*) as count from axolobot.mention where mention_id = '1453481950052700161')
SELECT count(*) as count from axolobot.mention where mention_id = '1454882508848435204'