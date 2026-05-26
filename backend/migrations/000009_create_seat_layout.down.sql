ALTER TABLE ticket_types DROP COLUMN IF EXISTS seat_section;
ALTER TABLE ticket_types DROP COLUMN IF EXISTS seat_rows;
DROP TABLE IF EXISTS seat_layouts;
