CREATE TABLE foods (
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  meal VARCHAR(100) NOT NULL, 
  name VARCHAR(100) NOT NULL,
  protein INT NOT NULL,
  carbohydrate INT NOT NULL, 
  fat INT NOT NULL, 
  calories INT NOT NULL, 
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  userId INT NOT NULL, 
  FOREIGN KEY (userId) REFERENCES users(id)
);

CREATE INDEX idx_foods_created ON foods(created_at);