package mediatr

import (
	"context"
	"reflect"
)

var ErrHandlerAlreadyExists = NewMediatrRegisterHandlerError("handler already exists")

var ErrHandlerNotFound = NewMediatrNotFoundHandlerError("handler not found")

var ErrConvertingHandler error = NewMediatrConvertingHandlerError("converting handler error")
var (
	handlers map[reflect.Type]interface{} = map[reflect.Type]interface{}{}
)

func Bind[TCommand any, TResponse any](handler CommandHandler[TCommand, TResponse]) error {
	var com TCommand
	comType := reflect.TypeOf(com)
	_, ok := handlers[comType]
	if ok {
		return ErrHandlerAlreadyExists
	}
	handlers[comType] = handler
	return nil
}

func Send[TCommand any, TResponse any](ctx context.Context, command TCommand) (TResponse, error) {
	comType := reflect.TypeOf(command)
	handler, ok := handlers[comType]
	if !ok {
		return *new(TResponse), ErrHandlerNotFound
	}
	handlerInstance, ok := handler.(CommandHandler[TCommand, TResponse])
	if !ok {
		return *new(TResponse), ErrConvertingHandler
	}
	resp, err := handlerInstance.Handle(ctx, command)
	if err != nil {
		return *new(TResponse), err
	}
	return resp, nil
}

func ClearCommands() {
	handlers = map[reflect.Type]interface{}{}
}
