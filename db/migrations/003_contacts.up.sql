create table if not exists contacts(
  id INTEGER NOT NULL PRIMARY KEY,
  user_id INTEGER,
  contact_id INTEGER,
  state INTEGER default 0,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(contact_id) REFERENCES users(id) 
);
