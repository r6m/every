package every

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func TestEvery_UnmarshalCaddyfile(t *testing.T) {
	data := []byte(`
	{
		every 2 days at 12:30 AM {
			run hello world
			user ubuntu
		}
		every day at 01:45 PM {
			run "gobackup perform -c /etc/gobackup.yml > /dev/null"
			user ubuntu
		}
	}`)

	blocks, err := caddyfile.Parse("Caddyfile", data)
	if err != nil {
		t.Error(err)
	}

	for _, s := range blocks[0].Segments {
		d := caddyfile.NewDispenser(s)
		config := new(Every)
		if err := config.UnmarshalCaddyfile(d); err != nil {
			t.Error("can't unmarshal every", err)
		}

	}
}

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
			"0 0 * Jun Sun,Fri",
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
