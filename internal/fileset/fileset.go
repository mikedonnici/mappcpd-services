/*
	Package fileset provides information about stored file resources
*/
package fileset

import (
	"database/sql"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/pkg/errors"
)

// FileSet represent a row from the fs_set table and describes a "set" of related files in cloud storage.
// ID - the fs_set.id value
// VolumeName - the name of the bucket
// SetPath - the path or pseudo path (S3) excluding the file name, eg '/cpd/', '/note/'
// EntityName - name of the db table containing the records to which files in this set are 'attached'
type FileSet struct {
	ID         int    `json:"id"`
	Volume string `json:"volume"`
	Path    string `json:"Path"`
	Entity string `json:"entity"`
}

// New fetches returns a pointer to a FileSet with values for the current file set for the given entity (e).
func New(e string) (*FileSet, error) {

	var fs FileSet

	query := "SELECT id, volume_name, set_path, entity_name FROM fs_set WHERE " +
		"active = 1 AND current = 1 AND entity_name = ?"
	err := datastore.MySQL.Session.QueryRow(query, e).Scan(
		&fs.ID,
		&fs.Volume,
		&fs.Path,
		&fs.Entity,
	)
	if err == sql.ErrNoRows {
		return &fs, errors.New("No file set found for entity " + e)
	}
	if err != nil {
		return &fs, errors.New("New() database error - " + err.Error())
	}

	return &fs, nil
}
