package eastmoney

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"github.com/letsfire/factory"
	"sync"
)

var pool *factory.Master

func init() {
	pool = factory.NewMaster(6, 2)
}

func newSpider(async bool) *colly.Collector {
	spider := colly.NewCollector()
	spider.Async = async
	spider.IgnoreRobotsTxt = false
	return spider
}

func assertError(exp bool, msg string) error {
	if exp {
		return nil
	}
	return errors.New(msg)
}

func callWithoutErr(err error, cbs ...func() error) error {
	if err != nil {
		return err
	}
	for _, cb := range cbs {
		if e := cb(); e != nil {
			return e
		}
	}
	return nil
}

func callWithoutErr2(err error, cbs ...func()) error {
	if err != nil {
		return err
	}
	for _, cb := range cbs {
		cb()
	}
	return nil
}

type Errors []error

var errorLocker = new(sync.Mutex)

func (es *Errors) add(errs ...error) {
	errorLocker.Lock()
	for _, err := range errs {
		if err != nil {
			*es = append(*es, err)
		}
	}
	errorLocker.Unlock()
}

func (es *Errors) first() error {
	for _, err := range *es {
		if err != nil {
			return err
		}
	}
	return nil
}

func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
