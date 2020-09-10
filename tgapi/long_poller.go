package tgapi

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

const (
	defaultPollTimeout   = 3 * time.Second
	shutdownPollInterval = time.Millisecond
)

// ErrorCallback is a function that is called when an error occurs during an HTTP request on get updates.
type ErrorCallback func(error)

type pollerOptions struct {
	pollTimeout         time.Duration
	listenErrorCallback ErrorCallback
}

func getDefaultPollerOptions() pollerOptions {
	return pollerOptions{
		pollTimeout: defaultPollTimeout,
		listenErrorCallback: func(err error) {
			log.Printf("listen updates: %v", err)
		},
	}
}

// LongPollerOption is used to customize polling behavior.
type LongPollerOption func(*pollerOptions)

// LongPollerErrorListener sets up a listener for polling errors.
func LongPollerErrorListener(listener ErrorCallback) LongPollerOption {
	return func(options *pollerOptions) {
		options.listenErrorCallback = listener
	}
}

// LongPollerErrorListener sets the timeout after an error occurs during polling.
func LongPollerPollTimeout(timeout time.Duration) LongPollerOption {
	return func(options *pollerOptions) {
		options.pollTimeout = timeout
	}
}

type Handler interface {
	HandleUpdate(context.Context, *Update)
}

type LongPoller struct {
	opts pollerOptions

	api     *API
	handler Handler
	ctx     context.Context
	cancel  context.CancelFunc
	// running goroutines.
	running int32
	// graceful stop. Closing a channel prevents requests for new updates.
	stop chan struct{}
}

func NewPoller(api *API, handler Handler, options ...LongPollerOption) *LongPoller {
	opts := getDefaultPollerOptions()
	for _, option := range options {
		option(&opts)
	}
	ctx, cancel := context.WithCancel(context.Background())
	poller := &LongPoller{
		api:     api,
		opts:    opts,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
		stop:    make(chan struct{}),
	}

	return poller
}

func (lp *LongPoller) Listen(updatesConfig *GetUpdates) {
	for {
		select {
		case <-lp.stop:
			return
		case <-lp.ctx.Done():
			return
		default:
		}

		updates, _, err := lp.api.GetUpdates(updatesConfig)
		if err != nil {
			if lp.opts.listenErrorCallback != nil {
				lp.opts.listenErrorCallback(err)
			}
			time.Sleep(lp.opts.pollTimeout)
			continue
		}

		for _, upd := range updates {
			upd := upd
			if upd.UpdateID >= updatesConfig.Offset {
				updatesConfig.Offset = upd.UpdateID + 1

				atomic.AddInt32(&lp.running, 1)
				cctx, cancel := context.WithCancel(lp.ctx)
				go func() {
					lp.handler.HandleUpdate(cctx, &upd)
					cancel()
					atomic.AddInt32(&lp.running, -1)
				}()
			}
		}
	}
}

// Shutdown a-la http.Server.
func (lp *LongPoller) Shutdown(ctx context.Context) error {
	close(lp.stop)
	defer lp.cancel()

	ticker := time.NewTicker(shutdownPollInterval)
	defer ticker.Stop()
	for {
		if atomic.LoadInt32(&lp.running) == 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// AcceptFunc is a function for validating incoming Update, similar to the path prefix in http.
// This function must be non-blocking.
type AcceptFunc func(*Update) bool

// HandlerFunc is a function that is called for an Update that satisfies all upstream AcceptFunc.
type HandlerFunc func(context.Context, *Update)

// CallTree is a type for handling the incoming updates.
type CallTree struct {
	accept  AcceptFunc
	handler HandlerFunc
	childs  []*CallTree
}

// NewCallTree creates a new instance of the call tree with the given handler.
func NewCallTree(defaultHandler HandlerFunc) *CallTree {
	return &CallTree{
		// first accept must be true.
		accept:  func(*Update) bool { return true },
		handler: defaultHandler,
	}
}

// NewChild adds a new child to the current node.
// The priority of the new child is lower than that of the previous child.
// Does not thread safe.
func (c *CallTree) NewChild(accept AcceptFunc, handler HandlerFunc) *CallTree {
	child := &CallTree{
		accept:  accept,
		handler: handler,
	}
	c.childs = append(c.childs, child)
	return child
}

var _ Handler = (*CallTree)(nil)

// HandleUpdate is th implementation method for the Handler interface.
func (c *CallTree) HandleUpdate(ctx context.Context, update *Update) {
	c.walk(ctx, update)
}

func (c *CallTree) walk(ctx context.Context, update *Update) (stop bool) {
	if !c.accept(update) {
		// no match on node, try next
		return false
	}

	for _, child := range c.childs {
		stop := child.walk(ctx, update)
		if stop {
			// handler found
			return true
		}
	}

	if c.handler != nil {
		// leaf
		c.handler(ctx, update)
		return true
	}
	// empty leaf
	return false
}
