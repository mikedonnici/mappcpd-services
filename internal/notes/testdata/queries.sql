-- name: create-test-schema
CREATE SCHEMA `%s`;

-- name: drop-test-schema
DROP DATABASE `%s`;

-- name: create-wf_note-table
CREATE TABLE IF NOT EXISTS `%s`.`wf_note` (
  `id`                 INT       NOT NULL AUTO_INCREMENT,
  `wf_note_type_id`    INT       NOT NULL,
  `ad_user_id_created` INT       NULL,
  `ad_user_id_updated` INT       NULL,
  `active`             TINYINT   NOT NULL DEFAULT 1,
  `created_at`         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`         TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
  `effective_on`       DATE      NOT NULL,
  `note`               TEXT      NOT NULL,
  PRIMARY KEY (`id`)
);

-- name: create-wf_note_association-table
CREATE TABLE IF NOT EXISTS `%s`.`wf_note_association` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `wf_note_id` INT NOT NULL,
  `member_id` INT NULL,
  `association_entity_id` INT NULL,
  `active` TINYINT NOT NULL DEFAULT 1,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL,
  `association` VARCHAR(100) NULL,
  PRIMARY KEY (`id`)
);

-- name: create-member-table
CREATE TABLE IF NOT EXISTS `%s`.`member` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `acl_member_role_id` INT NOT NULL,
  `a_name_prefix_id` INT NOT NULL,
  `country_id` INT NOT NULL,
  `ad_macro_id` INT NULL,
  `active` TINYINT(1) NOT NULL DEFAULT 1,
  `consent_directory` TINYINT(1) NOT NULL DEFAULT 0,
  `consent_contact` TINYINT(1) NOT NULL DEFAULT 0,
  `login` TINYINT NOT NULL DEFAULT 1,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',
  `last_login_at` TIMESTAMP NULL,
  `date_of_birth` DATE NULL,
  `date_of_entry` DATE NULL DEFAULT NULL,
  `gender` ENUM('M', 'F') NULL,
  `first_name` VARCHAR(45) NOT NULL,
  `middle_names` VARCHAR(100) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL,
  `suffix` VARCHAR(100) NULL,
  `qualifications_other` TEXT NULL,
  `mobile_phone` VARCHAR(45) NULL,
  `primary_email` VARCHAR(100) NULL,
  `secondary_email` VARCHAR(100) NULL,
  `password` VARCHAR(100) NOT NULL,
  `token` VARCHAR(45) NULL,
  `journal_number` VARCHAR(45) NULL,
  `bpay_number` VARCHAR(45) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `primary_email_UNIQUE` (`primary_email` ASC),
  UNIQUE INDEX `secondary_email_UNIQUE` (`secondary_email` ASC));

-- name: insert-wf_note-data
INSERT INTO `%s`.`wf_note`
(`id`,`wf_note_type_id`,`ad_user_id_created`,`ad_user_id_updated`,`active`,`created_at`,`updated_at`,`effective_on`,`note`)
VALUES
  (1, 1, 1, 1, 1, NOW(), NOW(), "1970-11-03", "Mike birthday"),
  (2, 1, 1, 1, 1, NOW(), NOW(), "1975-07-25", "Christie birthday"),
  (3, 1, 1, 1, 1, NOW(), NOW(), "2010-03-09", "Maia birthday"),
  (4, 1, 1, 1, 1, NOW(), NOW(), "2012-11-02", "Leo birthday");

-- name: insert-wf_note_association-data
INSERT INTO `%s`.`wf_note_association`
(`id`, `wf_note_id`, `member_id`, `association_entity_id`, `active`, `created_at`, `updated_at`, `association`)
VALUES
  ('1', '1', '1', '', '1', NOW(), NOW(), ''),
  ('2', '2', '1', '', '1', NOW(), NOW(), ''),
  ('3', '3', '1', '', '1', NOW(), NOW(), ''),
  ('4', '4', '1', '', '1', NOW(), NOW(), '');

-- name: insert-member-data
INSERT INTO `%s`.`member`
(`id`, `acl_member_role_id`, `a_name_prefix_id`, `country_id`, `ad_macro_id`, `active`, `consent_directory`,
 `consent_contact`, `login`, `created_at`, `updated_at`, `last_login_at`, `date_of_birth`, `date_of_entry`,
 `gender`, `first_name`, `middle_names`, `last_name`, `suffix`, `qualifications_other`, `mobile_phone`,
 `primary_email`, `secondary_email`, `password`, `token`, `journal_number`, `bpay_number`)
VALUES
  ('1', '1', '1', '1', '0', '1', '1', '1', '1', NOW(), NOW(), NOW(), '1970-11-03', '2000-01-01', 'M', 'Michael',
    'Peter', 'Donnici', '', '', '', 'michael@mesa.net.au', '', '', '', '', '');


