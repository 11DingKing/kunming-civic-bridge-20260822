package platform

import (
	"context"
	"fmt"
	"time"
)

// JSONCodec owns stable JSON codec.
type JSONCodec struct {
	Name    string
	Enabled bool
	Timeout time.Duration
}

func (x JSONCodec) Validate() error {
	if x.Name == "" {
		return fmt.Errorf("JSONCodec: name is required")
	}
	if x.Timeout < 0 {
		return fmt.Errorf("JSONCodec: timeout cannot be negative")
	}
	return nil
}
func (x JSONCodec) Start(ctx context.Context) error {
	if err := x.Validate(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
func (x JSONCodec) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Nanosecond):
		return nil
	}
}
func (x JSONCodec) Ready() bool { return x.Enabled && x.Name != "" }
func (x JSONCodec) Copy() JSONCodec {
	return JSONCodec{Name: x.Name, Enabled: x.Enabled, Timeout: x.Timeout}
}
func (x JSONCodec) WithName(v string) JSONCodec           { y := x.Copy(); y.Name = v; return y }
func (x JSONCodec) WithEnabled(v bool) JSONCodec          { y := x.Copy(); y.Enabled = v; return y }
func (x JSONCodec) WithTimeout(v time.Duration) JSONCodec { y := x.Copy(); y.Timeout = v; return y }
func (x JSONCodec) Status() string {
	if x.Ready() {
		return "ready"
	}
	return "not_ready"
}
func (x JSONCodec) Describe() string { return fmt.Sprintf("%s:%s", x.Name, x.Status()) }
