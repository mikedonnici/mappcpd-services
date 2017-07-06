package attachments

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"fmt"
	"github.com/34South/envr"
	"github.com/mappcpd/web-services/internal/platform/datastore"
	"github.com/pkg/errors"
	"os"
	"database/sql"
)

// Attachment contains data about an uploaded file (attachment)- that is
// its location in cloud storage and the entity with which it is associated.
type Attachment struct {
	FileSetID      int    `json:"fileSetId"`
	EntityID       int    `json:"entityId"`
	EntityName     string `json:"entityName"`
	CleanFilename  string `json:"cleanFilename"`
	CloudyFilename string `json:"cloudyFilename"`
}

func init() {
	envr.New("mappcpd-attachments", []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_REGION",
	}).Auto()
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

// Register records data about an uploaded file in the database. The record is inserted into a table that corresponds
// to the entity with which the file is being associated. There are a few differences between the tables, hence the switch.
func (a Attachment) Register() error {

	// Check for duplication
	id, err := a.Exists()
	if err != nil  {
		fmt.Println(sql.ErrNoRows)
		return errors.New(".Register() " + err.Error())
	}
	if id > 0 {
		msg := fmt.Sprintf(".Register() - an attachment for this file appears to already exist: %s.id = %d", a.EntityName, id)
		return errors.New(msg)
	}

	var query string

	switch a.EntityName {
	case "ce_m_activity_attachment":
		query = `INSERT INTO ce_m_activity_attachment ` +
			`(ce_m_activity_id, fs_set_id, active, created_at, updated_at, clean_filename, cloudy_filename) ` +
			`VALUES (%d, %d, 1, NOW(), NOW(), "%s", "%s")`
		query = fmt.Sprintf(query, a.EntityID, a.FileSetID, a.CleanFilename, a.CloudyFilename)
	default:
		return errors.New("Unknown entity name provided, cannot register the attachment")
	}

	fmt.Println("Registering the attachment with: ", query)
	_, err = datastore.MySQL.Session.Exec(query)
	if err != nil {
		return errors.New("Database error - " + err.Error())
	}

	return nil
}

// registerActivityAttachment registers an attachment for CPD activity, ie in the ce_m_activity_attachment table.
// It first checks if there is an existing record for the same file name, and if so, will do an update instead.
// Otherwise it will insert a new record
func (a Attachment) registerActivityAttachment() error {
	return nil
}

// Exists checks for an existing record for the attachment so we can prevent duplicate registrations.
func (a Attachment) Exists() (int, error) {

	var query string
	var id int

	switch a.EntityName {
	case "ce_m_activity_attachment":
		query = `SELECT id FROM ce_m_activity_attachment WHERE active = 1 AND ` +
			`ce_m_activity_id = %d AND fs_set_id = %d AND clean_filename = "%s" AND cloudy_filename = "%s" ` +
			`LIMIT 1`
		query = fmt.Sprintf(query, a.EntityID, a.FileSetID, a.CleanFilename, a.CloudyFilename)
	default:
		return id, errors.New(".Exists() - unknown entity name provided")
	}

	fmt.Println("Checking for an existing attachment record with: ", query)
	err := datastore.MySQL.Session.QueryRow(query).Scan(&id)
	// No rows is fine... that is what we are hoping for!
	if err != nil && err != sql.ErrNoRows {
		return id, errors.New("Database error - " + err.Error())
	}

	return id, nil
}

// GetFileSetID retrieves the current fs_set.id value that is current for a particular entity
func GetFileSetID(entityName string) (int, error) {

	var id int
	query := "SELECT id FROM fs_set WHERE active = 1 AND current = 1 AND entity_name = ?"
	err := datastore.MySQL.Session.QueryRow(query, entityName).Scan(&id)

	return id, err
}
