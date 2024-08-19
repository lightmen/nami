package alog

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/lightmen/nami/pkg/cast"
)

var _ Handler = (*handler)(nil)

type metaFunc func(context.Context) map[string]any

func defaultBuildMetadata(_ context.Context) (md map[string]any) {
	return md
}

type handler struct {
	level         atomic.Value
	bufPool       *sync.Pool
	fw            atomic.Value
	pid           []byte
	processName   string
	buildMetadata metaFunc
}

func NewFileHandler(appName, dir string, buildMeta metaFunc) *handler {
	dir = fmt.Sprintf("%s/%s", dir, appName)
	fw, err := newFileWriter(appName, dir)
	if err != nil {
		return nil
	}

	pool := &sync.Pool{New: func() any {
		buf := make([]byte, 0, 1024)
		bb := bytes.NewBuffer(buf)
		return bb
	}}

	pid := os.Getpid()

	h := &handler{
		bufPool:       pool,
		pid:           []byte(cast.ToString(pid)),
		processName:   getExeName(),
		buildMetadata: buildMeta,
	}
	h.fw.Store(fw)
	h.SetLevel(LevelInfo)

	return h
}

func (h *handler) SetLevel(level Level) {
	h.level.Store(level)
}

func (h *handler) Level() Level {
	return h.level.Load().(Level)
}

func (h *handler) Enabled(ctx context.Context, l Level) bool {
	return l >= h.Level()
}

func (h *handler) Handle(ctx context.Context, r Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	h.output(ctx, r)
	return nil
}

func (h *handler) output(ctx context.Context, r Record) {
	msg := fmt.Sprintf(r.Message, r.args...)
	info := r.info

	buffer := h.bufPool.Get().(*bytes.Buffer)
	defer h.bufPool.Put(buffer)
	buffer.Reset()

	md := h.buildMetadata(ctx)
	if md == nil {
		md = map[string]any{} //md不能为空，这样在打印日志转换为json的时候，可以输出 "{}"
	}

	//输出格式：[2006-01-02 15:04:05.000 -0700]	LogLevel pid appName-fileName:line msg {}
	buffer.WriteString(r.Time.Format("[2006-01-02 15:04:05.000 -0700]"))
	buffer.WriteByte('\t')
	buffer.WriteString(r.Level.String())
	buffer.WriteByte('\t')
	buffer.Write(h.pid)
	buffer.WriteByte('\t')
	buffer.WriteString(h.processName)
	buffer.WriteByte('-')
	buffer.WriteString(path.Base(info.fileName))
	buffer.WriteByte(':')
	buffer.WriteString(cast.ToString(info.line))
	buffer.WriteByte('\t')
	buffer.WriteString(info.funcName)
	buffer.WriteByte('\t')
	buffer.WriteString(msg)
	buffer.WriteByte('\t')
	buffer.WriteByte('-')
	buffer.WriteByte('\t')
	buffer.WriteString(cast.ToJson(md))
	buffer.WriteByte('\n')

	h.fw.Load().(*fileWriter).write(buffer.Bytes())
}

type fileWriter struct {
	appName string
	dir     string
	hour    int
	ticker  *time.Ticker
	fd      atomic.Value // *os.File
	watcher *fsnotify.Watcher
	sync.RWMutex
}

func newFileWriter(appName, dir string) (wr *fileWriter, err error) {
	st, err := os.Stat(dir)
	if err != nil {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}
	} else if !st.IsDir() {
		err = errors.New("dir " + dir + " illegal")
		return
	}

	wr = &fileWriter{
		appName: appName,
		dir:     dir,
		hour:    -1, // 小时存在0的情况，这里默认设置为-1
		ticker:  time.NewTicker(time.Second),
	}
	wr.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}

	wr.checkAndUpdate()
	go wr.watch()
	return
}

func (f *fileWriter) watch() {
	for {
		select {
		case evt := <-f.watcher.Events:
			f.handleFileEvent(evt)
		case <-f.ticker.C:
			f.checkAndUpdate()
		}
	}
}

func (f *fileWriter) handleFileEvent(evt fsnotify.Event) {
	if !evt.Op.Has(fsnotify.Chmod) {
		return
	}

	f.hour = 0 //当前文件被删除，置空内存记录的时间
	f.checkAndUpdate()
}

func (f *fileWriter) updateFd() {
	fn := f.getCurFileName()
	dir := path.Dir(fn)

	_, err := os.Stat(dir)
	if err != nil {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}
	}

	fd, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return
	}

	f.fd.Store(fd)
	wtList := f.watcher.WatchList()
	for _, name := range wtList {
		f.watcher.Remove(name)
	}
	f.watcher.Add(fn)
}

func (f *fileWriter) checkAndUpdate() {
	now := time.Now()
	if f.hour == now.Hour() {
		return
	}

	f.hour = now.Hour()
	f.updateFd()
}

func (f *fileWriter) write(bytes []byte) (n int, err error) {
	return f.fd.Load().(*os.File).Write(bytes)
}

func (f *fileWriter) getCurFileName() string {
	now := time.Now()
	return path.Clean(fmt.Sprintf("%s/%d%02d%02d/%s_%02d.log", f.dir, now.Year(), now.Month(), now.Day(), f.appName, now.Hour()))
}

func getExeName() (output string) {
	path, err := os.Executable()
	if err != nil {
		return ""
	}

	_, output = filepath.Split(path)
	return
}
