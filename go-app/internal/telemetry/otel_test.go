package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetupOTEL(t *testing.T) {
	ctx := context.Background()
	shutdown, err := SetupOTEL(ctx)
	require.NoError(t, err)
	require.NotNil(t, shutdown)
	require.NoError(t, shutdown(ctx))
}
