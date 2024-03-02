package mysql

import (
	"AituNews/pkg/models"
	"database/sql"
	"errors"
	"time"
)

type TicketModel struct {
	DB *sql.DB
}

func (tm *TicketModel) Create(userID, movieID int, movieTitle string, sessionTime time.Time) error {
	stmt := `INSERT INTO tickets (user_id, movie_id, movie_title, session_time) VALUES (?, ?, ?, ?)`
	_, err := tm.DB.Exec(stmt, userID, movieID, movieTitle, sessionTime)
	if err != nil {
		return err
	}
	return nil
}

func (tm *TicketModel) Update(userID, movieID int, movieTitle string, sessionTime time.Time) error {
	// Implement update logic if needed
	return errors.New("Update operation not implemented for tickets")
}

func (tm *TicketModel) Delete(ticketID int) error {
	stmt := `DELETE FROM tickets WHERE id = ?`
	_, err := tm.DB.Exec(stmt, ticketID)
	if err != nil {
		return err
	}
	return nil
}

func (tm *TicketModel) Get(ticketID int) (*models.Ticket, error) {
	stmt := `SELECT id, user_id, movie_id, movie_title, session_time FROM tickets WHERE id = ?`
	row := tm.DB.QueryRow(stmt, ticketID)

	ticket := &models.Ticket{}
	err := row.Scan(&ticket.ID, &ticket.UserID, &ticket.MovieID, &ticket.MovieTitle, &ticket.SessionTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return ticket, nil
}

func (tm *TicketModel) GetAll() ([]*models.Ticket, error) {
	stmt := `SELECT id, user_id, movie_id, movie_title, session_time FROM tickets`
	rows, err := tm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*models.Ticket
	for rows.Next() {
		ticket := &models.Ticket{}
		err := rows.Scan(&ticket.ID, &ticket.UserID, &ticket.MovieID, &ticket.MovieTitle, &ticket.SessionTime)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}
