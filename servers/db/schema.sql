create database if not exists db;

use db;

create table if not exists Users (
    ID int not null auto_increment primary key,
    Email varchar(320) not null,
    PassHash varchar(255) not null,
    UserName varchar(255) not null, 
    FirstName varchar(64) not null,
    LastName varchar(128) not null,
    PhotoURL varchar(128) not null,
    UNIQUE(Email),
    UNIQUE(UserName)
);

create table if not exists successful_logins (
    id int not null auto_increment primary key,
    user_id int not null,
    sign_in_time timestamp not null,
    login_ip varchar(32) not null
);