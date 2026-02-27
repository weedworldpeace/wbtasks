CREATE TABLE IF NOT EXISTS items_history (
    id SERIAL PRIMARY KEY,
    item_id UUID NOT NULL,
    action VARCHAR(10) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    user_role VARCHAR(20) NOT NULL,
    old_data JSONB,
    new_data JSONB,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_items_history_item_id ON items_history(item_id);
CREATE INDEX idx_items_history_changed_at ON items_history(changed_at);

CREATE OR REPLACE FUNCTION items_history_trigger()
RETURNS TRIGGER AS $$
DECLARE
    _user_id UUID := current_setting('app.current_user_id', true)::UUID;
    _user_role VARCHAR := current_setting('app.current_user_role', true);
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO items_history(item_id, action, user_id, user_role, new_data)
        VALUES (NEW.id, 'INSERT', _user_id, _user_role, row_to_json(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO items_history(item_id, action, user_id, user_role, old_data, new_data)
        VALUES (OLD.id, 'UPDATE', _user_id, _user_role, row_to_json(OLD), row_to_json(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO items_history(item_id, action, user_id, user_role, old_data)
        VALUES (OLD.id, 'DELETE', _user_id, _user_role, row_to_json(OLD));
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER items_history_trigger
    AFTER INSERT OR UPDATE ON items  
    FOR EACH ROW EXECUTE FUNCTION items_history_trigger();

CREATE TRIGGER items_history_before_delete_trigger
    BEFORE DELETE ON items
    FOR EACH ROW EXECUTE FUNCTION items_history_trigger();

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_items_updated_at
    BEFORE UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();