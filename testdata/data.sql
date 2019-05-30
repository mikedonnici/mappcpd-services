-- name: insert-data-a_meeting
INSERT INTO `%s`.`a_meeting` VALUES
  (1, 3, 1, NOW(), NOW(), '2015-09-11', 'Office', 'Applications - 11 September 2015', 'Considering applications...'),
  (1189, 3, 1, NOW(), NOW(), '2016-04-11', 'At the office', 'Sub - 11 April 2016', 'Review applications'),
  (1190, 3, 1, NOW(), NOW(), '2016-11-30', 'Somewhere', 'Membership Application Review', NULL);

-- name: insert-data-a_meeting_type
INSERT INTO `%s`.`a_meeting_type` VALUES
  (1, 1, '2015-08-30 17:10:09', '2015-08-30 17:10:09', 'AGM', 'Annual General Meeting'),
  (2, 1, '2015-08-30 17:10:09', '2015-08-30 17:10:09', 'Board', 'Board Meeting'),
  (3, 1, '2015-08-30 17:10:09', '2016-11-22 23:43:13', 'Sub Committee', 'Sub Committee');

-- name: insert-data-a_name_prefix
INSERT INTO `%s`.`a_name_prefix` VALUES
  (1, 1, '2013-06-03 17:29:45', NOW(), 'A/Prof'),
  (2, 1, '2013-06-03 17:29:46', NOW(), 'A/Prof Dame'),
  (3, 1, '2013-06-03 17:29:47', NOW(), 'Brig'),
  (4, 1, '2013-06-03 17:29:47', NOW(), 'Dame'),
  (5, 1, '2013-06-03 17:29:48', NOW(), 'Dr'),
  (6, 1, '2013-06-03 17:29:48', NOW(), 'Miss'),
  (7, 1, '2013-06-03 17:29:49', NOW(), 'Mr'),
  (8, 1, '2013-06-03 17:29:49', NOW(), 'Mrs'),
  (9, 1, '2013-06-03 17:29:50', NOW(), 'Ms'),
  (10, 1, '2013-06-03 17:29:51', NOW(), 'Professor'),
  (11, 1, '2013-06-03 17:29:51', NOW(), 'Professor Sir'),
  (12, 1, '2013-06-03 17:29:52', NOW(), 'Sir'),
  (13, 1, '2013-06-03 17:29:52', NOW(), 'Sister');

-- insert-data-acl_admin_resource

-- insert-data-acl_admin_role

-- insert-data-acl_admin_role_resource

-- insert-data-acl_member_resource

-- insert-data-acl_member_role

-- insert-data-acl_member_role_resource

-- insert-data-ad_macro

-- insert-data-ad_macro_transaction

-- insert-data-ad_permission

-- name: insert-data-ad_user
INSERT INTO `%s`.`ad_user` VALUES
  (1, 1, 1, 0, '2015-08-30 17:10:08', '2016-05-31 04:33:03', 'demo-admin', '41d0510a9067999b72f38ba0ce9f6195',
      'Demo Admin', 'demo', 'demo@noemail.com');

-- insert-data-ad_user_permission

-- name: insert-data-ce_activity
INSERT INTO `%s`.`ce_activity` VALUES
  (1, 1, 1, 0, 0, NOW(), NOW(), 'CE1', 'Conference session / workshop / course', '', 1.00, 50),
  (2, 1, 1, 0, 0, NOW(), NOW(), 'CE2', 'Reading, research, literature review', '', 1.00, 25),
  (3, 1, 1, 0, 0, NOW(), NOW(), 'CE3', 'Teaching - preperation and delivery', '', 1.00, 25),
  (4, 1, 1, 0, 0, NOW(), NOW(), 'CE4', 'Presentation', '', 1.00, 25),
  (5, 1, 1, 0, 0, NOW(), NOW(), 'CE5', 'Online content - other', '', 1.00, 25),
  (6, 3, 1, 0, 1, NOW(), NOW(), 'CE6', 'MappCPD online module', '', 1.00, 50),
  (7, 1, 2, 0, 0, NOW(), NOW(), 'PR1', 'Formal performance review or audit', '', 1.00, 20),
  (8, 1, 2, 0, 0, NOW(), NOW(), 'PR2', 'Informal peer review / meeting', '', 1.00, 20),
  (9, 1, 3, 0, 0, NOW(), NOW(), 'PQ1', 'Personal / professional / management course', '', 1.00, 20),
  (10, 1, 3, 0, 0, NOW(), NOW(), 'PQ2', 'Self-directed learning (professional qualities)', '', 1.00, 20),
  (20, 1, 10, 1, 0, NOW(), NOW(), 'RACP1', 'Practice Review & Improvement', '', 3.00, 50),
  (21, 1, 10, 1, 0, NOW(), NOW(), 'RACP2', 'Assessed Learning', '', 2.00, 50),
  (22, 1, 10, 1, 0, NOW(), NOW(), 'RACP3', 'Educational Development, Teaching & Research', '', 1.00, 50),
  (23, 1, 10, 1, 0, NOW(), NOW(), 'RACP4', 'Group Learning', '', 1.00, 50),
  (24, 1, 10, 1, 0, NOW(), NOW(), 'RACP5', 'Other Learning Activities', '', 1.00, 50);

-- name: insert-data-ce_activity_category
INSERT INTO `%s`.`ce_activity_category` VALUES
  (1, 1, '2015-08-30 17:10:10', '2015-08-30 17:10:10', 'Continuing Education',
   'Profesional knowledge and skills aquisition'),
  (2, 1, '2015-08-30 17:10:10', '2015-08-30 17:10:10', 'Performance Review',
   'Review of performance indicators, peer review, professional audit'),
  (3, 1, '2015-08-30 17:10:10', '2015-08-30 17:10:10', 'Professional Qualities', 'Personal skills, management, ethics'),
  (10, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'RACP', 'RACP CPD categories.');

-- name: insert-data-ce_activity_type
INSERT INTO `%s`.`ce_activity_type` VALUES
  (1, 20, 1, NOW(), NOW(), 'Practice audits/Clinical audits'),
  (2, 20, 1, NOW(), NOW(), 'Peer review'),
  (3, 20, 1, NOW(), NOW(), 'Patient satisfaction studies'),
  (4, 20, 1, NOW(), NOW(), 'Institution audits, e.g. hospital accreditation'),
  (5, 20, 1, NOW(), NOW(), 'Incident reporting/monitoring, e.g. morbidity & mortality meetings'),
  (6, 20, 1, NOW(), NOW(), 'Practice Review, e.g. Regular Practice Review'),
  (7, 20, 1, NOW(), NOW(), 'Multi Source Feedback (MSF)'),
  (8, 20, 1, NOW(), NOW(), 'Participation in the RACP Supervisor Professional Development Program (SPDP)'),
  (9, 20, 1, NOW(), NOW(), 'Other practice review & improvement activities'),
  (10, 21, 1, NOW(), NOW(), 'PhD studies'),
  (11, 21, 1, NOW(), NOW(), 'Formal postgraduate studies'),
  (12, 21, 1, NOW(), NOW(), 'Self-assessment programs'),
  (13, 21, 1, NOW(), NOW(), 'Courses to learn new techniques, e.g. Advanced Life Support (ALS)'),
  (14, 21, 1, NOW(), NOW(), 'Learner initiated and planned projects'),
  (15, 21, 1, NOW(), NOW(), 'Other assessed learning activities'),
  (16, 22, 1, NOW(), NOW(), 'Teaching, e.g. supervision, mentoring'),
  (17, 22, 1, NOW(), NOW(), 'Involvement in standards development'),
  (18, 22, 1, NOW(), NOW(), 'Reviewer'),
  (19, 22, 1, NOW(), NOW(), 'Writing examination questions'),
  (20, 22, 1, NOW(), NOW(), 'Examining'),
  (21, 22, 1, NOW(), NOW(), 'Publication (including preparation)'),
  (22, 22, 1, NOW(), NOW(), 'Presentation (including preparation)'),
  (23, 22, 1, NOW(), NOW(), 'Committee/working group/council involvement'),
  (24, 22, 1, NOW(), NOW(), 'Other educational development, teaching & research activities'),
  (25, 23, 1, NOW(), NOW(), 'Seminars'),
  (26, 23, 1, NOW(), NOW(), 'Conferences'),
  (27, 23, 1, NOW(), NOW(), 'Workshops'),
  (28, 23, 1, NOW(), NOW(), 'Grand rounds'),
  (29, 23, 1, NOW(), NOW(), 'Journal clubs'),
  (30, 23, 1, NOW(), NOW(), 'Hospital and other medical meetings'),
  (31, 23, 1, NOW(), NOW(), 'Other group learning activities'),
  (32, 24, 1, NOW(), NOW(), 'Reading journals and texts'),
  (33, 24, 1, NOW(), NOW(), 'Information searches, e.g. Medline'),
  (34, 24, 1, NOW(), NOW(), 'Audio/videotapes'),
  (35, 24, 1, NOW(), NOW(), 'Web-based learning'),
  (36, 24, 1, NOW(), NOW(), 'Other learning activities');

-- name: insert-data-ce_activity_unit
INSERT INTO `%s`.`ce_activity_unit` VALUES (1, 1, NOW(), NOW(), 1, 'hours', NULL),
  (2, 1, NOW(), NOW(), 0, 'item', 'Item, instance, single event.'),
  (3, 1, NOW(), NOW(), 0, 'module', 'Online learning module.');

-- insert-data-ce_audit

-- insert-data-ce_audit_m_activity

-- name: insert-data-ce_evaluation
INSERT INTO `%s`.`ce_evaluation` VALUES
  (1, 1, '2015-08-30 17:10:13', '2015-08-30 17:10:13', 0, 1, 1, 31, 12, 12, 100, 'Annual CPD', '12 month CPD period'),
  (2, 1, '2015-08-30 17:10:13', '2015-08-30 17:10:13', 0, 1, 1, 31, 12, 36, 250, 'CPD Triennium',
   '36 month CPD period');

-- insert-data-ce_event

-- name: insert-data-ce_m_activity
INSERT INTO `%s`.`ce_m_activity` VALUES
  (1, 1, 23, 25, NULL, NULL, 1, 1, NOW(), NOW(), '2018-02-03', 1.00, 1.00, 0, 'BJJ like Bruno Malfacine'),
  (2, 1, 23, 25, NULL, NULL, 1, 0, NOW(), NOW(), '2018-02-04', 1.00, 1.00, 0, 'Ate sausages and eggs'),
  (3, 1, 20, 1, NULL, NULL, 1, 0, NOW(), NOW(), '2018-02-05', 1.00, 3.00, 0, 'Baked bread');

-- name: insert-data-ce_m_activity_attachment
INSERT INTO `%s`.`ce_m_activity_attachment` VALUES
  (74, 247, 4, 1, '2018-02-06 00:04:35', '2018-02-06 00:04:35', 'Derwent1.png', '04abe653d926a3ccb122245671e6c064.png'),
  (75, 247, 4, 1, '2018-02-06 00:05:17', '2018-02-06 00:05:17', 'headerbg.jpg', '5886e6ab4b3b7b71ad112e56ef65ed66.jpg'),
  (77, 250, 4, 1, '2018-02-06 01:18:08', '2018-02-06 01:18:08', 'Derwent1.png', '04abe653d926a3ccb122245671e6c064.png'),
  (78, 249, 4, 1, '2018-02-06 02:55:29', '2018-02-06 02:55:29', 'Derwent1.png', '04abe653d926a3ccb122245671e6c064.png');

-- name: insert-data-ce_m_evaluation
INSERT INTO `%s`.`ce_m_evaluation` VALUES
  (1, 501, 1, 1, 1, '2015-08-31 04:27:24', '2017-01-02 09:05:50', 100, '2015-01-01', '2015-12-31', ''),
  (2, 503, 1, 1, 1, '2015-09-02 06:29:37', '2017-01-02 09:05:50', 100, '2015-01-01', '2015-12-31', ''),
  (3, 504, 1, 1, 1, '2015-09-28 03:56:53', '2017-01-02 09:05:50', 100, '2015-01-01', '2016-01-01', ''),
  (4, 505, 1, 1, 1, '2015-10-26 05:52:37', '2017-01-02 09:05:50', 100, '2015-01-01', '2015-12-01', ''),
  (5, 35, 1, 1, 1, '2015-10-29 03:57:18', '2017-01-02 09:05:35', 100, '2015-01-01', '2015-12-01', ''),
  (6, 506, 1, 1, 1, '2016-02-24 02:33:29', '2017-01-02 09:05:51', 100, '2016-01-01', '2017-01-01', '');

-- insert-data-cm_email_log

-- name: insert-data-cm_email_template
INSERT INTO `%s`.`cm_email_template` VALUES
  (1, 1, 1, '2015-08-30 17:10:15', '2015-08-30 17:10:15', 'Invoice 1st / Generic Notice',
      'Sent when a new invoice is generated, or whenever an admin user resends the invoice to the member via the \'Send Invoice to Member\' link',
      'MappCPD Subscription Invoice {{invoice_id}}', 'MappCPD', 'noreply@mappcpd.com',
      '', '1001,1002');

-- name: insert-data-cm_email_variable
INSERT INTO `%s`.`cm_email_variable` VALUES
  (1, 1, 1, '2015-08-30 17:10:16', '2015-08-30 17:10:16', 'member_id', 'Member\'s ID', 'Membership ID', 'MEMBER ID'),
  (2, 1, 1, '2015-08-30 17:10:16', '2015-08-30 17:10:16', 'member_prefix', 'Member\'s Name Prefix', 'Prefix of a member. eg: Mr, Dr etc', NULL),
  (3, 1, 1, '2015-08-30 17:10:16', '2015-08-30 17:10:16', 'member_firstname', 'Member\'s First Name', 'Member\'s First Name', 'MEMBER NAME'),
  (4, 1, 1, '2015-08-30 17:10:16', '2015-08-30 17:10:16', 'member_middlenames', 'Member\'s Second Name', 'Member\'s Second Name', NULL),
  (5, 1, 1, '2015-08-30 17:10:17', '2015-08-30 17:10:17', 'member_lastname', 'Member\'s Last Name', 'Member\'s Last Name', NULL),
  (6, 1, 1, '2015-08-30 17:10:17', '2015-08-30 17:10:17', 'member_suffix', 'Member\'s Name Suffix', 'Member\'s Name Suffix', NULL),
  (7, 1, 1, '2015-08-30 17:10:17', '2015-08-30 17:10:17', 'member_gender', 'Member\'s Gender', 'Member\'s Gender', NULL),
  (8, 1, 1, '2015-08-30 17:10:17', '2015-08-30 17:10:17', 'member_date_of_entry', 'Member\'s Date of Entry', 'Member\'s Date of Entry', NULL),
  (9, 1, 1, '2015-08-30 17:10:17', '2015-08-30 17:10:17', 'member_primary_email', 'Member\'s primary email', 'Member\'s primary email', NULL),
  (10, 1, 1, '2015-08-30 17:10:17', '2015-08-30 17:10:17', 'member_secondary_email', 'Member\'s secondary email', 'Member\'s secondary email', NULL),
  (11, 1, 1, '2015-08-30 17:10:18', '2015-08-30 17:10:18', 'member_country', 'Membership country', 'Member\'s membership country', 'MEMBER COUNTRY'),
  (12, 1, 1, '2015-08-30 17:10:18', '2015-08-30 17:10:18', 'member_mobile', 'Member\'s mobile number', 'Member\'s mobile number', NULL),
  (13, 1, 1, '2015-08-30 17:10:18', '2015-08-30 17:10:18', 'member_journal_number', 'Member\'s journal number', 'Member\'s journal number', NULL),
  (14, 1, 1, '2015-08-30 17:10:18', '2015-08-30 17:10:18', 'member_bpay_number', 'Member\'s BPAY number', 'Member\'s BPAY number', NULL),
  (16, 1, 1, '2015-08-30 17:10:18', '2015-08-30 17:10:18', 'member_url_login', 'Member log in', 'Outputs the raw URL of the log-in page', 'https://member.demo.mappcpd.com/login'),
  (17, 1, 1, '2015-08-30 17:10:19', '2015-08-30 17:10:19', 'membership_title', 'Member\'s current title', 'Member\'s current membership title', ''),
  (18, 1, 1, '2015-08-30 17:10:19', '2015-08-30 17:10:19', 'member_url_reset_password', 'Member password reset', 'Outputs the raw URL of the password reset page', 'https://member.demo.mappcpd.com/login/index/reset'),
  (1001, 1, 0, '2015-08-30 17:10:19', '2015-08-30 17:10:19', 'invoice_id', 'Invoice ID', 'Invoice ID', NULL),
  (1002, 1, 0, '2015-08-30 17:10:19', '2015-08-30 17:10:19', 'invoice_pay_link', 'Link to invoice payment page', 'Member Application Invoice Payment Link', NULL),
  (1003, 1, 0, '2015-08-30 17:10:19', '2015-08-30 17:10:19', 'password_reset_link', 'Reset password', 'Link to log user on with a token so they can reset their password', NULL),
  (1004, 1, 0, '2015-08-30 17:10:20', '2015-08-30 17:10:20', 'member_pending_issue_list', 'Member Pending Issue List', 'List of issues that are pending for a member to take action on.', NULL);

-- insert-data-cm_m_email

-- name: insert-data-country
INSERT INTO `%s`.`country` VALUES
  (14, 1, 1, 1000, NOW(), NOW(), 'AU', 'AUSTRALIA'),
  (159, 1, 1, 900, NOW(), NOW(), 'NZ', 'NEW ZEALAND'),
  (235, 1, 0, 500, NOW(), NOW(), 'GB', 'UNITED KINGDOM'),
  (236, 1, 0, 400, NOW(), NOW(), 'US', 'UNITED STATES');

-- name: insert-data-fn_inventory
INSERT INTO `%s`.`fn_inventory` VALUES
  (1, 1, '2015-08-30 17:10:20', '2015-08-30 17:10:20', 'Subs - Associate', 'Associate Membership Fee', 300.00, 'year',
   1),
  (2, 1, '2015-08-30 17:10:20', '2015-08-30 17:10:20', 'Subs - Ordinary', 'Ordinary Membership Fee', 600.00, 'year', 1),
  (3, 1, '2015-08-30 17:10:21', '2015-08-30 17:10:21', 'Subs - Fellow', 'Fellow Membership Fee', 900.00, 'year', 1),
  (4, 1, '2015-08-30 17:10:21', '2015-08-30 17:10:21', 'Misc', 'Miscellaneous item', 1.00, 'item', 1);

-- name: insert-data-fn_invoice_inventory
INSERT INTO `%s`.`fn_invoice_inventory` VALUES
  (1, 1, 1, 1, '2015-09-01 06:13:26', '2015-09-01 06:13:26', 'Associate Membership Fee', 0.33, 300.00, 10.00, 'GST'),
  (2, 2, 1, 1, '2015-09-01 06:46:13', '2015-09-01 06:46:13', 'Associate Membership Fee', 0.33, 300.00, 10.00, 'GST'),
  (3, 3, 1, 1, '2015-09-01 06:46:13', '2015-09-01 06:46:13', 'Associate Membership Fee', 0.33, 300.00, 10.00, 'GST'),
  (4, 4, 3, 1, '2015-09-01 06:46:13', '2015-09-01 06:46:13', 'Fellow Membership Fee', 0.33, 900.00, 10.00, 'GST'),
  (5, 5, 3, 1, '2015-09-01 06:46:13', '2015-09-01 06:46:13', 'Fellow Membership Fee', 0.33, 900.00, 10.00, 'GST'),
  (15564, 16577, 1, 1, '2016-04-04 03:26:08', '2016-04-04 03:26:08', 'Associate Membership Fee', 1.00, 300.00, 10.00, 'GST'),
  (15565, 16578, 1, 1, '2016-04-04 03:26:08', '2016-04-04 03:26:08', 'Associate Membership Fee', 1.00, 300.00, 10.00, 'GST'),
  (15566, 16579, 1, 1, '2016-04-04 03:26:08', '2016-04-04 03:26:08', 'Associate Membership Fee', 1.00, 300.00, 10.00, 'GST'),
  (15567, 16580, 3, 1, '2016-04-04 03:26:08', '2016-04-04 03:26:08', 'Fellow Membership Fee', 1.00, 900.00, 10.00, 'GST'),
  (15568, 16581, 3, 1, '2016-04-04 03:26:08', '2016-04-04 03:26:08', 'Fellow Membership Fee', 1.00, 900.00, 10.00, 'GST'),
  (15569, 16582, 1, 1, '2016-04-04 03:26:09', '2016-04-04 03:26:09', 'Associate Membership Fee', 0.75, 300.00, 10.00, 'GST'),
  (15570, 16583, 3, 1, '2016-04-04 03:30:20', '2016-04-04 03:30:20', 'Fellow Membership Fee', 0.75, 900.00, 10.00, 'GST'),
  (15571, 16584, 4, 1, '2017-01-02 04:39:48', '2017-01-02 04:39:48', 'Miscellaneous item', 1.00, 1.00, 10.00, 'GST'),
  (15572, 16585, 1, 1, '2017-01-02 09:05:22', '2017-01-02 09:05:22', 'Associate Membership Fee', 1.00, 300.00, 10.00, 'GST'),
  (15573, 16586, 1, 1, '2017-01-02 09:05:22', '2017-01-02 09:05:22', 'Associate Membership Fee', 1.00, 300.00, 10.00, 'GST'),
  (15574, 16587, 1, 1, '2017-01-02 09:05:22', '2017-01-02 09:05:22', 'Associate Membership Fee', 1.00, 300.00, 10.00, 'GST'),
  (15575, 16588, 1, 1, '2017-01-02 09:05:22', '2017-01-02 09:05:22', 'Associate Membership Fee', 1.00, 300.00, 10.00, 'GST'),
  (15576, 16589, 3, 1, '2017-01-02 09:05:22', '2017-01-02 09:05:22', 'Fellow Membership Fee', 1.00, 900.00, 10.00, 'GST'),
  (15577, 16590, 3, 1, '2017-01-02 09:05:22', '2017-01-02 09:05:22', 'Fellow Membership Fee', 1.00, 900.00, 10.00, 'GST'),
  (15578, 16591, 3, 1, '2017-01-02 09:05:23', '2017-01-02 09:05:23', 'Fellow Membership Fee', 1.00, 900.00, 10.00, 'GST'),
  (15579, 16592, 4, 1, '2017-01-02 21:31:05', '2017-01-02 21:31:05', 'Miscellaneous item', 1.00, 1.00, 10.00, 'GST'),
  (15580, 16593, 1, 1, '2017-01-04 21:30:36', '2017-01-04 21:30:36', 'Associate Membership Fee', 1.00, 300.00, 10.00,
          'GST');

-- name: insert-data-fn_m_invoice
INSERT INTO `%s`.`fn_m_invoice` VALUES
  (1, 1, 1, 1, 1, 0, NOW(), NOW(), NOW(), NOW(),
      '2018-01-01', '2018-01-15', '2018-01-01', '2018-12-31', 110.11, "Subs for 2018"),
  (2, 1, 1, 1, 1, 0, NOW(), NOW(), NOW(), NOW(),
      '2019-01-01', '2019-01-15', '2019-01-01', '2019-12-31', 220.22, "Subs for 2019");

-- name: insert-data-fn_invoice_payment
INSERT INTO `%s`.`fn_invoice_payment` VALUES 
  (1, 1, 1, 1, NOW(), NOW(), 108.95, "Allocation of payment id 1 to invoice id 1");

  -- name: insert-data-fn_payment
INSERT INTO `%s`.`fn_payment` VALUES
  (1, 1, 1, NULL, 1, NOW(), NOW(), '2015-09-01', 108.95, 'Payment for invoice id 1', 'data1 - 123', 'data2 - abc', 'data3 - 123abc', 'data4 - abc123'),
  (2, 4, 1, NULL, 1, NOW(), NOW(), '2016-09-01', 10.10, '', '', '', '', ''),
  (3, 3, 1, NULL, 1, NOW(), NOW(), '2017-09-01', 20.20, '', '', '', '', ''),
  (4, 2, 1, NULL, 1, NOW(), NOW(), '2018-09-01', 20.31, '', '', '', '', ''),
  (5, 1, 2, NULL, 1, NOW(), NOW(), '2018-09-01', 30.32, '', '', '', '', '');

-- name: insert-data-fn_m_subscription
INSERT INTO `%s`.`fn_m_subscription` VALUES
  (1, 1, 1, NULL, 1, 0, NOW(), NOW(), NULL, '2020-01-01', NULL);


-- name: insert-data-fn_payment_type
INSERT INTO `%s`.`fn_payment_type` VALUES
  (1, 1, '2013-07-24 20:08:56', '2013-07-24 20:08:56', 'BPAY', 'BPay - requires BPAy number for member.', '#1', '#2',
   '#3', '#4'), (2, 1, '2013-07-24 20:08:56', '2013-07-24 20:08:56', 'CC', 'Credit Card', '#1', '#2', '#3', '#4'),
  (3, 1, '2013-07-24 20:08:56', '2013-07-24 20:08:56', 'CHQ', 'Cheque', '#1', '#2', '#3', '#4'),
  (4, 1, '2013-07-24 20:08:56', '2013-07-24 20:08:56', 'EFT', 'Electronic Funds Transfer', '#1', '#2', '#3', '#4'),
  (5, 1, '2013-07-24 20:08:56', '2013-07-24 20:08:56', 'DIRECT', 'Direct Debit / Credit', '#1', '#2', '#3', '#4'),
  (6, 1, '2013-07-24 20:08:56', '2013-07-24 20:08:56', 'OTHER', 'Other payment type', '#1', '#2', '#3', '#4'),
  (7, 1, '2013-08-21 18:02:25', '2013-08-21 18:02:25', 'ONLINE AU', 'Online payment - AU Gateway', 'Card Type',
   'CC Number', 'Expiry', 'TRXN'),
  (8, 1, '2013-08-21 18:02:25', '2013-08-21 18:02:25', 'ONLINE NZ', 'Online payment - NZ Gateway', 'Card Type',
   'CC Number', 'Expiry', 'TRXN'),
  (9, 1, '2014-01-21 17:59:28', '2014-01-21 17:59:28', 'CREDIT', 'Pseudo payment for applying credit to Invoices ',
   NULL, NULL, NULL, NULL);

-- name: insert-data-fn_subscription
INSERT INTO `%s`.`fn_subscription` VALUES
  (1, 1, 1, '2015-08-30 17:10:22', '2015-08-30 17:10:22', 12, 'Associate Membership',
   'Annual subscription for Associates', 60, -30),
  (2, 1, 1, '2015-08-30 17:10:23', '2015-08-30 17:10:23', 12, 'Ordinary Membership',
   'Annual subscription for Ordinary Members', 60, -30),
  (3, 1, 1, '2015-08-30 17:10:23', '2015-08-30 17:10:23', 12, 'Fellowship', 'Annual subscription for Fellows', 60, -30);

-- name: insert-data-fn_subscription_inventory
INSERT INTO `%s`.`fn_subscription_inventory` VALUES
  (1, 1, 1, 1, '2015-08-26 04:05:31', '2015-08-26 04:05:31', 'Associate Membership Fee', 1.00),
  (2, 3, 3, 1, '2015-08-26 04:06:03', '2015-08-26 04:06:03', 'Fellow Membership Fee', 1.00),
  (3, 2, 2, 1, '2015-08-26 04:06:51', '2015-08-26 04:06:51', 'Ordinary Membership Fee', 1.00);

-- name: insert-data-fn_subscription_type
INSERT INTO `%s`.`fn_subscription_type` VALUES
  (1, 1, '2013-07-15 18:35:22', '2013-07-15 18:35:22', 1, 'Membership', 'Membership subscription');

-- name: insert-data-fn_tax
INSERT INTO `%s`.`fn_tax` VALUES
  (1, 14, 1, '2013-05-16 14:39:39', NOW(), 'GST', 'Goods & Services Tax (AU)', 10.00),
  (2, 159, 1, '2013-05-16 14:39:39', NOW(), 'GST', 'Goods & Services Tax (NZ)', 15.00);

-- name: insert-data-fs_set
INSERT INTO `%s`.`fs_set` VALUES
  (1, 1, 1, NOW(), NOW(), 'AWS-S3', '{"key": 1234}', 'test-volume', '/note/', 'wf_attachment'),
  (2, 1, 1, NOW(), NOW(), 'AWS-S3', '{"key": 1234}', 'test-volume', '/resource/', 'ol_resource_file'),
  (3, 1, 1, NOW(), NOW(), 'AWS-S3', '{"key": 1234}', 'test-volume', '/xml/', 'xml'),
  (4, 1, 1, NOW(), NOW(), 'AWS-S3', '{"key": 1234}', 'test-volume', '/cpd/', 'ce_m_activity_attachment');

-- name: insert-data-fs_url
INSERT INTO `%s`.`fs_url` VALUES
  (1, 1, 1, NOW(), NOW(), 10, 'https://cdn.test.com'),
  (2, 1, 1, NOW(), NOW(), 10, 'https://cdn.test.com'),
  (3, 1, 1, NOW(), NOW(), 10, 'https://cdn.test.com'),
  (4, 1, 1, NOW(), NOW(), 10, 'https://cdn.test.com');

-- insert-data-log_data_action

-- insert-data-log_data_field

-- insert-data-log_data_table

-- name: insert-data-member
INSERT INTO `%s`.`member` VALUES
  (1, 2, 0, 14, NULL, 1, 1, 1, 1, NOW(), NOW(), NULL, '1970-11-03', '2000-01-01', 'M', 'Michael', 'Peter', 'Donnici',
                                                NULL,
                                                NULL, '0402123123', 'michael@mesa.net.au', NULL,
   '5f4dcc3b5aa765d61d8327deb882cf99', NULL, NULL, NULL);

-- name: insert-data-mp_accreditation
INSERT INTO `%s`.`mp_accreditation` VALUES
  (1, 0, 1, '2015-08-30 17:10:25', '2015-08-30 17:10:25', 0, 'General Accreditation',
   'Industry accreditations that does NOT expire', '');

-- name: insert-data-mp_contact_type
INSERT INTO `%s`.`mp_contact_type` VALUES
  (1, 1, 1, 1, 100, NOW(), NOW(), 10, 'Mail',
   'Primary contact information for membership and billing.'),
  (2, 1, 1, 1, 50, NOW(), NOW(), 20, 'Directory',
   'The contact details that will show in the Member directory.'),
  (3, 1, 0, 1, 0, NOW(), NOW(), 30, 'Work', 'A place of work.'),
  (4, 1, 0, 1, 0, NOW(), NOW(), 40, 'Courier', 'Address details for courier.'),
  (5, 1, 0, 1, 0, NOW(), NOW(), 50, 'Home', 'Home contact information.'),
  (6, 1, 0, 1, 0, NOW(), NOW(), 100, 'Other', 'Any other contact type'),
  (7, 1, 0, 1, 10, NOW(), NOW(), 15, 'Journal', 'Postal address for hard-copy Journal.');

-- insert-data-mp_m_accreditation

-- name: insert-data-mp_m_contact
INSERT INTO `%s`.`mp_m_contact` VALUES
  (1, 1, 2, 14, 1, NOW(), NOW(), NULL, '+61 2 9999 9999', '+61 2 9999 8888',
      'info@here.com', 'www.here.com', '123 Some Street', 'Unit abc', 'back shed', 'Jervis Bay', 'NSW', '2540', NULL);

-- name: insert-data-mp_m_position
INSERT INTO `%s`.`mp_m_position` VALUES
  (1, 1, 1, 317, 1, NOW(), NOW(), '2017-12-03', '2017-12-03', NULL),
  (2, 1, 2, 95, 1, NOW(), NOW(), '2018-03-23', '2019-03-26', NULL),
  (3, 1, 3, 14, 1, NOW(), NOW(), NULL, NULL, NULL);

-- name: insert-data-mp_m_qualification
INSERT INTO `%s`.`mp_m_qualification` VALUES
  (1, 1, 21, 0, 1, NOW(), NOW(), 1992, 'BBio', NULL);

-- name: insert-data-mp_m_speciality
INSERT INTO `%s`.`mp_m_speciality` VALUES
  (1, 1, 1, 1, NOW(), NOW(), NULL, NULL, 'This is a comment'),
  (2, 1, 2, 1, NOW(), NOW(), NULL, NULL, 'This is another comment');

-- insert-data-mp_m_tag

-- name: insert-data-mp_position
INSERT INTO `%s`.`mp_position` VALUES
  (1, 0, 1, '2015-08-30 17:10:25', '2015-08-30 17:10:25', 0, 'AFFIL', 'Affiliate',
   '... has an affiliation with an organisation'),
  (2, 0, 1, '2015-08-30 17:10:26', '2015-08-30 17:10:26', 0, 'MEMBER', 'Member', '... is a member of an organisation'),
  (3, 0, 1, '2015-08-30 17:10:26', '2015-08-30 17:10:26', 0, 'CHAIR', 'Chair',
   '... is the chair of a group / committee'),
  (4, 0, 1, '2015-08-30 17:10:26', '2015-08-30 17:10:26', 0, 'MANAGER', 'Manager', 'A management position.');

-- name: insert-data-mp_qualification
INSERT INTO `%s`.`mp_qualification` VALUES
  (1, 0, 1, NOW(), NOW(), 0, 'BCh', 'Bachelor of Chemistry', ''),
  (2, 0, 1, NOW(), NOW(), 0, 'MBBS', 'Bachelor of Medicine, Bachelor of Surgery', ''),
  (3, 0, 1, NOW(), NOW(), 0, 'BSc', 'Bachelor of Science', ''),
  (4, 0, 1, NOW(), NOW(), 0, 'ChB', 'Bachelor of Surgery', ''),
  (5, 0, 1, NOW(), NOW(), 0, 'DDU', 'Diploma of Diagnostic Ultrasound', ''),
  (6, 0, 1, NOW(), NOW(), 0, 'DSc', 'Doctor of Science', ''),
  (7, 0, 1, NOW(), NOW(), 0, 'FACC', 'Fellow of the American College of Cardiology', ''),
  (8, 0, 1, NOW(), NOW(), 0, 'FACS', 'Fellow of the American College of Surgeons', ''),
  (9, 0, 1, NOW(), NOW(), 0, 'FAHA', 'Fellow of the American Heart Association', ''),
  (10, 0, 1, NOW(), NOW(), 0, 'FCCP', 'Fellow of the American College of Chest Physicians', ''),
  (11, 0, 1, NOW(), NOW(), 0, 'FCSANZ', 'Fellow of the Cardiac Society of Australia and New Zealand', ''),
  (12, 0, 1, NOW(), NOW(), 0, 'FRACP', 'Fellow of the Royal Australasian College of Physicians', ''),
  (13, 0, 1, NOW(), NOW(), 0, 'FRACS', 'Fellow of the Royal Australiasian College of Surgeons', ''),
  (14, 0, 1, NOW(), NOW(), 0, 'FRCP', 'Fellow of the Royal College of Physicians', ''),
  (15, 0, 1, NOW(), NOW(), 0, 'FRCS', 'Fellow of the Royal College of Surgeons', ''),
  (16, 0, 1, NOW(), NOW(), 0, 'FTSE', 'Fellow of the Austalian Academy of Technological Sciences and Engineering', ''),
  (17, 0, 1, NOW(), NOW(), 0, 'MB', 'Bachelor of Medicine', ''),
  (18, 0, 1, NOW(), NOW(), 0, 'MD', 'Doctor of Medicine', ''),
  (19, 0, 1, NOW(), NOW(), 0, 'MRCP', 'Member of the Royal College of Physicians', ''),
  (20, 0, 1, NOW(), NOW(), 0, 'MSc', 'Master of Science', ''),
  (21, 0, 1, NOW(), NOW(), 0, 'PhD', 'PhD', ''),
  (22, 0, 1, NOW(), NOW(), 0, 'BS', 'Bachelor of Surgery', ''),
  (23, 0, 1, NOW(), NOW(), 0, 'BEd', 'Bachelor of Education', ''),
  (24, 0, 1, NOW(), NOW(), 0, 'MEd', 'Master of Education', ''),
  (25, 0, 1, NOW(), NOW(), 0, 'MA', 'Master of Arts', ''),
  (26, 0, 1, NOW(), NOW(), 0, 'MBA', 'Master of Business Administration', ''),
  (27, 0, 1, NOW(), NOW(), 0, 'BSc (Med)', 'Bachelor of Science (Medicine)', ''),
  (28, 0, 1, NOW(), NOW(), 0, 'BMedSc', 'Bachelor of Medical Sciences', ''),
  (29, 0, 1, NOW(), NOW(), 0, 'BM', 'Bachelor of Medicine', '');

-- name: insert-data-mp_speciality
INSERT INTO `%s`.`mp_speciality` VALUES
  (1, 1, NOW(), NOW(), 'Cardiac Care Nurse (Medical)', '', 0),
  (2, 1, NOW(), NOW(), 'Cardiac Cath Lab Nurse', '', 0),
  (3, 1, NOW(), NOW(), 'Cardiac Technologist', '', 0),
  (4, 1, NOW(), NOW(), 'Cardiovascular Genetic Diseases', '', 0),
  (5, 1, NOW(), NOW(), 'Cardiovascular Surgery', '', 0);

-- name: insert-data-mp_tag
INSERT INTO `%s`.`mp_tag` VALUES
  (1, 1, NOW(), NOW(), 'Allied Health', ''),
  (2, 1, NOW(), NOW(), 'Nurse', ''),
  (3, 1, NOW(), NOW(), 'Surgeon', ''),
  (4, 1, NOW(), NOW(), 'Advanced Trainee', ''),
  (5, 1, NOW(), NOW(), 'Retired / Comp', 'Complimentary subscription');

-- name: insert-data-ms_m_application
INSERT INTO `%s`.`ms_m_application` VALUES
  (1, 502, 388, 440, 2, NULL, 1, '2015-09-01 04:30:30', '2015-09-01 04:44:01', '2015-09-01', 1, 'first application'),
  (2, 482, 72, 440, 2, 1, 1, '2015-09-01 06:28:15', '2015-09-01 06:43:22', '2015-09-01', 1, NULL),
  (3, 488, 389, 456, 2, 1, 1, '2015-09-01 06:29:00', '2015-09-01 06:43:23', '2015-09-01', 1, NULL),
  (4, 499, 423, 195, 4, 2, 1, '2015-09-01 06:29:56', '2015-09-01 06:45:34', '2015-09-01', 1, NULL),
  (5, 485, 168, 317, 4, 2, 1, '2015-09-01 06:31:16', '2015-09-01 06:45:34', '2015-09-01', 1, NULL),
  (6, 502, 388, 440, 4, NULL, 1, '2017-09-01 04:30:30', '2017-09-01 04:44:01', '2017-07-01', 1, 'second application'),
  (7, 502, 388, 440, 4, NULL, 0, '2017-09-01 04:30:30', '2017-09-01 04:44:01', '2017-07-01', 1, 'soft-deleted record');

-- name: insert-data-ms_m_application_meeting
INSERT INTO `%s`.`ms_m_application_meeting` VALUES
  (1, 1, NULL, 1, 1, '2015-09-01 04:32:16', '2015-09-01 04:32:16', -1, NULL);

-- insert-data-ms_m_permission

-- name: insert-data-ms_m_status
INSERT INTO `%s`.`ms_m_status` VALUES
  (1, 1, 1, 0, 1, 1, '2015-08-30 17:10:57', NOW(), NULL);

-- name: insert-data-ms_m_title
INSERT INTO `%s`.`ms_m_title` VALUES
  (1, 1, 2, NULL, 1, 1, '2015-08-30 17:10:57', NOW(), '2015-08-31',
   'Test data does will not show historic titles');


-- insert-data-ms_permission


-- name: insert-data-ms_status
INSERT INTO `%s`.`ms_status` VALUES
  (1, 1, 1, '2013-01-21 14:52:28', '2013-06-18 15:11:35', 1, 1, 1, 1, 1, 'Active', ''),
  (10003, 1, 0, '2013-01-29 11:41:19', '2013-01-29 11:41:19', 0, 0, 0, 0, 0, 'Pending', ''),
  (10004, 1, 0, '2013-01-29 11:41:33', '2013-06-18 15:11:42', 0, 0, 0, 0, 0, 'Lapsed', ''),
  (10005, 1, 0, '2013-01-29 11:41:53', '2013-06-18 15:11:50', 0, 0, 0, 0, 0, 'Suspended', ''),
  (10006, 1, 0, '2013-01-29 11:42:12', '2013-06-18 15:11:57', 0, 0, 0, 0, 0, 'Inactive', ''),
  (10007, 1, 0, '2013-05-15 20:12:00', '2013-06-18 15:12:04', 0, 0, 0, 0, 0, 'Resigned', ''),
  (10008, 1, 0, NOW(), '2013-06-18 15:12:11', 0, 0, 0, 0, 0, 'Deceased', '');

-- name: insert-data-ms_title
INSERT INTO `%s`.`ms_title` VALUES
  (1, 0, 1, '2013-01-17 23:17:04', '2013-07-08 19:44:45', 0, 0, 0, 0, 0, 'Applicant', ''),
  (2, 0, 1, '2013-01-17 23:17:05', '2013-10-02 16:37:57', 1, 1, 1, 1, 1, 'Associate', ''),
  (3, 0, 1, '2013-01-29 11:40:21', '2013-01-29 11:40:21', 1, 1, 1, 1, 1, 'Ordinary', ''),
  (4, 0, 1, '2013-01-29 11:39:01', '2013-10-02 16:38:09', 1, 1, 1, 1, 1, 'Fellow', '');

-- name: insert-data-ol_category
INSERT INTO `%s`.`ol_category` VALUES
  (1, 1, '2015-08-30 17:10:30', '2015-08-30 17:10:30', 'GMED', 'General Medicine', ''),
  (2, 1, '2015-08-30 17:10:30', '2015-08-30 17:10:30', 'CARD', 'Cardiology', ''),
  (3, 1, '2015-08-30 17:10:31', '2015-08-30 17:10:31', 'OPTH', 'Opthamology', ''),
  (4, 1, '2015-08-30 17:10:31', '2015-08-30 17:10:31', 'PARA', 'Paramedicine', ''),
  (5, 1, '2015-08-30 17:10:31', '2015-08-30 17:10:31', 'MARINE', 'Marine', ''),
  (6, 1, '2017-01-25 06:21:51', '2017-01-25 06:21:51', 'DAN', 'Dance', '');

-- name: insert-data-ol_m_category
INSERT INTO `%s`.`ol_m_category` VALUES
  (1, 501, 4, 1, '2015-08-31 03:47:10', '2015-08-31 03:47:10'),
  (2, 7821, 4, 1, '2016-07-19 01:22:29', '2016-07-19 01:22:29'),
  (3, 7821, 6, 1, '2017-01-25 06:23:16', '2017-01-25 06:23:16');

-- insert-data-ol_m_module

-- insert-data-ol_m_module_rating

-- insert-data-ol_m_module_slide

-- insert-data-ol_m_module_slide_option

-- insert-data-ol_module

-- insert-data-ol_module_category

-- insert-data-ol_module_cpd

-- insert-data-ol_module_rating

-- insert-data-ol_module_resource

-- insert-data-ol_option

-- name: insert-data-ol_resource
INSERT INTO `%s`.`ol_resource` VALUES
  (6576, 80, 1, 1, '2018-02-14 23:34:22', '2018-02-15 05:25:36', '2017-12-01', 2017, 12, 1,
         'Gavage of Fecal Samples From Patients With Colorectal Cancer Promotes Intestinal Carcinogenesis in Germ-Free and Conventional Mice.',
   'Altered gut microbiota is implicated in development of colorectal cancer (CRC). Some intestinal bacteria have been reported to potentiate intestinal carcinogenesis by producing genotoxins, altering the immune response and intestinal microenvironment, and activating oncogenic signaling pathways. We investigated whether stool from patients with CRC could directly induce colorectal carcinogenesis in mice.',
   'Carcinogenesis,Colon Cancer,Germ-Free,Stool Transplantation,Animals,Azoxymethane,Case-Control Studies,Cell Proliferation,Cell Transformation Neoplastic,Colon,Colonic Polyps,Colorectal Neoplasms,Disease Models Animal,Feces,Gastrointestinal Microbiome,Gene Expression Regulation Neoplastic,Germ-Free Life,Host-Pathogen Interactions,Humans,Inflammation Mediators,Ki-67 Antigen,Lymphocytes Tumor-Infiltrating,Male,Mice Inbred C57BL,Th1 Cells,Th17 Cells,Wong SH,Zhao L,Zhang X,Nakatsu G,Han J,Xu W,Xiao X,Kwong TNY,Tsoi H,Wu WKK,Zeng B,Chan FKL,Sung JJY,Wei H,Yu J,28823860',
   'https://doi.org/10.1053/j.gastro.2017.08.022', 'http://localhost:8080/r6576', '', NULL,
   '{\"category\":\"\",\"free\":false,\"public\":false,\"source\":\"pubmed\",\"sourceId\":\"28823860\",\"sourceName\":\"Gastroenterology\",\"sourceNameAbbrev\":\"Gastroenterology\",\"sourcePubDate\":\"2017 Dec\",\"sourceVolume\":\"153\",\"sourceIssue\":\"6\",\"sourcePages\":\"1621-1633.e6\"}'),
  (6577, 80, 1, 1, '2018-02-14 23:34:22', '2018-02-15 05:25:36', '2017-12-01', 2017, 12, 1,
         'NETSstudy: development of a Hirschsprung\'s disease core outcome set.',
   'The objective of this study was to develop a Hirschsprung\'s disease (HD) core outcome set (COS).',
   'Core Outcome Set,Gastroenterology,Hirschsprung’s Disease,Paediatric Surgery,Adolescent,Child,Child Preschool,Delphi Technique,Developed Countries,Hirschsprung Disease,Humans,Infant,Infant Newborn,Patient Reported Outcome Measures,Severity of Illness Index,Stakeholder Participation,Treatment Outcome,Allin BSR,Bradnock T,Kenny S,Kurinczuk JJ,Walker G,Knight M,,28784616',
   'https://doi.org/10.1136/archdischild-2017-312901', 'http://localhost:8080/r6577', '', NULL,
   '{\"category\":\"\",\"free\":false,\"public\":false,\"source\":\"pubmed\",\"sourceId\":\"28784616\",\"sourceName\":\"Archives of disease in childhood\",\"sourceNameAbbrev\":\"Arch Dis Child\",\"sourcePubDate\":\"2017 Dec\",\"sourceVolume\":\"102\",\"sourceIssue\":\"12\",\"sourcePages\":\"1143-1151\"}'),
  (6578, 80, 1, 1, '2018-02-14 23:34:22', '2018-02-15 05:25:37', '2017-12-01', 2017, 12, 1,
         'Thymus transplantation for complete DiGeorge syndrome: European experience.',
   'Thymus transplantation is a promising strategy for the treatment of athymic complete DiGeorge syndrome (cDGS).',
   'DiGeorge syndrome,athymia,thymus transplantation,Autoimmune Diseases,Cells Cultured,Child,Child Preschool,DiGeorge Syndrome,Europe,Female,Humans,Immune Reconstitution,Infant,Male,Organ Culture Techniques,Organ Transplantation,Postoperative Complications,T-Lymphocytes,Thymus Gland,Transplantation Homologous,Treatment Outcome,Davies EG,Cheung M,Gilmour K,Maimaris J,Curry J,Furmanski A,Sebire N,Halliday N,Mengrelis K,Adams S,Bernatoniene J,Bremner R,Browning M,Devlin B,Erichsen HC,Gaspar HB,Hutchison L,Ip W,Ifversen M,Leahy TR,McCarthy E,Moshous D,Neuling K,Pac M,Papadopol A,Parsley KL,Poliani L,Ricciardelli I,Sansom DM,Voor T,Worth A,Crompton T,Markert ML,Thrasher AJ,28400115',
   'https://doi.org/10.1016/j.jaci.2017.03.020', 'http://localhost:8080/r6578', '', NULL,
   '{\"category\":\"\",\"free\":false,\"public\":false,\"source\":\"pubmed\",\"sourceId\":\"28400115\",\"sourceName\":\"The Journal of allergy and clinical immunology\",\"sourceNameAbbrev\":\"J Allergy Clin Immunol\",\"sourcePubDate\":\"2017 Dec\",\"sourceVolume\":\"140\",\"sourceIssue\":\"6\",\"sourcePages\":\"1660-1670.e16\"}');

-- insert-data-ol_resource_attribute

-- insert-data-ol_resource_attribute_value

-- insert-data-ol_resource_file

-- insert-data-ol_resource_filetype

-- insert-data-ol_resource_type

-- insert-data-ol_slide

-- insert-data-ol_slide_resource

-- name: insert-data-organisation
INSERT INTO `%s`.`organisation` VALUES
  (1, NULL, 1, 1, 1, NOW(), NOW(), 'ABC', 'ABC Organisation', '', '', '', '', '', '', '', '', '', ''),
  (2, NULL, 1, 1, 1, NOW(), NOW(), 'DEF', 'DEF Organisation', '', '', '', '', '', '', '', '', '', ''),
  (3, 1, 2, 1, 1, NOW(), NOW(), 'ABC-1', 'ABC Sub1', '', '', '', '', '', '', '', '', '', ''),
  (4, 1, 3, 1, 1, NOW(), NOW(), 'ABC-2', 'ABC Sub2', '', '', '', '', '', '', '', '', '', ''),
  (5, 1, 4, 1, 1, NOW(), NOW(), 'ABC-3', 'ABC Sub3', '', '', '', '', '', '', '', '', '', '');

-- name: insert-data-organisation_type
INSERT INTO `%s`.`organisation_type` VALUES
  (1, 1, '2013-06-11 20:45:11', '2013-06-11 20:45:31', 'CSANZ Council', 'Group / type for CSANZ councils'),
  (2, 1, NOW(), NOW(), 'CSANZ Working Group', 'CSANZ Working Group'),
  (3, 1, NOW(), NOW(), 'Committee', 'Committee'),
  (4, 1, NOW(), NOW(), 'Sub-committee', 'Sub-committee'),
  (5, 1, NOW(), '2013-08-19 15:57:34', 'CSANZ Board', 'Group / type for CSANZ Board'),
  (6, 1, NOW(), NOW(), 'Other', 'Other'),
  (7, 1, '2013-07-14 15:37:02', '2013-07-14 15:37:13', 'Institute / Hospital', ''),
  (8, 1, '2013-07-14 17:32:31', '2013-07-14 17:32:31', 'University / Education', ''),
  (9, 1, '2014-07-23 12:00:14', '2014-07-23 12:00:14', 'Heart Foundation', 'Heart Foundation of Australia'),
  (10, 1, '2014-08-04 19:21:08', '2014-08-04 19:21:08', 'Government', '');

-- name: insert-data-wf_attachment
INSERT INTO `%s`.`wf_attachment` VALUES
  (1, 1, 1, 1, 1, NOW(), NOW(), 'filename.ext');

-- name: insert-data-wf_issue
INSERT INTO `%s`.`wf_issue` VALUES
  (1, 1, NULL, NULL, NULL, 1, 1, 1, '2015-09-01 06:13:27', '2015-09-01 06:57:09', '2015-09-01',
   'A new invoice has been raised and is pending payment. (INV0001)',
   'Members can pay online or by alternate methods specified on the invoice.'),
  (2, 4, NULL, NULL, NULL, 1, 0, 0, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01',
   'All active members require a membership subscription, even if that subscription is complimentary.',
   'Assign the appropriate membership subscription.'),
  (3, 4, NULL, NULL, NULL, 1, 0, 0, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01',
   'All active members require a membership subscription, even if that subscription is complimentary.',
   'Assign the appropriate membership subscription.');

-- name: insert-data-wf_issue_association
INSERT INTO `%s`.`wf_issue_association` VALUES
  (1, 1, 502, 1, 1, '2015-09-01 06:13:27', '2015-09-01 06:13:27', 'invoice'),
  (2, 2, 1, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', NULL),
  (3, 3, 2, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', NULL);

-- name: insert-data-wf_issue_category
INSERT INTO `%s`.`wf_issue_category` VALUES
  (1, 1, '2013-07-04 14:28:03', '2013-07-04 14:34:15', 'General Administration',
   'Miscellaneous, uncategorised and general administrative Issues.'),
  (2, 1, '2013-07-04 14:34:58', '2013-07-04 14:34:58', 'Membership Applications',
   'Issues relating to Membership Applications.'),
  (3, 1, '2013-07-04 14:35:39', '2013-07-04 14:35:39', 'CPD', 'Issues relating to CPD & Online Learning.'),
  (4, 1, '2013-07-04 14:36:09', '2013-07-04 14:36:09', 'Finance',
   'Issues relating to Subscriptions, Invoicing and Payments.'),
  (5, 1, '2013-07-04 14:36:36', '2013-07-04 14:36:36', 'Data Integrity', 'Issues relating to missing or broken data.');

-- name: insert-data-wf_issue_type
INSERT INTO `%s`.`wf_issue_type` VALUES 
(1,4,NULL,1,1,1,1,'2013-07-05 10:37:45','2013-11-28 15:53:14','Invoice Raised','A new invoice has been raised and is pending payment.','Members can pay online or by alternate methods specified on the invoice.',NULL),
(2,4,NULL,1,1,1,1,'2013-07-05 10:38:26','2013-11-28 15:54:15','Invoice Past Due','Invoice is past due and pending payment.','Please pay online or by alternate methods specified on the invoice.',NULL),
(3,5,NULL,1,1,1,0,'2013-09-04 12:32:00','2013-09-10 10:16:40','Missing Primary Email','All active members require a primary email to access the system and for communications.','Add a primary email address to member profile.',NULL),
(4,5,NULL,1,1,1,0,'2013-09-04 12:34:11','2013-09-10 10:15:49','Missing Membership Subscription','All active members require a membership subscription, even if that subscription is complimentary.','Assign the appropriate membership subscription.',NULL),
(5,5,NULL,1,1,1,0,'2013-09-04 12:36:54','2013-09-10 10:18:07','Empty Contact Card','All active members require at least one field to be filled in for each persistent contact card.','Edit member contact information and provide at least one value (or \'n/a\') in each of the persistent contact cards.',NULL),
(6,4,NULL,1,1,1,1,'2013-11-08 09:48:01','2013-11-28 15:55:38','Invoice Final Notice','Invoice is past due and a final notice has been issued by email.','Please note no further notices will be emailed.  If your subscription is not paid immediately, your Membership of the Society will be cancelled.',NULL),
(8,1,NULL,1,1,0,0,'2015-04-09 18:38:32','2015-04-09 18:38:32','Email Communication Failure','A recent email communication failed for some reason. ','Check the specific messages in the Members communication tab for clues as to the appropriate follow up.',NULL),
(9,4,NULL,1,1,0,0,'2016-03-21 14:12:09','2016-03-21 14:12:09','Invoice Overpaid','Total of payments allocated to invoice exceeds the invoice total. ','Require manual intervention to remove payment allocations as well as refund if applicable.',NULL),
(10,2,NULL,1,1,0,0,'2019-03-12 10:45:07','2019-03-12 10:45:07','Online Application','Online applications pending acceptance.','Check supplied information, assign appropriate title and status, allocate to meetings.',NULL),
(10000,1,NULL,1,0,0,0,'2013-09-11 17:06:29','2013-09-12 11:53:12','General Admin','-','-',NULL);

-- name: insert-data-wf_note
INSERT INTO `%s`.`wf_note` VALUES
  (1, 10001, 1, 1, 1, '2015-09-01 04:32:49', '2015-09-01 04:32:49', '2015-09-01', 'Application note'),
  (2, 1, NULL, NULL, 1, '2015-09-01 06:13:27', '2015-09-01 06:13:27', '2015-09-01', 'Issue raised'),
  (3, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised'),
  (4, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised');

-- name: insert-data-wf_note_association
INSERT INTO `%s`.`wf_note_association` VALUES
  (1, 1, 1, 1, 1, NOW(), NOW(), 'application'),
  (2, 2, 1, 1, 1, NOW(), NOW(), 'issue'),
  (3, 3, 1, 2, 1, NOW(), NOW(), 'issue');

-- name: insert-data-wf_note_type
INSERT INTO `%s`.`wf_note_type` (`id`,`active`,`system`,`created_at`,`updated_at`,`name`,`description`) VALUES
  (1,1,1,'2013-11-08 09:53:15','2013-11-08 09:53:15','System','Note added by the system housekeeping functions'),
  (10001,1,0,'2013-05-07 02:18:02','2013-06-27 12:43:21','General',''),
  (10002,1,0,'2013-06-27 12:43:33','2013-06-27 12:43:33','Award',''),
  (10003,1,0,'2013-06-27 15:15:10','2013-06-27 15:15:10','Call',''),
  (10004,1,0,'2013-06-27 15:38:56','2013-06-27 15:38:56','Deceased',''),
  (10005,1,0,'2013-06-27 16:12:30','2013-06-27 16:12:30','Email',''),
  (10006,1,0,'2013-06-27 16:22:04','2013-06-27 16:22:04','File Note',''),
  (10007,1,0,'2013-06-27 16:22:04','2013-06-27 16:22:04','History',''),
  (10008,1,0,'2013-06-27 16:22:04','2013-06-27 16:22:04','Lapsed',''),
  (10009,1,0,'2013-06-27 16:22:04','2013-07-31 16:28:04','Lead Extraction',''),
  (10010,1,0,'2013-06-27 16:22:04','2013-06-27 16:22:04','Named Lecture',''),
  (10011,1,0,'2013-06-27 16:22:04','2013-06-27 16:22:04','Prize',''),
  (10012,1,0,'2013-06-27 16:22:04','2013-06-27 16:22:04','Removed',''),
  (10013,1,0,'2013-06-28 11:02:47','2013-06-28 11:02:47','Research',''),
  (10014,1,0,'2013-06-28 11:03:06','2013-06-28 11:03:06','Resigned',''),
  (10015,1,0,'2013-06-28 11:03:06','2013-06-28 11:03:06','Retired',''),
  (10016,1,0,'2013-06-28 11:03:06','2013-06-28 11:03:06','Suspend',''),
  (10017,1,0,'2013-07-16 12:28:28','2013-07-16 12:28:28','Communication','Define type of communication, eg email'),
  (10018,1,0,'2013-07-16 16:22:40','2013-07-16 16:22:40','Advanced Trainee','Name of Hospital'),
  (10019,1,0,'2013-08-08 11:45:48','2013-08-08 11:45:48','Other / Scholarship','ASM Scholarship, Travelling Scholarship to ACC AHA and ESC, Research Scholarships'),
  (10020,1,0,'2013-10-28 08:25:37','2013-10-28 08:25:37','Personal','');

