package main

/*

CREATE TABLE msg (
    id int(11) NOT NULL AUTO_INCREMENT,
    sender_id int(11) NOT NULL COMMENT '发送者',
    receiver_id int(11) NOT NULL COMMENT '接收者',
    content varchar(256) NOT NULL COMMENT '内容',
    status tinyint(4) NOT NULL,
    createtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
  );

  CREATE TABLE mail (
    id int(11) NOT NULL AUTO_INCREMENT,
    sender_id int(11) NOT NULL,
    receiver_id int(11) NOT NULL,
    title varchar(128) NOT NULL,
    content varchar(1024) NOT NULL,
    status tinyint(4) NOT NULL,
    createtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
  );

  CREATE TABLE notice (
  id int(11) NOT NULL AUTO_INCREMENT,
  No int(11) NOT NULL,
  sender_id int(11) NOT NULL COMMENT '发送者',
  receiver_id int(11) NOT NULL COMMENT '接收者',
  content varchar(256) CHARACTER SET utf8mb4 NOT NULL COMMENT '内容',
  status tinyint(4) NOT NULL,
  createtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id,No)
);


cassandra -----------------
CREATE KEYSPACE space_for_back WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}  AND durable_writes = true;

CREATE TABLE num_log_2s (
    server_id int,
    create_time timestamp,
    id uuid,
    num int,
    PRIMARY KEY (server_id, create_time)
) WITH CLUSTERING ORDER BY (create_time DESC);

CREATE TABLE num_log_100ms (
    server_id int,
    create_time timestamp,
    id uuid,
    num int,
    PRIMARY KEY (server_id, create_time)
) WITH CLUSTERING ORDER BY (create_time DESC);

CREATE TABLE num_log_4s (
    server_id int,
    create_time timestamp,
    id uuid,
    num int,
    PRIMARY KEY (server_id, create_time)
) WITH CLUSTERING ORDER BY (create_time DESC);

CREATE TABLE num_log_10ms (
    server_id int,
    create_time timestamp,
    id uuid,
    num int,
    PRIMARY KEY (server_id, create_time)
) WITH CLUSTERING ORDER BY (create_time DESC);

*/
