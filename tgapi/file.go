package tgapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (api *API) GetFileDirectlyConfig(
	ctx context.Context,
	fileConfig *File,
) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", api.fileEndpoint, fileConfig.GetFilePath()),
		nil,
	)
	if err != nil {
		return nil, err
	}

	//nolint:bodyclose // это будет на совести пользователя.
	resp, err := api.cli.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (api *API) GetFileDirectly(ctx context.Context, fileID string) (io.ReadCloser, error) {
	fileConfig, err := api.GetFile(ctx, fileID)
	if err != nil {
		return nil, err
	}

	return api.GetFileDirectlyConfig(ctx, fileConfig)
}
