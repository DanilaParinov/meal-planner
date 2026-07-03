-- Migration 004: add weight_g column to meals
ALTER TABLE meals ADD COLUMN IF NOT EXISTS weight_g INT;
