package every

import (
	"testing"
)

func TestEvery_CronExpr(t *testing.T) {
	tests := []struct {
		name    string
		every   *Every
		want    string
		wantErr bool
	}{
		{
			"every-2-minutes",
			&Every{Every: "2 minutes"},
			"*/2 * * * *",
			false,
		},
		{
			"every-2-minutes-on-Fri",
			&Every{Every: "2 minutes on Fri"},
			"*/2 * * * Fri",
			false,
		},
		{
			"every-day-at-3pm",
			&Every{Every: "day at 3 pm"},
			"0 15 * * *",
			false,
		},
		{
			"every-hour-in-Jun-on-Sun,Fri",
			&Every{Every: "hour in Jun on Sun,Fri"},
			"0 * * Jun Sun,Fri",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Every{
				Every: tt.every.Every,
				User:  tt.every.User,
				Run:   tt.every.Run,
			}
			got, err := e.CronExpr()
			if (err != nil) != tt.wantErr {
				t.Errorf("Every.CronExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Every.CronExpr() = %v, want %v", got, tt.want)
			}
		})
	}
}
