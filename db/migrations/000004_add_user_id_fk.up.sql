ALTER TABLE orders 
    ADD COLUMN user_id BIGINT NOT NULL,
    ADD CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);
