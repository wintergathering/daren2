package sqlite

import (
	"database/sql"

	daren "github.com/wintergathering/daren2"
)

type dareService struct {
	db *sql.DB
}

func NewDareService(db *sql.DB) daren.DareService {
	return &dareService{
		db: db,
	}
}

//methods -------------

func (ds dareService) CreateDare(d *daren.Dare) (int, error) {
	qry := `
		INSERT INTO dares (title, dare_text, added_by) VALUES (?, ?, ?)	
	`

	stmt, err := ds.db.Prepare(qry)

	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(d.Title, d.Text, d.AddedBy)

	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func (ds dareService) GetDareByID(id int) (*daren.Dare, error) {
	dare := &daren.Dare{}

	err := ds.db.QueryRow("SELECT * FROM dares WHERE dare_id = ?", id).Scan(&dare.ID, &dare.Title, &dare.Text, &dare.AddedBy, &dare.Seen, &dare.CreatedAt)

	if err != nil {
		return nil, err
	}

	return dare, nil
}

func (ds dareService) GetRandomDare() (*daren.Dare, error) {
	qry := `
		SELECT *
		FROM dares
		WHERE seen = 0	
		ORDER BY RANDOM()
		LIMIT 1;
	`
	dare := &daren.Dare{}

	err := ds.db.QueryRow(qry).Scan(&dare.ID, &dare.Title, &dare.Text, &dare.AddedBy, &dare.Seen, &dare.CreatedAt)

	if err != nil {
		return nil, err
	}

	return dare, nil
}

func (ds dareService) GetAllDares() ([]*daren.Dare, error) {
	//TODO
	return nil, nil
}
