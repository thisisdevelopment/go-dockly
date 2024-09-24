package iferr

import (
	"net/http"
	"os"
	"reflect"

	"dario.cat/mergo"
	"github.com/logrusorgru/aurora"
	errors "github.com/thisisdevelopment/go-dockly/v2/xerrors"
	"github.com/thisisdevelopment/go-dockly/v2/xlogger"
)

var Default *IfErr

func init() {
	var err error

	log, err := xlogger.New(&xlogger.Config{
		CallerDepth: 4,
	})

	if err != nil {
		log.Errorf("Failed to initialize the default IfErr %+v\n", err)
	}

	Default, err = New(WithLogger(log))

	if err != nil {
		log.Errorf("Failed to initialize the default IfErr %+v\n", err)
	}
}

type IfErr struct {
	*Options
}

type Options struct {
	Log     *xlogger.Logger
	Verbose *bool
}

type OptionFunc func(*Options)

func New(optFns ...OptionFunc) (*IfErr, error) {
	opts := &Options{}

	for _, optFn := range optFns {
		optFn(opts)
	}

	if opts.Verbose == nil {
		true := true
		opts.Verbose = &true
	}

	if opts.Log == nil {
		log, err := xlogger.New(&xlogger.Config{
			CallerDepth: 3,
		})

		if err != nil {
			return nil, err
		}

		opts.Log = log
	}

	return &IfErr{
		Options: opts,
	}, nil
}

func WithLogger(log *xlogger.Logger) OptionFunc {
	return func(o *Options) {
		o.Log = log
	}
}

func Verbose(flag bool) OptionFunc {
	return func(o *Options) {
		o.Verbose = &flag
	}
}

func Warn(err error) { Default.Warn(err) }
func (ie *IfErr) Warn(err error) {
	if err != nil {
		if *ie.Verbose {
			ie.Log.Warnf("%+v\n", aurora.BrightRed(err))
		} else {
			ie.Log.Warnf("%v\n", aurora.BrightRed(err))
		}
	}
}

func Exit(err error, ctx ...string) { Default.Exit(err, ctx) }
func (ie *IfErr) Exit(err error, ctx []string) {
	if err != nil {
		var context = "no recover: "
		if len(ctx) > 0 {
			context = ctx[0]
		}
		if *ie.Verbose {
			ie.Log.Error(aurora.Sprintf("%s %+v", aurora.Yellow(context), aurora.BrightRed(err)))
		} else {
			ie.Log.Error(aurora.Sprintf("%s %v", aurora.Yellow(context), aurora.BrightRed(err)))
		}
		os.Exit(-1)
	}
}

func Panic(err error, ctx ...string) { Default.Panic(err, ctx) }
func (ie *IfErr) Panic(err error, ctx []string) {
	if err != nil {
		var message string
		var context = "panic: "
		if len(ctx) > 0 {
			context = ctx[0]
		}

		if *ie.Verbose {
			message = aurora.Sprintf("%s %+v", aurora.Yellow(context), aurora.BrightRed(err))
		} else {
			message = aurora.Sprintf("%s %v", aurora.Yellow(context), aurora.BrightRed(err))
		}

		panic(message)
	}
}

type Fataler interface {
	Fatalf(format string, args ...interface{})
}

func Fail(f Fataler, err error) { Default.Fail(f, err) }
func (ie *IfErr) Fail(f Fataler, err error) {
	if err != nil {
		f.Fatalf("%+v\n", aurora.BrightRed(err))
	}
}

type ResponseOpts struct {
	Code    int
	Message string
	Depth   int
}

func Respond(w http.ResponseWriter, err error, opts ...*ResponseOpts) bool {
	return Default.Respond(w, err, opts...)
}
func (ie *IfErr) Respond(w http.ResponseWriter, err error, opts ...*ResponseOpts) bool {
	if err == nil {
		return false
	}

	opt := &ResponseOpts{
		Code: http.StatusInternalServerError,
	}

	if len(opts) > 0 {
		err := Merge(opt, opts[0])
		ie.Warn(err)
	}

	if opt.Message == "" {
		opt.Message = errors.Message(err, opt.Depth)
	}

	if opt.Message == "" {
		opt.Message = err.Error()
	}

	if *ie.Verbose {
		ie.Log.Error("HTTP Error: ", err)
	} else {
		ie.Log.Error("HTTP Error: ", err)
	}

	http.Error(w, opt.Message, opt.Code)
	return true
}

func Merge(dst, src interface{}) error {
	if src != nil && !reflect.ValueOf(src).IsNil() {
		return mergo.Merge(dst, src, mergo.WithOverride)
	}

	return nil
}
