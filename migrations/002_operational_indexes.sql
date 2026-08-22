CREATE INDEX IF NOT EXISTS idx_audit_object ON audit_logs(object_type, object_id, created_at);
CREATE INDEX IF NOT EXISTS idx_events_suggestion ON suggestion_events(suggestion_id, created_at);
CREATE INDEX IF NOT EXISTS idx_assignments_lease ON assignments(status, lease_until);
