package helper

import "testing"

func TestGoLiteralValIntArr2d(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case1",
			args: args{
				val: "1#2#3|4#5#6",
			},
			want: "[][]int{{1,2,3},{4,5,6}}",
		},
		{
			name: "case2",
			args: args{
				val: "",
			},
			want: "[][]int",
		},
		{
			name: "case3",
			args: args{
				val: "1#2",
			},
			want: "[][]int{{1,2}}",
		},
		{
			name: "case4",
			args: args{
				val: "1#2|",
			},
			want: "[][]int{{1,2}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GoLiteralValIntArr2d(tt.args.val); got != tt.want {
				t.Errorf("GoLiteralValIntArr2d() = %v, want %v", got, tt.want)
			}
		})
	}
}
