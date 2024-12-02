-- migrate:up

CREATE TABLE comment (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL DEFAULT '',
  contact TEXT NOT NULL DEFAULT '',
  subject TEXT NOT NULL DEFAULT '',
  cc TEXT NOT NULL DEFAULT '',
  message TEXT NOT NULL DEFAULT '',
  ip TEXT NOT NULL DEFAULT '',
  user_agent TEXT NOT NULL DEFAULT '',
  referrer TEXT NOT NULL DEFAULT '',
  host_page TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
  modified_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TRIGGER comment_update_timestamp
AFTER UPDATE ON comment
FOR EACH ROW
-- Only trigger if modified_at wasn't manually set
WHEN NEW.modified_at = OLD.modified_at
BEGIN
UPDATE comment
SET modified_at = current_timestamp
WHERE id = NEW.id ;
END ;

-- migrate:down
DROP TRIGGER comment_update_timestamp ;
DROP TABLE comment ;
