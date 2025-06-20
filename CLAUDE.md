# Instructions

## PostgreSQL
- Always use GENERATED ALWAYS AS IDENTITY instead of SERIAL or BIGSERIAL types
- Always use TIMESTAMPTZ for datetime columns
- Always use JSONB for JSON data (if you need to)
- Always use TEXT instead of VARCHAR(N)

## Golang
- Always use camelCase for JSON fields
