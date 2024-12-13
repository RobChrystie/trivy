#!/bin/bash

docker run --name mysql -v ./data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=password -d mysql:9.1.0
docker exec -it mysql bash

# mysql -u root -p

# create database recordings;

# use recordings

#DROP TABLE IF EXISTS album;
#CREATE TABLE album (
#  id         INT AUTO_INCREMENT NOT NULL,
#  title      VARCHAR(128) NOT NULL,
#  artist     VARCHAR(255) NOT NULL,
#  price      DECIMAL(5,2) NOT NULL,
#  PRIMARY KEY (`id`)
#);

#INSERT INTO album
#  (title, artist, price)
#VALUES
#  ('Blue Train', 'John Coltrane', 56.99),
#  ('Giant Steps', 'John Coltrane', 63.99),
#  ('Jeru', 'Gerry Mulligan', 17.99),
#  ('Sarah Vaughan', 'Sarah Vaughan', 34.98);
