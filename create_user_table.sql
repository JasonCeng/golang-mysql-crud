CREATE DATABASE IF NOT EXISTS `mpc` DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
USE `mpc`;
DROP TABLE IF EXISTS `User`;
CREATE TABLE IF NOT EXISTS `User` (
    `Id` int not null auto_increment,
    `Name` varchar(200) not null,
    `Amt` float null,
    primary key(Id)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;
INSERT INTO `User` (`Id`, `Name`, `Amt`) VALUES (1, '小明', 9.9);
INSERT INTO `User` (`Id`, `Name`, `Amt`) VALUES (2, 'liujun', 19.9);
INSERT INTO `User` (`Id`, `Name`, `Amt`) VALUES (3, '曹亮', 88.88);
select id,name,amt from User;
select id,name,amt from User where id = "1";