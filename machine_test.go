package machine_test

import (
	"context"
	"fmt"
	"github.com/autom8ter/machine"
	"testing"
	"time"
)

func Test(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	m, err := machine.New(ctx, &machine.Opts{
		MaxRoutines: 100,
		Debug:       true,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	for x := 0; x < 1000; x++ {
		m.Go(func(ctx context.Context) error {
			i := x
			m.Cache().Set(fmt.Sprint(i), map[string]interface{}{
				"index": fmt.Sprint(i),
				"current": m.Current(),
			})
			time.Sleep(200 * time.Millisecond)
			return nil
		})
	}
	m.Cache().Range(func(id string, data map[string]interface{}) bool {
		t.Logf("id = %v data = %v\n", id, data)
		return true
	})
	t.Logf("stats = %v\n", m.Stats())
	if errs := m.Wait(); len(errs) > 0 {
		for _, err := range errs {
			t.Logf("workerPool error: %s", err)
		}
	}
	t.Logf("after: %v", m.Current())

}