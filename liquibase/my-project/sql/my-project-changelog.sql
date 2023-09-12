--liquibase formatted sql

--changeset lamboktulus1379:1 labels:my-project-label context:my-project-context
--comment: my-project comment
create table user (
    id int primary key auto_increment not null,
    name varchar(50) not null,
    user_name varchar(50),
    password varchar(50),
    created_by varchar(50),
    updated_by varchar(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
)
--rollback DROP TABLE user;

--changeset lamboktulus1379:2 labels:my-project-label context:my-project-context
--comment: my-project comment
ALTER TABLE `user` 
    MODIFY `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
--rollback ALTER TABLE `user` MODIFY `updated_at` TIMESTAMP;