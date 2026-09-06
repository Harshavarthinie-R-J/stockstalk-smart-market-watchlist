package change

import "testing"

func TestDetect(t *testing.T) {
	tests := []struct {
		name         string
		previous     float64
		current      float64
		weekHigh     float64
		weekLow      float64
		wantType     Type
		wantSeverity Severity
	}{
		{
			name:         "normal movement",
			previous:     100,
			current:      101,
			weekHigh:     120,
			weekLow:      80,
			wantType:     PriceMovement,
			wantSeverity: Normal,
		},
		{
			name:         "notable movement",
			previous:     100,
			current:      103,
			weekHigh:     120,
			weekLow:      80,
			wantType:     PriceMovement,
			wantSeverity: Notable,
		},
		{
			name:         "significant movement",
			previous:     100,
			current:      106,
			weekHigh:     120,
			weekLow:      80,
			wantType:     PriceMovement,
			wantSeverity: Significant,
		},
		{
			name:         "new high",
			previous:     100,
			current:      120,
			weekHigh:     120,
			weekLow:      80,
			wantType:     NewHigh,
			wantSeverity: Significant,
		},
		{
			name:         "new low",
			previous:     100,
			current:      80,
			weekHigh:     120,
			weekLow:      80,
			wantType:     NewLow,
			wantSeverity: Significant,
		},
	}

	service := &Service{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.Detect(
				tt.previous,
				tt.current,
				tt.weekHigh,
				tt.weekLow,
			)

			if got.Type != tt.wantType {
				t.Fatalf(
					"type = %s, want %s",
					got.Type,
					tt.wantType,
				)
			}

			if got.Severity != tt.wantSeverity {
				t.Fatalf(
					"severity = %s, want %s",
					got.Severity,
					tt.wantSeverity,
				)
			}
		})
	}
}

func TestVolumeAnomaly(t *testing.T) {
	service := &Service{}

	c := service.Detect(
		100,
		103,
		120,
		80,
	)

	service.EnrichWithVolume(
		nil,
		&c,
		20_000_000,
		8_000_000,
	)

	if !c.VolumeAnomaly {
		t.Fatal("expected volume anomaly")
	}

	if c.VolumeRatio < 2.4 || c.VolumeRatio > 2.6 {
		t.Fatalf(
			"volume ratio = %f, want approximately 2.5",
			c.VolumeRatio,
		)
	}
}

func TestEventFusion(t *testing.T) {
	service := &Service{}

	c := service.Detect(
		100,
		103,
		120,
		80,
	)

	service.EnrichWithVolume(
		nil,
		&c,
		20_000_000,
		8_000_000,
	)

	c = service.Fuse(c)

	if c.Type != FusedEvent {
		t.Fatalf(
			"type = %s, want %s",
			c.Type,
			FusedEvent,
		)
	}
}

func TestAttribution(t *testing.T) {
	service := &Service{}

	c := service.Detect(
		100,
		103,
		120,
		80,
	)

	service.EnrichWithVolume(
		nil,
		&c,
		20_000_000,
		8_000_000,
	)

	c = service.Fuse(c)
	c = service.AddAttribution(c)

	if c.Attribution == "" {
		t.Fatal("expected attribution")
	}

	if c.AttributionConfidence == NoClearContext {
		t.Fatal("expected attribution confidence")
	}
}
