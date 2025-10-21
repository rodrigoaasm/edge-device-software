package plc_runtime

import (
	"context"
	"log"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/stretchr/testify/assert"
)

type MockOPCUAClient struct {
	F_mch_progress_1 chan bool
	RespCh           chan interface{}
}

func (c MockOPCUAClient) ReadVar(opcNs uint16, path string, vary string) (interface{}, error) {
	return <-c.F_mch_progress_1, nil
}

func (c MockOPCUAClient) WriteVar(opcNs uint16, path string, vars map[string]interface{}) error {
	c.RespCh <- vars
	return nil
}

func TestRuntime(t *testing.T) {
	mockSelfClient := MockOPCUAClient{F_mch_progress_1: make(chan bool, 2), RespCh: make(chan interface{}, 12)}
	mockLineClient := MockOPCUAClient{F_mch_progress_1: make(chan bool, 2), RespCh: make(chan interface{}, 12)}
	log := stdr.New(log.Default())
	ctx := logr.NewContext(context.Background(), log)
	defer ctx.Done()

	plcRunner := NewPLCRunner(100)

	assert.Equal(t, plcRunner.tWindow, 10)
	assert.Equal(t, plcRunner.injectScrewPosition, 9)
	assert.Equal(t, plcRunner.manu1ScrewPosition, 45)
	assert.Equal(t, plcRunner.manu2ScrewPosition, 267)
	assert.Equal(t, plcRunner.returnScrewPosition, 284)
	assert.Equal(t, plcRunner.startInjection, 17)
	assert.Equal(t, plcRunner.maxStep1Injection, 20)
	assert.Equal(t, plcRunner.maxStep2Injection, 34)
	assert.Equal(t, plcRunner.repressInjection, 267)
	assert.Equal(t, plcRunner.down1Pressure, 292)
	assert.Equal(t, plcRunner.down2Pressure, 300)
	assert.Equal(t, plcRunner.startMaxPressure, 34)
	assert.Equal(t, plcRunner.endMaxPressure, 42)
	assert.Equal(t, plcRunner.maxTemp, 34)

	go plcRunner.Run(mockSelfClient, mockLineClient, log)
	mockLineClient.F_mch_progress_1 <- true

	var maxI uint64 = 0
	for _ = range mockSelfClient.RespCh {
		if maxI == 499 {
			mockLineClient.F_mch_progress_1 <- false
		} else {
			mockLineClient.F_mch_progress_1 <- true
		}
		if maxI == 500 {
			plcRunner.Stop()
			break
		}
		maxI++
	}
}
