CREATE table user(
    id integer primary key autoincrement,
    sn varchar(8) not null unique,
    name varchar(6) not null unique,
    email varchar(50),
    hashed_password char(60),
    -- created datetime default current_timestamp
    created timestamp default (datetime('now','localtime'))
);

CREATE TABLE role(
       id integer primary key autoincrement,
       name varchar(20) NOT NULL unique,
       description varchar(50)
);
