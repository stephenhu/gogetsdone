create table if not exists groups(
  id INTEGER NOT NULL PRIMARY KEY,
  owner_id INTEGER,
  name VARCHAR NOT NULL,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(owner_id) REFERENCES users(id)
);

create table if not exists user_groups(
  id INTEGER NOT NULL PRIMARY KEY,
  group_id INTEGER,
  user_id INTEGER,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(group_id) REFERENCES groups(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);
