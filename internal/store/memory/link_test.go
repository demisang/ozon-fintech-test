package memory_test

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/demisang/ozon-fintech-test/internal/models"
	"github.com/demisang/ozon-fintech-test/internal/store/memory"
	"github.com/demisang/ozon-fintech-test/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	log := logger.GetLogger()
	ctx := context.Background()

	length := 100
	storage := memory.New(log)
	localStorage := make([]models.Link, length)
	wg := sync.WaitGroup{}

	t.Run("multiple simultaneous write", func(t *testing.T) {
		for i := 0; i < length; i++ {
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				tmp, err := storage.Create(ctx, models.CreateLinkDto{URL: "stroka" + strconv.Itoa(i)})
				require.NoError(t, err)
				localStorage[i] = tmp
			}()
		}
		wg.Wait()
	})

	t.Run("Multiple simultaneous read", func(t *testing.T) {
		for _, link := range localStorage {
			link := link
			wg.Add(1)
			go func() {
				defer wg.Done()
				tmp, err := storage.GetByCode(ctx, link.Code)
				require.NoError(t, err)
				require.Equal(t, link.URL, tmp.URL)
			}()
		}
	})
}
