package discount

import "testing"

func TestDescountPercentage(t *testing.T) {
	type args struct {
		port    string
		product int32
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{"testing", args{":50051", 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DescountPercentage(tt.args.port, tt.args.product); got < tt.want {
				t.Errorf("DescountPercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}
