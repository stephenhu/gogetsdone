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
  token VARCHAR,
  salt VARCHAR NOT NULL,
  icon VARCHAR,
  rank_id INTEGER DEFAULT 1,
  registered BOOLEAN DEFAULT 0,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(rank_id) REFERENCES ranks(id)
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
  origin_id INTEGER,
  visibility INTEGER DEFAULT 0,
  task VARCHAR(1024) NOT NULL,
  meta VARCHAR,
  deferred BOOLEAN DEFAULT 0,
  actual DATETIME,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(owner_id) REFERENCES users(id),
  FOREIGN KEY(delegate_id) REFERENCES users(id),
  FOREIGN KEY(state_id) REFERENCES states(id),
  FOREIGN KEY(priority_id) REFERENCES priorities(id),
  FOREIGN KEY(origin_id) REFERENCES tasks(id)
);

create table if not exists hashtags(
  id INTEGER NOT NULL PRIMARY KEY,
  hashtag VARCHAR NOT NULL UNIQUE,
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
  comment VARCHAR(1024) NOT NULL,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(task_id) REFERENCES tasks(id),
  FOREIGN KEY(user_id) REFERENCES users(id)  
);

INSERT into ranks(rank, count) VALUES("trainee", 0);
INSERT into ranks(rank, count) VALUES("neophyte", 10);
INSERT into ranks(rank, count) VALUES("novice", 20);
INSERT into ranks(rank, count) VALUES("trained", 50);
INSERT into ranks(rank, count) VALUES("qualified", 100);
INSERT into ranks(rank, count) VALUES("experienced", 200);
INSERT into ranks(rank, count) VALUES("professional", 500);
INSERT into ranks(rank, count) VALUES("maestro", 1000);

INSERT into states(state) VALUES("open");
INSERT into states(state) VALUES("deferred");
INSERT into states(state) VALUES("closed");

INSERT into priorities(priority) VALUES("low");
INSERT into priorities(priority) VALUES("medium");
INSERT into priorities(priority) VALUES("high");
INSERT into priorities(priority) VALUES("urgent");
