package validator_test

import (
	"context"
	"testing"

	"github.com/fakovacic/editor/internal/app/write/validator"
)

func TestValidator(t *testing.T) {
	ctx := context.Background()

	v := validator.New()

	err := v.AddClient(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}

	isReady := v.IsReady(ctx)
	if !isReady {
		t.Fatal("expected ready")
	}

	v.Clear(ctx)

	isReady = v.IsReady(ctx)
	if !isReady {
		t.Fatal("expected ready")
	}

	err = v.RemoveClient(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}

	isReady = v.IsReady(ctx)
	if isReady {
		t.Fatal("expected not ready")
	}
}

func TestValidatorMultipleUsers(t *testing.T) {
	ctx := context.Background()

	v := validator.New()

	err := v.AddClient(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}

	err = v.AddClient(ctx, "2")
	if err != nil {
		t.Fatal(err)
	}

	isReady := v.IsReady(ctx)
	if isReady {
		t.Fatal("expected not ready")
	}

	err = v.ReadyClient(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}

	isReady = v.IsReady(ctx)
	if isReady {
		t.Fatal("expected not ready")
	}

	err = v.ReadyClient(ctx, "2")
	if err != nil {
		t.Fatal(err)
	}

	isReady = v.IsReady(ctx)
	if !isReady {
		t.Fatal("expected ready")
	}

	err = v.UnreadyClient(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}

	isReady = v.IsReady(ctx)
	if isReady {
		t.Fatal("expected not ready")
	}

	err = v.ReadyClient(ctx, "1")
	if err != nil {
		t.Fatal(err)
	}

	v.Clear(ctx)

	isReady = v.IsReady(ctx)
	if isReady {
		t.Fatal("expected not ready")
	}
}

func TestValidatorError(t *testing.T) {
	ctx := context.Background()

	v := validator.New()

	err := v.RemoveClient(ctx, "1")
	if err == nil {
		t.Fatal("expected error")
	}

	err = v.ReadyClient(ctx, "1")
	if err == nil {
		t.Fatal("expected error")
	}

	err = v.UnreadyClient(ctx, "1")
	if err == nil {
		t.Fatal("expected error")
	}
}
