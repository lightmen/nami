package tracing

import (
	"context"
	"fmt"
	"reflect"

	"github.com/lightmen/nami/pkg/aerror"
	"github.com/lightmen/nami/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type Tracer struct {
	tracer trace.Tracer
	kind   trace.SpanKind
	opt    *options
}

func NewTracer(kind trace.SpanKind, opts ...Option) *Tracer {
	opt := &options{
		name:       "nami",
		prefix:     []string{"x-md"},
		propagator: propagation.NewCompositeTextMapPropagator(Metadata{}, propagation.Baggage{}, propagation.TraceContext{}),
	}
	for _, o := range opts {
		o(opt)
	}

	if opt.provider != nil {
		otel.SetTracerProvider(opt.provider)
	}

	if kind != trace.SpanKindClient && kind != trace.SpanKindServer {
		panic(fmt.Sprintf("unsupported span kind: %v", kind))
	}

	tr := &Tracer{
		tracer: otel.Tracer(opt.name),
		kind:   kind,
		opt:    opt,
	}

	return tr
}

func (t *Tracer) Start(ctx context.Context, spanName string, carrier propagation.TextMapCarrier) (context.Context, trace.Span) {
	if t.kind == trace.SpanKindServer {
		ctx = t.opt.propagator.Extract(ctx, carrier)
	}
	ctx, span := t.tracer.Start(ctx, spanName, trace.WithSpanKind(t.kind))

	if t.kind == trace.SpanKindClient {
		t.opt.propagator.Inject(ctx, carrier)
	}

	return ctx, span
}

func (t *Tracer) End(ctx context.Context, span trace.Span, msg any, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Key("rpc.status_code").Int64(int64(aerror.Code(err))))
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "OK")
	}

	var attrs []attribute.KeyValue
	var tr transport.Transporter
	var key string

	if t.kind == trace.SpanKindServer {
		tr, _ = transport.FromServerContext(ctx)
		key = "send_msg.size"
	} else {
		tr, _ = transport.FromClientContext(ctx)
		key = "recv_msg.size"
	}

	if p, ok := msg.(proto.Message); ok && !reflect.ValueOf(p).IsNil() {
		attrs = append(attrs, attribute.Key(key).Int(proto.Size(p)))
	}

	if tr != nil {
		header := tr.Header()
		keys := header.Keys()
		for _, key := range keys {
			value := header.Get(key)
			if value == "" {
				continue
			}
			if t.opt.hasPrefix(key) {
				attrs = append(attrs, attribute.Key(key).String(value))
			}
		}
	}

	span.SetAttributes(attrs...)

	span.End()
}
