-- A lesson of a type format 8 removed: only a write past lesson.type's CHECK leaves one, and
-- quick_check on a read-only connection lets it pass.
PRAGMA ignore_check_constraints = ON;
UPDATE lesson SET type = 'type_cli' WHERE uuid = 'quiz';
