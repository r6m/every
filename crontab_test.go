package every

import (
	"strings"
	"testing"
)

func TestCrontab_writeCrontab(t *testing.T) {
	crontab := `*/2 * * * * command >/dev/null 2>&1
*/2 * * * * command >/dev/null 2>&1

# Begin every generated file for /Users/reza/Everyfile
*/2 * * * * command >/dev/null 2>&1
*/2 * * * * command >/dev/null 2>&1
# End every generated file for /Users/reza/Everyfile`

	type args struct {
		configPath string
		items      []*Every
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "update-existing",
			args: args{
				configPath: "~/Everyfile",
				items: []*Every{
					{Every: "2 minutes", Run: "command >/dev/null 2>&1"},
					{Every: "2 minutes", Run: "command >/dev/null 2>&1"},
					{Every: "2 minutes", Run: "command >/dev/null 2>&1"},
				},
			},
			want:    crontab,
			wantErr: false,
		},
		{
			name: "write-crontab",
			args: args{
				configPath: "/tmp/Everyfile",
				items: []*Every{
					{Every: "2 minutes", Run: "command >/dev/null 2>&1"},
					{Every: "2 minutes", Run: "command >/dev/null 2>&1"},
				},
			},
			want: "# Begin every generated jobs for /tmp/Everyfile" + "\n" +
				"*/2 * * * * command >/dev/null 2>&1" + "\n" +
				"*/2 * * * * command >/dev/null 2>&1" + "\n" +
				"# End every generated jobs for /tmp/Everyfile",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := updateCrontab(crontab, tt.args.configPath, tt.args.items...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crontab.writeCrontab() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("Crontab.writeCrontab() = %v\nwant:\n%v", got, tt.want)
			}
		})
	}
}
