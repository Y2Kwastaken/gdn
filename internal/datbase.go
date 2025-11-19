package internal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

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

	_, err = conn.Exec(query, imageIdBytes, metadata.Title, metadata.ImageType, metadata.Description)
	if err != nil {
		return nil, err
	}

	trsn, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	query = `INSERT INTO image_tags (id, tag) VALUES ( ?, ? )`
	for i := range metadata.Tags {
		tag := metadata.Tags[i]
		println(tag)
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
	query := `SELECT id, image_name, image_type, description FROM image_meta WHERE id = ?`

	inBytes, err := inUUID.MarshalBinary()
	if err != nil {
		return nil, err
	}

	row := conn.QueryRow(query, inBytes)

	meta := &ImageMeta{}
	var uuidBlob []byte // this read is useless, but idk if I can just dev/null with scanner
	err = row.Scan(&uuidBlob, &meta.ImageName, &meta.ImageType, &meta.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	meta.Id, err = uuid.FromBytes(uuidBlob[:])
	if err != nil {
		return nil, err
	}

	query = `SELECT tag FROM image_tags WHERE id = ?`
	rows, err := conn.Query(query, uuidBlob)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
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

func (db *Database) QueryIds(limit int, offset int) ([]uuid.UUID, error) {
	if limit <= 0 || offset < 0 {
		return nil, fmt.Errorf("limit or offset out of bounds limit: %d, offset: %d", limit, offset)
	}

	conn := db.conn
	query := `SELECT id FROM image_meta LIMIT ? OFFSET ?`

	rows, err := conn.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uuids []uuid.UUID
	for rows.Next() {
		var tmp []byte
		if err := rows.Scan(&tmp); err != nil {
			return nil, err
		}

		uuid, err := uuid.FromBytes(tmp)
		if err != nil {
			return nil, err
		}

		uuids = append(uuids, uuid)
	}

	return uuids, nil
}

func (db *Database) CountEntries() (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM image_meta").Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}
