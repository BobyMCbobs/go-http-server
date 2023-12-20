package metrics

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMetrics_Handle(t *testing.T) {
	type fields struct {
		Enabled bool
		Port    string
	}
	tests := []struct {
		name          string
		noQuitChannel bool
		fields        fields
	}{
		{
			name: "basic",
			fields: fields{
				Enabled: true,
				Port:    fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000),
			},
		},
		{
			name: "not enabled",
			fields: fields{
				Enabled: false,
				Port:    fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Metrics{
				Enabled: tt.fields.Enabled,
				Port:    tt.fields.Port,
			}
			quitChan := make(chan bool, 1)
			if tt.noQuitChannel {
				quitChan = nil
			}
			go m.Handle(quitChan)
			time.Sleep(1 * time.Second)
			if !tt.noQuitChannel {
				defer func() {
					quitChan <- true
				}()
			}
		})
	}
}
