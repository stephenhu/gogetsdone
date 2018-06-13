create table if not exists ranks(
  id INTEGER NOT NULL PRIMARY KEY,
  rank VARCHAR,
  count INTEGER,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP 
);

create table if not exists users(
  id INTEGER NOT NULL PRIMARY KEY,
  email VARCHAR NOT NULL UNIQUE,
  name VARCHAR NOT NULL UNIQUE,
  mobile VARCHAR,
  password VARCHAR NOT NULL,
  salt VARCHAR,
  icon VARCHAR,
  rank_id INTEGER,
  registered BOOLEAN,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(rank_id) REFERENCES ranks(id)
);

create table if not exists follows(
  id INTEGER NOT NULL PRIMARY KEY,
  user_id INTEGER,
  follow_id INTEGER,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)  
);

create table if not exists states(
  id INTEGER NOT NULL PRIMARY KEY,
  state VARCHAR,
  description VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP 
);

create table if not exists priorities(
  id INTEGER NOT NULL PRIMARY KEY,
  priority VARCHAR,
  description VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP 
);

create table if not exists tasks(
  id INTEGER NOT NULL PRIMARY KEY,
  owner_id INTEGER,
  delegate_id INTEGER,
  state_id INTEGER,
  priority_id INTEGER,
  visibility INTEGER DEFAULT 0,
  task VARCHAR NOT NULL,
  estimate DATETIME DEFAULT CURRENT_TIMESTAMP,
  actual DATETIME DEFAULT CURRENT_TIMESTAMP,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(owner_id) REFERENCES users(id),
  FOREIGN KEY(delegate_id) REFERENCES delegates(id),
  FOREIGN KEY(state_id) REFERENCES states(id),
  FOREIGN KEY(priority_id) REFERENCES priorities(id)  
);

create table if not exists hashtags(
  id INTEGER NOT NULL PRIMARY KEY,
  hashtag VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP 
);

create table if not exists task_hashtags(
  id INTEGER NOT NULL PRIMARY KEY,
  task_id INTEGER,
  hashtag_id INTEGER,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(task_id) REFERENCES tasks(id),
  FOREIGN KEY(hashtag_id) REFERENCES hashtags(id)  
);

create table if not exists comments(
  id INTEGER NOT NULL PRIMARY KEY,
  task_id INTEGER,
  user_id INTEGER,
  comment VARCHAR NOT NULL,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(task_id) REFERENCES tasks(id),
  FOREIGN KEY(user_id) REFERENCES users(id)  
);

INSERT into ranks(rank, count) VALUES("white", 0);
INSERT into ranks(rank, count) VALUES("yellow", 10);
INSERT into ranks(rank, count) VALUES("orange", 20);
INSERT into ranks(rank, count) VALUES("green", 50);
INSERT into ranks(rank, count) VALUES("blue", 100);
INSERT into ranks(rank, count) VALUES("brown", 200);
INSERT into ranks(rank, count) VALUES("black", 500);
INSERT into ranks(rank, count) VALUES("red", 1000);

INSERT into states(state) VALUES("open");
INSERT into states(state) VALUES("deferred");
INSERT into states(state) VALUES("closed");

INSERT into priorities(priority) VALUES("low");
INSERT into priorities(priority) VALUES("medium");
INSERT into priorities(priority) VALUES("high");
INSERT into priorities(priority) VALUES("urgent");
