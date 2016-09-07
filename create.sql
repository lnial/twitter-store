CREATE TABLE tweet_media
(
    tweet_id character varying(500) NOT NULL,
    url character varying(500) NOT NULL,
    rt_user_id int,
    PRIMARY KEY (tweet_id, url)
);
