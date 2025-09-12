package internal

import (
	"database/sql"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type ImageMeta struct {
	Id          uuid.UUID
	ImageName   string
	ImageType   string
	Description string
	Tags        []string
}

type Database struct {
	conn *sql.DB
}

func NewDBConnection(constr string) (*Database, error) {
	conn, err := sql.Open("sqlite", constr)
	if err != nil {
		return nil, err
	}

	return &Database{conn: conn}, nil
}

func (db *Database) SetupTables() error {
	conn := db.conn
	query := `CREATE TABLE IF NOT EXISTS image_meta (
		id BLOB PRIMARY KEY,
		image_name TEXT NOT NULL,
		image_type TEXT NOT NULL,
		description TEXT

	)`

	_, err := conn.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS image_tags (
		id BLOB,
		tag TEXT
	)`

	_, err = conn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UploadImageMeta(metadata *Metadata) (*uuid.UUID, error) {
	conn := db.conn
	imageId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO image_meta (id , image_name, image_type, description) VALUES( ?, ?, ?, ? )`
	imageIdBytes, err := imageId.MarshalBinary()
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(query, imageIdBytes, metadata.Name, metadata.ImageType, metadata.Description)
	if err != nil {
		return nil, err
	}

	trsn, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	query = `INSERT INTO image_tags (id, tag) VALUES ( ?, ? )`
	for tag := range metadata.Tags {
		_, err = trsn.Exec(query, imageIdBytes, tag)
		if err != nil {
			return nil, err
		}
	}

	err = trsn.Commit()
	if err != nil {
		return nil, err
	}

	return &imageId, nil
}

func (db *Database) DeleteImage(uuid uuid.UUID) error {
	conn := db.conn

	uuidBytes, err := uuid.MarshalBinary()
	if err != nil {
		return err
	}

	query := `DELETE FROM image_meta WHERE id = ?`
	_, err = conn.Exec(query, uuidBytes)
	if err != nil {
		return err
	}

	query = `DELETE FROM image_tags where id = ?`
	_, err = conn.Exec(query, uuidBytes)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) QueryImage(inUUID uuid.UUID) (*ImageMeta, error) {
	conn := db.conn
	query := `SELECT * FROM image_meta WHERE id = ?`

	inBytes, err := inUUID.MarshalBinary()
	if err != nil {
		return nil, err
	}

	row := conn.QueryRow(query, inBytes)
	meta := &ImageMeta{}
	var uuidBlob [16]byte // this read is useless, but idk if I can just dev/null with scanner
	err = row.Scan(&uuidBlob, &meta.ImageName, &meta.ImageType, &meta.Description)
	if err != nil {
		return nil, err
	}

	meta.Id, err = uuid.FromBytes(uuidBlob[:])
	if err != nil {
		return nil, err
	}

	query = `SELECT * FROM image_tags WHERE id = ?`
	rows, err := conn.Query(query, inBytes)
	if err != nil {
		return nil, err
	}

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(nil, &tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	meta.Tags = tags
	return meta, nil
}
