CREATE TABLE projects (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    slug VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP
);