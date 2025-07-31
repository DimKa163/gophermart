package mediatr

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegister_Should_Be_Success(t *testing.T) {
	defer cleanup()
	err := Bind[*RegisterShouldBeSuccessCommand, *Data](&RegisterShouldBeSuccessHandler{})
	assert.Nil(t, err)
}

func TestRegister_Second_Time_Should_Be_Error(t *testing.T) {
	defer cleanup()
	_ = Bind[*RegisterShouldBeSuccessCommand, *Data](&RegisterShouldBeSuccessHandler{})
	err := Bind[*RegisterShouldBeSuccessCommand, *Data](&RegisterShouldBeSuccessHandler{})
	assert.NotNil(t, err)
}

func TestSend_Should_Be_Success(t *testing.T) {
	defer cleanup()
	_ = Bind[*CallCommandShouldBeSuccessCommand, *Data](&CallCommandShouldBeSuccessHandler{})
	d, err := Send[*CallCommandShouldBeSuccessCommand, *Data](context.Background(), &CallCommandShouldBeSuccessCommand{})
	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, Data("test"), *d)
}

type RegisterShouldBeSuccessCommand struct{}

type Data string

type RegisterShouldBeSuccessHandler struct{}

func (rssh *RegisterShouldBeSuccessHandler) Handle(_ context.Context, _ *RegisterShouldBeSuccessCommand) (*Data, error) {
	data := Data("test")
	return &data, nil
}

type CallCommandShouldBeSuccessCommand struct{}
type CallCommandShouldBeSuccessHandler struct{}

func (ccsbsh *CallCommandShouldBeSuccessHandler) Handle(_ context.Context, _ *CallCommandShouldBeSuccessCommand) (*Data, error) {
	data := Data("test")
	return &data, nil
}

func cleanup() {
	ClearCommands()
}
