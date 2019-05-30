package note

// Queries is a map containing common queries for the package
var queries = map[string]string{
	"select-note":               selectNote,
	"select-note-by-id":         selectNoteByID,
	"select-notes-by-member-id": selectNotesByMemberID,
	"select-attachments":        selectAttachments,
	"insert-note":               insertNote,
	"insert-note-association":   insertNoteAssociation,
}

const selectNote = `
SELECT
  wn.id                         AS ID,
  wnt.name                      AS Type,
  wnt.id                        AS TypeID,
  COALESCE(wna.association, '') AS Association,
  wna.association_entity_id     AS AssociationID,
  m.id                          AS MemberID,
  wn.created_at                 AS CreatedAt,
  wn.updated_at                 AS UpdatedAt,
  wn.effective_on               AS Date,
  wn.note                       AS Note
FROM wf_note wn
  LEFT JOIN wf_note_type wnt ON wn.wf_note_type_id = wnt.id
  LEFT JOIN wf_note_association wna ON wn.id = wna.wf_note_id
  LEFT JOIN member m ON wna.member_id = m.id
WHERE wn.active = 1 `

const selectNoteByID = selectNote + ` AND wn.id = ?`

const selectNotesByMemberID = selectNote + ` AND m.id = ?`

const selectAttachments = `
SELECT
  wa.id AS ID,
  wa.clean_filename AS FileName,
  CONCAT(u.base_url, s.set_path, wa.wf_note_id, "/", wa.id, "-", wa.clean_filename) AS FileUrl
FROM wf_attachment wa
  LEFT JOIN fs_set s ON wa.fs_set_id = s.id
  LEFT JOIN fs_url u ON s.id = u.fs_set_id`

const insertNote = `
INSERT INTO wf_note (
	wf_note_type_id, 
	updated_at, 
	effective_on, 
	note 
) VALUES (?, NOW(), NOW(), ?)`

const insertNoteAssociation = `
INSERT INTO wf_note_association (
	wf_note_id, 
	member_id, 
	association_entity_id, 
	updated_at, 
	association
) VALUES (?, ?, ?, NOW(), ?)`
