CREATE TABLE `components` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(256) DEFAULT '',
  `status` int(11) DEFAULT NULL,
  `version` varchar(128) DEFAULT NULL,
  `user` varchar(256) DEFAULT NULL,
  `create_date` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=362 DEFAULT CHARSET=utf8;