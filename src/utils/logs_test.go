// Copyright Christian Przybulinski
// All Rights Reserved

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
		{"Debug Mode", args{"debug"}, log.DebugLevel},
		{"Error Mode", args{"error"}, log.ErrorLevel},
		{"Warn Mode", args{"warn"}, log.WarnLevel},
		{"Info Mode", args{"info"}, log.InfoLevel},
		{"Empty String", args{""}, log.InfoLevel},
		{"Wrong string", args{"jasihddshnaiohdasuihda"}, log.InfoLevel},
		{"Warning Mode", args{"warning"}, log.WarnLevel},
		{"Err Mode", args{"err"}, log.ErrorLevel},
		{"Using Numbers", args{"12312421521"}, log.InfoLevel},
		{"Using a not setted envvar", args{os.Getenv("not exists variable")}, log.InfoLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitLog(tt.args.logLevel)

			gotResult := log.GetLevel()
			if gotResult != tt.wantResult {
				t.Errorf("logs.InitLog() name = %v failed = %v, want %v", tt.name, gotResult, tt.wantResult)
			}
		})
	}
}
