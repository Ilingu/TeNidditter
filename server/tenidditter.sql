-- @block
create table Account (
  account_id INT AUTO_INCREMENT

  username VARCHAR(255) NOT NULL UNIQUE
  password CHAR(128) NOT NULL

  created_at DATETIME NOT NULL DEFAULT GETDATE()

  PRIMARY KEY (account_id)
);

-- @block
CREATE INDEX username_index ON Account(username);