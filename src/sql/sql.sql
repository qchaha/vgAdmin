CREATE DATABASE fordealOps CHARSET utf8mb4;
USE fordealOps;
CREATE TABLE user_info(
    id DOUBLE PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(32),
    password VARCHAR(32),
    role VARCHAR(16)
)CHARACTER SET utf8mb4;
INSERT INTO user_info VALUES(null, 'root', md5('root123'), 'super');
ALTER TABLE user_info ADD COLUMN status TINYINT(1) COMMENT '0:lock, 1:open';
ALTER TABLE user_info ADD COLUMN last_modify_time VARCHAR(19);
UPDATE user_info SET status = 1;
UPDATE user_info set last_modify_time = '2018-12-02 15:22:11';
ALTER TABLE user_info ADD COLUMN email VARCHAR(32);
UPDATE user_info SET email = 'admin@fordeal.com';
UPDATE user_info SET email = 'lansi@fordeal.com' WHERE username = 'lansi';
ALTER TABLE user_info MODIFY username VARCHAR(32) UNIQUE NOT NULL;


CREATE TABLE token_pool(
    id DOUBLE PRIMARY KEY AUTO_INCREMENT,
    token VARCHAR(32) UNIQUE NOT NULL,
    username VARCHAR(32) UNIQUE NOT NULL
)CHARACTER SET utf8mb4;
