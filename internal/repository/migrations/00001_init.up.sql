create table recordings (
    id                  uuid        primary key,
    login               text        not null,
    created_at          timestamp   not null,
    ready               boolean     not null,
    title               text        not null,
    category            text        not null,
    language            text        not null,
    viewers             integer     not null,
    stream_started_at   timestamp   not null
);

create index idx_recordings_login on recordings(login);
