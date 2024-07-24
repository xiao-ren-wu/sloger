# sloger

```go
func TestSloger(t *testing.T) {
	f, err := os.Create("test.log")
	if err != nil {
		t.Fatalf(err.Error())
	}

	slogerHandler := NewHandler(
		&slog.HandlerOptions{AddSource: true},
		Json,
		[]string{"x-tt-logId"},
		os.Stdout, f)

	logger := slog.New(slogerHandler)
	slog.SetDefault(logger)

	slog.Info("hello slog")
	ctx := context.WithValue(context.Background(), "x-tt-logId", "12986329")
	slog.InfoContext(ctx, "hello slog")
}

```