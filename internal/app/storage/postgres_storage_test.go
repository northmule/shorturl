package storage

import (
	"database/sql"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) ExecContext(url models.URL) error {
	return nil
}
func (m *MockDB) FindByShortURL(shortURL string) (*models.URL, error) {
	return nil, nil
}

func TestPostgresStorage_Add(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		url models.URL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				DB: tt.fields.DB,
			}
			if err := p.Add(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresStorage_FindByShortURL(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.URL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				DB: tt.fields.DB,
			}
			got, err := p.FindByShortURL(tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByShortURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresStorage_FindByURL(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.URL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				DB: tt.fields.DB,
			}
			got, err := p.FindByURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresStorage_Ping(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				DB: tt.fields.DB,
			}
			if err := p.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresStorage_createTable(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PostgresStorage{
				DB: tt.fields.DB,
			}
			if err := p.createTable(); (err != nil) != tt.wantErr {
				t.Errorf("createTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
