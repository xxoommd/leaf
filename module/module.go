package module

import (
	"sync"

	"github.com/xxoommd/leaf/conf"
	"github.com/xxoommd/leaf/log"
	"go.uber.org/zap"
)

type Module interface {
	OnInit()
	BeforeDestroy()
	OnDestroy()
	Run(closeSig chan bool)
}

type module struct {
	mi       Module
	closeSig chan bool
	wg       sync.WaitGroup
}

var mods []*module

func Register(mi Module) {
	m := new(module)
	m.mi = mi
	m.closeSig = make(chan bool, 1)

	mods = append(mods, m)
}

func Init() {
	for i := 0; i < len(mods); i++ {
		mods[i].mi.OnInit()
	}

	for i := 0; i < len(mods); i++ {
		m := mods[i]
		m.wg.Add(1)
		go run(m)
	}
}

func Destroy() {
	for i := len(mods) - 1; i >= 0; i-- {
		m := mods[i]
		m.mi.BeforeDestroy()
		m.closeSig <- true
		m.wg.Wait()
		destroy(m)
	}
}

func run(m *module) {
	m.mi.Run(m.closeSig)
	m.wg.Done()
}

func destroy(m *module) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				// buf := make([]byte, conf.LenStackBuf)
				// l := runtime.Stack(buf, false)
				// log.Error("%v: %s", r, buf[:l])
				log.ZapLogger.Error("recovery stack", zap.Stack("recovery"))
			} else {
				log.ZapLogger.Error("recovery error", zap.Any("error", r))
			}
		}
	}()

	m.mi.OnDestroy()
}
