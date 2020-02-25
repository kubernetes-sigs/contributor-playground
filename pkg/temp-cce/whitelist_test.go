package temp_cce

import (
	"context"
	"net/http"
	"testing"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
)

func TestClient_CheckWhiteList(t *testing.T) {
	type args struct {
		ctx         context.Context
		featureType FeatureType
		signOpt     *bce.SignOption
	}
	tests := []struct {
		name    string
		envs    []*testEnvConfig
		args    args
		want    bool
		wantErr bool
	}{
		// All test cases.
		{
			name: "normal case",
			envs: []*testEnvConfig{
				{
					uri:          "/v1/cluster/check_white_list",
					method:       "GET",
					statusCode:   http.StatusOK,
					responseBody: []byte(`{"isExist":true}`),
				},
			},
			args: args{
				ctx:         context.TODO(),
				featureType: EnableClusterRBAC,
			},
			want: true,
		},
		{
			name: "bad json case",
			envs: []*testEnvConfig{
				{
					uri:          "/v1/cluster/check_white_list",
					method:       "GET",
					statusCode:   http.StatusOK,
					responseBody: []byte(`"isExist":true`),
				},
			},
			args: args{
				ctx:         context.TODO(),
				featureType: EnableClusterRBAC,
			},
			wantErr: true,
		},
		{
			name: "bad request case",
			envs: []*testEnvConfig{
				{
					uri:        "/v1/cluster/check_white_list",
					method:     "GET",
					statusCode: http.StatusBadRequest,
				},
			},
			args: args{
				ctx:         context.TODO(),
				featureType: EnableClusterRBAC,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestEnv(tt.envs)
			defer tearDownTestEnv()

			c := cceClient
			got, err := c.CheckWhiteList(tt.args.ctx, tt.args.featureType, tt.args.signOpt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CheckWhiteList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.CheckWhiteList() = %v, want %v", got, tt.want)
			}
		})
	}
}
