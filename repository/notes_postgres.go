package repository

import (
	"database/sql"
	"fmt"

	"example.com/kode-notes/core"
)

type NotesPostgres struct {
	db *sql.DB
}

func NewNotesPostgres(db *sql.DB) *NotesPostgres {
	return &NotesPostgres{db: db}
}

func (r *NotesPostgres) CreateNote(note *core.Note) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (user_id, title, description) VALUES ($1, $2, $3) RETURNING id", NOTES_TABLE)

	row := r.db.QueryRow(query, note.User_ID, note.Title,
		note.Description)

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *NotesPostgres) GetAllNotes(user_id int) ([]*core.Note, error) {
	var notes []*core.Note

	query := fmt.Sprintf("SELECT * FROM %s WHERE "+
		"%s.user_id = $1", NOTES_TABLE, NOTES_TABLE)

	rows, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n core.Note
		if err := rows.Scan(&n.User_ID, &n.ID, &n.Title,
			&n.Description); err != nil {
			return notes, err
		}
		notes = append(notes, &n)
	}

	return notes, nil
}
