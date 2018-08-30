create table if not exists contact_states(
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR NOT NULL UNIQUE,
  description VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP
);


create table if not exists contacts(
  id INTEGER NOT NULL PRIMARY KEY,
  user_id INTEGER,
  contact_id INTEGER,
  contact_state_id INTEGER,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(contact_id) REFERENCES users(id),
  FOREIGN KEY(contact_state_id) REFERENCES contact_states(id)
);

INSERT into contact_states(name) VALUES("pending");
INSERT into contact_states(name) VALUES("requested");
INSERT into contact_states(name) VALUES("accepted");
INSERT into contact_states(name) VALUES("declined");
