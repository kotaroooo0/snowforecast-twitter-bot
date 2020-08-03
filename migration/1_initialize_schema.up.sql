drop table if exists snow_resorts;

create table snow_resorts
(
    id auto_increment,
    name varchar(40),
    search_key varchar(40),
    primary key (id),
);

load data infile 'data.csv' into table snow_resorts fields terminated by ',';
