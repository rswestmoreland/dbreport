DROP TABLE IF EXISTS user_agent_info;
DROP TABLE IF EXISTS auth_attempts;
DROP TABLE IF EXISTS user_accounts;

CREATE TABLE user_accounts (
  id INT PRIMARY KEY,
  full_name VARCHAR(120) NOT NULL,
  username VARCHAR(80) NOT NULL UNIQUE,
  email VARCHAR(180) NOT NULL UNIQUE
);

CREATE TABLE auth_attempts (
  id INT PRIMARY KEY,
  user_id INT NOT NULL,
  login_time DATETIME NOT NULL,
  result ENUM('success', 'failure') NOT NULL,
  FOREIGN KEY (user_id) REFERENCES user_accounts(id)
);

CREATE TABLE user_agent_info (
  id INT PRIMARY KEY,
  auth_attempt_id INT NOT NULL,
  browser VARCHAR(80) NOT NULL,
  country VARCHAR(80) NOT NULL,
  FOREIGN KEY (auth_attempt_id) REFERENCES auth_attempts(id)
);
