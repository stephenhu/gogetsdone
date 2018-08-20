create table if not exists stats(
  id INTEGER NOT NULL PRIMARY KEY,
  user_id INTEGER,
  total_tasks INTEGER,
  total_deferred INTEGER,
  total_delegated INTEGER,
  total_open INTEGER,
  total_time_to_complete INTEGER,
  daily_rate FLOAT,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP 
);
