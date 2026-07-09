package config

import "testing"

func TestLoadReadsEnvironmentValues(t *testing.T) {
	t.Setenv("MONGO_URI", "mongodb://localhost:27017")
	t.Setenv("MONGO_DB", "testdb")
	t.Setenv("MONGO_COLLECTION", "testusers")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.MongoURI != "mongodb://localhost:27017" {
		t.Fatalf("expected mongo URI from env, got %q", cfg.MongoURI)
	}

	if cfg.DBName != "testdb" {
		t.Fatalf("expected db name from env, got %q", cfg.DBName)
	}

	if cfg.CollectionName != "testusers" {
		t.Fatalf("expected collection name from env, got %q", cfg.CollectionName)
	}
}
