package members

const selectMemberActivityQuery =
	`SELECT
    	cma.id AS 'memberActivityId',
    	cma.member_id AS 'memberId',
    	cma.activity_on AS 'memberActivityDate',
    	COALESCE(cma.description, '') AS 'memberActivityDescription',
    	(cma.quantity * cma.points_per_unit) AS 'activityCredit',
    	cma.quantity AS 'quantity',
    	COALESCE(cau.name, '') AS 'unit',
    	cma.points_per_unit AS 'creditPerUnit',
    	cac.id AS 'categoryId',
    	COALESCE(cac.name, '') AS 'categoryName',
    	COALESCE(cac.description, '') AS 'categoryDescription',
    	ca.id AS 'activityId',
    	COALESCE(ca.code, '') AS 'activityCode',
    	COALESCE(ca.name, '') AS 'activityName',
    	COALESCE(ca.description, '') AS 'activityDescription',
    	cat.id AS 'typeId',
    	COALESCE(cat.name, '') AS 'typeName'
	FROM
    	ce_m_activity cma
        	LEFT JOIN
    	ce_activity ca ON cma.ce_activity_id = ca.id
        	LEFT JOIN
    	ce_activity_unit cau ON ca.ce_activity_unit_id = cau.id
        	LEFT JOIN
    	ce_activity_category cac ON ca.ce_activity_category_id = cac.id
        	LEFT JOIN
    	ce_activity_type cat ON cma.ce_activity_type_id = cat.id`
