package smspartner_test

import (
	"testing"

	"github.com/hoflish/smspartner-go/v1"
)

func TestDate(t *testing.T) {
	tests := [...]struct {
		date       smspartner.Date
		wantSDD    string
		wantHour   int
		wantMinute int
		wantErr    bool
	}{
		{smspartner.NewDate(2018, 8, 16, 17, 45), "16/08/2018", 17, 45, false},
		{smspartner.NewDate(2018, 13, 40, 17, 4), "09/02/2019", 17, 0, true},
	}

	for i, tt := range tests {
		if tt.date.ScheduledDeliveryDate() != tt.wantSDD {
			t.Errorf("#%d. got: %s, want: %s", i, tt.date.ScheduledDeliveryDate(), tt.wantSDD)
		}
		if tt.date.Time.Hour() != tt.wantHour {
			t.Errorf("#%d. got: %d, want: %d", i, tt.date.Time.Hour(), tt.wantHour)
		}

		gotMinute, err := tt.date.MinuteToSendSMS()
		if tt.wantErr != (err != nil) {
			t.Errorf("#%d. expected a non-nil error", i)
		}
		if gotMinute != tt.wantMinute {
			t.Errorf("#%d. got: %d, want: %d", i, gotMinute, tt.wantMinute)
		}
	}
}
