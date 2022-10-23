-- @block
create table Account (
  account_id INT AUTO_INCREMENT

  username VARCHAR(255) NOT NULL UNIQUE
  password CHAR(128) NOT NULL

  created_at DATETIME NOT NULL DEFAULT GETDATE()

  PRIMARY KEY (account_id)
);
CREATE INDEX username_index ON Account(username);

-- @block
create table Subteddits (
  subteddit_id INT AUTO_INCREMENT

  subname VARCHAR(255) NOT NULL UNIQUE
  PRIMARY KEY (subteddit_id)
);
CREATE INDEX subname_idx ON Subteddits(subname);

-- @block
create table Twittos (
  twittos_id INT AUTO_INCREMENT

  username VARCHAR(255) NOT NULL UNIQUE
  PRIMARY KEY (twittos_id)
);
CREATE INDEX username_idx ON Twittos(username);

-- @block
create table Teship (
  follower_id INT
  subteddit_id INT
);
CREATE INDEX follower_id_idx ON Teship(follower_id);


-- @block
create table Twiship (
  follower_id INT
  twittos_id INT
);
CREATE INDEX follower_id_idx ON Twiship(follower_id);
