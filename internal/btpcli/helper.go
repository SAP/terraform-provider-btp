package btpcli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// firstElementOrDefault returns the first element of a slice or if not available the given defaultValue
func firstElementOrDefault[T any](slice []T, defaultValue T) T {
	return nthElementOrDefault(slice, 0, defaultValue)
}

// nthElementOrDefault returns the n-th element of a slice or if not available the given defaultValue
func nthElementOrDefault[T any](slice []T, n int, defaultValue T) T {
	if n >= len(slice) {
		return defaultValue
	}

	return slice[n]
}

func doExecute[T interface{}](cliClient *v2Client, ctx context.Context, req *CommandRequest, options ...CommandOptions) (T, CommandResponse, error) {
	var obj T

	res, err := cliClient.Execute(ctx, req, options...)

	if err != nil {
		return obj, res, err
	}

	defer func() {
		if tempErr := res.Body.Close(); tempErr != nil {
			if err != nil {
				err = fmt.Errorf("original error: %v; error while closing response body: %v", err, tempErr)
			} else {
				err = tempErr
			}
		}
	}()

	if err = json.NewDecoder(res.Body).Decode(&obj); err == nil || err == io.EOF {
		return obj, res, nil
	} else {
		return obj, res, err
	}
}
