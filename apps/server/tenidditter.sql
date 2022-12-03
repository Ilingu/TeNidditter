-- @block
create table Account (
  account_id INT AUTO_INCREMENT

  username VARCHAR(255) NOT NULL UNIQUE
  password CHAR(128) NOT NULL CHECK (LENGTH(password) >= 8)
  recovery_codes VARCHAR(108)

  created_at DATETIME NOT NULL DEFAULT GETDATE()

  PRIMARY KEY (account_id)
);
CREATE INDEX username_index ON Account(username);
ALTER TABLE Account
DROP recovery_codes;
ALTER TABLE Account
ADD recovery_codes VARCHAR(108);


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
ALTER TABLE Teship
ADD ship_id INT AUTO_INCREMENT PRIMARY KEY;

-- @block
create table Twiship (
  follower_id INT
  twittos_id INT
);
CREATE INDEX follower_id_idx ON Twiship(follower_id);
ALTER TABLE Twiship
ADD ship_id INT AUTO_INCREMENT PRIMARY KEY;

-- @block
create table NitterLists (
  list_id INT AUTO_INCREMENT

  account_id INT
  title TEXT NOT NULL

  PRIMARY KEY (list_id)
);
CREATE INDEX account_id_idx ON NitterLists(account_id);

-- @block
create table Neets (
  neet_id VARCHAR(255) NOT NULL UNIQUE CHECK (LENGTH(neet_id) = 19)
  neet_data JSON NOT NULL

  PRIMARY KEY (neet_id)
);
CREATE INDEX neet_id_idx ON Neets(neet_id);

-- @block
create table ListToNeet (
  list_id INT
  neet_id VARCHAR(255) NOT NULL UNIQUE CHECK (LENGTH(neet_id) = 19)
);
CREATE INDEX list_id_idx ON ListToNeet(list_id);
