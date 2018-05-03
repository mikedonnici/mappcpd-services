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
  (1, 1, '2013-06-03 17:29:45', '0000-00-00 00:00:00', 'A/Prof'),
  (2, 1, '2013-06-03 17:29:46', '0000-00-00 00:00:00', 'A/Prof Dame'),
  (3, 1, '2013-06-03 17:29:47', '0000-00-00 00:00:00', 'Brig'),
  (4, 1, '2013-06-03 17:29:47', '0000-00-00 00:00:00', 'Dame'),
  (5, 1, '2013-06-03 17:29:48', '0000-00-00 00:00:00', 'Dr'),
  (6, 1, '2013-06-03 17:29:48', '0000-00-00 00:00:00', 'Miss'),
  (7, 1, '2013-06-03 17:29:49', '0000-00-00 00:00:00', 'Mr'),
  (8, 1, '2013-06-03 17:29:49', '0000-00-00 00:00:00', 'Mrs'),
  (9, 1, '2013-06-03 17:29:50', '0000-00-00 00:00:00', 'Ms'),
  (10, 1, '2013-06-03 17:29:51', '0000-00-00 00:00:00', 'Professor'),
  (11, 1, '2013-06-03 17:29:51', '0000-00-00 00:00:00', 'Professor Sir'),
  (12, 1, '2013-06-03 17:29:52', '0000-00-00 00:00:00', 'Sir'),
  (13, 1, '2013-06-03 17:29:52', '0000-00-00 00:00:00', 'Sister');

-- name: insert-data-acl_admin_resource
INSERT INTO `%s`.`acl_admin_resource` VALUES
  (1, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Home', 'Home page'),
  (2, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Member', 'Manage Member Records'),
  (3, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Workflow', 'Manage administration workflows'),
  (4, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Finance', 'Finance Module'),
  (5, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Events', 'Manage events'),
  (6, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Learning', 'Learning Module'),
  (7, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Communication', 'Communication Module'),
  (8, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Macro', 'Manage Macros'),
  (9, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Report', 'Report Module'),
  (10, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Setup', 'Setup system details'),
  (11, 10, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Setup Admin User', 'Setup Admin User'),
  (12, 10, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Setup ACL', 'Setup ACL'),
  (10000, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'function', 'Log', 'Log system activities'),
  (10001, 10000, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'function', 'Log Data Changes',
   'Log changes made to only tables defined in \'log_data_table\' db table'),
  (10002, 6, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'function', 'Learning CPD',
   'Assign CPD point on completing a learning module');

-- name: insert-data-acl_admin_role
INSERT INTO `%s`.`acl_admin_role` VALUES
  (1, 1, 0, '2013-10-02 16:08:09', '2013-10-02 16:08:09', 'Super User', 'System administrator account'),
  (2, 1, 1, '2013-10-02 16:08:09', '2013-10-02 16:08:09', 'Admin', 'General account for admin level users'),
  (3, 1, 0, '2014-09-21 17:43:12', '2014-09-21 17:43:12', 'Limited', 'Limited access to Member search feature.');

-- name: insert-data-acl_admin_role_resource
INSERT INTO `%s`.`acl_admin_role_resource` VALUES
  (9, 2, 1, 1, '2013-10-02 16:35:47', '2013-10-02 16:35:47', NULL),
  (10, 2, 2, 1, '2013-10-02 16:35:47', '2013-10-02 16:35:47', NULL),
  (11, 2, 3, 1, '2013-10-02 16:35:47', '2013-10-02 16:35:47', NULL),
  (12, 2, 4, 1, '2013-10-02 16:35:47', '2013-10-02 16:35:47', NULL),
  (13, 2, 5, 1, '2013-10-02 16:35:47', '2013-10-02 16:35:47', NULL),
  (14, 2, 6, 1, '2013-10-02 16:35:47', '2013-10-02 16:35:47', NULL),
  (15, 2, 7, 1, '2013-10-02 16:35:47', '2013-10-02 16:35:47', NULL),
  (16, 2, 8, 1, '2013-10-02 16:35:47', '2013-10-02 16:35:47', NULL),
  (17, 2, 9, 1, '2013-10-02 16:35:47', '2013-10-02 16:35:47', NULL),
  (19, 2, 10, 1, '2013-11-21 16:11:59', '2013-11-21 16:11:59', NULL),
  (20, 3, 2, 1, '2014-09-21 17:47:04', '2014-09-21 17:47:04', NULL),
  (21, 3, 1, 1, '2014-09-21 17:53:44', '2014-09-21 17:53:44', NULL);

-- name: insert-data-acl_member_resource
INSERT INTO `%s`.`acl_member_resource` VALUES
  (1, 0, 0, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Home', 'Home Page'),
  (2, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Learning', 'Learning Module'),
  (3, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'CPD', 'CPD Module'),
  (4, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Journals', 'Journals Page'),
  (5, 0, 0, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'Curriculum', 'Curriculum Page'),
  (6, 0, 0, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'FAQ', 'FAQ page'),
  (7, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'page', 'MyBilling', 'My Billing Page'),
  (8, 10003, 1, '2013-11-14 19:06:33', '2013-11-14 19:06:33', 'page', 'Resource Library', 'Resource library page'),
  (9, 3, 1, '2013-11-21 15:58:04', '2013-11-21 15:58:04', 'page', 'CPD Evaluation Period', 'Defined CPD Evaluation Periods for reporting and auditing'),
  (10000, 2, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'function', 'Learning CPD', 'Allocate CPD points on completing learning module'),
  (10001, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'function', 'Log', 'Log system activities'),
  (10002, 10001, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'function', 'Log Data',
   'Changes made to only tables defined in \'log_data_table\' db table will be logged'),
  (10003, 0, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'function', 'Library', 'Resource library');

-- name: insert-data-acl_member_role
INSERT INTO `%s`.`acl_member_role` VALUES
  (1, 1, 1, '2013-09-26 20:40:25', '2013-09-26 20:40:25', 'Full Access',
   'Member account has full access to all application areas.'),
  (2, 1, 0, '2013-10-14 14:19:17', '2013-10-14 14:24:24', 'Restricted Access',
   'At this stage member access is limited to <a href=\"/member/index/summary\">Profile</a>  and <a href=\"/finance/invoice\">myBilling</a>.');

-- name: insert-data-acl_member_role_resource
INSERT INTO `%s`.`acl_member_role_resource` VALUES
  (9, 1, 1, 1, '2013-09-26 20:40:43', '2013-09-26 20:40:43', NULL),
  (10, 1, 2, 1, '2013-09-26 20:40:43', '2013-09-26 20:40:43', NULL),
  (11, 1, 3, 1, '2013-09-26 20:40:43', '2013-09-26 20:40:43', NULL),
  (13, 1, 5, 1, '2013-09-26 20:40:43', '2013-09-26 20:40:43', NULL),
  (14, 1, 6, 1, '2013-09-26 20:40:43', '2013-09-26 20:40:43', NULL),
  (15, 1, 7, 1, '2013-09-26 20:40:43', '2013-09-26 20:40:43', NULL),
  (17, 2, 7, 1, '2013-10-14 14:19:36', '2013-10-14 14:19:36', NULL),
  (19, 1, 4, 1, '2013-11-24 14:35:03', '2013-11-24 14:35:03', NULL),
  (20, 1, 8, 1, '2013-12-04 11:35:31', '2013-12-04 11:35:31', NULL),
  (21, 1, 9, 1, '2014-03-06 17:14:49', '2014-03-06 17:14:49', NULL);

-- name: insert-data-ad_macro
INSERT INTO `%s`.`ad_macro` VALUES
  (1, 1, 1, '2015-09-01 06:43:22', '2015-09-01 06:43:22', 'Application', NULL),
  (2, 1, 1, '2015-09-01 06:45:34', '2015-09-01 06:45:34', 'Application', NULL),
  (66, 1, 1, '2016-04-04 03:21:29', '2016-04-04 03:21:29', 'Application', 'All good!'),
  (67, 1, 1, '2016-04-04 03:29:17', '2016-04-04 03:29:17', 'Application', NULL);

-- name: insert-data-ad_macro_transaction
INSERT INTO `%s`.`ad_macro_transaction` VALUES
  (1, 1, 482, 1, 1, '2015-09-01 06:43:23', '2015-09-01 06:43:23',
   'Application ID 2\nSet Title to Associate\nSet Date of Election to 01 Sep 2015\nSet Membership Subscription to Associate Membership\nSet Access Level to Full Access'),
  (2, 1, 488, 1, 1, '2015-09-01 06:43:23', '2015-09-01 06:43:23',
   'Application ID 3\nSet Title to Associate\nSet Date of Election to 01 Sep 2015\nSet Membership Subscription to Associate Membership\nSet Access Level to Full Access'),
  (3, 2, 499, 1, 1, '2015-09-01 06:45:34', '2015-09-01 06:45:34',
   'Application ID 4\nSet Title to Fellow\nSet Date of Election to 01 Sep 2015\nSet Membership Subscription to Fellowship\nSet Access Level to Full Access'),
  (4, 2, 485, 1, 1, '2015-09-01 06:45:34', '2015-09-01 06:45:34',
   'Application ID 5\nSet Title to Fellow\nSet Date of Election to 01 Sep 2015\nSet Membership Subscription to Fellowship\nSet Access Level to Full Access'),
  (373, 66, 7822, 1, 1, '2016-04-04 03:21:29', '2016-04-04 03:21:29',
   'Application ID 6536\nResult of following meeting(s) set to Accepted: Sub - 11 April 2016\nSet Title to Associate\nSet Date of Election to 04 Apr 2016\nSet Membership Subscription to Associate Membership\nSet Access Level to Full Access'),
  (374, 67, 7824, 1, 1, '2016-04-04 03:29:18', '2016-04-04 03:29:18',
   'Application ID 6538\nResult of following meeting(s) set to Accepted: Sub - 11 April 2016\nSet Title to Fellow\nSet Date of Election to 04 Apr 2016\nSet Membership Subscription to Fellowship\nSet Access Level to Full Access');

-- insert-data-ad_permission


-- name: insert-data-ad_user
INSERT INTO `%s`.`ad_user` VALUES
  (1, 1, 1, 0, '2015-08-30 17:10:08', '2016-05-31 04:33:03', 'demo-admin', '41d0510a9067999b72f38ba0ce9f6195',
      'Demo Admin', 'demo', 'demo@noemail.com');

-- insert-data-ad_user_permission

-- name: insert-data-ce_activity
INSERT INTO `%s`.`ce_activity` VALUES
  (1, 1, 1, 0, 0, '2015-08-30 17:10:11', '2015-08-30 17:10:11', 'CE1', 'Conference session / workshop / course', '',
      1.00, 50),
  (2, 1, 1, 0, 0, '2015-08-30 17:10:11', '2015-08-30 17:10:11', 'CE2', 'Reading, research, literature review', '', 1.00, 25),
  (3, 1, 1, 0, 0, '2015-08-30 17:10:11', '2016-09-16 14:11:19', 'CE3', 'Teaching - preperation and delivery', '', 1.00, 25),
  (4, 1, 1, 0, 0, '2015-08-30 17:10:11', '2015-08-30 17:10:11', 'CE4', 'Presentation', '', 1.00, 25),
  (5, 1, 1, 0, 0, '2015-08-30 17:10:11', '2015-08-30 17:10:11', 'CE5', 'Online content - other', '', 1.00, 25),
  (6, 3, 1, 0, 1, '2015-08-30 17:10:12', '2015-08-31 03:42:22', 'CE6', 'MappCPD online module', '', 1.00, 50),
  (7, 1, 2, 0, 0, '2015-08-30 17:10:12', '2015-08-30 17:10:12', 'PR1', 'Formal performance review or audit', '', 1.00, 20),
  (8, 1, 2, 0, 0, '2015-08-30 17:10:12', '2015-08-30 17:10:12', 'PR2', 'Informal peer review / meeting', '', 1.00, 20),
  (9, 1, 3, 0, 0, '2015-08-30 17:10:12', '2015-08-30 17:10:12', 'PQ1', 'Personal / professional / management course', '', 1.00, 20),
  (10, 1, 3, 0, 0, '2015-08-30 17:10:12', '2015-08-30 17:10:12', 'PQ2', 'Self-directed learning (professional qualities)', '', 1.00, 20),
  (20, 1, 10, 1, 0, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'RACP1', 'Practice Review & Improvement', '', 3.00, 50),
  (21, 1, 10, 1, 0, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'RACP2', 'Assessed Learning', '', 2.00, 50),
  (22, 1, 10, 1, 0, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'RACP3',
       'Educational Development, Teaching & Research', '', 1.00, 50),
  (23, 1, 10, 1, 0, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'RACP4', 'Group Learning', '', 1.00, 50),
  (24, 1, 10, 1, 0, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'RACP5', 'Other Learning Activities', '', 1.00, 50);

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
  (1, 20, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Practice audits/Clinical audits'),
  (2, 20, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Peer review'),
  (3, 20, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Patient satisfaction studies'),
  (4, 20, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Institution audits, e.g. hospital accreditation'),
  (5, 20, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Incident reporting/monitoring, e.g. morbidity & mortality meetings'),
  (6, 20, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Practice Review, e.g. Regular Practice Review'),
  (7, 20, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Multi Source Feedback (MSF)'),
  (8, 20, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Participation in the RACP Supervisor Professional Development Program (SPDP)'),
  (9, 20, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Other practice review & improvement activities'),
  (10, 21, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'PhD studies'),
  (11, 21, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Formal postgraduate studies'),
  (12, 21, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Self-assessment programs'),
  (13, 21, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Courses to learn new techniques, e.g. Advanced Life Support (ALS)'),
  (14, 21, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Learner initiated and planned projects'),
  (15, 21, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Other assessed learning activities'),
  (16, 22, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Teaching, e.g. supervision, mentoring'),
  (17, 22, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Involvement in standards development'),
  (18, 22, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Reviewer'),
  (19, 22, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Writing examination questions'),
  (20, 22, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Examining'),
  (21, 22, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Publication (including preparation)'),
  (22, 22, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Presentation (including preparation)'),
  (23, 22, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Committee/working group/council involvement'),
  (24, 22, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Other educational development, teaching & research activities'),
  (25, 23, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Seminars'),
  (26, 23, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Conferences'),
  (27, 23, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Workshops'),
  (28, 23, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Grand rounds'),
  (29, 23, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Journal clubs'),
  (30, 23, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Hospital and other medical meetings'),
  (31, 23, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Other group learning activities'),
  (32, 24, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Reading journals and texts'),
  (33, 24, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Information searches, e.g. Medline'),
  (34, 24, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Audio/videotapes'),
  (35, 24, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Web-based learning'),
  (36, 24, 1, '2018-02-04 03:03:07', '2018-02-04 03:03:07', 'Other learning activities');

-- name: insert-data-ce_activity_unit
INSERT INTO `%s`.`ce_activity_unit` VALUES (1, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 1, 'hours', NULL),
  (2, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 0, 'item', 'Item, instance, single event.'),
  (3, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 0, 'module', 'Online learning module.');

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
  (5839, 7821, 23, 25, NULL, NULL, 1, 1, '2018-04-30 04:29:13', '2018-04-30 04:29:13', '2018-04-30', 1.00, 1.00, 0,
   'BJJ like Bruno Malfacine... control hands and sweep'),
  (5840, 7821, 23, 25, NULL, NULL, 1, 0, '2018-04-30 07:09:48', '2018-04-30 07:09:48', '2018-04-30', 1.00, 1.00, 0,
   'Test activity'),
  (5841, 7821, 20, 1, NULL, NULL, 1, 0, '2018-04-30 07:15:24', '2018-04-30 07:15:24', '2018-04-30', 1.00, 3.00, 0,
   'sdasd');

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
  (14, 1, 1, 1000, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'AU', 'AUSTRALIA'),
  (159, 1, 1, 900, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'NZ', 'NEW ZEALAND'),
  (235, 1, 0, 500, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'GB', 'UNITED KINGDOM'),
  (236, 1, 0, 400, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'US', 'UNITED STATES');

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

-- name: insert-data-fn_invoice_payment
INSERT INTO `%s`.`fn_invoice_payment` VALUES (1, 1, 1, 1, '2015-09-01 06:56:48', '2015-09-01 06:56:48', 108.90, NULL);

-- name: insert-data-fn_m_invoice
INSERT INTO `%s`.`fn_m_invoice` VALUES
  (1, 502, 1, 1, 1, 0, '2015-09-01 06:13:26', '2015-09-01 06:56:48', '2015-09-01 06:13:27', '2015-09-01 06:13:33',
      '2015-09-01', '2015-09-15', '2015-09-01', '2015-12-31', 108.90, NULL),
  (2, 482, 1, 1, 0, 0, '2015-09-01 06:46:13', '2017-01-02 09:05:27', '2015-09-01 06:46:13', '2017-01-02 09:05:27',
      '2015-09-01', '2015-09-15', '2015-09-01', '2015-12-31', 108.90, NULL),
  (3, 488, 1, 1, 0, 0, '2015-09-01 06:46:13', '2017-01-02 09:05:27', '2015-09-01 06:46:13', '2017-01-02 09:05:27',
      '2015-09-01', '2015-09-15', '2015-09-01', '2015-12-31', 108.90, NULL),
  (4, 499, 3, 1, 0, 0, '2015-09-01 06:46:13', '2017-01-02 09:05:27', '2015-09-01 06:46:13', '2017-01-02 09:05:27',
      '2015-09-01', '2015-09-15', '2015-09-01', '2015-12-31', 326.70, NULL),
  (5, 485, 3, 1, 0, 0, '2015-09-01 06:46:13', '2017-01-02 09:05:28', '2015-09-01 06:46:13', '2017-01-02 09:05:28',
      '2015-09-01', '2015-09-15', '2015-09-01', '2015-12-31', 326.70, NULL);

-- name: insert-data-fn_m_subscription
INSERT INTO `%s`.`fn_m_subscription` VALUES
  (1, 502, 1, NULL, 1, 0, '2015-09-01 04:46:58', '2017-01-02 09:05:22', NULL, '2018-01-01', NULL),
  (2, 482, 1, 1, 1, 0, '2015-09-01 06:43:23', '2017-01-02 09:05:22', NULL, '2018-01-01', NULL),
  (3, 488, 1, 1, 1, 0, '2015-09-01 06:43:23', '2017-01-02 09:05:22', NULL, '2018-01-01', NULL),
  (4, 499, 3, 2, 1, 0, '2015-09-01 06:45:34', '2017-01-02 09:05:22', NULL, '2018-01-01', NULL),
  (5, 485, 3, 2, 1, 0, '2015-09-01 06:45:34', '2017-01-02 09:05:22', NULL, '2018-01-01', NULL);

-- name: insert-data-fn_payment
INSERT INTO `%s`.`fn_payment` VALUES
  (1, 5, 502, NULL, 1, '2015-09-01 06:56:48', '2015-09-01 06:56:48', '2015-09-01', 108.90, '', '', '', '', '');

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
  (1, 14, 1, '2013-05-16 14:39:39', '0000-00-00 00:00:00', 'GST', 'Goods & Services Tax (AU)', 10.00),
  (2, 159, 1, '2013-05-16 14:39:39', '0000-00-00 00:00:00', 'GST', 'Goods & Services Tax (NZ)', 15.00);

-- insert-data-fs_set

-- insert-data-fs_url

-- insert-data-log_data_action

-- insert-data-log_data_field

-- insert-data-log_data_table

-- name: insert-data-member
INSERT INTO `%s`.`member` VALUES
  (1, 2, 0, 14, NULL, 1, 1, 0, 1, '2015-08-30 17:10:36', '0000-00-00 00:00:00', NULL, '1970-06-17', NULL, NULL, 'Dolan',
                                                                                '', 'Webster', NULL, NULL, NULL,
   'ac@malesuadaut.net', NULL, '', NULL, NULL, NULL),
  (2, 2, 0, 14, NULL, 1, 1, 0, 1, '2015-08-30 17:10:36', '0000-00-00 00:00:00', NULL, '1978-05-12', NULL, NULL, 'Velma',
                                                                                '', 'Whitley', NULL, NULL, NULL,
   'sapien.cursus@gravidamauris.ca', NULL, '', NULL, NULL, NULL),
  (3, 2, 0, 14, NULL, 1, 1, 0, 1, '2015-08-30 17:10:36', '0000-00-00 00:00:00', NULL, '1965-10-07', NULL, NULL,
                                                                                'Jakeem', '', 'Sellers', NULL, NULL,
                                                                                NULL, 'Vivamus@semperetlacinia.edu',
   NULL, '', NULL, NULL, NULL),
  (4, 2, 0, 14, NULL, 1, 1, 0, 1, '2015-08-30 17:10:36', '0000-00-00 00:00:00', NULL, '1954-10-22', NULL, NULL,
                                                                                'Nicole', '', 'Colon', NULL, NULL, NULL,
   'risus.Nunc.ac@magnaLoremipsum.com', NULL, '', NULL, NULL, NULL);

-- name: insert-data-mp_accreditation
INSERT INTO `%s`.`mp_accreditation` VALUES
  (1, 0, 1, '2015-08-30 17:10:25', '2015-08-30 17:10:25', 0, 'General Accreditation',
   'Industry accreditations that does NOT expire', '');

-- name: insert-data-mp_contact_type
INSERT INTO `%s`.`mp_contact_type` VALUES
  (1, 1, 1, 1, 100, '2013-05-16 14:43:12', '0000-00-00 00:00:00', 10, 'Mail',
   'Primary contact information for membership and billing.'),
  (2, 1, 1, 1, 50, '2013-05-16 14:43:12', '0000-00-00 00:00:00', 20, 'Directory',
   'The contact details that will show in the Member directory.'),
  (3, 1, 0, 1, 0, '2013-05-16 14:43:13', '0000-00-00 00:00:00', 30, 'Work', 'A place of work.'),
  (4, 1, 0, 1, 0, '2013-05-16 14:43:14', '0000-00-00 00:00:00', 40, 'Courier', 'Address details for courier.'),
  (5, 1, 0, 1, 0, '2013-05-16 14:43:15', '0000-00-00 00:00:00', 50, 'Home', 'Home contact information.'),
  (6, 1, 0, 1, 0, '2013-06-04 20:12:29', '0000-00-00 00:00:00', 100, 'Other', 'Any other contact type'),
  (7, 1, 0, 1, 10, '2014-07-15 14:56:38', '2014-07-15 14:56:44', 15, 'Journal',
   'Postal address for hard-copy Journal.');

-- insert-data-mp_m_accreditation

-- name: insert-data-mp_m_contact
INSERT INTO `%s`.`mp_m_contact` VALUES
  (1, 7821, 2, 14, 1, '2016-03-03 20:33:37', '2018-03-26 20:31:34', NULL, '+61 2 9999 9999', '+61 2 9999 8888',
      'info@here.com', 'www.here.com', '123 Some Street', 'Unit abc', 'back shed', 'Jervis Bay', 'NSW', '2540', NULL);

-- name: insert-data-mp_m_position
INSERT INTO `%s`.`mp_m_position` VALUES
  (1, 501, 2, 317, 1, '2017-12-02 20:31:22', '2017-12-02 20:31:22', '2017-12-03', '2017-12-03', NULL),
  (2, 7821, 3, 95, 1, '2018-03-26 22:24:48', '2018-03-26 23:01:59', '2018-03-23', '2019-03-26', NULL),
  (3, 7821, 1, 14, 1, '2018-03-26 23:04:29', '2018-03-26 23:04:29', NULL, NULL, NULL);

-- name: insert-data-mp_m_qualification
INSERT INTO `%s`.`mp_m_qualification` VALUES
  (1, 501, 89, 229, 1, '2017-12-02 20:30:57', '2017-12-02 20:30:57', 1992, 'BBio', NULL);

-- name: insert-data-mp_m_speciality
INSERT INTO `%s`.`mp_m_speciality` VALUES
  (1, 7821, 19, 1, '2018-03-26 23:52:12', '2018-03-26 23:52:12', NULL, NULL, 'Test this one.'),
  (2, 7821, 11, 1, '2018-03-27 00:06:15', '2018-03-27 00:06:15', NULL, NULL, 'Likes computers');

-- insert-data-mp_m_tag

-- name: insert-data-mp_position
INSERT INTO `%s`.`mp_position` VALUES
  (1, 0, 1, '2015-08-30 17:10:25', '2015-08-30 17:10:25', 0, 'AFFIL', 'Affiliation',
   '... has an affiliation with an organisation'),
  (2, 0, 1, '2015-08-30 17:10:26', '2015-08-30 17:10:26', 0, 'MEMBER', 'Member', '... is a member of an organisation'),
  (3, 0, 1, '2015-08-30 17:10:26', '2015-08-30 17:10:26', 0, 'CHAIR', 'CHAIR',
   '... is the chair of a group / committee'),
  (4, 0, 1, '2015-08-30 17:10:26', '2015-08-30 17:10:26', 0, 'MANAGER', 'Manager', 'A management position.');

-- name: insert-data-mp_qualification
INSERT INTO `%s`.`mp_qualification` VALUES
  (1, 0, 1, '2013-06-10 16:30:29', '2013-06-10 16:30:29', 0, 'BCh', 'Bachelor of Chemistry', ''),
  (2, 0, 1, '2013-06-10 16:32:15', '2013-06-10 16:32:15', 0, 'MBBS', 'Bachelor of Medicine, Bachelor of Surgery', ''),
  (3, 0, 1, '2013-06-10 16:32:51', '2013-06-10 16:32:51', 0, 'BSc', 'Bachelor of Science', ''),
  (4, 0, 1, '2013-06-10 16:33:37', '2013-06-10 16:33:37', 0, 'ChB', 'Bachelor of Surgery', ''),
  (5, 0, 1, '2013-06-10 16:34:13', '2013-06-10 16:34:13', 0, 'DDU', 'Diploma of Diagnostic Ultrasound', ''),
  (6, 0, 1, '2013-06-10 16:34:33', '2013-06-10 16:34:33', 0, 'DSc', 'Doctor of Science', ''),
  (7, 0, 1, '2013-06-10 16:35:04', '2013-06-10 16:35:04', 0, 'FACC', 'Fellow of the American College of Cardiology', ''),
  (8, 0, 1, '2013-06-10 16:35:27', '2013-06-10 16:35:27', 0, 'FACS', 'Fellow of the American College of Surgeons', ''),
  (9, 0, 1, '2013-06-10 16:35:44', '2013-06-10 16:35:44', 0, 'FAHA', 'Fellow of the American Heart Association', ''),
  (10, 0, 1, '2013-06-10 16:36:03', '2013-06-10 16:36:03', 0, 'FCCP', 'Fellow of the American College of Chest Physicians', ''),
  (11, 0, 1, '2013-06-10 16:36:25', '2013-06-10 16:36:25', 0, 'FCSANZ', 'Fellow of the Cardiac Society of Australia and New Zealand', ''),
  (12, 0, 1, '2013-06-10 16:37:04', '2013-06-10 16:37:04', 0, 'FRACP', 'Fellow of the Royal Australasian College of Physicians', ''),
  (13, 0, 1, '2013-06-10 16:37:23', '2013-10-03 20:22:18', 0, 'FRACS', 'Fellow of the Royal Australiasian College of Surgeons', ''),
  (14, 0, 1, '2013-06-10 16:37:44', '2013-06-10 16:37:44', 0, 'FRCP', 'Fellow of the Royal College of Physicians', ''),
  (15, 0, 1, '2013-06-10 16:38:11', '2013-06-10 16:38:11', 0, 'FRCS', 'Fellow of the Royal College of Surgeons', ''),
  (16, 0, 1, '2013-06-10 16:38:29', '2013-06-10 16:38:29', 0, 'FTSE', 'Fellow of the Austalian Academy of Technological Sciences and Engineering', ''),
  (17, 0, 1, '2013-06-10 16:39:08', '2013-06-10 16:39:08', 0, 'MB', 'Bachelor of Medicine', ''),
  (18, 0, 1, '2013-06-10 16:39:31', '2013-06-10 16:39:31', 0, 'MD', 'Doctor of Medicine', ''),
  (19, 0, 1, '2013-06-10 16:39:54', '2013-06-10 16:39:54', 0, 'MRCP', 'Member of the Royal College of Physicians', ''),
  (20, 0, 1, '2013-06-10 16:40:11', '2013-06-10 19:59:00', 0, 'MSc', 'Master of Science', ''),
  (21, 0, 1, '2013-06-10 16:42:59', '2013-06-10 16:42:59', 0, 'PhD', 'PhD', ''),
  (22, 0, 1, '2013-06-10 16:43:22', '2013-06-10 16:43:22', 0, 'BS', 'Bachelor of Surgery', ''),
  (23, 0, 1, '2013-06-10 20:02:52', '2013-06-10 20:02:52', 0, 'BEd', 'Bachelor of Education', ''),
  (24, 0, 1, '2013-06-10 20:03:11', '2013-06-10 20:03:11', 0, 'MEd', 'Master of Education', ''),
  (25, 0, 1, '2013-06-10 20:03:56', '2013-06-10 20:03:56', 0, 'MA', 'Master of Arts', ''),
  (26, 0, 1, '2013-06-10 20:21:47', '2013-06-10 20:21:47', 0, 'MBA', 'Master of Business Administration', ''),
  (27, 0, 1, '2013-06-10 20:29:57', '2013-06-10 20:29:57', 0, 'BSc (Med)', 'Bachelor of Science (Medicine)', ''),
  (28, 0, 1, '2013-06-10 20:30:45', '2013-06-10 20:30:45', 0, 'BMedSc', 'Bachelor of Medical Sciences', ''),
  (29, 0, 1, '2013-06-10 20:34:03', '2013-06-10 20:34:03', 0, 'BM', 'Bachelor of Medicine', '');

-- name: insert-data-mp_speciality
INSERT INTO `%s`.`mp_speciality` VALUES
  (1, 1, '2013-06-11 20:26:46', '2013-06-11 20:27:44', 'Cardiac Care Nurse (Medical)', '', 0),
  (2, 1, '2013-06-11 20:27:05', '2013-06-11 20:27:05', 'Cardiac Cath Lab Nurse', '', 0),
  (3, 1, '2013-06-11 20:28:01', '2013-06-11 20:28:01', 'Cardiac Technologist', '', 0),
  (4, 1, '2013-06-11 20:28:23', '2013-06-11 20:28:23', 'Cardiovascular Genetic Diseases', '', 0),
  (5, 1, '2013-06-11 20:31:29', '2013-06-11 20:31:29', 'Cardiovascular Surgery', '', 0);

-- name: insert-data-mp_tag
INSERT INTO `%s`.`mp_tag` VALUES
  (1, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'Allied Health', ''),
  (2, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'Nurse', ''),
  (3, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'Surgeon', ''),
  (4, 1, '2013-06-12 16:27:22', '2013-06-12 16:27:22', 'Advanced Trainee', ''),
  (5, 1, '2013-07-15 14:19:37', '2013-07-15 14:19:37', 'Retired / Comp', 'Complimentary subscription');

-- name: insert-data-ms_m_application
INSERT INTO `%s`.`ms_m_application` VALUES
  (1, 502, 388, 440, 2, NULL, 1, '2015-09-01 04:30:30', '2015-09-01 04:44:01', '2015-09-01', 1, 'all good!'),
  (2, 482, 72, 440, 2, 1, 1, '2015-09-01 06:28:15', '2015-09-01 06:43:22', '2015-09-01', 1, NULL),
  (3, 488, 389, 456, 2, 1, 1, '2015-09-01 06:29:00', '2015-09-01 06:43:23', '2015-09-01', 1, NULL),
  (4, 499, 423, 195, 4, 2, 1, '2015-09-01 06:29:56', '2015-09-01 06:45:34', '2015-09-01', 1, NULL),
  (5, 485, 168, 317, 4, 2, 1, '2015-09-01 06:31:16', '2015-09-01 06:45:34', '2015-09-01', 1, NULL);

-- name: insert-data-ms_m_application_meeting
INSERT INTO `%s`.`ms_m_application_meeting` VALUES
  (1, 1, NULL, 1, 1, '2015-09-01 04:32:16', '2015-09-01 04:32:16', -1, NULL);

-- insert-data-ms_m_permission

-- name: insert-data-ms_m_status
INSERT INTO `%s`.`ms_m_status` VALUES (1, 1, 1, 0, 1, 0, '2015-08-30 17:10:57', '0000-00-00 00:00:00', NULL),
  (2, 2, 1, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', NULL),
  (3, 3, 1, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', NULL),
  (4, 4, 1, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', NULL),
  (5, 5, 1, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', NULL),
  (6, 6, 1, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', NULL),
  (7, 7, 1, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', NULL),
  (8, 8, 1, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', NULL),
  (9, 9, 1, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', NULL);

-- name: insert-data-ms_m_title
INSERT INTO `%s`.`ms_m_title` VALUES
  (1, 1, 2, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', '2015-08-31',
   'Test data does will not show historic titles'),
  (2, 2, 2, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', '2015-08-31',
   'Test data does will not show historic titles'),
  (3, 3, 2, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', '2015-08-31',
   'Test data does will not show historic titles'),
  (4, 4, 2, NULL, 1, 1, '2015-08-30 17:10:57', '0000-00-00 00:00:00', '2015-08-31',
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
  (10008, 1, 0, '0000-00-00 00:00:00', '2013-06-18 15:12:11', 0, 0, 0, 0, 0, 'Deceased', '');

-- name: insert-data-ms_title
INSERT INTO `%s`.`ms_title`
VALUES (1, 0, 1, '2013-01-17 23:17:04', '2013-07-08 19:44:45', 0, 0, 0, 0, 0, 'Applicant', ''),
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
  (3, 1, 1, 1, 1, NOW(), NOW(), 'ABC-1', 'ABC Sub1', '', '', '', '', '', '', '', '', '', ''),
  (4, 1, 1, 1, 1, NOW(), NOW(), 'ABC-2', 'ABC Sub2', '', '', '', '', '', '', '', '', '', ''),
  (5, 1, 1, 1, 1, NOW(), NOW(), 'ABC-3', 'ABC Sub3', '', '', '', '', '', '', '', '', '', '');


-- name: insert-data-organisation_type
INSERT INTO `%s`.`organisation_type`
VALUES (1, 1, '2013-06-11 20:45:11', '2013-06-11 20:45:31', 'CSANZ Council', 'Group / type for CSANZ councils'),
  (2, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'CSANZ Working Group', 'CSANZ Working Group'),
  (3, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'Committee', 'Committee'),
  (4, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'Sub-committee', 'Sub-committee'),
  (5, 1, '0000-00-00 00:00:00', '2013-08-19 15:57:34', 'CSANZ Board', 'Group / type for CSANZ Board'),
  (6, 1, '0000-00-00 00:00:00', '0000-00-00 00:00:00', 'Other', 'Other'),
  (7, 1, '2013-07-14 15:37:02', '2013-07-14 15:37:13', 'Institute / Hospital', ''),
  (8, 1, '2013-07-14 17:32:31', '2013-07-14 17:32:31', 'University / Education', ''),
  (9, 1, '2014-07-23 12:00:14', '2014-07-23 12:00:14', 'Heart Foundation', 'Heart Foundation of Australia'),
  (10, 1, '2014-08-04 19:21:08', '2014-08-04 19:21:08', 'Government', '');

-- insert-data-wf_attachment

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
  (1, 4, NULL, 1, 1, 1, 1, '2013-07-04 14:37:45', '2013-11-27 17:53:14', 'Invoice Raised',
      'A new invoice has been raised and is pending payment.',
   'Members can pay online or by alternate methods specified on the invoice.', NULL),
  (2, 4, NULL, 1, 1, 1, 1, '2013-07-04 14:38:26', '2013-11-27 17:54:15', 'Invoice Past Due',
      'Invoice is past due and pending payment.', 'Please pay online or by alternate methods specified on the invoice.',
   NULL), (3, 5, NULL, 1, 1, 1, 0, '2013-09-03 16:32:00', '2013-09-09 14:16:40', 'Missing Primary Email',
              'All active members require a primary email to access the system and for communications.',
           'Add a primary email address to member profile.', NULL),
  (4, 5, NULL, 1, 1, 1, 0, '2013-09-03 16:34:11', '2013-09-09 14:15:49', 'Missing Membership Subscription',
      'All active members require a membership subscription, even if that subscription is complimentary.',
   'Assign the appropriate membership subscription.', NULL),
  (5, 5, NULL, 1, 1, 1, 0, '2013-09-03 16:36:54', '2013-09-09 14:18:07', 'Empty Contact Card',
      'All active members require at least one field to be filled in for each persistent contact card.',
   'Edit member contact information and provide at least one value (or \'n/a\') in each of the persistent contact cards.',
   NULL), (6, 4, NULL, 1, 1, 1, 1, '2013-11-07 11:48:01', '2013-11-27 17:55:38', 'Invoice Final Notice',
              'Invoice is past due and a final notice has been issued by email.',
           'Please note no further notices will be emailed.  If your subscription is not paid immediately, your Membership of the Society will be cancelled.',
           NULL), (8, 1, NULL, 1, 1, 0, 0, '2015-04-08 22:38:32', '2015-04-08 22:38:32', 'Email Communication Failure',
                      'A recent email communication failed for some reason. ',
                   'Check the specific messages in the Member\'s communication tab for clues as to the appropriate follow up.',
                   NULL), (9, 4, NULL, 1, 1, 0, 0, '2016-03-20 16:13:08', '2016-03-20 16:13:08', 'Invoice Overpaid',
                              'Total of payments allocated to invoice exceeds the invoice total. ',
                           'Require manual intervention to remove payment allocations as well as refund if applicable.',
                           NULL),
  (10000, 1, NULL, 1, 0, 0, 0, '2013-09-10 21:06:29', '2013-09-11 15:53:12', 'General Admin', '-', '-', NULL);

-- name: insert-data-wf_note
INSERT INTO `%s`.`wf_note` VALUES
  (1, 10001, 1, 1, 1, '2015-09-01 04:32:49', '2015-09-01 04:32:49', '2015-09-01', 'Application scanned and attached.'),
  (2, 1, NULL, NULL, 1, '2015-09-01 06:13:27', '2015-09-01 06:13:27', '2015-09-01', 'Issue raised.'),
  (3, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised.'),
  (4, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised.'),
  (5, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised.'),
  (6, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised.'),
  (7, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised.'),
  (8, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised.'),
  (9, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised.'),
  (10, 1, NULL, NULL, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', '2015-09-01', 'Issue raised.');

-- name: insert-data-wf_note_association
INSERT INTO `%s`.`wf_note_association`
VALUES (1, 1, 1, 1, 1, '2015-09-01 04:32:49', '2015-09-01 04:32:49', 'application'),
  (2, 2, 1, 1, 1, '2015-09-01 06:13:27', '2015-09-01 06:13:27', 'issue'),
  (3, 3, 1, 2, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', 'issue'),
  (4, 4, 2, 3, 1, '2015-09-01 06:35:15', '2015-09-01 06:35:15', 'issue');

-- name: insert-data-wf_note_type
INSERT INTO `%s`.`wf_note_type`
VALUES (1, 1, 1, '2013-11-07 11:53:15', '2013-11-07 11:53:15', 'System', 'Note added by system housekeeping'),
  (10001, 1, 0, '2013-05-06 06:18:02', '2013-06-26 16:43:21', 'General', ''),
  (10002, 1, 0, '2013-05-06 06:18:02', '2013-06-26 16:43:21', 'Account', ''),
  (10003, 1, 0, '2013-06-26 16:43:33', '2013-06-26 16:43:33', 'Contact', ''),
  (10004, 1, 0, '2013-06-26 20:22:04', '2013-06-26 20:22:04', 'Prize / Award', '');

