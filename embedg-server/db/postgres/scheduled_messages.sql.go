// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: scheduled_messages.sql

package postgres

import (
	"context"
	"database/sql"
	"time"
)

const deleteScheduledMessage = `-- name: DeleteScheduledMessage :exec
DELETE FROM scheduled_messages WHERE id = $1 AND guild_id = $2
`

type DeleteScheduledMessageParams struct {
	ID      string
	GuildID string
}

func (q *Queries) DeleteScheduledMessage(ctx context.Context, arg DeleteScheduledMessageParams) error {
	_, err := q.db.ExecContext(ctx, deleteScheduledMessage, arg.ID, arg.GuildID)
	return err
}

const getDueScheduledMessages = `-- name: GetDueScheduledMessages :many
SELECT id, creator_id, guild_id, channel_id, message_id, saved_message_id, cron_expression, trigger_at, trigger_once, enabled, created_at, updated_at FROM scheduled_messages WHERE trigger_at <= $1
`

func (q *Queries) GetDueScheduledMessages(ctx context.Context, triggerAt time.Time) ([]ScheduledMessage, error) {
	rows, err := q.db.QueryContext(ctx, getDueScheduledMessages, triggerAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ScheduledMessage
	for rows.Next() {
		var i ScheduledMessage
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.GuildID,
			&i.ChannelID,
			&i.MessageID,
			&i.SavedMessageID,
			&i.CronExpression,
			&i.TriggerAt,
			&i.TriggerOnce,
			&i.Enabled,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getScheduledMessage = `-- name: GetScheduledMessage :one
SELECT id, creator_id, guild_id, channel_id, message_id, saved_message_id, cron_expression, trigger_at, trigger_once, enabled, created_at, updated_at FROM scheduled_messages WHERE id = $1 AND guild_id = $2
`

type GetScheduledMessageParams struct {
	ID      string
	GuildID string
}

func (q *Queries) GetScheduledMessage(ctx context.Context, arg GetScheduledMessageParams) (ScheduledMessage, error) {
	row := q.db.QueryRowContext(ctx, getScheduledMessage, arg.ID, arg.GuildID)
	var i ScheduledMessage
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.GuildID,
		&i.ChannelID,
		&i.MessageID,
		&i.SavedMessageID,
		&i.CronExpression,
		&i.TriggerAt,
		&i.TriggerOnce,
		&i.Enabled,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getScheduledMessages = `-- name: GetScheduledMessages :many
SELECT id, creator_id, guild_id, channel_id, message_id, saved_message_id, cron_expression, trigger_at, trigger_once, enabled, created_at, updated_at FROM scheduled_messages WHERE guild_id = $1 ORDER BY updated_at DESC
`

func (q *Queries) GetScheduledMessages(ctx context.Context, guildID string) ([]ScheduledMessage, error) {
	rows, err := q.db.QueryContext(ctx, getScheduledMessages, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ScheduledMessage
	for rows.Next() {
		var i ScheduledMessage
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.GuildID,
			&i.ChannelID,
			&i.MessageID,
			&i.SavedMessageID,
			&i.CronExpression,
			&i.TriggerAt,
			&i.TriggerOnce,
			&i.Enabled,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertScheduledMessage = `-- name: InsertScheduledMessage :one
INSERT INTO scheduled_messages (id, creator_id, guild_id, channel_id, message_id, saved_message_id, cron_expression, trigger_at, trigger_once, enabled, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, creator_id, guild_id, channel_id, message_id, saved_message_id, cron_expression, trigger_at, trigger_once, enabled, created_at, updated_at
`

type InsertScheduledMessageParams struct {
	ID             string
	CreatorID      string
	GuildID        string
	ChannelID      string
	MessageID      sql.NullString
	SavedMessageID string
	CronExpression sql.NullString
	TriggerAt      time.Time
	TriggerOnce    bool
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (q *Queries) InsertScheduledMessage(ctx context.Context, arg InsertScheduledMessageParams) (ScheduledMessage, error) {
	row := q.db.QueryRowContext(ctx, insertScheduledMessage,
		arg.ID,
		arg.CreatorID,
		arg.GuildID,
		arg.ChannelID,
		arg.MessageID,
		arg.SavedMessageID,
		arg.CronExpression,
		arg.TriggerAt,
		arg.TriggerOnce,
		arg.Enabled,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i ScheduledMessage
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.GuildID,
		&i.ChannelID,
		&i.MessageID,
		&i.SavedMessageID,
		&i.CronExpression,
		&i.TriggerAt,
		&i.TriggerOnce,
		&i.Enabled,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateScheduledMessage = `-- name: UpdateScheduledMessage :one
UPDATE scheduled_messages SET channel_id = $3, message_id = $4, saved_message_id = $5, cron_expression = $6, trigger_at = $7, trigger_once = $8, enabled = $9, updated_at = $10 WHERE id = $1 AND guild_id = $2 RETURNING id, creator_id, guild_id, channel_id, message_id, saved_message_id, cron_expression, trigger_at, trigger_once, enabled, created_at, updated_at
`

type UpdateScheduledMessageParams struct {
	ID             string
	GuildID        string
	ChannelID      string
	MessageID      sql.NullString
	SavedMessageID string
	CronExpression sql.NullString
	TriggerAt      time.Time
	TriggerOnce    bool
	Enabled        bool
	UpdatedAt      time.Time
}

func (q *Queries) UpdateScheduledMessage(ctx context.Context, arg UpdateScheduledMessageParams) (ScheduledMessage, error) {
	row := q.db.QueryRowContext(ctx, updateScheduledMessage,
		arg.ID,
		arg.GuildID,
		arg.ChannelID,
		arg.MessageID,
		arg.SavedMessageID,
		arg.CronExpression,
		arg.TriggerAt,
		arg.TriggerOnce,
		arg.Enabled,
		arg.UpdatedAt,
	)
	var i ScheduledMessage
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.GuildID,
		&i.ChannelID,
		&i.MessageID,
		&i.SavedMessageID,
		&i.CronExpression,
		&i.TriggerAt,
		&i.TriggerOnce,
		&i.Enabled,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}