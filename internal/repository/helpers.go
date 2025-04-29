package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToPG(in uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: in, Valid: true}
}

func NilUUIDToPG(in *uuid.UUID) pgtype.UUID {
	if in == nil {
		return pgtype.UUID{}
	}

	return pgtype.UUID{Bytes: *in, Valid: true}
}

func PGToUUID(in pgtype.UUID) uuid.UUID {
	if !in.Valid {
		return uuid.Nil
	}

	return in.Bytes
}

func TimeToPG(in time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:             in.UTC(),
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}

func NilTimeToPG(in *time.Time) pgtype.Timestamp {
	if in == nil {
		return pgtype.Timestamp{}
	}

	return pgtype.Timestamp{
		Time:             in.UTC(),
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}

func PGToTime(in pgtype.Timestamp) time.Time {
	return in.Time.UTC()
}

func NilStringToPG(in *string) pgtype.Text {
	if in == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{String: *in, Valid: true}
}

func NilBoolToPG(in *bool) pgtype.Bool {
	if in == nil {
		return pgtype.Bool{}
	}

	return pgtype.Bool{Bool: *in, Valid: true}
}

func NilIntToPG(in *int) pgtype.Int4 {
	if in == nil {
		return pgtype.Int4{}
	}

	return pgtype.Int4{Int32: int32(*in), Valid: true}
}
