create database if not exists `rating_db`;
create database if not exists `clinic_db`;
create database if not exists `scheduling_db`;

create user 'root'@'localhost' identified by 'verysecretyes';
grant all on *.*  to 'root'@'%';