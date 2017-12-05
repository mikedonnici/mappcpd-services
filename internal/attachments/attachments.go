/*
	Package attachments handles uploading of files to cloud storage and registration of the files in the database.
*/
package attachments

import (
	"fmt"
	"os"
	"time"

	"database/sql"

	"github.com/34South/envr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"

	"github.com/mappcpd/web-services/internal/fileset"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"strconv"
)

// Attachment contains data about an uploaded file (attachment)- that is
// its location in cloud storage and the entity with which it is associated.
// The relevant fields will depend on the type of upload, and member or admin.
type Attachment struct {

	// ID is the identifier of the attachment record
	ID int64 `json:"id"`

	// EntityID represents is id of the record with which the attachment is associated
	EntityID int64 `json:"entityId"`

	// UserID is a stored with attachment records when they are added by an admin user.
	UserID int64 `json:"userId""`

	// Clean filename is a sanitised version of the original filename
	CleanFilename string `json:"cleanFilename"`

	// Cloudy Filename is an obfuscated MD5 of the original filename and, when present, the file is stored with this name
	CloudyFilename string `json:"cloudyFilename"`

	// URL is the public URL to access this attachment
	// todo can we use signed urls to request access?
	URL string `json:"url"`

	// An Attachment always has an associated FileSet which represents how it is stored
	FileSet fileset.FileSet
}

// attachmentQueries stores relevant sql for a particular attachment type.
type attachmentQueries struct {
	exists   string
	register string
}

func init() {
	envr.New("mappcpd-attachments", []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_REGION",
	}).Auto()
}

// New returns a pointer to an Attachment value.
func New() *Attachment {
	return &Attachment{}
}

// Validate checks the Attachment for the required values prior to registration
func (a *Attachment) Validate() error {

	// First, and most important value is the FileSet, as this determines some of the other requirements.
	if err := a.FileSet.CheckFields(); err != nil {
		return err
	}

	// For an Attachment to be registered also need EntityID and CleanFilename
	if a.EntityID == 0 {
		return errors.New("Attachment.EntityID (int) has a zero value")
	}
	if a.CleanFilename == "" {
		return errors.New("Attachment.CleanFilename is an empty string")
	}

	// These depend on the entity
	switch a.FileSet.Entity {
	case "ce_m_activity_attachment":
		// activity attachments need a CloudyFilename
		if a.CloudyFilename == "" {
			return errors.New("Attachment.CloudyFilename is an empty string - value required for activity attachments")
		}
	case "wf_attachment":
		// UserID for note attachments to identify the admin user
		if a.UserID == 0 {
			return errors.New("Attachment.UserID has a zero value - admin (user) ID is required for note attachments")
		}
	}

	return nil
}

// Register records data about an uploaded file in the database. The record is inserted into a table that corresponds
// to the entity with which the file is being associated. There are a few differences between the tables, hence the switch.
// If the record already exists it cannot be registered again, however this is not really an error so Register()
// can just populate the fields and carry on. From the caller's perspective this makes no difference - there is an uploaded
// file and details about it are being recorded in the database. HOWEVER, what if the user ID is changes? In this case we could
// force an update of the record. For now it will populate the fields which makes this operation IDEMPOTENT.
// 'flags' is a hack to pass in an optional setting - at this stage just to set thumbnail = 1 for a resource file.
func (a *Attachment) Register(flags ...string) error {

	// Validate first
	if err := a.Validate(); err != nil {
		return errors.New("Attachment validation error - " + err.Error())
	}

	// Check if already registered, if so we will set the ID and URL and return without error
	if err := a.Exists(); err != nil {
		return errors.New("Error checking for existing registration - " + err.Error())
	}
	if a.ID > 0 { // ID was set so it DOES exist
		if err := a.setURL(); err != nil {
			return errors.New("Error setting URL for attachment - " + err.Error())
		}
		return nil
	}

	// If we're here the attachment is NOT already registered, so register it
	var query string

	switch a.FileSet.Entity {
	case "ce_m_activity_attachment":
		query = `INSERT INTO ce_m_activity_attachment ` +
			`(ce_m_activity_id, fs_set_id, active, created_at, updated_at, clean_filename, cloudy_filename) ` +
			`VALUES (%d, %d, 1, NOW(), NOW(), "%s", "%s")`
		query = fmt.Sprintf(query, a.EntityID, a.FileSet.ID, a.CleanFilename, a.CloudyFilename)

	case "wf_attachment":
		query = `INSERT INTO wf_attachment ` +
			`(wf_note_id, ad_user_id, fs_set_id, active, created_at, updated_at, clean_filename) ` +
			`VALUES (%d, %d, %d, 1, NOW(), NOW(), "%s")`
		query = fmt.Sprintf(query, a.EntityID, a.UserID, a.FileSet.ID, a.CleanFilename)

	case "ol_resource_file":
		var thumbnail int
		if flags[0] == "thumbnail" {
			thumbnail = 1
		}
		query = `INSERT INTO ol_resource_file ` +
			`(ol_resource_id, ad_user_id, fs_set_id, active, thumbnail, created_at, updated_at, clean_filename, cloudy_filename) ` +
			`VALUES (%d, %d, %d, 1, %d, NOW(), NOW(), "%s", "%s")`
		query = fmt.Sprintf(query, a.EntityID, a.UserID, a.FileSet.ID, thumbnail, a.CleanFilename, a.CloudyFilename)

	default:
		return errors.New("Error registering attachment - unknown entity name")
	}

	result, err := datastore.MySQL.Session.Exec(query)
	if err != nil {
		return errors.New("Database error - " + err.Error())
	}

	// Set the ID
	id, err := result.LastInsertId()
	if err != nil {
		return errors.New("Error getting last insert id - " + err.Error())
	}
	a.ID = int64(id)

	return nil
}

// Exists checks for an existing record for the attachment so we can prevent duplicate registrations.
// Note that it returns an error, except for sql.ErrorNoRows, which indicates that the record does not exist.
// If the record is found then the fields are populated.
func (a *Attachment) Exists() error {

	var query string
	var id int

	switch a.FileSet.Entity {
	case "ce_m_activity_attachment":
		query = `SELECT id FROM ce_m_activity_attachment WHERE active = 1 AND ` +
			`ce_m_activity_id = %d AND fs_set_id = %d AND clean_filename = "%s" AND cloudy_filename = "%s" ` +
			`LIMIT 1`
		query = fmt.Sprintf(query, a.EntityID, a.FileSet.ID, a.CleanFilename, a.CloudyFilename)

	case "wf_attachment":
		query = `SELECT id FROM wf_attachment WHERE active = 1 AND ` +
			`wf_note_id = %d AND fs_set_id = %d AND clean_filename = "%s" ` +
			`LIMIT 1`
		query = fmt.Sprintf(query, a.EntityID, a.FileSet.ID, a.CleanFilename)

	case "ol_resource_file":
		query = `SELECT id FROM ol_resource_file WHERE active = 1 AND ` +
			`ol_resource_id = %d AND fs_set_id = %d AND clean_filename = "%s" AND cloudy_filename = "%s" ` +
			`LIMIT 1`
		query = fmt.Sprintf(query, a.EntityID, a.FileSet.ID, a.CleanFilename, a.CloudyFilename)

	default:
		return errors.New("Unknown entity: " + a.FileSet.Entity)
	}

	err := datastore.MySQL.Session.QueryRow(query).Scan(&id)
	// No rows is not an error here
	if err == sql.ErrNoRows {
		return nil
	}
	// some other error
	if err != nil {
		return err
	}

	// Record exists so set ID
	a.ID = int64(id)
	return nil
}

// setURL sets the public URL for an attachment by looking up fs_url record
func (a *Attachment) setURL() error {

	var url string

	query := "SELECT base_url FROM fs_url WHERE active = 1 AND fs_set_id = %d ORDER BY priority ASC LIMIT 1"
	query = fmt.Sprintf(query, a.FileSet.ID)
	err := datastore.MySQL.Session.QueryRow(query).Scan(&url)
	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("No fs_url record found for file_set.id = %d - %s", a.FileSet.ID, err.Error())
		return errors.New(msg)
	}
	if err != nil {
		msg := fmt.Sprintf("Error looking up base url for file_set.id = %d - %s", a.FileSet.ID, err.Error())
		return errors.New(msg)
	}

	a.URL = url + a.FileSet.Path + strconv.FormatInt(a.EntityID, 10) + "/" + a.CloudyFilename
	return nil
}

// S3PutRequest issues a signed URL that allows for a PUT to an Amazon S3 bucket. It receives the
// key (full path to file including file name', and the name of the bucket.
// The aws package ASSUMES the presence of AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY env vars so they have been
// added to init() above. AWS_REGION was added by me.
func S3PutRequest(key, bucket string) (string, error) {

	sess := session.Must(session.NewSession())
	svc := s3.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		//Body:   strings.NewReader("EXPECTED CONTENTS"),
	})

	return req.Presign(15 * time.Minute)
}
