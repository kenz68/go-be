package wpool

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

var (
	errDefault = errors.New("wrong argument type")
	descriptor = JobDescriptor{
		ID:    JobID("1"),
		JType: jobType("anyType"),
		Metadata: jobMetadata{
			"foo": "foo",
			"bar": "bar",
		},
	}

	execFn = func(ctx context.Context, args any) (any, error) {
		argVal, ok := args.(int)
		if !ok {
			return nil, errDefault
		}

		return argVal * 2, nil
	}
)

func Test_Job_Execute(t *testing.T) {
	ctx := context.TODO()
	type fields struct {
		descriptor JobDescriptor
		execFn     ExecutionFn
		args       any
	}

	tests := []struct {
		name   string
		fields fields
		want   Result
	}{
		{
			name: "job execution success",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				args:       10,
			},
			want: Result{
				Value:      20,
				Descriptor: descriptor,
			},
		},
		{
			name: "job execution failure",
			fields: fields{
				descriptor: descriptor,
				execFn:     execFn,
				args:       "20",
			},
			want: Result{
				Err:        errDefault,
				Descriptor: descriptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := Job{
				Descriptor: tt.fields.descriptor,
				ExecFn:     tt.fields.execFn,
				Args:       tt.fields.args,
			}
			got := j.execute(ctx)

			if tt.want.Err != nil {
				if !reflect.DeepEqual(got.Err.Error(), tt.want.Err.Error()) {
					t.Errorf("execute() = %v, want %v", got.Err, tt.want.Err)
				}
				return
			}

			if !reflect.DeepEqual(got.Value, tt.want.Value) {
				t.Errorf("execute() = %v, want %v", got.Value, tt.want.Value)
			}
		})
	}
}
