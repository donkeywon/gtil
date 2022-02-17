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

    WithLogger(self Service, logger *zap.Logger)
    WithContext(ctx context.Context)
    Context() context.Context

    Closed() <-chan struct{}
    CloseCh()
    ListenAndClose(self Service)

    Children() map[string]Service
    OpenChildren() error
    CloseChildren()
    WaitChildrenClose()
    ChildrenLastError() error
    WithChildContext(ctx context.Context, cancel context.CancelFunc)
    AppendService(name string, svc Service)

    Statistics() map[string]float64
    AppendError(err ...error)
    LastError() error

    Debug(msg string, fields ...zap.Field)
    Info(msg string, fields ...zap.Field)
    Warn(msg string, fields ...zap.Field)
    Error(msg string, fields ...zap.Field)
    Panic(msg string, fields ...zap.Field)
    DPanic(msg string, fields ...zap.Field)
    Fatal(msg string, fields ...zap.Field)
}

func DoOpen(self Service, ctx context.Context, logger *zap.Logger) error {
    self.WithLogger(self, logger)
    self.WithContext(ctx)
    self.WithChildContext(context.WithCancel(context.Background()))

    self.Info("Opening")
    defer self.Info("Opened")

    err := self.OpenChildren()
    if err != nil {
        return errors.Wrapf(err, ErrOpen, self.Name())
    }

    err = self.Open()
    if err != nil {
        return errors.Wrapf(err, ErrOpen, self.Name())
    }

    go self.ListenAndClose(self)
    return nil
}

func DoClose(self Service) error {
    self.Info("Closing")
    defer self.Info("Closed")

    select {
    case <-self.Closed():
        return nil
    default:
        defer self.CloseCh()

        self.CloseChildren()
        self.WaitChildrenClose()

        err := self.ChildrenLastError()
        if err != nil {
            err = errors.Wrapf(err, ErrCloseChildren)
        }

        err = self.Close()
        if err != nil {
            err = errors.Wrapf(err, ErrClose, self.Name())
        }

        return err
    }
}

func DoShutdown(self Service) error {
    self.Info("Shutting down")
    defer self.Info("Shutdown")
    select {
    case <-self.Closed():
        return nil
    default:
        defer self.CloseCh()

        var err error
        for _, child := range self.Children() {
            err = multierr.Append(err, DoShutdown(child))
        }
        if err != nil {
            err = errors.Wrapf(err, ErrShutdownChildren)
        }

        err = self.Shutdown()
        if err != nil {
            err = errors.Wrapf(err, ErrShutdown, self.Name())
        }

        return err
    }
}

type BaseService struct {
    logger      *zap.Logger
    ctx         context.Context
    err         error
    closed      chan struct{}
    children    map[string]Service
    childCtx    context.Context
    childCancel context.CancelFunc
}

func NewBase() *BaseService {
    return &BaseService{
        closed:   make(chan struct{}),
        children: make(map[string]Service),
    }
}

func (bs *BaseService) WithLogger(self Service, logger *zap.Logger) {
    bs.logger = logger.Named(self.Name())
}

func (bs *BaseService) WithContext(ctx context.Context) {
    bs.ctx = ctx
}

func (bs *BaseService) Context() context.Context {
    return bs.ctx
}

func (bs *BaseService) Closed() <-chan struct{} {
    return bs.closed
}

func (bs *BaseService) CloseCh() {
    close(bs.closed)
}

func (bs *BaseService) ListenAndClose(self Service) {
    <-bs.ctx.Done()
    bs.logger.Debug("Receive cancel, start close")
    bs.AppendError(DoClose(self))
}

func (bs *BaseService) Children() map[string]Service {
    return bs.children
}

func (bs *BaseService) OpenChildren() error {
    var err error
    for _, child := range bs.children {
        err = multierr.Append(err, DoOpen(child, bs.childCtx, bs.logger))
    }
    if err != nil {
        err = errors.Wrapf(err, ErrOpenChildren)
    }
    return err
}

func (bs *BaseService) CloseChildren() {
    bs.childCancel()
}

func (bs *BaseService) WaitChildrenClose() {
    for _, child := range bs.children {
        <-child.Closed()
    }
}

func (bs *BaseService) ChildrenLastError() error {
    var err error
    for _, child := range bs.children {
        err = multierr.Append(err, child.LastError())
    }
    return err
}

func (bs *BaseService) WithChildContext(ctx context.Context, cancel context.CancelFunc) {
    bs.childCtx = ctx
    bs.childCancel = cancel
}

func (bs *BaseService) AppendService(name string, svc Service) {
    if name == "" || svc == nil {
        return
    }

    bs.children[name] = svc
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

func (bs *BaseService) Debug(msg string, fields ...zap.Field) {
    bs.logger.Debug(msg, fields...)
}

func (bs *BaseService) Info(msg string, fields ...zap.Field) {
    bs.logger.Info(msg, fields...)
}

func (bs *BaseService) Warn(msg string, fields ...zap.Field) {
    bs.logger.Warn(msg, fields...)
}

func (bs *BaseService) Error(msg string, fields ...zap.Field) {
    bs.logger.Error(msg, fields...)
}

func (bs *BaseService) Panic(msg string, fields ...zap.Field) {
    bs.logger.Panic(msg, fields...)
}

func (bs *BaseService) DPanic(msg string, fields ...zap.Field) {
    bs.logger.DPanic(msg, fields...)
}

func (bs *BaseService) Fatal(msg string, fields ...zap.Field) {
    bs.logger.Fatal(msg, fields...)
}
