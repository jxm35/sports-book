-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE `competition` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL,
  `year` int(11) NOT NULL,
  `us_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1

CREATE TABLE `sports-book`.`team` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 NOT NULL,
  `us_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1

CREATE TABLE `sports-book`.`player` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 NOT NULL,
  `position` varchar(15) NOT NULL,
  `us_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1

CREATE TABLE `sports-book`.`match` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `home_team` int(11) NOT NULL,
  `away_team` int(11) NOT NULL,
  `competition` int(11) NOT NULL,
  `home_goals` int(11) NOT NULL,
  `away_goals` int(11) NOT NULL,
  `home_expected_goals` decimal(4,2) DEFAULT NULL,
  `away_expected_goals` decimal(4,2) DEFAULT NULL,
  `us_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `away_team` FOREIGN KEY (`away_team`) REFERENCES `team` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `competition` FOREIGN KEY (`competition`) REFERENCES `competition` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `home_team` FOREIGN KEY (`home_team`) REFERENCES `team` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1


CREATE TABLE `sports-book`.`appearance` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `player` int(11) NOT NULL,
  `match` int(11) NOT NULL,
  `team` int(11) NOT NULL,
  `goals` int(11) NOT NULL,
  `expected_goals` decimal(4,2) DEFAULT NULL,
  `expected_goals_chain` decimal(4,2) DEFAULT NULL,
  `expected_goals_buildup` decimal(4,2) DEFAULT NULL,
  `assists` int(11) NOT NULL,
  `expected_assists` decimal(4,2) DEFAULT NULL,
  `key_passes` int(11) NOT NULL,
  `yellow_cards` int(11) NOT NULL,
  `red_cards` int(11) NOT NULL,
  `minutes` int(11) NOT NULL,
  `us_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `match` FOREIGN KEY (`match`) REFERENCES `match` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `player` FOREIGN KEY (`player`) REFERENCES `player` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `team` FOREIGN KEY (`team`) REFERENCES `team` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1