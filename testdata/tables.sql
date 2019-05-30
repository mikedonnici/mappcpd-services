-- name: create-table-ad_user
CREATE TABLE IF NOT EXISTS `%s`.`ad_user` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `acl_admin_role_id` INT NOT NULL COMMENT 'The role (permissions group) into which this admin user is assigned.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `locked` TINYINT NOT NULL DEFAULT 0 COMMENT 'Account locked after X failed login attempts.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `username` VARCHAR(45) NOT NULL COMMENT 'admin user\'s username',
  `password` VARCHAR(45) NOT NULL COMMENT 'Admin user\'s password',
  `name` VARCHAR(100) NOT NULL COMMENT 'Admin user full name',
  `short_name` VARCHAR(16) NOT NULL COMMENT 'Short name is for display in lists, e.g. can use initials, first or nick name.',
  `email` VARCHAR(100) NULL COMMENT 'Contact email for admin - not used for anything at present but may be used for alerts etc.',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name` (`username` ASC),
  UNIQUE INDEX `short_name_UNIQUE` (`short_name` ASC),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC))
  ENGINE = InnoDB
  COMMENT = 'Admin user records';


-- name: create-table-ad_user_permission
CREATE TABLE IF NOT EXISTS `%s`.`ad_user_permission` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ad_user_id` INT NOT NULL COMMENT 'The id of the admin user that has been granted the permission.',
  `ad_permission_id` INT NOT NULL COMMENT 'The id of the permission that has been granted',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Association between admin user and a permission, that is, stores the granting of permissions to specific admin users.';


-- name: create-table-ce_activity
CREATE TABLE IF NOT EXISTS `%s`.`ce_activity` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ce_activity_unit_id` INT NOT NULL COMMENT 'The unit of measurement for the activity type.',
  `ce_activity_category_id` INT NOT NULL,
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete.',
  `system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'An activity type that is managed by the application and not modifiable by the end users.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `code` VARCHAR(10) NOT NULL COMMENT 'An arbitrary short code assigned to the activity type.',
  `name` VARCHAR(100) NOT NULL COMMENT 'The descriptive name of the activity type / category.',
  `description` TEXT NOT NULL COMMENT 'Description of the type of CPD activity and / or instruction to the user such as what documentation needs to be provided as proof.',
  `points_per_unit` DECIMAL(5,2) NOT NULL COMMENT 'The amount of points allocated per unit of this type of activity.',
  `annual_points_cap` TINYINT UNSIGNED NOT NULL COMMENT 'Specifies an annual maximum cap for the activity. Need to standardise to a year so we can apply appropriate capping to evaluation periods of varying length.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines the various CPD activities or categories of activity, that members can undertake in order to satisfy their CPD requirements.';


-- name: create-table-log_data_action
CREATE TABLE IF NOT EXISTS `%s`.`log_data_action` (
  `id` INT(10) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `log_data_table_id` INT(10) NOT NULL,
  `record_id` INT NOT NULL COMMENT 'The id of the record that was changed, in the table identified by log_data_table_id.',
  `user_id` INT(10) NOT NULL COMMENT 'The id of the user (either admin or member) that performed the action.',
  `member_id` INT NULL COMMENT 'If relevant, the id of the member to which the data change relates.',
  `active` TINYINT NOT NULL DEFAULT 1,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created at - ie time and date that the action was performed.',
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `user_type` VARCHAR(20) NOT NULL COMMENT 'Defines if the user was \'admin\' or \'member\'.',
  `action` ENUM('insert','update','delete') NOT NULL COMMENT 'Describes the action that was performed on the data',
  `message` TEXT NOT NULL COMMENT 'Further details about the event.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Logs various system events and data changes over time.';


-- name: create-table-member
CREATE TABLE IF NOT EXISTS `%s`.`member` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `acl_member_role_id` INT NOT NULL COMMENT 'The ACL role / group to which this member has been assigned.',
  `a_name_prefix_id` INT NOT NULL COMMENT 'Name prefix - eg Dr, Mr etc.',
  `country_id` INT NOT NULL COMMENT 'This is the membership country of the member and has nothing to do with their residence or contact details. It is used for billing, tax codes etc. ',
  `ad_macro_id` INT NULL COMMENT 'If this record was inserted or updated by a macro we store the macro id here.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `consent_directory` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'A flag to specify if the user has given con set to appear in the directory. Directory listing is controlled by both title and status, but this flag is the final one required for a directory listing. So a member would need a 1 in all three flags - we default this to 1 so it is not an impediment by default.',
  `consent_contact` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Member gives constant for contact information to be given to third parties where appropriate.',
  `login` TINYINT NOT NULL DEFAULT 1 COMMENT 'A flag to allow user to login. User also requires a membership title and status that allow login.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `last_login_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'The last login for the member.',
  `date_of_birth` DATE NULL COMMENT 'Date of birth of the member',
  `date_of_entry` DATE NULL DEFAULT NULL COMMENT 'Date of entry of the member into the organisations - this field is a bit of a legacy field for data migration - strictly speaking the membership title history should cover this.',
  `gender` ENUM('M', 'F') NULL,
  `first_name` VARCHAR(45) NOT NULL COMMENT 'Member\'s first name\n',
  `middle_names` VARCHAR(100) NOT NULL COMMENT 'One or more middle names.',
  `last_name` VARCHAR(45) NOT NULL COMMENT 'Member\'s surname / family name\n',
  `suffix` VARCHAR(100) NULL,
  `qualifications_other` TEXT NULL COMMENT 'This is used to store qualifications (post nominals) other than those we have normalised into the mp_qualifications table. This field was introduced primarily to accommodate the import of existing data that was difficult to normalise.',
  `mobile_phone` VARCHAR(45) NULL COMMENT 'Mobile phone number.',
  `primary_email` VARCHAR(100) NULL COMMENT 'Primary email address, also used for authentication to the member system.',
  `secondary_email` VARCHAR(100) NULL COMMENT 'Secondary email is used in case the primary email become inactive or gets forgotten. The user can also login with this email.',
  `password` VARCHAR(100) NOT NULL COMMENT 'Password is stored as an MD5 hash.',
  `token` VARCHAR(45) NULL COMMENT 'A temporary authentication token used to log the user in via a link so they can reset their password. This token should be cleared immediately as part of the login process.',
  `journal_number` VARCHAR(45) NULL COMMENT 'Journal number is given to the member as a reference for their subscription to a primary journal publication. (This should move to an external table later)',
  `bpay_number` VARCHAR(45) NULL COMMENT 'BPay number is generated by admin and allocated for Australian members only - for direct deposit of funds. (This should move to an external table later)',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `primary_email_UNIQUE` (`primary_email` ASC),
  UNIQUE INDEX `secondary_email_UNIQUE` (`secondary_email` ASC))
  ENGINE = InnoDB
  COMMENT = 'Member is the central entity of the system and stores basic information about the person (member) including details to login.';


-- name: create-table-ce_m_activity
CREATE TABLE IF NOT EXISTS `%s`.`ce_m_activity` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'The member who completed the activity.',
  `ce_activity_id` INT NOT NULL COMMENT 'The activity (type) that was undertaken.',
  `ce_activity_type_id` INT NULL COMMENT 'The activity type.',
  `ce_event_id` INT NULL COMMENT 'Optional link to an Event, also used to record Member attendance at Events. This is here instead of a junction table between event and member.',
  `ad_macro_id` INT NULL COMMENT 'If this record was inserted or updated by a macro we store the macro id here.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete - hard delete is allowed also under some circumstances.',
  `evidence` TINYINT NULL COMMENT 'A flag to indicate if the member has evidence available to support the claim of this activity - e.g. Document. NULL is unknown, 0 is NO and 1 is YES.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `activity_on` DATE NULL DEFAULT NULL COMMENT 'The date that the activity was performed. This could be a start or end date for a multi-day activity. This date is used to capture activity within  defined evaluation period.',
  `quantity` DECIMAL(5,2) NOT NULL COMMENT 'The number of units of the activity that were completed. e.g. 4 x hours.',
  `points_per_unit` DECIMAL(5,2) NOT NULL COMMENT 'Points for each unit is copied from the ce_activity definition table at the time the activity is recorded. This is in case the value for the activity is changed at some stage in the future.\n\nWe copy the current value from the ce_activity table each time a new activity is entered, or each time the evaluation period report is generated for an OPEN EP.\n\nFor a closed EP we will NOT reset this value so the historical values are maintained. \n\nThis means that the value for an activity MAY change over time for the user. This is part of the rules and the final value will be the current value at the time the EP is closed.',
  `annual_points_cap` SMALLINT NOT NULL DEFAULT 0 COMMENT 'Standardised (per year) points cap for the activity. As for points_per_unit we copy the current value from the ce_activity table each time a new activity is entered, or each time the evaluation period report is generated for an OPEN EP.\n\nFor a closed EP we will NOT reset this value so the historical values are maintained. \n\nIn both cases we can use ANY value for the same activity type (they should all be the same anyway) for the applications of caps. Yes, this is very redundant data BUT we decide was better to do it this way as it saved us managing a separate table for the same purpose.',
  `description` TEXT NULL COMMENT 'Optional descriptive text about the activity.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'A record of a particular CPD activity undertaken by a member.';


-- name: create-table-ms_m_title
CREATE TABLE IF NOT EXISTS `%s`.`ms_m_title` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'Link the member',
  `ms_title_id` INT NOT NULL COMMENT 'Link to the membership (title)',
  `ad_macro_id` INT NULL COMMENT 'If this record was inserted or updated by a macro we store the macro id here.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `current` TINYINT NOT NULL DEFAULT 0 COMMENT 'Defines the current (i.e. latest) title and is used to make table joins easier when creating lists of members.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `granted_on` DATE NULL DEFAULT NULL COMMENT 'The date the member received this level of membership. Would generally coincide with a board meeting date.',
  `comment` TEXT NULL COMMENT 'Optional comment',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Membership title record history for a member.';


-- name: create-table-ad_permission
CREATE TABLE IF NOT EXISTS `%s`.`ad_permission` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'record last updated',
  `name` VARCHAR(80) NOT NULL COMMENT 'Name of the permission.',
  `description` TEXT NOT NULL COMMENT 'Description of the permission.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  AUTO_INCREMENT = 6
  COMMENT = 'Stores a list of the various permissions that may be granted to admin users.';


-- name: create-table-organisation_type
CREATE TABLE IF NOT EXISTS `%s`.`organisation_type` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name` VARCHAR(45) NOT NULL COMMENT 'Name of organisation type',
  `description` VARCHAR(100) NOT NULL COMMENT 'Descriptive text',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines types of organisations such as Universities, Colleges, etc etc.';


-- name: create-table-mp_m_qualification
CREATE TABLE IF NOT EXISTS `%s`.`mp_m_qualification` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique Identifier',
  `member_id` INT NOT NULL COMMENT 'The member',
  `mp_qualification_id` INT NOT NULL COMMENT 'The qualification',
  `organisation_id` INT NULL DEFAULT NULL COMMENT 'Optional link to an Organisation. ',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `year` INT(4) UNSIGNED NULL DEFAULT NULL COMMENT 'The year the qualification was obtained… NULL if unknown.',
  `qualification_suffix` VARCHAR(45) NULL DEFAULT NULL COMMENT 'Used to further describe the qualification for thinks like (Hons), (First Class Hons) and so on. Will appear immediately to the right of the qualification short name.',
  `comment` TEXT NULL COMMENT 'Optional comment',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores a qualification that has been obtained by a member.';


-- name: create-table-mp_qualification
CREATE TABLE IF NOT EXISTS `%s`.`mp_qualification` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `organisation_id` INT NOT NULL DEFAULT 0 COMMENT 'Optional link to the Organisation (qualification provider). If 0 then this qualification is general and may come from more than one organisation. Eg. B Sc.',
  `active` TINYINT NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation date and time',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Last updated date and time',
  `display_priority` SMALLINT NOT NULL DEFAULT 0 COMMENT 'An optional value that will allow records to be ordered when displayed. Higher priority is displayed first.',
  `short_name` VARCHAR(45) NOT NULL COMMENT 'Accepted short name or abbreviation. Eg MBBS.',
  `name` VARCHAR(100) NOT NULL COMMENT 'The full name for the qualification',
  `description` TEXT NOT NULL COMMENT 'Descriptive text',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Store a list of recognised qualifications that can be assigned to members. Each qualification may optionally be linked to a specific organisation.';


-- name: create-table-ce_m_evaluation
CREATE TABLE IF NOT EXISTS `%s`.`ce_m_evaluation` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'The member',
  `ce_evaluation_id` INT NOT NULL COMMENT 'The evaluation period \'type\'',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `closed` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to indicate that this evaluation period has ended and is closed to any further modification. This will also prevent the recording on new CPD Activity with dates falling within this period.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `cpd_points_required` INT(5) NOT NULL COMMENT 'The number of CPD points that are required to satisfy this evaluation period.',
  `start_on` DATE NULL DEFAULT NULL COMMENT 'The start date for this evaluation period.',
  `end_on` DATE NULL DEFAULT NULL COMMENT 'The end date for this evaluation period.',
  `comment` TEXT NOT NULL COMMENT 'A comment about this instanc of the evluation period type. Eg if it gets modified.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines an evaluation period for an individual member. An evaluation period is an arbitrary stretch of time over which the member cpd activity will be assessed.';


-- name: create-table-organisation
CREATE TABLE IF NOT EXISTS `%s`.`organisation` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'If present, denotes that the record is a group within an existing organisation.',
  `parent_organisation_id` INT NULL DEFAULT NULL,
  `organisation_type_id` INT NOT NULL COMMENT 'Further classifies the organisation type. Eg University, Government Authority etc.',
  `country_id` INT NOT NULL COMMENT 'The country of origin / location for the organisation.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Last updated',
  `short_name` VARCHAR(45) NOT NULL COMMENT 'Short name or accepted abbreviation for the Organisation. Eg Uni Syd,',
  `name` VARCHAR(255) NOT NULL COMMENT 'Full name of Organisation',
  `address1` VARCHAR(100) NOT NULL,
  `address2` VARCHAR(100) NOT NULL,
  `address3` VARCHAR(100) NOT NULL,
  `locality` VARCHAR(100) NOT NULL,
  `state` VARCHAR(45) NOT NULL,
  `postcode` VARCHAR(20) NOT NULL,
  `phone` VARCHAR(45) NOT NULL,
  `fax` VARCHAR(45) NOT NULL,
  `email` VARCHAR(100) NOT NULL,
  `web` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC))
  ENGINE = InnoDB
  COMMENT = 'Organisation table stores a list of any type or organisation that may be relevant to the system. This might include Universities, committees, industry or government bodies. \nIt is simply for the purposes of associating various profile attributes with a relevant body.';


-- name: create-table-ce_evaluation
CREATE TABLE IF NOT EXISTS `%s`.`ce_evaluation` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `audit_compulsory` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to indicate that this type of evaluation should always trigger an audit once it has expired.',
  `start_date` TINYINT NULL COMMENT 'Default start date (numeric day of month) of the evaluation period.',
  `start_month` TINYINT NULL COMMENT 'Default start month (numeric) of the evaluation period.',
  `end_date` TINYINT NULL COMMENT 'Default end date (numeric day of month) of the evaluation period.',
  `end_month` TINYINT NULL COMMENT 'Default end month (numeric) of the evaluation period.',
  `duration_months` TINYINT NULL COMMENT 'The duration of the evaluation period in months e.g. 12, 24, 36 etc.',
  `points_required` INT(5) NOT NULL COMMENT 'The number of cpd points required to satisfy the requirements of this evaluation period.',
  `name` VARCHAR(45) BINARY NOT NULL COMMENT 'A name for the evaluation period. For example: \'Standard Evaluation\'.',
  `description` TEXT NOT NULL COMMENT 'A description of the evaluation period.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines a standard evaluation period type which may have a pre-defined length, points requirements etc.\n';


-- name: create-table-wf_issue_type
CREATE TABLE IF NOT EXISTS `%s`.`wf_issue_type` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `wf_issue_category_id` INT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'The category of the issue - simply for grouping and display purposes..',
  `acl_member_role_id` INT NULL COMMENT 'If present this value indicates that an active issue of this type should impose restrictions on the member as though they belonged to this member_role group. If multiple issues are active then the system will take the intersection of all the role restrictions and impose them on the member after they login. This does not change the members actual assigned role and will revert once the issues are resolved.',
  `active` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `system_raise` TINYINT NOT NULL DEFAULT 0 COMMENT 'Is system only issue - that is only the system is able to raise an instance of this typ eof issue. If this is set to 0 then admin users are able to raise an issue of this type.',
  `system_resolve` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Is system only issue - that is only the system is able to RESOLVE an instance of this type of issue. If set to 1 we will hide the input to resolve the issue from the admin users. this will avoid admin users manually closing issues which will then be automatically raised again by the housekeeping scripts. If this is set to 0 then admin users are able to rresolve an issue of this type.',
  `member_visible` TINYINT NOT NULL DEFAULT 0 COMMENT 'Flag to indicate if this issue type is visible to the member, defaults to NO (0).',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated.',
  `name` VARCHAR(100) NOT NULL COMMENT 'Short name for issue',
  `description` TEXT NOT NULL COMMENT 'Longer descriptive text about issue type.',
  `required_action` TEXT NOT NULL COMMENT 'A friendly message to show members that explains the required action they should take in relation to this issue.',
  `dependencies` VARCHAR(100) NULL COMMENT 'Defines a comma separated list of issue ids that are impediments to this issue. That is if those issues exist and are open this issue cannot be resolved.\n',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC))
  ENGINE = InnoDB
  AUTO_INCREMENT = 10000
  COMMENT = 'Defines the various types of administrative issues, both system and admin types.';


-- name: create-table-wf_issue
CREATE TABLE IF NOT EXISTS `%s`.`wf_issue` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `wf_issue_type_id` INT NOT NULL COMMENT 'The issue type.',
  `ad_user_id_created` INT NULL COMMENT 'Admin user who created the issue, NULL will be a system issue.',
  `ad_user_id_updated` INT NULL COMMENT 'Admin user who last updated the issue, NULL for system.',
  `ad_user_id_assigned` INT NULL COMMENT 'Admin user to whom the issue has been assigned, NULL for nobody or a system issue.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `resolved` TINYINT NOT NULL DEFAULT 0 COMMENT 'Resolved flag.',
  `member_visible` TINYINT NOT NULL DEFAULT 0 COMMENT 'Flag to indicate if this issue type is visible to the member. Inherits value from the wf_issue_type parent record but may be overridden in this instance. The rules for changing this value are implemented in the application. ',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `live_on` DATE NULL DEFAULT NULL COMMENT 'The date the issue gets raised and becomes live.',
  `description` TEXT NOT NULL COMMENT 'The description of the issue, pre-filled from wf_issue_type.description and can then be edited for each instance of the issue.',
  `required_action` TEXT NOT NULL COMMENT 'The required action by either admin or member, pre-filled from wf_issue_type.required_action and can then be edited for each instance of the issue.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'An instance of an issue that MAY relate to a member. If member_id is 0 the issue is defines as global in context.';


-- name: create-table-ce_audit
CREATE TABLE IF NOT EXISTS `%s`.`ce_audit` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ce_m_evaluation_id` INT NOT NULL COMMENT 'The evaluation period that is subject to audit audited.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `completed_on` DATE NULL DEFAULT NULL COMMENT 'Date the audit was finalised',
  `result` TINYINT NOT NULL COMMENT 'The pass / fail status of the audit. All audits will start as 0 (pending) the be either passed or failed. May change to enum (\'PENDING\', \'PASS\', \'FAILED\')\n',
  `audited_by` VARCHAR(45) NOT NULL COMMENT 'The name of the person who did the audit. Note this is NOT a link to an admin user as audits may be carried out by people other than admin users.',
  `comment` TEXT NULL COMMENT 'An optional comment about the audit.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'An audit record defines an evaluation period for which all of the claimed CPD activity will be verified by an admin user.\n';


-- name: create-table-ce_audit_m_activity
CREATE TABLE IF NOT EXISTS `%s`.`ce_audit_m_activity` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ce_audit_id` INT NOT NULL COMMENT 'The audit that this activity is part of. Note: This reference may be redundant but is probably here for convenience.',
  `ce_m_activity_id` INT NOT NULL COMMENT 'The member cps activity record to be verified.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created date',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `verified` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to signal if the activity record has been verified. 1 = verified.',
  `comment` TEXT NULL COMMENT 'An optional comment.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'This table stores a reference to all of the activity records that fall within the date range of the evaluation period being subject to audit. Each CPD activity will be verified by checking supporting evidence.';


-- name: create-table-fn_inventory
CREATE TABLE IF NOT EXISTS `%s`.`fn_inventory` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name` VARCHAR(100) NOT NULL COMMENT 'Name of the item.',
  `description` TEXT NOT NULL COMMENT 'Description of the item.',
  `unit_charge` DECIMAL(10,2) NOT NULL COMMENT 'Charge per unit, ex Tax.',
  `unit_name` VARCHAR(45) NOT NULL COMMENT 'The type of unites - eg item, hours, kg.',
  `tax` TINYINT(1) NOT NULL COMMENT 'Does tax apply to this item? e.g. GST or VAT.\n',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Common line items that may appear on invoices. Defaults are set here but can be overridden in the actual invoice.';


-- name: create-table-fn_m_invoice
CREATE TABLE IF NOT EXISTS `%s`.`fn_m_invoice` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier, will also be the invoice number so will need to start from a value that exceeds the current highest issued invoices from previous systems.',
  `member_id` INT NOT NULL COMMENT 'The member to whom the invoice has been issued.',
  `fn_subscription_id` INT(11) NULL COMMENT 'The original subscription that the invoice was generated from. The invoice will come from the fn_m_subscription record tied to the member however this may change over time so we keep a record of the original subscription.',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `paid` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Flag to indicate the invoice has been paid - used to speed up invoices repots.',
  `system_send` TINYINT NOT NULL DEFAULT 1 COMMENT 'A flag to trigger sending of the invoice by the housekeeping script. Defaults to 1 but can be overridden by the admin set when they set up a new invoice.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `completed_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Time stamp that signals when the invoice has been marked as complete and will be visible to the member.',
  `last_sent_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'The date and time that the invoice was last sent to the member as an email and attached invoice.',
  `invoiced_on` DATE NULL DEFAULT NULL COMMENT 'The invoice date.',
  `due_on` DATE NULL DEFAULT NULL COMMENT 'Invoice due date',
  `start_on` DATE NULL DEFAULT NULL COMMENT 'Defines the start date for the billing or subscription period of this invoice, when applicable.',
  `end_on` DATE NULL DEFAULT NULL COMMENT 'Defines the end date for the billing or subscription period of this invoice, when applicable.',
  `invoice_total` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT 'The total of the invoice (calculated for convenience wrt on screen reporting)	',
  `comment` TEXT NULL COMMENT 'An invoice comment if required.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Invoices issued to members. An invoice is always issued to a member, but the payments may be received from organisations.';


-- name: create-table-fn_invoice_payment
CREATE TABLE IF NOT EXISTS `%s`.`fn_invoice_payment` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `fn_m_invoice_id` INT NOT NULL COMMENT 'The invoice to which the payment has been allocated.',
  `fn_payment_id` INT NOT NULL COMMENT 'The payment record from which the allocation has been deducted.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record created at.',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated at.',
  `amount` DECIMAL(10,2) NOT NULL COMMENT 'The amount that has been allocated from the total payment amount. May be all or part of the payment amount.',
  `comment` TEXT NULL COMMENT 'A comment about the allocation?\n',
  PRIMARY KEY (`id`),
  INDEX `invoice_id_idx1` (`fn_m_invoice_id` ASC),
  INDEX `payment_id_idx1` (`fn_payment_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'Records the allocation of an amount of a payment (all or part) to an invoice.';


-- name: create-table-fn_m_subscription
CREATE TABLE IF NOT EXISTS `%s`.`fn_m_subscription` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'The member who owns this subscription.',
  `fn_subscription_id` INT(11) NOT NULL COMMENT 'The subscription template',
  `ad_macro_id` INT NULL COMMENT 'If this record was inserted or updated by a macro we store the macro id here.',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete.',
  `complimentary` TINYINT NOT NULL DEFAULT 0 COMMENT 'Flag to tell housekeeping to skip invoice generation when processing this subscription renewal. All other steps will be taken - i.e. the renewal date will be updated etc, but the invoke will not be generated.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated.',
  `pro_rata_bill_on` DATE NULL COMMENT 'Option pro-rata date. Ideally this should be NULL if it is ',
  `renew_on` DATE NULL DEFAULT NULL COMMENT 'The next date for full period renewal of the subscription and generation of the invoice.',
  `comment` TEXT NULL COMMENT 'General comment.',
  PRIMARY KEY (`id`),
  INDEX `fk_member_subscription_subscription1_idx` (`fn_subscription_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'A members subscription.';


-- name: create-table-fn_payment
CREATE TABLE IF NOT EXISTS `%s`.`fn_payment` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `fn_payment_type_id` INT NOT NULL COMMENT 'Specifies the type of payment',
  `member_id` INT NOT NULL COMMENT 'Set if the payment was received from a member, else 0.',
  `organisation_id` INT NULL COMMENT 'Set if the payment was received from an organisation, else set to 0.',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `payment_on` DATE NULL DEFAULT NULL COMMENT 'Payment received on date.',
  `amount_received` DECIMAL(10,2) NOT NULL COMMENT 'The total amount received.',
  `comment` TEXT NULL COMMENT 'An optional comment about the payment.',
  `field1_data` VARCHAR(45) NULL COMMENT 'The descriptive data for general field specified in fn_payment_type table.',
  `field2_data` VARCHAR(45) NULL COMMENT 'The descriptive data for general field specified in fn_payment_type table.',
  `field3_data` VARCHAR(45) NULL COMMENT 'The descriptive data for general field specified in fn_payment_type table.',
  `field4_data` VARCHAR(45) NULL COMMENT 'The descriptive data for general field specified in fn_payment_type table.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Records the receipt of a sum of money from a member OR an organisation. Thus this table has no _m_ in its name as payments may optionally be specified as being from an organisation. Payments are made and must be allocated against one or more invoices.';


-- name: create-table-wf_note
CREATE TABLE IF NOT EXISTS `%s`.`wf_note` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `wf_note_type_id` INT NOT NULL COMMENT 'The type of note.',
  `ad_user_id_created` INT NULL COMMENT 'Admin user who created the note. NULL values were allowed for import of historic data but new records will be set to logged in user.',
  `ad_user_id_updated` INT NULL COMMENT 'Admin user who last updated the note. NULL values were allowed for import of historic data but new records will be set to logged in user.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `effective_on` DATE NULL DEFAULT NULL COMMENT 'The relevant date of the event or item referred to by the note. Allows for user to specify the date of something even if the note is added (created_at_ much later.',
  `note` TEXT NOT NULL,
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores administrative notes which may be linked to a member record, or linked to an issue record that is global. That is, an Issue record that does not have a link to a member. \n\nIf the note is linked to a member the member_id field will be filled. The note may also be associated with an issue or an application record, in which case it should STILL have a member_id value.\n\nThe only time the member_id is 0 is when the note is linked to a gloabl Issue, which itself does not have a member_id.\n';


-- name: create-table-mp_m_accreditation
CREATE TABLE IF NOT EXISTS `%s`.`mp_m_accreditation` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'The member who has attained this accreditation.',
  `mp_accreditation_id` INT NOT NULL COMMENT 'The accredittion that is held by the member.',
  `organisation_id` INT NULL COMMENT 'The organisation that grants this accreditation.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `start_on` DATE NULL COMMENT 'Date on which the accreditation was obtained',
  `end_on` DATE NULL DEFAULT NULL COMMENT 'Date on which the accreditation expires, if it has limited tenure',
  `comment` TEXT NULL COMMENT 'Optional comment field.',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `member_accreditation_id_UNIQUE` (`member_id` ASC, `mp_accreditation_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'Accreditations held by individual members. Maintained as a history table.';


-- name: create-table-mp_accreditation
CREATE TABLE IF NOT EXISTS `%s`.`mp_accreditation` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `organisation_id` INT NOT NULL COMMENT 'The organisation that is tasked with granting a particular accreditation.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created on',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `display_priority` SMALLINT NOT NULL DEFAULT 0 COMMENT 'An optional value that will allow records to be ordered when displayed. Higher priority is displayed first.',
  `short_name` VARCHAR(45) NOT NULL COMMENT 'Accepted short name / abbrev. for the accreditation.',
  `name` VARCHAR(80) NOT NULL COMMENT 'Full name of the accreditation',
  `description` TEXT NOT NULL COMMENT 'Description of the accreditation.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Accreditations are industry-specific qualifications or acknowledgements, generally with limited tenure and with very specify scope. Otherwise they are identical in most respects to qualifications. This table defines the types of accreditations.';


-- name: create-table-mp_m_position
CREATE TABLE IF NOT EXISTS `%s`.`mp_m_position` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'The member who holds this position.',
  `mp_position_id` INT NOT NULL COMMENT 'The position held by the member.',
  `organisation_id` INT NULL COMMENT 'Link to the organisation, may be inherited value via position record. Or specified by the admin.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `start_on` DATE NULL COMMENT 'The date on which the member was appointed to this position.',
  `end_on` DATE NULL DEFAULT NULL COMMENT 'If limited tenure, the date on which this position will expire. This will be used to raise expiring type housekeeping issues.',
  `comment` TEXT NULL COMMENT 'Optional comment field.',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `member_position_organisation_starton_id_UNIQUE` (`member_id` ASC, `mp_position_id` ASC, `organisation_id` ASC, `start_on` ASC))
  ENGINE = InnoDB
  COMMENT = 'Position help by individual members. A member may hold multiple positions, with multiple organisations, but not the same position with same organisation (need to define a unique index here).';


-- name: create-table-ms_m_status
CREATE TABLE IF NOT EXISTS `%s`.`ms_m_status` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'Link to member',
  `ms_status_id` INT NOT NULL COMMENT 'Link to status',
  `ad_macro_id` INT NULL COMMENT 'If this record was inserted or updated by a macro we store the macro id here.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `current` TINYINT NOT NULL DEFAULT 0 COMMENT 'Defines the current (i.e. latest) status and is used to make table joins easier when creating lists of members.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `comment` TEXT NULL COMMENT 'Optional comment.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'The membership status history of a member.'
  PACK_KEYS = Default;


-- name: create-table-ms_title
CREATE TABLE IF NOT EXISTS `%s`.`ms_title` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ms_subscription_id_default` INT NOT NULL COMMENT 'Default subscription template for this membership.',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `login` TINYINT NOT NULL DEFAULT 0 COMMENT 'Can login',
  `directory` TINYINT NOT NULL DEFAULT 0 COMMENT 'Will be shown in directory',
  `subscription` TINYINT NOT NULL DEFAULT 0 COMMENT 'Members with this membership should have / are subject to the application of fees via a subscription of type MEMBERSHIP. If set to 1 the member record should have an active membership subscription. If set to 0 they can have any other type of subscription but an existing membership subscription will be ignored when renewal invoices are run.',
  `application` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to indicate whether a Title is eligible to be applied for via an Application. This flag is used to limit the lists where apt as well.',
  `cpd` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to indicate if the membership title is subject to the CPD component rules.',
  `name` VARCHAR(45) NOT NULL COMMENT 'The \'title\' of the Membership',
  `description` VARCHAR(255) NOT NULL COMMENT 'Descriptive text',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines the Membership names or titles as defined by the organisation. Eg Associate, Fellow, Non-member etc. In combination with status table also defined some privileges for the member within the system.';


-- name: create-table-mp_m_contact
CREATE TABLE IF NOT EXISTS `%s`.`mp_m_contact` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'Member that owns the contact card.',
  `mp_contact_type_id` INT NOT NULL COMMENT 'Type of contact card',
  `country_id` INT NULL COMMENT 'Country for contact card address.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `mobile` VARCHAR(45) NULL COMMENT 'Optional mobile number',
  `phone` VARCHAR(45) NULL COMMENT 'Option phone (landline)',
  `fax` VARCHAR(45) NULL COMMENT 'Optional fax number',
  `email` VARCHAR(100) NULL COMMENT 'Optional email',
  `web` VARCHAR(100) NULL COMMENT 'Optional web site url.',
  `address1` VARCHAR(45) NULL COMMENT 'Optional address information',
  `address2` VARCHAR(45) NULL,
  `address3` VARCHAR(45) NULL,
  `locality` VARCHAR(45) NULL COMMENT 'Optional town / suburb',
  `state` VARCHAR(45) NULL COMMENT 'Optional state',
  `postcode` VARCHAR(45) NULL COMMENT 'Optional post code',
  `comment` TEXT NULL COMMENT 'Optional comment to explain the card.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores an individual contact card type for a member.';


-- name: create-table-mp_contact_type
CREATE TABLE IF NOT EXISTS `%s`.`mp_contact_type` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `persistent` TINYINT NOT NULL COMMENT 'If set to \'1\' this card may NOT be deleted. That is, all member records should have a card of this type, even if all the fields are empty.',
  `allow_member_edit` TINYINT NOT NULL DEFAULT 0 COMMENT 'Flag to allow the member to edit this type of contact record. This only applies when persistent = 1 and is ignored for records where persistent = 0. In other words, admin can specify if a member is allowed to edit a persistent contact card type, but members can always add / edit non-persistent contact card types. ',
  `tax_lookup_priority` TINYINT NOT NULL DEFAULT 0 COMMENT 'Defines the priority with which this type of contact card (address) is used for tax calculations. That is, the Country value in this card is used to lookup tax tables for that country. The system will look for a country value in the cards 9in order of priority) then will use the first country value it finds to lookup taxes. If it fails to find any Country values it will fall back to the membership country.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `order` TINYINT NOT NULL,
  `name` VARCHAR(45) NOT NULL COMMENT 'A name for the card type',
  `description` VARCHAR(255) NOT NULL COMMENT 'A description of the card type',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines the type of contact card e.g. Primary, Courier, Home etc.';


-- name: create-table-mp_position
CREATE TABLE IF NOT EXISTS `%s`.`mp_position` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `organisation_id` INT NOT NULL COMMENT 'Optional link to the Organisation if the position is unique to an organisation.',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `display_priority` SMALLINT NOT NULL DEFAULT 0 COMMENT 'An optional value that will allow records to be ordered when displayed. Higher priority is displayed first.',
  `short_name` VARCHAR(45) NOT NULL COMMENT 'Short name or abbreviation.',
  `name` VARCHAR(100) NOT NULL COMMENT 'The name of the position',
  `description` TEXT NOT NULL COMMENT 'Descriptive text about the position.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores all the possible positions that a member may hold in various organisations. Eg President, Chair, Member, etc.';


-- name: create-table-ms_m_application
CREATE TABLE IF NOT EXISTS `%s`.`ms_m_application` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'Link to the member who submitted the application',
  `member_id_nominator` INT NULL COMMENT 'The Member who nominated this member\'s application… allowed NULL to migrate old data ',
  `member_id_seconder` INT NULL COMMENT 'The Member who seconded this member\'s application … allowed NULL to migrate old data ',
  `ms_title_id` INT NOT NULL COMMENT 'The Membership (title) being applied for. ',
  `ad_macro_id` INT NULL COMMENT 'If this record was inserted or updated by a macro we store the macro id here.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `applied_on` DATE NULL DEFAULT NULL COMMENT 'Date the application was received / acknowledged / processed.',
  `result` TINYINT NOT NULL DEFAULT -1 COMMENT 'The outcome of the application -  value of -1 is the default case in which the outcome of the application is unknown or PENDING. An explicit 0 is REJECTED and an explicit 1 is ACCEPTED. We can add other explicit values later if needed.\n\nEach of the ms_m_application_meeting records relating to an application also have this flag. Thus we can create a boolean workflow where if all meeting records show a 1 we can infer that the application should also be a 1. This also enables us to introduce other types of events in the workflow.',
  `comment` TEXT NULL,
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores records of applications from members to attain a specific membership (title).';


-- name: create-table-a_meeting
CREATE TABLE IF NOT EXISTS `%s`.`a_meeting` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `a_meeting_type_id` INT NOT NULL COMMENT 'The meeting type.',
  `active` TINYINT(4) NOT NULL DEFAULT '1' COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `meeting_on` DATE NULL DEFAULT NULL COMMENT 'Date of meeting',
  `location` VARCHAR(100) NULL COMMENT 'Location of meeting',
  `name` VARCHAR(255) NOT NULL COMMENT 'Name or descriptive title of the meeting.',
  `comment` TEXT NULL COMMENT 'General comments about the meeting.',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `a_meeting_type_on_location_UNIQUE` (`a_meeting_type_id` ASC, `meeting_on` ASC, `location` ASC))
  ENGINE = InnoDB
  COMMENT = 'Stored general info about board meetings which are linked to';


-- name: create-table-wf_attachment
CREATE TABLE IF NOT EXISTS `%s`.`wf_attachment` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `wf_note_id` INT NOT NULL COMMENT 'The note that this file is attached to / with.',
  `ad_user_id` INT NOT NULL COMMENT 'The admin user who added the file.',
  `fs_set_id` INT NOT NULL COMMENT 'Tells us in which set we will find file system information that will allow us to locate the file.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `clean_filename` VARCHAR(255) NOT NULL COMMENT 'The filename of the original document prior to upload - stored for defence.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Documents table provides a way to attach files to notes. A document is always attached via a note and, therefore, can only be associated with member records. As notes can be associated with applications and issues, documents can also be associated with these entities.';


-- name: create-table-ms_status
CREATE TABLE IF NOT EXISTS `%s`.`ms_status` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `system` TINYINT NOT NULL DEFAULT 0 COMMENT 'System flag prevents the record from being editable via the application. If set to 1 we will ensure it does not get modified.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `login` TINYINT NOT NULL COMMENT 'Can login',
  `directory` TINYINT NOT NULL COMMENT 'Will be shown in directory',
  `subscription` TINYINT NOT NULL COMMENT 'Members with this membership should have / are subject to the application of fees via a subscription of type MEMBERSHIP. If set to 1 the member record should have an active membership subscription. If set to 0 they can have any other type of subscription but an existing membership subscription will be ignored when renewal invoices are run.',
  `workflow` TINYINT NOT NULL DEFAULT 0 COMMENT 'This flag tells the system to include members with this status as part of the workflow functionality. So if a member status has 1 here that member can have issues raised against their record. A zero value here means that no new issues will be raise, and any existing issues will be closed by the system housekeeping scripts.',
  `cpd` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to indicate if the membership title is subject to the CPD component rules.',
  `name` VARCHAR(45) NOT NULL COMMENT 'The name of the status value',
  `description` VARCHAR(255) NOT NULL COMMENT 'Descriptive text',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  AUTO_INCREMENT = 10000
  COMMENT = 'Status table describes the membership or members status and, in combination with the membership title, various attributes of the membership.';


-- name: create-table-fn_invoice_inventory
CREATE TABLE IF NOT EXISTS `%s`.`fn_invoice_inventory` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `fn_m_invoice_id` INT NOT NULL COMMENT 'The invoice on which this line item appears.',
  `fn_inventory_id` INT NOT NULL COMMENT 'The inventory record represented by this line item.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `description` TEXT NOT NULL COMMENT 'The description of the line item, taken from inventory but can be edited.',
  `quantity` DECIMAL(8,2) NOT NULL DEFAULT 1 COMMENT 'Quantity / multiplier of units.',
  `unit_charge` DECIMAL(8,2) NOT NULL DEFAULT 0 COMMENT 'The unit charge for this line item, taken as a snapshot from the inventory record, excluding tax. Can be edited.',
  `tax_rate` DECIMAL(5,2) NOT NULL DEFAULT 0 COMMENT 'The applicable tax for this line item. This will be calculated based on a set of business rules.',
  `tax_name` VARCHAR(45) NULL COMMENT 'The name of the tax applied from the fn_tax table. Stored with record in case tax table changes.',
  PRIMARY KEY (`id`),
  INDEX `invoice_id_idx1` (`fn_m_invoice_id` ASC),
  INDEX `inventory_id_idx1` (`fn_inventory_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'Defines the line items that appear on an invoice. These will generally come from the inventory table. ';


-- name: create-table-ol_module
CREATE TABLE IF NOT EXISTS `%s`.`ol_module` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_module_id_original` INT NOT NULL COMMENT 'References the original module id and is used for versioning. When a new module is created with field will be set to the current id.  All subsequent revisions of a module will refer to a single original, and the version number will be used to determine their order. (Could just use created_at really?)',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `current` TINYINT NOT NULL DEFAULT 0 COMMENT 'Flag the current revision',
  `revision` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'A revision number used when modules are updated. Any module that has been attempted in the past cannot be modified and instead a new revision is created and the previous one de-activated.',
  `feedback` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Feedback flag for the module. 0 - no feedback step, 1 - feedback step is there but optional, 2 - Feedback step is compulsory.',
  `started` MEDIUMINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'The number of times this module has been started. A start is an entry in the ol_m_module table. Stored to make it easier to determine if modules can be edited etc.',
  `finished` MEDIUMINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'The number of times this module has been finished.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `published_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Date and time the module was marked as complete.',
  `name` TEXT NOT NULL COMMENT 'The module name or descriptive title',
  `description` TEXT NOT NULL COMMENT 'Description of the module',
  `objective` TEXT NOT NULL COMMENT 'The learning outcome / objective of the module',
  `instruction` TEXT NOT NULL COMMENT 'Instructions to the user on how to complete the module.',
  `pass_percentage` TINYINT UNSIGNED NULL COMMENT 'The nominal percentage required for this to be a pass. If NULL then scores are not relevant.',
  `estimated_total_mins` SMALLINT UNSIGNED NOT NULL COMMENT 'A guide for the user to indicate how long the entire module should take to complete. Given in minutes it is not used in any calculations.',
  `attempt_limit` SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Defines the number of times a user is allowed to attempt this module within the attempt_limit_reset_days.',
  `attempt_limit_reset_days` SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'After this many days the user is able to attempt the module again as many times as is defined in the attempt_limit field.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'A Module is an individual learning unit comprising one or more media resources and an optional set of questions. A module may be assigned to one or more relevant subtopics.';


-- name: create-table-ol_m_module
CREATE TABLE IF NOT EXISTS `%s`.`ol_m_module` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'The member who has undertaken the module',
  `ol_module_id` INT NOT NULL COMMENT 'The module being attempted',
  `ce_m_activity_id` INT NULL DEFAULT NULL COMMENT 'If the module has CPD points awarded on completion, we link to the actual member cpd activity record so we are able to display the historical CPD value of the module to the user. It may also be useful for reporting and auditing as we have data proof that the module was actually attempted and have a direct connection between the activity record and the module attempt. ',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created at, also flag the start of the module.',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated, also reflects the completion date and time if the module is complete.',
  `slides_completed_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'The date and time the module questions were completed.',
  `module_completed_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'The date and time the entire module was completed including feedback step. This is the flag that closes the module off from further editing.',
  `resource_seconds` SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Cumulative seconds that the user has the resources page opened. Incremented using an ajax request.',
  `feedback` TEXT NULL COMMENT 'Member can give feedback about content once it is completed. This is freeform text feedback, as opposed to the rating type of feedback.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'A learning module undertaken by a member.';


-- name: create-table-ol_module_resource
CREATE TABLE IF NOT EXISTS `%s`.`ol_module_resource` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_module_id` INT NOT NULL COMMENT 'The module...',
  `ol_resource_id` INT NOT NULL COMMENT 'The media resource',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `ol_module_resource_UNIQUE` (`ol_module_id` ASC, `ol_resource_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'A resource associated with a learning module.';


-- name: create-table-ol_slide
CREATE TABLE IF NOT EXISTS `%s`.`ol_slide` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_module_id` INT NOT NULL COMMENT 'The module the question relates to.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `sequence` TINYINT NOT NULL COMMENT 'The order in which the questions will appear.',
  `timer_max_minutes` TINYINT NOT NULL DEFAULT 5 COMMENT 'The maximum time in minutes that we will increment the module timer feature. This is set to avoid the situation where the user leaves the module to attend to other things. In this case we will only increment the module timer by this amount.',
  `type` ENUM('info','question') NOT NULL COMMENT 'A slide can be of type QUESTION (which requires answers) or of type INFO which can contain info only.',
  `summary` TEXT NOT NULL COMMENT 'Slide summary is the short content for the slide- for a question the actual question copy goers into this field, and for an info slide this is used as a title or excerpt explaining the slide content. This will be plain text only and thus can appear in the slide list view as well as the results page.',
  `content` TEXT NULL COMMENT 'The slide content field is used to store more complex content to be shown on the slide, below the summary - for example HTML content, embedded images or a table. This is optional for both info and question slides however it will generally be used on an info slide.\n',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'A slide can be of type question or info and is used to progress the user through the module.';


-- name: create-table-ol_m_module_slide_option
CREATE TABLE IF NOT EXISTS `%s`.`ol_m_module_slide_option` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_m_module_slide_id` INT NOT NULL COMMENT 'The question (slide) record being answered, in the context of the members attempt on the module - i.e. the instance of this question (slide) within the current attempt on the module. This relationship is required because each question is allowed to have more than one answer given.\n\n',
  `ol_option_id` INT NOT NULL COMMENT 'The option selected by the member.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'A members answer to a module they have undertaken.';


-- name: create-table-ol_option
CREATE TABLE IF NOT EXISTS `%s`.`ol_option` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_slide_id` INT NOT NULL COMMENT 'The question this answer relates to.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `sequence` TINYINT NOT NULL COMMENT 'Will specify the order that the answers will appear in the multiple choice format.',
  `positive_value` TINYINT NOT NULL DEFAULT 0 COMMENT 'Flag to specify that this answer is correct, and the relative value / score / mark for selecting this answer. ',
  `negative_value` TINYINT NOT NULL DEFAULT 0 COMMENT 'If positive_value is 0 then the answer is incorrect and this value (score / mark) will be DEDUCTED from the score. If positive_value is > 0 then the answer is correct and this value (score / mark) will be DEDUCTED from the score in the event that the user DOES NOT select it.\n\n\n',
  `option` TEXT NOT NULL COMMENT 'The answer text.',
  `explanation` TEXT NULL COMMENT 'Option to provide an explanatory note for answers, including the incorrect ones. This gets displayed on the final results page.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Each question can have one-to-many answers. A true false question would have two possible answers (true / false) but questions can have any number of answers for the multiple choice format.';


-- name: create-table-ol_resource
CREATE TABLE IF NOT EXISTS `%s`.`ol_resource` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_resource_type_id` INT NOT NULL COMMENT 'The \'type\' of resource, e.g. image, video, sound, document.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `primary` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to indicate that the resource is significant, substantial, complete etc. That is, the resource contains a reference to education material that could be used on its own, outside the contact of a module. Examples might be external web sites, video lectures, published papers. This flag should be set to 0 when the resource is a kind of fragment that is used in context with questions - such as a labelled image, and short sound file  and so on. \n\nPrimary resources will be featured on the Member login page - so this flag allows us to endure trivial resources can be left out if required.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `presented_on` DATE NULL DEFAULT NULL COMMENT 'Optional date that specifies when the resource item was presented, published or any other relevant date attribute.',
  `presented_year` MEDIUMINT(4) NULL COMMENT 'The year the resource was created / published / presented.',
  `presented_month` TINYINT NULL COMMENT 'The month the resource was created / published / presented',
  `presented_date` TINYINT NULL COMMENT 'The date (day) the resource was created / published / presented.',
  `name` TEXT NOT NULL COMMENT 'A  name for the resource, must be unique. NB have used TEXT here as some resource names are very long.\n',
  `description` TEXT NOT NULL COMMENT 'A description of the resource',
  `keywords` TEXT NOT NULL COMMENT 'Comma separated keywords to assist with resource library search.',
  `resource_url` VARCHAR(255) NOT NULL COMMENT 'The resource location… either a URL in the system CDN or an external URL.',
  `short_url` VARCHAR(255) NULL,
  `thumbnail_url` VARCHAR(255) NOT NULL COMMENT 'A url for thumbnail image used as a cover sheet or link to the actual resource.',
  `ol_resource_related_ids` VARCHAR(255) NULL COMMENT 'Comma separated list of related resource ids - hack for now but allows us to show related resources.',
  `attributes` TEXT NULL,
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'A resource record is an individual piece of media content, such as an image, a video or sound file, a document or an external website. The collection of resources makes up the media library to support all of the individual learning modules.';


-- name: create-table-ol_category
CREATE TABLE IF NOT EXISTS `%s`.`ol_category` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `code` VARCHAR(10) NOT NULL COMMENT 'Short category code',
  `name` VARCHAR(255) NOT NULL COMMENT 'Category name',
  `description` TEXT NOT NULL COMMENT 'Category description',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'The top level topic categories.';


-- name: create-table-ol_resource_type
CREATE TABLE IF NOT EXISTS `%s`.`ol_resource_type` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` VARCHAR(45) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name` VARCHAR(45) NOT NULL COMMENT 'Resource type name',
  `description` VARCHAR(255) NOT NULL COMMENT 'Resource type description',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines the type of resource. Eg. Image, Video, Audi etc.';


-- name: create-table-ce_event
CREATE TABLE IF NOT EXISTS `%s`.`ce_event` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `start_on` DATE NULL DEFAULT NULL COMMENT 'The date the event was held or the start date for a multi-day event.',
  `end_on` DATE NULL DEFAULT NULL COMMENT 'The date the event ends, optional for a single day event.',
  `location` VARCHAR(255) NOT NULL COMMENT 'The location of the event.',
  `name` VARCHAR(255) NOT NULL COMMENT 'The name of the event.',
  `description` TEXT NOT NULL COMMENT 'A description of the event.',
  `information_url` TEXT NULL COMMENT 'Allows us to create a link to a website url that has information about the event.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'A record of CPD related events. These records are used to assist with bulk recording of CPD via macros.';


-- name: create-table-cm_email_template
CREATE TABLE IF NOT EXISTS `%s`.`cm_email_template` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `system` TINYINT NOT NULL DEFAULT 0 COMMENT 'Defines if template is used by system or created by admin. Value of 1 indicates a system email template and cannot be deleted, but can be modified by admins.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created at',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated at',
  `name` VARCHAR(100) NOT NULL COMMENT 'A friendly name for the template, eg Monthly Newsletter',
  `description` VARCHAR(255) NOT NULL COMMENT 'A more detailed description of the template',
  `subject` VARCHAR(255) NOT NULL COMMENT 'The default email subject for the template',
  `from_name` VARCHAR(100) NOT NULL COMMENT 'The text name for the From: email header ',
  `from_email` VARCHAR(100) NOT NULL COMMENT 'The email address for the From: header',
  `body_html` TEXT NOT NULL COMMENT 'The html version of the email body',
  `cm_email_variable_ids` TEXT NULL COMMENT 'A comma separated list of the variable ids (from cm_email_variable table) that are available in the email template, other than those marked as persistent in cm_email_variable table (which are always available).',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  AUTO_INCREMENT = 10000
  COMMENT = 'Stores standard  template emails that may be sent frequently. Simply used as a starting point for the creation of a new email.'
  PACK_KEYS = Default;


-- name: create-table-cm_email
CREATE TABLE IF NOT EXISTS `%s`.`cm_email` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unqiue identifier',
  `cm_email_template_id` INT NULL COMMENT 'The (optional) template that was used as a starting point for this email.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `ready` TINYINT NOT NULL DEFAULT 0 COMMENT 'Simple flag to ensure the user is happy for the email to be broadcast at the broadcast_at time date.',
  `all_sent` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag that is set by a cron job when there are no cm_m_email recipients with sent = 0. This is so that the sending cron script can easily locate emails for which it should check ay pending deliveries rather than checking them all every time.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created.',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated.',
  `broadcast_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'The date the email is scheduled, or was scheduled to be broadcast - i.e. can be future or past date. Cron jobs use this date to decide if they need to process emails for this blast.',
  `name` VARCHAR(100) NOT NULL COMMENT 'A name for this communication which is descriptive of the purpose.',
  `comment` TEXT NULL COMMENT 'If required, a comment about the purpose of the communication.',
  `subject` VARCHAR(255) NOT NULL COMMENT 'The subject for the email',
  `from_name` VARCHAR(100) NOT NULL COMMENT 'Name portion of the From: header',
  `from_email` VARCHAR(100) NOT NULL COMMENT 'Email portion of the From: header',
  `body_html` TEXT NOT NULL COMMENT 'The HTML version of the email body.',
  `cm_email_variable_ids` TEXT NULL COMMENT 'Comma separated id values of the variable that are / were available for this email - A copy of the same field from cm_email_template table so we have an instance of the variables that were available in a particular template.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'An email that is blasted to one or more members.';


-- name: create-table-cm_m_email
CREATE TABLE IF NOT EXISTS `%s`.`cm_m_email` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'The member that received the email communication.',
  `cm_email_id` INT NOT NULL COMMENT 'The email that was sent.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created, ie when email was sent.',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `pickup_at` TIMESTAMP NULL COMMENT 'This field gets set to current time stamp when the email sending script picks up this record as part of its next batch of emails to broadcast. This is effectively a flag to say that the script started to process this email, and if the processing is completed, the process field gets set to 1. It provides a way for us to retry sending emails if the script did not complete for some reason - i.e. a start and end flag.',
  `processed` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to indicate that the system has made an attempt to form an email and send it to this recipient. If, for example, the member email address was missing or was malformed process would be set to 1 but sent would be set to 0. Cron job will use this flag rather than sent so that it does not indefinitely try to resend to members without email address or some similar permanent error.',
  `sent` TINYINT NOT NULL DEFAULT 0 COMMENT 'Sent by system - that is, the application sent the email but we don\'t know what happened to it after that. This flag is used by the cron jobs that send emails via the third party system. This flag is set after processed is set to 1 (that is all the bits are in place) and the php mail command has been executed.',
  `system_message` VARCHAR(255) NOT NULL COMMENT 'A message from the system to expand on the result of the scripts mail command. For example - Could not send because the user has no email address.',
  `sent_to` VARCHAR(255) NOT NULL COMMENT 'Store the recipients email for the sake of data clarity e.g. if email changes.',
  `parameters` TEXT NULL COMMENT 'Serialised values for personalisation variables in the email template.',
  `mx_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'The date and time that the current status was reported from the mail exchanger, via cm_email_log table data.',
  `mx_status` VARCHAR(100) NULL COMMENT 'The latest status of the email as reported from the mail exchanger, via cm_email_log table data.',
  `mx_description` TEXT NULL COMMENT 'A description of the latest email status reported from the mail exchanger, via cm_email_log table data.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Records all instances of email communications sent to individual members.';


-- name: create-table-acl_member_resource
CREATE TABLE IF NOT EXISTS `%s`.`acl_member_resource` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `parent_id` INT(11) NULL COMMENT 'Self-referencing id to indicate parent record.',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `type` VARCHAR(45) NOT NULL COMMENT 'Type of resource - e.g. page or action.',
  `name` VARCHAR(255) NOT NULL COMMENT 'The name of the resource',
  `description` TEXT NULL DEFAULT NULL COMMENT 'A description of the resource',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  AUTO_INCREMENT = 10000
  COMMENT = 'Defined ACL for resources in the application.';


-- name: create-table-acl_member_role
CREATE TABLE IF NOT EXISTS `%s`.`acl_member_role` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT 'Specified that this role is the default assigned to a member',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `name` VARCHAR(255) NOT NULL COMMENT 'Name of role (group)',
  `description` TEXT NULL DEFAULT NULL COMMENT 'Description of role (group)',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines roles of groups into which a member is assigned. The member with then inherit all of the permissions associated with this role.';


-- name: create-table-acl_member_role_resource
CREATE TABLE IF NOT EXISTS `%s`.`acl_member_role_resource` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `acl_member_role_id` INT(11) NOT NULL COMMENT 'The member role from which this resource may be accessed',
  `acl_member_resource_id` INT(11) NOT NULL COMMENT 'The member resource that may be accessed from this role',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `privileges` VARCHAR(100) NULL DEFAULT NULL COMMENT 'Privileges (not in use yet)',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `acl_member_role_resource_id_UNIQUE` (`acl_member_role_id` ASC, `acl_member_resource_id` ASC))
  ENGINE = InnoDB
  AUTO_INCREMENT = 9;


-- name: create-table-ad_macro
CREATE TABLE IF NOT EXISTS `%s`.`ad_macro` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ad_user_id` INT NOT NULL COMMENT 'The id of the admin user that executed the macro.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created date - the date and time the macro was run',
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `name` VARCHAR(45) NOT NULL COMMENT 'The name of the macro that was run - maybe set by the system on run.',
  `comment` TEXT NULL COMMENT 'An optional description of the macro / reason for running.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores  records of Macros that are executed by admin users.';


-- name: create-table-ol_slide_resource
CREATE TABLE IF NOT EXISTS `%s`.`ol_slide_resource` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_slide_id` INT NOT NULL,
  `ol_resource_id` INT NOT NULL,
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `ol_slide_resource_id_UNIQUE` (`ol_slide_id` ASC, `ol_resource_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'Resources associated with a specific slide.';


-- name: create-table-ol_module_category
CREATE TABLE IF NOT EXISTS `%s`.`ol_module_category` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_module_id` INT NOT NULL COMMENT 'The module in question',
  `ol_category_id` INT NOT NULL COMMENT 'The category this module belongs to.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'record last updated',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `ol_module_category_UNIQUE` (`ol_module_id` ASC, `ol_category_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'Assignment of individual learning modules to one or more categories.';


-- name: create-table-ol_module_rating
CREATE TABLE IF NOT EXISTS `%s`.`ol_module_rating` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `sequence` TINYINT NOT NULL COMMENT 'The order the questions might appear.',
  `question` TEXT NOT NULL COMMENT 'The question or prompt. Eg. \'How would you rate this modules relevance to….\'',
  `score_max` TINYINT NOT NULL COMMENT 'Max score for this rating question / prompt. ',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores the questions or prompts that illicit a rating (e.g. from 1-5) from the member in relation to the module.';


-- name: create-table-ol_m_module_rating
CREATE TABLE IF NOT EXISTS `%s`.`ol_m_module_rating` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_m_module_id` INT NOT NULL COMMENT 'The member\'s instance of the module being rated.',
  `ol_module_rating_id` INT NOT NULL COMMENT 'The module rating (question) being scored.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated.',
  `score` TINYINT UNSIGNED NOT NULL COMMENT 'The score given.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'The rating / score given for each rating question, from a member for a particular module.';


-- name: create-table-fn_payment_type
CREATE TABLE IF NOT EXISTS `%s`.`fn_payment_type` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated	',
  `name` VARCHAR(45) NOT NULL COMMENT 'Name of the payment type',
  `description` VARCHAR(255) NOT NULL COMMENT 'Description of the payment type	',
  `field1_name` VARCHAR(45) NULL COMMENT 'Name of generalised description field in fn_payment table.',
  `field2_name` VARCHAR(45) NULL COMMENT 'Name of generalised description field in fn_payment table.',
  `field3_name` VARCHAR(45) NULL COMMENT 'Name of generalised description field in fn_payment table.',
  `field4_name` VARCHAR(45) NULL COMMENT 'Name of generalised description field in fn_payment table.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines the different types of payments. Have used generalised descriptive column names in this table for the corresponding data fields in the fn_payment table. This is not great design but there is no point creating an excessively generalised design for this.';


-- name: create-table-wf_note_type
CREATE TABLE IF NOT EXISTS `%s`.`wf_note_type` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `system` TINYINT NOT NULL DEFAULT 0 COMMENT 'Denotes a note type used by system which cannot be deleted / modified by admin. Set auto_increment to high value like 10000 so all system ids can be placed before this value.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created at',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name` VARCHAR(45) NOT NULL COMMENT 'Name for the note type.',
  `description` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  AUTO_INCREMENT = 10000
  COMMENT = 'Descriptive of the type of note e.g.. Contact, General etc.';


-- name: create-table-wf_note_association
CREATE TABLE IF NOT EXISTS `%s`.`wf_note_association` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'unique identifer',
  `wf_note_id` INT NOT NULL COMMENT 'The note record',
  `member_id` INT NULL COMMENT 'Defines the note as member-specific, i.e. this is the Member the note relates to. May be empty if the note relates to a \"global\" issue.  Ff this value is NULL the note should have a value for association_entity_id. That is, either member_id OR association_entity_id OR BOTH should have values.',
  `association_entity_id` INT NULL COMMENT 'The id of the record in the associated entity. if this value is NULL the note should have a value for member_id. That is either member_id OR association_entity_id OR BOTH should have values.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `association` VARCHAR(100) NULL COMMENT 'The name of the entity (table) with which we are associating the note. This may end up being the literal table name such as mp_m_application, or a friendly name like: application which is then mapped to the table name in a global config.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Maintains details of a Note record relationship to other entities.';


-- name: create-table-wf_issue_association
CREATE TABLE IF NOT EXISTS `%s`.`wf_issue_association` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `wf_issue_id` INT NOT NULL COMMENT 'The issue for which the association is being created.',
  `member_id` INT NULL COMMENT 'The member associated with the issue, if empty the issue is deemed to be global.',
  `association_entity_id` INT NULL COMMENT 'The id of the record being associated with.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `association` VARCHAR(100) NULL COMMENT 'The name of the entity (table) with which we are associating the issue. This may end up being the literal table name such as mp_m_application, or a friendly name like: application which is then mapped to the table name in a global config.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines associations between Issues and other entities. The main association will be with members via member_id, however an iassue can also be associated with other entities using this method.';


-- name: create-table-ce_activity_unit
CREATE TABLE IF NOT EXISTS `%s`.`ce_activity_unit` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `specify_quantity` TINYINT NOT NULL COMMENT 'Flag that allows the user to specify the quantity of the cpd activity. If this is set to 1 the user cn specify, if it is set to 0 then the quantity will be forced to  a value of 1 for a single instanc or item type of activity.',
  `name` VARCHAR(45) NOT NULL COMMENT 'The name of the unit, eg hours, days, item, event.',
  `description` VARCHAR(255) NULL COMMENT 'Optional description',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores the types of units used to measure CPD activity such as hours, days, items, rounds and whatever. ';


-- name: create-table-mp_tag
CREATE TABLE IF NOT EXISTS `%s`.`mp_tag` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name` VARCHAR(45) NOT NULL COMMENT 'The tag itself',
  `description` VARCHAR(255) NULL COMMENT 'Explanation of the tag.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB;


-- name: create-table-mp_m_tag
CREATE TABLE IF NOT EXISTS `%s`.`mp_m_tag` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'The member who is tagged',
  `mp_tag_id` INT NOT NULL COMMENT 'The tag being applauds',
  `active` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `member_tag_id_UNIQUE` (`member_id` ASC, `mp_tag_id` ASC))
  ENGINE = InnoDB;


-- name: create-table-wf_issue_category
CREATE TABLE IF NOT EXISTS `%s`.`wf_issue_category` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created at',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Rrecord last updated.',
  `name` VARCHAR(45) NOT NULL COMMENT 'Name of group, must be short as it is used to dynamically generate tabs in the application.',
  `description` VARCHAR(255) NOT NULL COMMENT 'Description of the issue group.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Grouping for issue types for display and reporting purposes. this is NOT managed by the users and is setup in the initial rollout.';


-- name: create-table-country
CREATE TABLE IF NOT EXISTS `%s`.`country` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `membership` TINYINT NOT NULL COMMENT 'Defines if the country is a membership country for the organisation. Used for organisations that span multiple countries and have some kind of administrative division or jurisdiction based on country. For example may be used to determine the appropriate tax names and rates for membership subscriptions.',
  `display_priority` INT NOT NULL COMMENT 'Used to control the order for displaying the countries, higher is first, then alphabetical.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `code` VARCHAR(10) NOT NULL COMMENT 'Country code, eg AU, NZ.',
  `name` VARCHAR(100) NOT NULL COMMENT 'Full Country name',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores a list of countries for use in the system.';


-- name: create-table-ol_resource_filetype
CREATE TABLE IF NOT EXISTS `%s`.`ol_resource_filetype` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_resource_type_id` INT NOT NULL COMMENT 'A more specific sub category of type that explains the file type specifically, e.g. jpg, pdf, mpeg etc.',
  `active` VARCHAR(45) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name` VARCHAR(45) NOT NULL COMMENT 'Resource type name',
  `description` VARCHAR(255) NOT NULL COMMENT 'Resource type description',
  `extension` VARCHAR(4) NOT NULL COMMENT 'Accepted file extension within the system.',
  `other_extensions` TEXT NOT NULL COMMENT 'A comma separated list of other possible extensions that may be used to identify the file.',
  `mime_type` VARCHAR(45) NOT NULL COMMENT 'Mime type.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'NOT IN USE - intended to specifically define file types so we can handle them in a specific way when the user want to view the resource file.';


-- name: create-table-ol_m_category
CREATE TABLE IF NOT EXISTS `%s`.`ol_m_category` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'Link to member record.',
  `ol_category_id` INT NOT NULL COMMENT 'Link to category record',
  `active` TINYINT UNSIGNED NOT NULL COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `member_category_UNIQUE` (`member_id` ASC, `ol_category_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'Defines the learning categories of interest to the member.';


-- name: create-table-ol_m_module_slide
CREATE TABLE IF NOT EXISTS `%s`.`ol_m_module_slide` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_m_module_id` INT NOT NULL COMMENT 'The member\'s instance of the module being undertaken.',
  `ol_slide_id` INT NOT NULL COMMENT 'The question being referenced ',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `elapsed_time_seconds` MEDIUMINT NOT NULL DEFAULT 0 COMMENT 'Incremental time count in seconds that user is on this question page.',
  `score` TINYINT NULL COMMENT 'The score the member has received for their answers to the question (slide). Could be calculated but stored here to make life easier. Info slides will have a NULL value here.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'A copy of the question records attempted by a member for a module. Stored here because each question may have more than one answer given.';


-- name: create-table-ol_module_cpd
CREATE TABLE IF NOT EXISTS `%s`.`ol_module_cpd` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier.',
  `ol_module_id` INT NOT NULL COMMENT 'The module to which this CPD data applies. NOTE this is unique because we have a one-to-one relationship.',
  `ce_activity_id` INT NOT NULL COMMENT 'Links to the relevant CPD activity for which we will record the CPD points for completion of the module.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created.',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated.',
  `allocate_points` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to specify if CPD points are applicable for the module. Overrides all other settings.  Set to FALSE by default so must be explicitly set for a module.',
  `allocate_on_pass` TINYINT NOT NULL DEFAULT 0 COMMENT 'Only allocate CPD points if the  user attains the required pass percentage (ol_module.pass_percentage).',
  `allocate_instance_limit` SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'The number of times (instances) CPD points will be allocated for this module, per allocate_instance_limit_reset_days.',
  `allocate_instance_limit_reset_days` SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'The days after which the CPD allocation instance limit will be reset. That is, after this many days the user is able to repeat the module and gain points for the number of times specified in allocate_instance_limit.',
  `cpd_index` DECIMAL(5,2) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'This is an index that represents the relative weight of each module for the purposes of awarding cpd activity points. This is analogous to the QUANTITY specified by the user (e.g.number of hours) that is then multiplied by the ce_activity.points_per_unit value to calculate the final points recorded for the completion of the module.',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `ol_module_id_UNIQUE` (`ol_module_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'Stores all CPD-related information about a module. Even though this is a one-to-one relationship it has been separated out to simplify the management of CPD points allocation for learning modules. If there exists in this table the appropriate data for allocating CPD points then the system will do so, otherwise CPD points will be ignored.';


-- name: create-table-mp_speciality
CREATE TABLE IF NOT EXISTS `%s`.`mp_speciality` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `name` VARCHAR(255) NOT NULL COMMENT 'Name of the speciality or professional area',
  `description` TEXT NOT NULL COMMENT 'Description of the speciality or professional area',
  `display_priority` SMALLINT NOT NULL DEFAULT 0 COMMENT 'An optional value that will allow records to be ordered when displayed. Higher priority is displayed first.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Professional areas or specialities for members. A member may have one or more of these.';


-- name: create-table-mp_m_speciality
CREATE TABLE IF NOT EXISTS `%s`.`mp_m_speciality` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `member_id` INT NOT NULL COMMENT 'The member who has the speciality',
  `mp_speciality_id` INT UNSIGNED NOT NULL COMMENT 'The speciality that the member has',
  `active` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `preference` TINYINT NULL,
  `start_on` DATE NULL COMMENT 'If relevant, the date that the member started practising this speciality.',
  `comment` TEXT NULL COMMENT 'Option comment about the member and this speciality.',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `member_speciality_id_UNIQUE` (`member_id` ASC, `mp_speciality_id` ASC))
  ENGINE = InnoDB
  COMMENT = '	';


-- name: create-table-ms_permission
CREATE TABLE IF NOT EXISTS `%s`.`ms_permission` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name` VARCHAR(45) NOT NULL COMMENT 'Name of the permission.',
  `description` VARCHAR(255) NOT NULL COMMENT 'Description of the permission.',
  `default_value` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Default value for this permission when adding to mp_m_permission table.',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC))
  ENGINE = InnoDB;


-- name: create-table-ms_m_permission
CREATE TABLE IF NOT EXISTS `%s`.`ms_m_permission` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier.',
  `active` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `member_id` INT NOT NULL,
  `ms_permission_id` INT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`))
  ENGINE = InnoDB;


-- name: create-table-ce_activity_category
CREATE TABLE IF NOT EXISTS `%s`.`ce_activity_category` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name` VARCHAR(255) NOT NULL COMMENT 'The name of the category / group.',
  `description` TEXT NOT NULL COMMENT 'Option al description of the activity category /  group.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Used to store CPD activity categories or groups for collecting related CPD activity types and showing grouped lists etc.';


-- name: create-table-fs_set
CREATE TABLE IF NOT EXISTS `%s`.`fs_set` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL COMMENT 'Soft delete',
  `current` TINYINT NOT NULL COMMENT 'This is the current volume / set for this component. Tells the system to USE this set.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `storage_type` ENUM('LOCAL','AWS-S3','RS-CF') NOT NULL COMMENT 'This determines the type of uploader we will use to transfer the files to their primary storage location. Should be enumerated, e.g. AWS-S3 (Amazon S3), RS-CF (Rackspace Cloud Files).',
  `storage_credentials` TEXT NOT NULL COMMENT 'A serialised set of access credentials - these will be specific to the storage type.',
  `volume_name` VARCHAR(100) NOT NULL COMMENT 'The name of the volume (or container / bucket in cloud environment) e.g. CSANZ-MAPP-R00 . ',
  `set_path` VARCHAR(255) NOT NULL COMMENT 'This is the base path, relative to the volume (container / bucket) below which the dynamically created objects will be referenced. For example an S3 bucket (volume) called CSANZ-MAPP-R00 is used to store a set of Resource files, we can add these resource files into a directory (or simulated directory for cloud storage) called /resource/ to make things easier to locate. From here the rest of the path to the file will be determined by our upload handler.',
  `entity_name` VARCHAR(100) NOT NULL COMMENT 'This defines the name of the table (of file records) for which this storage set is intended. For example, if a set record is going to be used for note attachments then this value would be wf_attachment, for resource file storage it would be ol_resource_file. ',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines file sets which are separate volumes and locations for files associated with components of the application. This table defines file sets that are accessible via the local file system.';


-- name: create-table-fs_url
CREATE TABLE IF NOT EXISTS `%s`.`fs_url` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `fs_set_id` INT NOT NULL COMMENT 'The set for which this URL dispatch method is defined.',
  `active` TINYINT NOT NULL COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `priority` SMALLINT NOT NULL COMMENT 'The priority of this URL for serving files in a set, highest priority is tested first.',
  `base_url` VARCHAR(255) NOT NULL COMMENT 'Defines the base url for accessing the file, the rest of the information is gathered from the set record. This will generally contain a CDN base url, or CNAME for the same. This url will correspond to a cloud files container e.g. MAPP-CSANZ-R00, and is appended with the relevant details from the fs_set record to make a path (actually an object) name.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines file urls which may (optionally) be used for downloading files from a CDN, cloud store or similar - that is, accessing files by URLs.';


-- name: create-table-ol_resource_attribute
CREATE TABLE IF NOT EXISTS `%s`.`ol_resource_attribute` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT(4) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `order` INT(4) NOT NULL COMMENT 'Use to retrieve / display attributes in a given order. ',
  `name` VARCHAR(45) NOT NULL COMMENT 'The name (label) of the attribute - this is a code-friendly name.',
  `display_name` VARCHAR(100) NOT NULL COMMENT 'The interface display name for the attribute - can be longer and more descriptive.',
  `description` VARCHAR(255) NULL COMMENT 'Explanation of the attribute which can be used as an input note as well',
  `placeholder` VARCHAR(45) NULL COMMENT 'Placeholder text to for the form input to guide the user as to the required format.',
  `category` ENUM('citation') NULL COMMENT 'A category or group name that is used to combine related attributes.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  AUTO_INCREMENT = 4
  COMMENT = 'Defines attribute names that can be assigned to resources as required.';


-- name: create-table-ol_resource_attribute_value
CREATE TABLE IF NOT EXISTS `%s`.`ol_resource_attribute_value` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_resource_id` INT NOT NULL COMMENT 'The resource that has the value',
  `ol_resource_attribute_id` INT NOT NULL COMMENT 'The attribute assigned to the resource.',
  `active` TINYINT(4) NOT NULL DEFAULT '1' COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `value` TEXT NULL DEFAULT NULL COMMENT 'The value of the attribute in question.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  AUTO_INCREMENT = 4
  COMMENT = 'Defines the values of attributes assigned to a resource. ';


-- name: create-table-ol_resource_file
CREATE TABLE IF NOT EXISTS `%s`.`ol_resource_file` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ol_resource_id` INT NOT NULL,
  `ad_user_id` INT NOT NULL COMMENT 'The admin user who added the file.',
  `fs_set_id` INT NOT NULL COMMENT 'Tells us in which set we will find file system information that will allow us to locate the file.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `thumbnail` TINYINT NOT NULL DEFAULT 0 COMMENT 'A flag to specify that this file resource is a thumbnail image.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `clean_filename` VARCHAR(255) NOT NULL COMMENT 'The filename of the original file cleaned up and stored for reference.',
  `cloudy_filename` VARCHAR(255) NOT NULL COMMENT 'Obfuscated filename used to store the files so that the url does not give away anything about the resource. Only relevant to a learning environment so that the name does not give any clues about the answer1',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores information about the actual resource files that have been uploaded and managed within the application. ';


-- name: create-table-fn_subscription
CREATE TABLE IF NOT EXISTS `%s`.`fn_subscription` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `fn_subscription_type_id` INT(11) NOT NULL,
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `recurrence_months` TINYINT(4) NOT NULL COMMENT 'Defines the recurrence interval in months. This value is added to the current renew_on value once the invoice has been generated, so that the invoice will be regenerated at a date defined by this interval.',
  `name` VARCHAR(45) NOT NULL COMMENT 'A short name / desc for the subscription template',
  `description` VARCHAR(255) NOT NULL COMMENT 'A longer description if required.',
  `prorata_threshold_days` TINYINT NOT NULL DEFAULT 0 COMMENT 'If pro-rata renewal of subscription is within this many days from a full renewal, then the system will opt for the full renewal instead of the pro-rata. NOTE: the full renewal billing period start is still defined by fn_m_subscription.renew_on date. So effectively it will just be an early full renewal and works exactly the same as the renewal_offset_days value. ',
  `renewal_offset_days` TINYINT NOT NULL DEFAULT 0 COMMENT 'Allows the processing of the subscription renewal before or after the fn_m_subscription.renew_on date. A negative value will force the renewal early, and a positive renewal will be x days after renew_on. The actual renewal process is the same as if the pro-rata renewal falls within the threshold days to a full renewal. That is, the Invoice is generated NOW and has current date as invoice date, but the billing period start is still defined by renew_on date.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines Subscription templates used to generate renewal invoices.';


-- name: create-table-fn_subscription_inventory
CREATE TABLE IF NOT EXISTS `%s`.`fn_subscription_inventory` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `fn_subscription_id` INT NOT NULL COMMENT 'The subscription to which this line item belongs',
  `fn_inventory_id` INT NOT NULL COMMENT 'The inventory item on this subscription.',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `description` TEXT NOT NULL COMMENT 'Description of the inventory line item, can be edited.',
  `quantity` DECIMAL(8,2) NOT NULL COMMENT 'The quantity of the line item.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines the inventory items that appear on subscription rene';


-- name: create-table-fn_subscription_type
CREATE TABLE IF NOT EXISTS `%s`.`fn_subscription_type` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT(4) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `limit_per_member` TINYINT(4) NOT NULL COMMENT 'This subscription type is limited to this number per member (0 is no limit).',
  `name` VARCHAR(45) NOT NULL COMMENT 'Name of the subscription type.',
  `description` VARCHAR(255) NOT NULL COMMENT 'Description of subscription type.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB;


-- name: create-table-fn_tax
CREATE TABLE IF NOT EXISTS `%s`.`fn_tax` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `country_id` INT NOT NULL COMMENT 'The country for which this tax applies.',
  `active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated.',
  `name` VARCHAR(45) NOT NULL COMMENT 'Name of the tax.',
  `description` VARCHAR(255) NOT NULL COMMENT 'Description of tax.',
  `rate` DECIMAL(5,2) NOT NULL COMMENT 'The rate of the tax as a percentage.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Tax table that defines taxes for each country.';


-- name: create-table-log_data_field
CREATE TABLE IF NOT EXISTS `%s`.`log_data_field` (
  `id` INT(10) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `log_data_action_id` INT(10) NOT NULL,
  `active` TINYINT NOT NULL DEFAULT 1,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created at - ie time and date that the action was performed.',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `field_name` VARCHAR(20) NOT NULL COMMENT 'Defines if the user was \'admin\' or \'member\'.',
  `value_before` TEXT NOT NULL COMMENT 'Describes the action that was taken by the user.',
  `value_after` TEXT NOT NULL COMMENT 'Further details about the event.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Logs various system events and data changes over time.';


-- name: create-table-log_data_table
CREATE TABLE IF NOT EXISTS `%s`.`log_data_table` (
  `id` INT(10) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created at - ie time and date that the action was performed.',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `table_name` VARCHAR(20) NOT NULL COMMENT 'Defines if the user was \'admin\' or \'member\'.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines the table names for which we will log data changes. These values need to be added as part of setup.';


-- name: create-table-a_name_prefix
CREATE TABLE IF NOT EXISTS `%s`.`a_name_prefix` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated.',
  `name` VARCHAR(45) NOT NULL COMMENT 'The short name or abbreviation for the name prefix - e.g. Dr, Mr, Mrs, Professor etc. as we want it to appear before the first name.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Name prefix is for title or honorific - e.g. Dr, Mr, Professor etc.';


-- name: create-table-a_meeting_type
CREATE TABLE IF NOT EXISTS `%s`.`a_meeting_type` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name` VARCHAR(255) NOT NULL COMMENT 'Name or descriptive title of the TYPE of meeting - e.g. AGM, Quarterly etc etc.',
  `description` TEXT NOT NULL COMMENT 'Description of type of meeting.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Meeting types are used to group meetings';


-- name: create-table-ms_m_application_meeting
CREATE TABLE IF NOT EXISTS `%s`.`ms_m_application_meeting` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ms_m_application_id` INT NOT NULL COMMENT 'The application record that will be reviewed at the meeting.',
  `ad_macro_id` INT NULL COMMENT 'If this record was inserted or updated by a macro we store the macro id here.',
  `a_meeting_id` INT NOT NULL COMMENT 'The meeting at which the application will be reviewed.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `result` TINYINT NOT NULL DEFAULT -1 COMMENT 'The outcome of the application AT THIS MEETING - value of -1 is the default case in which the outcome of the application is unknown or PENDING. An explicit 0 is REJECTED and an explicit 1 is ACCEPTED. The final outcome of the application is stored in the ms_m_application table in the result column.',
  `comment` TEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `ms_m_application_meeting_UNIQUE` (`ms_m_application_id` ASC, `a_meeting_id` ASC))
  ENGINE = InnoDB
  COMMENT = 'One or more associated meetings relevant to the processing of a member application. Provides a workflow or review or approval processes.';


-- name: create-table-ad_macro_transaction
CREATE TABLE IF NOT EXISTS `%s`.`ad_macro_transaction` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ad_macro_id` INT NOT NULL COMMENT 'The macro to which this transaction belongs.',
  `member_id` INT NULL COMMENT 'The member against whose record the transaction was taken. This is set to NULL in the event that the member id was malformed or not found.',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT 'Status indicates success (1) or fail (0 - default) of this transaction.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created date - the date and time the macro was run',
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `log` TEXT NULL COMMENT 'Details of the transaction.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores  records of the individual transactions / events in a macro.';


-- name: create-table-acl_admin_role
CREATE TABLE IF NOT EXISTS `%s`.`acl_admin_role` (
  `id` INT(11) NOT NULL COMMENT 'Unique identifier',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT 'Specified that this role is the default assigned to admin user',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `name` VARCHAR(255) NOT NULL COMMENT 'Name of role (group)',
  `description` TEXT NULL DEFAULT NULL COMMENT 'Description of role (group)',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Defines roles of groups into which a member is assigned. The member with then inherit all of the permissions associated with this role.';


-- name: create-table-acl_admin_resource
CREATE TABLE IF NOT EXISTS `%s`.`acl_admin_resource` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `parent_id` INT(11) NULL COMMENT 'Self-referencing id to indicate parent record.',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `type` VARCHAR(45) NOT NULL COMMENT 'Type of resource: page or function - page represents access to a top-level component whereas function represents the availability (globally) of that function - this can then be fine tuned in acl_admin_role_resource.',
  `name` VARCHAR(255) NOT NULL COMMENT 'The name of the resource',
  `description` TEXT NULL DEFAULT NULL COMMENT 'A description of the resource',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  AUTO_INCREMENT = 10000
  COMMENT = 'Defined ACL for resources in the application.';


-- name: create-table-acl_admin_role_resource
CREATE TABLE IF NOT EXISTS `%s`.`acl_admin_role_resource` (
  `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `acl_admin_role_id` INT(11) NOT NULL COMMENT 'The role that will be allowed to access the resource',
  `acl_admin_resource_id` INT(11) NOT NULL COMMENT 'The resource that can be accessed by the role',
  `active` TINYINT(1) NOT NULL DEFAULT '1' COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `privileges` VARCHAR(100) NULL DEFAULT NULL COMMENT 'Privileges (not in use yet)',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `acl_admin_role_resource_id_UNIQUE` (`acl_admin_role_id` ASC, `acl_admin_resource_id` ASC))
  ENGINE = InnoDB
  AUTO_INCREMENT = 9;


-- name: create-table-cm_email_variable
CREATE TABLE IF NOT EXISTS `%s`.`cm_email_variable` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `persistent` TINYINT NOT NULL DEFAULT 0 COMMENT 'If set to 1 this variable will always be available in all email templates.',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `name_strict` VARCHAR(100) NOT NULL COMMENT 'Unique variable-style name e.g. member_fullname, must be unique.',
  `name_friendly` VARCHAR(100) NOT NULL COMMENT 'A friendlier short name for the variable eg Member Full Name.',
  `description` TEXT NOT NULL COMMENT 'A short description of the variable and it\'s use.',
  `fallback` TEXT NULL COMMENT 'The fallback value is used if the variable value is not available, missing or empty.',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_strict_UNIQUE` (`name_strict` ASC))
  ENGINE = InnoDB
  COMMENT = 'Defines the variables available for use in email templates.';


-- name: create-table-cm_email_log
CREATE TABLE IF NOT EXISTS `%s`.`cm_email_log` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `cm_email_id` INT NOT NULL COMMENT 'The parent email record that was sent - i.e. the email that was the source for the broadcast. This is redundant as it can be looked up using cm_m_email_id however is included for implementation ease. This value is initially sent the the mail exchanger as an additional email header and is then posted back from the remote system to reconcile email events. ',
  `cm_m_email_id` INT NOT NULL COMMENT 'The email that the logged event relates to. This value is initially sent the the mail exchanger as an additional email header and is then posted back from the remote system to reconcile email events. ',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record updated',
  `mx_at` TIMESTAMP NULL COMMENT 'The date and time that the status was recorded at the mail exchanger.',
  `mx_status` VARCHAR(100) NULL COMMENT 'The status / event name reported by the mail exchanger - extracted from the raw data.',
  `mx_description` TEXT NULL COMMENT 'The status / event description - extracted from the raw data.',
  `mx_data` TEXT NULL COMMENT 'The serialised raw data posted back / retrieved from the mail exchanger.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Stores information about the status of emails that are relayed via an external mail exchanger. This data is posted back to our application by web hooks or gather via the remote system API, and then written to this table. It is used to keep a history of the email transactions and to populate the cm_m_email table with the current status of each email.';


-- name: create-table-session_admin
CREATE TABLE IF NOT EXISTS `%s`.`session_admin` (
  `id` CHAR(32) NOT NULL,
  `modified` INT(11) NULL DEFAULT NULL,
  `lifetime` INT(11) NULL DEFAULT NULL,
  `data` MEDIUMTEXT NULL,
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Session table - structure required by Zend Framework';


-- name: create-table-session_member
CREATE TABLE IF NOT EXISTS `%s`.`session_member` (
  `id` CHAR(32) NOT NULL,
  `modified` INT(11) NULL DEFAULT NULL,
  `lifetime` INT(11) NULL DEFAULT NULL,
  `data` MEDIUMTEXT NULL,
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Session table - structure required by Zend Framework';


-- name: create-table-ce_m_activity_attachment
CREATE TABLE IF NOT EXISTS `%s`.`ce_m_activity_attachment` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ce_m_activity_id` INT NOT NULL COMMENT 'The member cpd activity with which this file is associated. ',
  `fs_set_id` INT NOT NULL COMMENT 'Tells us in which set we will find file system information that will allow us to locate the file.',
  `active` TINYINT NOT NULL DEFAULT 1 COMMENT 'Soft delete',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'Record last updated',
  `clean_filename` VARCHAR(255) NOT NULL COMMENT 'The filename of the original document prior to upload - stored for defence.',
  `cloudy_filename` VARCHAR(255) NOT NULL COMMENT 'The obfuscated cloud file name.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'Documents table provides a way to attach files to notes. A document is always attached via a note and, therefore, can only be associated with member records. As notes can be associated with applications and issues, documents can also be associated with these entities.';

-- name: create-table-ce_activity_type
CREATE TABLE IF NOT EXISTS `%s`.`ce_activity_type` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Unique identifier',
  `ce_activity_id` INT NOT NULL COMMENT 'The activity to which this type relates.',
  `active` TINYINT NOT NULL DEFAULT 1,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created',
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `name` VARCHAR(255) NOT NULL COMMENT 'Descriptive name for the activity type.',
  PRIMARY KEY (`id`))
  ENGINE = InnoDB
  COMMENT = 'This table was added to allow for prescriptive activity descriptions.';

