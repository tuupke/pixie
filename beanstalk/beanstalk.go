package beanstalk

import (
	"errors"
	"sync"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/rs/zerolog/log"
)

// Connect can be used to connect to a beanstalk instance and listen on some tubes.
func Connect(network, addr string, tubes ...string) Connection {
	cc := &Conn{
		Tubes: tubes,
		Pool:  sync.Pool{},
	}

	cc.Pool.New = func() interface{} {
		conn, err := dial(network, addr, cc.Tubes...)
		log.Err(err).Strs("tubes", cc.Tubes).Msg("dialed beanstalk")
		if conn == nil {
			return (*beanstalk.Conn)(nil)
		}

		return conn
	}

	return cc
}

func dial(network, addr string, tubes ...string) (*beanstalk.Conn, error) {
	bs, err := beanstalk.Dial(network, addr)
	if err != nil {
		return (*beanstalk.Conn)(nil), err
	}

	bs.TubeSet = *beanstalk.NewTubeSet(bs, tubes...)
	return bs, err
}

type Connection interface {
	Close() error
	Reserve(timeout time.Duration) (id uint64, body []byte, err error)
	Put(tube string, body []byte, pri uint32, delay, ttr time.Duration) (id uint64, err error)
	Delete(id uint64) error
	Release(id uint64, pri uint32, delay time.Duration) error
	Bury(id uint64, pri uint32) error
	KickJob(id uint64) error
	Touch(id uint64) error
	Peek(id uint64) (body []byte, err error)
	Stats() (map[string]string, error)
	StatsJob(id uint64) (map[string]string, error)
	ListTubes() ([]string, error)
}

type Conn struct {
	sync.Pool
	Tubes []string
}

var ErrClosed = errors.New("TODO MARTA")

func TimeoutErr(err error) bool {
	cErr, ok := err.(beanstalk.ConnError)
	if !ok {
		return false
	}

	return cErr.Err == beanstalk.ErrTimeout
}

func (c *Conn) Close() error {
	c.Pool.New = func() interface{} {
		return (*beanstalk.Conn)(nil)
	}

	for ci := c.Pool.Get(); ci != nil; ci = c.Pool.Get() {
		bc, ok := ci.(*beanstalk.Conn)
		if !ok {
			log.Warn().Msg("could not cast connection while closing")
			continue
		}

		log.Err(bc.Close()).Msg("closed connection")
	}

	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.Close()
}

func (c *Conn) Reserve(timeout time.Duration) (id uint64, body []byte, err error) {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return 0, nil, ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.TubeSet.Reserve(timeout)
}

func (c *Conn) Put(tube string, body []byte, pri uint32, delay, ttr time.Duration) (id uint64, err error) {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return 0, ErrClosed
	}

	defer c.Pool.Put(bsc)
	bsc.Tube.Name = tube
	return bsc.Put(body, pri, delay, ttr)
}

func (c *Conn) Delete(id uint64) error {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.Delete(id)
}

func (c *Conn) Release(id uint64, pri uint32, delay time.Duration) error {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.Release(id, pri, delay)
}

func (c *Conn) Bury(id uint64, pri uint32) error {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.Bury(id, pri)
}

func (c *Conn) KickJob(id uint64) error {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.KickJob(id)
}

func (c *Conn) Touch(id uint64) error {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.Touch(id)
}

func (c *Conn) Peek(id uint64) (body []byte, err error) {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return nil, ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.Peek(id)
}

func (c *Conn) Stats() (map[string]string, error) {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return nil, ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.Stats()
}

func (c *Conn) StatsJob(id uint64) (map[string]string, error) {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return nil, ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.StatsJob(id)
}

func (c *Conn) ListTubes() ([]string, error) {
	bsc := c.Get().(*beanstalk.Conn)
	if bsc == nil {
		return nil, ErrClosed
	}

	defer c.Pool.Put(bsc)
	return bsc.ListTubes()
}
