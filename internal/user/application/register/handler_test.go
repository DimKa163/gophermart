package register

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestRegisterNewUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

}
