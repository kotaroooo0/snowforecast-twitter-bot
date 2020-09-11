drop database if exists snowforecast_twitter_bot;
create database snowforecast_twitter_bot;
use snowforecast_twitter_bot;

drop table if exists snow_resorts;
create table snow_resorts (
    id int not null auto_increment primary key,
    name varchar(128),
    search_key varchar(128)
);
