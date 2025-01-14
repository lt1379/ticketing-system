--liquibase formatted sql

--changeset lamboktulussimamora:1 labels:initialize context:dev
--comment: creating ticket table
CREATE TABLE `ticket`
(
    `id`         INT          NOT NULL AUTO_INCREMENT,
    `title`      VARCHAR(100) NOT NULL,
    `message`    TEXT         NOT NULL,
    `user_id`    INT          NOT NULL,
    `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX        `index_user_id` (`user_id` ASC) VISIBLE,
    INDEX        `index_created_at` (`created_at` ASC) VISIBLE
);

--rollback DROP TABLE `ticket`;

--changeset lamboktulussimamora:2 labels:modify context:dev
--comment: added column status
ALTER TABLE `ticket`
    ADD COLUMN `status` VARCHAR(3) NOT NULL DEFAULT 'opn' AFTER `user_id`;

--rollback ALTER TABLE `ticketing_system`.`ticket` DROP COLUMN `status`;
