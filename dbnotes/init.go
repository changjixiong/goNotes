package main

import "goNotes/dbnotes/dbhelper"

/*

CREATE TABLE msg (
    id int(11) NOT NULL AUTO_INCREMENT,
    sender_id int(11) NOT NULL,
    receiver_id int(11) NOT NULL,
    content varchar(256) NOT NULL,
    status tinyint(4) NOT NULL,
    createtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
  )

  CREATE TABLE mail (
    id int(11) NOT NULL AUTO_INCREMENT,
    sender_id int(11) NOT NULL,
    receiver_id int(11) NOT NULL,
    title varchar(128) NOT NULL,
    content varchar(1024) NOT NULL,
    status tinyint(4) NOT NULL,
    createtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
  )

*/

func init() {

	dbhelper.GetDB("127.0.0.1", 3306, "dbnote", "root", "123456")

}
