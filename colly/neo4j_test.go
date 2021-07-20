package neo4j

import (
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestStorage_Init(t *testing.T) {
	type fields struct {
		Database string
		URI      string
		Username string
		Password string
		driver   *neo4j.Driver
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// Test cases.
		{name: "Return Error", fields: fields{URI: "URI_with_no_protocol", Username: "neo4j", Password: "random_pass"}, wantErr: true},
		{name: "Return Non-Error (bolt)", fields: fields{URI: "bolt://localhost:7687", Username: "neo4j", Password: "random_pass"}, wantErr: false},
		{name: "Return Non-Error (neo4j+s)", fields: fields{URI: "neo4j+s://localhost:7687", Username: "neo4j", Password: "random_pass"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				Database: tt.fields.Database,
				URI:      tt.fields.URI,
				Username: tt.fields.Username,
				Password: tt.fields.Password,
				driver:   tt.fields.driver,
			}
			if err := s.Init(); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Tests if Database is correctly set
	t.Run("Set Custom DB name", func(t *testing.T) {
		s := &Storage{
			Database: "testing_db",
			URI:      "bolt://localhost:7687",
			Username: "neo4j",
			Password: "random_pass",
		}
		if err := s.Init(); err != nil {
			t.Errorf("Storage.Init() error = %v", err)
		}
		if s.Database != "testing_db" {
			t.Errorf("Storage Database not properly set")
		}
	})

}
