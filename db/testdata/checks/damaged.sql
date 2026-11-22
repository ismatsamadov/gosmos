-- A NULL where the schema says NOT NULL, which only a write around the schema leaves and
-- quick_check finds.
PRAGMA writable_schema = ON;
UPDATE sqlite_schema SET sql = replace(sql, 'title         TEXT NOT NULL,', 'title         TEXT,') WHERE name = 'lesson';
PRAGMA writable_schema = RESET;
UPDATE lesson SET title = NULL WHERE uuid = 'hello';
PRAGMA writable_schema = ON;
UPDATE sqlite_schema SET sql = replace(sql, 'title         TEXT,', 'title         TEXT NOT NULL,') WHERE name = 'lesson';
PRAGMA writable_schema = RESET;
