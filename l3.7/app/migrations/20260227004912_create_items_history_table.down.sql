DROP TRIGGER IF EXISTS items_history_trigger ON items;
DROP TRIGGER IF EXISTS update_items_updated_at ON items;
DROP TRIGGER IF EXISTS items_history_before_delete_trigger ON items;
DROP TABLE IF EXISTS items_history;
DROP FUNCTION IF EXISTS items_history_trigger();
DROP FUNCTION IF EXISTS update_updated_at_column();