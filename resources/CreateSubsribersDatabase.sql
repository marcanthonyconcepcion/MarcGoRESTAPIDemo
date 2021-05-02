/*
Author: Marc Concepcion
Copyright 2021, Marc Concepcion
Email: marcanthonyconcepcion@gmail.com
*/

drop database if exists `subscribers_database`;
create database if not exists `subscribers_database`;
use `subscribers_database`;
drop table if exists `subscribers`;
create table if not exists `subscribers` (
	`index`				int				primary key auto_increment,
    `email_address`		varchar(255)	not null unique,
    `last_name`			varchar(50),		
    `first_name`		varchar(100),
    `activation_flag`	tinyint			default 0 not null
);
