-- +goose Up
-- +goose StatementBegin

-- Define the global utility function
CREATE OR REPLACE FUNCTION public.set_updated_at() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS public.set_updated_at() CASCADE;
-- +goose StatementEnd
