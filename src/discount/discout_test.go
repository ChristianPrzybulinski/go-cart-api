// Copyright Christian Przybulinski
// All Rights Reserved

package discount

import "testing"

//The test only covers if the return isnt lower than zero
//Since its already a mocked gRPC, only testing the return in in case its not reachable is OK
func TestDescountPercentage(t *testing.T) {
	type args struct {
		port    string
		product int32
		timeout int
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{"Test case 1", args{":50051", 1, 5}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DescountPercentage(tt.args.port, tt.args.product, tt.args.timeout); got < tt.want {
				t.Errorf("DescountPercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}
