CREATE DATABASE IF NOT EXISTS `test`;
GRANT ALL PRIVILEGES ON *.* TO 'test'@'%';

USE `test`;

CREATE TABLE `user` (
  `name` VARCHAR(191) NOT NULL
) ENGINE=InnoDB;

INSERT INTO `user` (`name`) values ('test')
