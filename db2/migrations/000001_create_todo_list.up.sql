CREATE TABLE list (
    id int PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE item (
    id int PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    done boolean not null default 0,
    list_id INT NOT NULL,
    FOREIGN KEY (list_id) REFERENCES list (id) ON DELETE CASCADE
   
);  