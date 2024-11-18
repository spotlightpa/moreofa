-- migrate:up
CREATE TABLE comment (
  id INTEGER PRIMARY KEY, name TEXT NOT NULL DEFAULT '',
  contact TEXT NOT NULL DEFAULT '', message TEXT NOT NULL DEFAULT '',
  ip TEXT NOT NULL DEFAULT '', user_agent TEXT NOT NULL DEFAULT '',
  referrer TEXT NOT NULL DEFAULT '', host_page TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER comment_update_timestamp
AFTER
UPDATE
  ON comment FOR EACH ROW -- Only trigger if updated_at wasn't manually set
  WHEN NEW.updated_at = OLD.updated_at BEGIN
UPDATE
  comment
SET
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = NEW.id;
END;

-- migrate:down
DROP
  TRIGGER comment_update_timestamp;
DROP
  TABLE comment;
