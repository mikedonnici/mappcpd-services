// Packages notes provides access to Notes data
package notes

import "github.com/mappcpd/web-services/internal/platform/datastore"

// Note is just that - a note recorded in the system which may be linked to a member,
// and more specifically to another entity such as an application, or an issue
type Note struct {
	ID            int          `json:"id" bson:"id"`
	MemberID      int          `json:"memberId" bson:"memberId"`
	DateCreated   string       `json:"dateCreated" bson:"dateCreated"`
	DateUpdated   string       `json:"dateUpdated" bson:"dateUpdated"`
	DateEffective string       `json:"dateEffective" bson:"dateEffective"`
	Content       string       `json:"content" bson:"content"`
	Attachments   []Attachment `json:"attachments" bson:"attachments"`
}

type Notes []Note

// Attachment is a file linked to a note
type Attachment struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

// Attachments fetches the attachments for a Note
func (n *Note) SetAttachments() error {

	query := `SELECT
	wa.id,
	wa.clean_filename,
	CONCAT(u.base_url, s.set_path, wa.wf_note_id, "/", wa.id, "-", wa.clean_filename)
	FROM wf_attachment wa
	LEFT JOIN fs_set s ON wa.fs_set_id = s.id
	LEFT JOIN fs_url u ON s.id = u.fs_set_id
	WHERE wa.wf_note_id = ?`

	rows, err := datastore.MySQL.Session.Query(query, n.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {

		a := Attachment{}
		rows.Scan(
			&a.ID,
			&a.Name,
			&a.URL,
		)

		n.Attachments = append(n.Attachments, a)
	}

	return nil
}

// NoteByID fetches a Note from the default datastore
func NoteByID(id int) (Note, error) {
	return noteByID(id, datastore.MySQL)
}

// NoteByIDStore fetches a Note from the specified datastore - used for testing
func NoteByIDStore(id int, conn datastore.MySQLConnection) (Note, error) {
	return noteByID(id, conn)
}

// noteByID fetches a Note record from the specified data store
func noteByID(id int, connection datastore.MySQLConnection) (Note, error) {

	// Create Note value
	n := Note{ID: id}

	// Coalesce any NULL-able fields
	query := `SELECT m.id, wn.created_at, wn.updated_at, wn.effective_on, wn.note
		FROM wf_note wn
		LEFT JOIN wf_note_association wna ON wn.id = wna.wf_note_id
		LEFT JOIN member m ON wna.member_id = m.id
		WHERE wn.id = ?`

	err := connection.Session.QueryRow(query, id).Scan(
		&n.MemberID,
		&n.DateCreated,
		&n.DateUpdated,
		&n.DateEffective,
		&n.Content,
	)
	if err != nil {
		return n, err
	}

	// Add attachments
	//err = n.SetAttachments()

	return n, err
}

// MemberNotes fetches all the notes linked to a Member
func MemberNotes(memberID int) (*Notes, error) {

	ns := Notes{}

	// Coalesce any NULL-able fields
	query := `SELECT
		wn.id,
		m.id,
		wn.created_at,
		wn.updated_at,
		wn.effective_on,
		wn.note
		FROM wf_note wn
		LEFT JOIN wf_note_association wna ON wn.id = wna.wf_note_id
		LEFT JOIN member m ON wna.member_id = m.id
		WHERE m.id = ?
		ORDER BY wn.effective_on DESC`

	rows, err := datastore.MySQL.Session.Query(query, memberID)
	if err != nil {
		return &ns, err
	}
	defer rows.Close()

	for rows.Next() {

		n := Note{}
		rows.Scan(
			&n.ID,
			&n.MemberID,
			&n.DateCreated,
			&n.DateUpdated,
			&n.DateEffective,
			&n.Content,
		)

		// Attachments, if any
		n.SetAttachments()

		ns = append(ns, n)
	}

	return &ns, nil
}
