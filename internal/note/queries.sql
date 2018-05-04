-- name: select-note
SELECT
  wn.id           AS ID,
  wnt.name        AS Type,
  m.id            AS MemberID,
  wn.created_at   AS CreatedAt,
  wn.updated_at   AS UpdatedAt,
  wn.effective_on AS Date,
  wn.note         AS Note
FROM wf_note wn
  LEFT JOIN wf_note_type wnt ON wn.wf_note_type_id = wnt.id
  LEFT JOIN wf_note_association wna ON wn.id = wna.wf_note_id
  LEFT JOIN member m ON wna.member_id = m.id


-- name: select-attachments
SELECT
  wa.id AS ID,
  wa.clean_filename AS FileName,
  CONCAT(u.base_url, s.set_path, wa.wf_note_id, "/", wa.id, "-", wa.clean_filename) AS FileUrl
FROM wf_attachment wa
  LEFT JOIN fs_set s ON wa.fs_set_id = s.id
  LEFT JOIN fs_url u ON s.id = u.fs_set_id
