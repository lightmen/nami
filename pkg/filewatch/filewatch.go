// 监控文件修改后自动 reload
package filewatch

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/pkg/aerror"
	"github.com/lightmen/nami/pkg/safe"
)

type fileNode struct {
	load func(f string) error
}

type watcher struct {
	*fsnotify.Watcher
	fileMap sync.Map
}

var gw *watcher

func init() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	gw = &watcher{
		Watcher: w,
	}

	go func() {
		for {
			select {
			case evt := <-gw.Events:
				if !evt.Op.Has(fsnotify.Write) {
					break
				}
				alog.Info("file changed: %s, op: %s", evt.Name, evt.Op.String())

				val, ok := gw.fileMap.Load(evt.Name)
				if !ok {
					alog.Error("cannot found load function: %s", evt.Name)
					break
				}

				fn := val.(*fileNode)

				safe.Func(func() {
					if err = fn.load(evt.Name); err != nil {
						// 业务使用者会打印日志的，所以这里不需要
						return
					}
				})

			case e := <-gw.Errors:
				alog.Fatal("file watcher error: %s", e.Error())
			}
		}
	}()
}

// @name 为文件绝对路径
func Add(name string, load func(f string) error) (err error) {
	if name == "" {
		err = aerror.InvalidParam
		return
	}

	if err = gw.Add(name); err != nil {
		return
	}

	if err = load(name); err != nil {
		return
	}

	fn := &fileNode{
		load: load,
	}

	gw.fileMap.Store(name, fn)
	return nil
}

// @name 为文件绝对路径
func Remove(name string) error {
	gw.fileMap.Delete(name)
	return gw.Remove(name)
}
