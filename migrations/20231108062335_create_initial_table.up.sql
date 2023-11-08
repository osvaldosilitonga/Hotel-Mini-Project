CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR NOT NULL UNIQUE,
  password VARCHAR NOT NULL,
  saldo INT NOT NULL
);

CREATE TABLE rooms(
  id SERIAL PRIMARY KEY,
  category VARCHAR(50) NOT NULL,
  room_number VARCHAR NOT NULL,
  status VARCHAR(50) NOT NULL DEFAULT 'ready',
  price INT NOT NULL
);

CREATE TABLE orders(
  id SERIAL PRIMARY KEY,
  user_id INT,
  room_id INT,
  adult INT NOT NULL DEFAULT 1,
  child INT DEFAULT 0,
  check_in DATE NOT NULL,
  check_out DATE NOT NULL,
  status VARCHAR NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (room_id) REFERENCES rooms(id)
);

CREATE TABLE payments(
  id SERIAL PRIMARY KEY,
  order_id INT,
  method VARCHAR(20),
  amount INT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'unpaid',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_id) REFERENCES orders(id)
);