package app

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

const T_WEAK = 0
const T_STRONG = 1

type Entry struct {
	UID       int
	FirstName string
	LastName  string
}

type EntryRepository struct {
	db *sql.DB
}

func (r *EntryRepository) CreateOrUpdate(e Entry) error {
	const fn = "app.repository.CreateOrUpdate"

	stmt, err := r.db.Prepare(`
		INSERT INTO sdn_entries (uid, first_name, last_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (uid) DO UPDATE
		SET first_name = excluded.first_name,
		    last_name = excluded.last_name,
			updated_at = current_timestamp;
	`)
	if err != nil {
		return fmt.Errorf("%s :%w", fn, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.UID, e.FirstName, e.LastName)

	// fmt.Printf("from chanel: %#v, err=%#v\n", e, err)

	return fmt.Errorf("%s :%w", fn, err)
}

func (r *EntryRepository) ClearAll() error {
	const fn = "app.repository.ClearAll"

	_, err := r.db.Exec(`TRUNCATE TABLE sdn_entries;`)
	if err != nil {
		return fmt.Errorf("%s :%w", fn, err)
	}
	return nil
}

func (r *EntryRepository) HasRecords() (bool, error) {
	const fn = "app.repository.hasRecords"

	var exists bool
	if err := r.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM sdn_entries);`).Scan(&exists); err != nil {
		return false, fmt.Errorf("%s :%w", fn, err)
	}
	return exists, nil
}

func (r *EntryRepository) Search(term string, flag int) ([]Entry, error) {
	const fn = "app.repository.Search"

	term = strings.TrimSpace(term)
	term = regexp.QuoteMeta(term)
	entries := make([]Entry, 0)

	if len(term) == 0 {
		return entries, nil
	}

	switch flag {
	case T_STRONG:
		term = fmt.Sprintf("^%s$", term)
	default:
		term = fmt.Sprintf("(%s)", strings.Join(strings.Split(term, " "), "|"))
	}

	stmt, err := r.db.Prepare(`
		SELECT uid, first_name, last_name 
		FROM sdn_entries 
		WHERE concat_ws(' ', first_name, last_name) ~* $1
	`)
	if err != nil {
		return nil, fmt.Errorf("%s :%w", fn, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(term)
	if err != nil {
		return nil, fmt.Errorf("%s :%w", fn, err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.UID, &entry.FirstName, &entry.LastName); err != nil {
			return nil, fmt.Errorf("%s :%w", fn, err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func NewEntryRepository(db *sql.DB) *EntryRepository {
	return &EntryRepository{db: db}
}
