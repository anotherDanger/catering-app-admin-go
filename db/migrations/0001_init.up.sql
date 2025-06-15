create database catering;
use catering;

create table products(
    id char(36),
    name varchar(100),
    description text,
    price int,
    stock int,
    created_at TIMESTAMP,
    modified_at TIMESTAMP
);

create table users(
    id int auto_increment primary key ,
    username varchar(100),
    password varchar(255)
);

create table admin(
    id char(36) primary key ,
    username varchar(100),
    password varchar(255)
);

create table orders(
    id char(36) primary key ,
    product_name varchar(100),
    username varchar(100),
    quantity int,
    total double,
    status varchar(20) default 'pending',
    created_at timestamp

);

insert into admin values ('e7b8a9d4-3f5a-4c82-b7e2-2c3f49b0e9c1', "admin", "123");