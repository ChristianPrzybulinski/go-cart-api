package utils

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestInitLog(t *testing.T) {
	type args struct {
		logLevel string
	}
	tests := []struct {
		name       string
		args       args
		wantResult log.Level
	}{
		{"debug mode", args{"debug"}, log.DebugLevel},
		{"error mode", args{"error"}, log.ErrorLevel},
		{"warn mode", args{"warn"}, log.WarnLevel},
		{"info mode", args{"info"}, log.InfoLevel},
		{"empty string", args{""}, log.InfoLevel},
		{"wrong string", args{"jasihddshnaiohdasuihda"}, log.InfoLevel},
		{"warning mode", args{"warning"}, log.WarnLevel},
		{"err mode", args{"err"}, log.ErrorLevel},
		{"numbers", args{"12312421521"}, log.InfoLevel},
		{"null", args{os.Getenv("variavel que nao existe")}, log.InfoLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitLog(tt.args.logLevel)

			gotResult := log.GetLevel()
			if gotResult != tt.wantResult {
				t.Errorf("logs.InitLog() failed = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
