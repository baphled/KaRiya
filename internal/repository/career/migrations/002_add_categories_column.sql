-- +goose Up
-- Add categories column to support event categorization
ALTER TABLE career_events ADD COLUMN categories TEXT;

-- +goose Down
-- NOTE: SQLite does not support DROP COLUMN natively
-- This migration cannot be fully rolled back without recreating the table
-- The column will remain but be unused if rolled back
-- For production rollback, manual table recreation would be required
