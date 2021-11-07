CREATE SCHEMA IF NOT EXISTS axolobot;
USE axolobot;
DROP TABLE IF EXISTS mention;
CREATE TABLE mention (
  mention_id varchar(64) NOT NULL,
  PRIMARY KEY (mention_id)
);