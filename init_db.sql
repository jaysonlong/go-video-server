drop database if exists go_video_server;
create database go_video_server;
use go_video_server;

create table users (
  id int primary key auto_increment,
  user_name varchar(64) unique,
  pwd varchar(64)
);

create table video_info (
  id varchar(64) primary key,
  author_id int unsigned,
  name varchar(64),
  display_ctime varchar(64),
  create_time datetime default CURRENT_TIMESTAMP
);

create table comments (
  id int primary key auto_increment,
  author_id int unsigned,
  video_id varchar(64),
  content varchar(255),
  display_ctime varchar(64),
  create_time datetime default CURRENT_TIMESTAMP
);

create table sessions (
  session_id varchar(64) primary key,
  login_name varchar(64),
  TTL int
);

create table video_del_rec (
  video_id varchar(64) primary key
);