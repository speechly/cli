package cmd

import (
	salv1 "github.com/speechly/api/go/speechly/sal/v1"
)

type ValidateWriter struct {
	appId  string
	stream salv1.Compiler_ValidateClient
}

type CompileWriter struct {
	appId  string
	stream salv1.Compiler_CompileClient
}

func createAppsource(appId string, data []byte) *salv1.AppSource {
	contentType := salv1.AppSource_CONTENT_TYPE_TAR
	return &salv1.AppSource{AppId: appId, DataChunk: data, ContentType: contentType}
}

func (u ValidateWriter) Write(data []byte) (n int, err error) {
	if err = u.stream.Send(createAppsource(u.appId, data)); err != nil {
		return 0, err
	}
	return len(data), nil
}

func (u CompileWriter) Write(data []byte) (n int, err error) {
	if err = u.stream.Send(createAppsource(u.appId, data)); err != nil {
		return 0, err
	}
	return len(data), nil
}