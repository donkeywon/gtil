package service

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Service interface {
	Name() string
	Open() error
	Close() error
	Shutdown() error

	withLogger(self Service, logger *zap.Logger)
	withContext(ctx context.Context)
	Context() context.Context

	Closed() <-chan struct{}
	closeCh()
	listenAndClose(self Service)

	GetChildrenSvc(name string) Service
	ChildrenSvcs() []Service
	ChildrenLastError() error
	AppendService(name string, svc Service)
	openChildren() error
	closeChildren()
	waitChildrenClose()
	withChildContext(ctx context.Context, cancel context.CancelFunc)

	Statistics() map[string]float64
	AppendError(err ...error)
	LastError() error
}

func DoOpen(self Service, ctx context.Context, logger *zap.Logger) error {
	self.withLogger(self, logger)
	self.withContext(ctx)
	self.withChildContext(context.WithCancel(context.Background()))

	err := self.openChildren()
	if err != nil {
		return errors.Wrapf(err, ErrOpenSvc, self.Name())
	}

	err = self.Open()
	if err != nil {
		return errors.Wrapf(err, ErrOpenSvc, self.Name())
	}

	go self.listenAndClose(self)
	return nil
}

func DoClose(self Service) error {
	select {
	case <-self.Closed():
		return nil
	default:
		defer self.closeCh()

		self.closeChildren()
		self.waitChildrenClose()

		err := multierr.Combine(self.ChildrenLastError(), self.Close())
		if err != nil {
			return errors.Wrapf(err, ErrCloseSvc, self.Name())
		}

		return nil
	}
}

func DoShutdown(self Service) error {
	select {
	case <-self.Closed():
		return nil
	default:
		defer self.closeCh()

		var err error
		for _, child := range self.ChildrenSvcs() {
			err = multierr.Append(err, DoShutdown(child))
		}
		err = multierr.Append(err, self.Shutdown())

		if err != nil {
			return errors.Wrapf(err, ErrShutdownSvc, self.Name())
		}

		return nil
	}
}

type BaseService struct {
	*zap.Logger
	ctx         context.Context
	err         error
	closed      chan struct{}
	children    map[string]Service
	childrenArr []Service
	childCtx    context.Context
	childCancel context.CancelFunc
}

func NewBase() *BaseService {
	return &BaseService{
		closed:      make(chan struct{}),
		children:    make(map[string]Service),
		childrenArr: make([]Service, 0, 1),
	}
}

func (bs *BaseService) withLogger(self Service, logger *zap.Logger) {
	bs.Logger = logger.Named(self.Name())
}

func (bs *BaseService) withContext(ctx context.Context) {
	bs.ctx = ctx
}

func (bs *BaseService) Context() context.Context {
	return bs.ctx
}

func (bs *BaseService) Closed() <-chan struct{} {
	return bs.closed
}

func (bs *BaseService) closeCh() {
	close(bs.closed)
}

func (bs *BaseService) listenAndClose(self Service) {
	<-bs.ctx.Done()
	bs.Debug("Receive cancel, start close")
	bs.AppendError(DoClose(self))
}

func (bs *BaseService) GetChildrenSvc(name string) Service {
	return bs.children[name]
}

func (bs *BaseService) ChildrenSvcs() []Service {
	return bs.childrenArr
}

func (bs *BaseService) openChildren() error {
	var err error
	for _, child := range bs.childrenArr {
		err = multierr.Append(err, DoOpen(child, bs.childCtx, bs.Logger))
	}
	return err
}

func (bs *BaseService) closeChildren() {
	bs.childCancel()
}

func (bs *BaseService) waitChildrenClose() {
	for _, child := range bs.childrenArr {
		<-child.Closed()
	}
}

func (bs *BaseService) ChildrenLastError() error {
	var err error
	for _, child := range bs.childrenArr {
		err = multierr.Append(err, child.LastError())
	}
	return err
}

func (bs *BaseService) withChildContext(ctx context.Context, cancel context.CancelFunc) {
	bs.childCtx = ctx
	bs.childCancel = cancel
}

func (bs *BaseService) AppendService(name string, svc Service) {
	if name == "" || svc == nil {
		return
	}

	bs.children[name] = svc
	bs.childrenArr = append(bs.childrenArr, svc)
}

func (bs *BaseService) AppendError(err ...error) {
	bs.err = multierr.Append(bs.err, multierr.Combine(err...))
}

func (bs *BaseService) Statistics() map[string]float64 {
	return nil
}

func (bs *BaseService) LastError() error {
	return bs.err
}
