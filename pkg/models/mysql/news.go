package mysql

import (
	"AituNews/pkg/models"
	"database/sql"
	"errors"
)

type NewsModel struct {
	DB *sql.DB
}

func (m *NewsModel) Get(id int) (*models.News, error) {

	stmt := `SELECT id, title, content, details, created, category FROM news
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.News{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Details, &s.Created, &s.Category)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil

}

func (m *NewsModel) Latest() ([]*models.News, error) {
	stmt := `SELECT id, title, content, details, created, category FROM news ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	news := []*models.News{}
	for rows.Next() {
		s := &models.News{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Details, &s.Created, &s.Category)
		if err != nil {
			return nil, err
		}
		news = append(news, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return news, nil
}

func (m *NewsModel) LatestByCategory(category string) ([]*models.News, error) {
	stmt := `SELECT id, title, content, details, created, category FROM news WHERE category = ? ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	news := []*models.News{}
	for rows.Next() {
		s := &models.News{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Details, &s.Created, &s.Category)
		if err != nil {
			return nil, err
		}
		news = append(news, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return news, nil
}

func (m *NewsModel) Insert(title, content, details, category string) (int, error) {
	stmt := `INSERT INTO news ( title, content, details, created, category)
    VALUES( ?, ?, ?, UTC_TIMESTAMP(), ? )`

	result, err := m.DB.Exec(stmt, title, content, details, category)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}
