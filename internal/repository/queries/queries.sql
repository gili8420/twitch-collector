-- name: CreateRecording :exec
insert into recordings (id, login, created_at, ready, title, category, language, viewers, stream_started_at)
values (@id, @login, @created_at, @ready, @title, @category, @language, @viewers, @stream_started_at);

-- name: GetStreamerRecordings :many
select *
from recordings
where login = @login;

-- name: UpdateRecording :one
update recordings
set login               = coalesce(sqlc.narg(login), login),
    created_at          = coalesce(sqlc.narg(created_at), created_at),
    ready               = coalesce(sqlc.narg(ready), ready),
    title               = coalesce(sqlc.narg(title), title),
    category            = coalesce(sqlc.narg(category), category),
    language            = coalesce(sqlc.narg(language), language),
    viewers             = coalesce(sqlc.narg(viewers), viewers),
    stream_started_at   = coalesce(sqlc.narg(stream_started_at), stream_started_at)
where id = @recording_id
returning id;