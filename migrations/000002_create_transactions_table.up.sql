CREATE TABLE IF NOT EXISTS transactions(
   id VARCHAR(36) PRIMARY KEY,
   sender_id VARCHAR(36) NOT NULL,
   receiver_id VARCHAR(36) NOT NULL,
   amount BIGINT NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (sender_id) REFERENCES users(id),
   FOREIGN KEY (receiver_id) REFERENCES users(id)
);
