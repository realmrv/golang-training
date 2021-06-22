create table if not exists articles
(
    id      int auto_increment
        primary key,
    title   varchar(100) not null,
    brief   varchar(255) not null,
    content text         not null
);
