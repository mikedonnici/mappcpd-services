package activity

// Queries is a map containing common queries for the package
var Queries = map[string]string {
  "select-activities": selectActivities,
}

const selectActivities = `SELECT
  a.id                      AS ActivityID,
  a.code                    AS ActivityCode,
  a.name                    AS ActivityName,
  a.description             AS ActivityDescription,
  a.ce_activity_category_id AS ActivityCategoryID,
  c.name                    AS ActivityCategoryName,
  a.ce_activity_unit_id     AS ActivityUnitID,
  u.name                    AS ActivityUnitName,
  a.points_per_unit         AS CreditPerUnit,
  a.annual_points_cap       AS MaxCredit
FROM
  ce_activity a
  LEFT JOIN
  ce_activity_category c ON a.ce_activity_category_id = c.id
  LEFT JOIN
  ce_activity_unit u ON a.ce_activity_unit_id = u.id`

