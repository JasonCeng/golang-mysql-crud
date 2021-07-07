CREATE DATABASE IF NOT EXISTS `mpc` DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
USE `mpc`;
DROP TABLE IF EXISTS `User`;
CREATE TABLE IF NOT EXISTS `User` (`Id` varchar(20), `Name` varchar(20));
INSERT INTO `User` (`Id`, `Name`) VALUES ("1", "xiaoming");
select id,name from User where id = "1";