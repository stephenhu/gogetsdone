create table if not exists stats(
  id INTEGER NOT NULL PRIMARY KEY,
  user_id INTEGER,
  tc INTEGER,
  dc INTEGER,
  oc INTEGER,
  attc FLOAT,
  drio FLOAT,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP 
);
