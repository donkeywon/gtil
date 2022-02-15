package service

import (
    "context"
    "github.com/donkeywon/gtil/logger"
    "testing"
    "time"
)

type testService struct {
    *BaseService
}

func (t *testService) Name() string {
    return "testService"
}

func (t *testService) Open() error {
    return nil
}

func (t *testService) Close() error {
    return nil
}

func (t *testService) Shutdown() error {
    return nil
}

func NewTestService() *testService {
    return &testService{
        BaseService: NewBase(),
    }
}

type testServiceChildA struct {
    *BaseService
}

func NewTestServiceA() *testServiceChildA {
    return &testServiceChildA{
        BaseService: NewBase(),
    }
}

func (t *testServiceChildA) Name() string {
    return "testServiceChildA"
}

func (t *testServiceChildA) Open() error {
    return nil
}

func (t *testServiceChildA) Close() error {
    time.Sleep(time.Second * 2)
    return nil
}

func (t *testServiceChildA) Shutdown() error {
    time.Sleep(time.Second * 2)
    return nil
}

type testServiceChildB struct {
    *BaseService
}

func NewTestServiceB() *testServiceChildB {
    return &testServiceChildB{
        BaseService: NewBase(),
    }
}

func (t *testServiceChildB) Name() string {
    return "testServiceChildB"
}

func (t *testServiceChildB) Open() error {
    return nil
}

func (t *testServiceChildB) Close() error {
    time.Sleep(time.Second * 3)
    return nil
}

func (t *testServiceChildB) Shutdown() error {
    time.Sleep(time.Second * 3)
    return nil
}

func TestBaseService(t *testing.T) {
    ts := NewTestService()

    l, _ := logger.FromConfig(logger.DefaultConsoleConfig())
    ctx, cancel := context.WithCancel(context.Background())

    tsa := NewTestServiceA()
    tsb := NewTestServiceB()

    ts.AppendService(tsa.Name(), tsa)
    ts.AppendService(tsb.Name(), tsb)

    DoOpen(ts, ctx, l)

    go func() {
        time.Sleep(time.Second * 3)
        cancel()
    }()

    <-ts.Closed()
}
