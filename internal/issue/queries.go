package issue

var queries = map[string]string{
	"insert-issue":             insertIssue,
	"insert-issue-association": insertIssueAssociation,
	"select-issue-by-id":       selectIssueByID,
	"select-issue-type-by-id": selectIssueTypeByID,
}

const insertIssue = `
INSERT INTO wf_issue (
	wf_issue_type_id, 
	updated_at, 
	live_on, 
	description, 
	required_action
) VALUES (%d, NOW(), NOW(), %q, %q)`

const insertIssueAssociation = `
INSERT INTO wf_issue_association (
	wf_issue_id, 
	member_id, 
	association_entity_id, 
	updated_at, 
	association
) VALUES (%d, %d, %d, NOW(), %q)`

const selectIssue = `
SELECT 
    i.id AS IssueID,
    i.resolved AS IssueResolved,
    i.member_visible AS VisibleToMember,
    i.description AS Description,
    i.required_action AS Action,
    ia.member_id AS MemberID,
    ia.association AS AssocEntity,
    ia.association_entity_id AS AssocEntityID,
    it.id AS IssueTypeID,
	it.name AS IssueType,
	it.Description as IssueTypeDescription,
	ic.id as IssueCategoryID,
	ic.name AS IssueCategory,
	ic.description AS IssueCategoryDescription
FROM
    wf_issue i
        LEFT JOIN
    wf_issue_type it ON i.wf_issue_type_id = it.id
        LEFT JOIN
    wf_issue_category ic ON it.wf_issue_category_id = ic.id
        LEFT JOIN
	wf_issue_association ia ON i.id = ia.wf_issue_id
WHERE 1`

const selectIssueByID = selectIssue + ` AND i.id = ?`

const selectIssueType = `
SELECT 
    t.id AS TypeID,
    t.name AS Type,
    t.description AS Description,
    t.required_action AS Action,
    c.id AS CategoryID,
    c.name AS CategoryName,
    c.description AS CategoryDescription
FROM
    wf_issue_type t
LEFT JOIN wf_issue_category c ON t.wf_issue_category_id = c.id    
WHERE
	t.active = 1`

const selectIssueTypeByID = selectIssueType + ` AND t.id = ?`

