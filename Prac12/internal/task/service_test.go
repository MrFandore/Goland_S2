package task

import "testing"

func TestServiceCreateAndUpdate(t *testing.T) {
	service := NewService()

	created, err := service.Create(CreateInput{
		Title:       "Сравнить REST и GraphQL",
		Description: "Практическая работа №12",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.ID != "t_003" {
		t.Fatalf("Create() id = %s, want t_003", created.ID)
	}

	done := true
	updated, err := service.Update(created.ID, UpdateInput{Done: &done})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if !updated.Done {
		t.Fatal("Update() did not set done flag")
	}
}
