package entrypoint

import (
	"context"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	vi          *viper.Viper
	mu          sync.Mutex
	reloadCh    = make(chan struct{})
	shutdownCtx context.Context
	cancelFn    context.CancelFunc
	ep          *EntryPoint
	wd          string
)

const prefix = "app.entrypoint"

func init() {
	shutdownCtx, cancelFn = context.WithCancel(context.Background())
}

// Viper returns instance of Viper
func Viper() *viper.Viper {
	mu.Lock()
	defer mu.Unlock()
	if vi != nil {
		return vi
	}
	vi = viper.GetViper()
	return vi
}

// Initialize returns instance of entry point singleton manager, workDir uses to attach templates
func Initialize(workDir string, v *viper.Viper) (*EntryPoint, error) {
	mu.Lock()
	defer mu.Unlock()
	if ep != nil {
		return ep, nil
	}
	if len(workDir) > 0 {
		wd = workDir
	} else {
		wd, _ = os.Getwd()
	}
	vi, ep = v, &EntryPoint{}
	return ep, nil
}

// OnShutdown subscribe on shutdown event for gracefully exit via context.
func OnShutdown() context.Context {
	return shutdownCtx
}

// OnReload subscribe on reload event.
func OnReload() <-chan struct{} {
	return reloadCh
}

// EntryPoint manager of single point of application
type EntryPoint struct {
}

// Shutdown raise shutdown event.
func Shutdown(ctx context.Context, code int) {
	mu.Lock()
	defer mu.Unlock()
	cancelFn()
	if _, ok := ctx.Deadline(); ok {
		<-ctx.Done()
	}
}

// Reload raise reload event.
func (e *EntryPoint) Reload() {
	mu.Lock()
	defer mu.Unlock()
	ch := reloadCh
	reloadCh = make(chan struct{})
	close(ch)
}

// WorkDir returns current work directory
func WorkDir() string {
	return wd
}
