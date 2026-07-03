-- Migration 006: add is_drink column to meals
ALTER TABLE meals ADD COLUMN IF NOT EXISTS is_drink BOOLEAN NOT NULL DEFAULT false;
