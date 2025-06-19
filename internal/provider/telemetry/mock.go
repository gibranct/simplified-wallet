package telemetry

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

// MockTelemetry is a mock implementation of the Telemetry interface
type MockTelemetry struct {
	mock.Mock
}

// Start mocks the Start method of the Telemetry interface
func (m *MockTelemetry) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	args := m.Called(ctx, name, opts)
	return args.Get(0).(context.Context), args.Get(1).(Span)
}

// Shutdown mocks the Shutdown method of the Telemetry interface
func (m *MockTelemetry) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockSpan is a mock implementation of the Span interface
type MockSpan struct {
	mock.Mock
	embedded.Span
}

// End mocks the End method of the trace.Span interface
func (m *MockSpan) End(options ...trace.SpanEndOption) {
	m.Called(options)
}

// AddEvent mocks the AddEvent method of the trace.Span interface
func (m *MockSpan) AddEvent(name string, options ...trace.EventOption) {
	m.Called(name, options)
}

// AddLink mocks the AddLink method of the trace.Span interface
func (m *MockSpan) AddLink(link trace.Link) {
	m.Called(link)
}

// IsRecording mocks the IsRecording method of the trace.Span interface
func (m *MockSpan) IsRecording() bool {
	args := m.Called()
	return args.Bool(0)
}

// RecordError mocks the RecordError method of the trace.Span interface
func (m *MockSpan) RecordError(err error, options ...trace.EventOption) {
	m.Called(err, options)
}

// SpanContext mocks the SpanContext method of the trace.Span interface
func (m *MockSpan) SpanContext() trace.SpanContext {
	args := m.Called()
	return args.Get(0).(trace.SpanContext)
}

// SetStatus mocks the SetStatus method of the trace.Span interface
func (m *MockSpan) SetStatus(code codes.Code, description string) {
	m.Called(code, description)
}

// SetName mocks the SetName method of the trace.Span interface
func (m *MockSpan) SetName(name string) {
	m.Called(name)
}

// SetAttributes mocks the SetAttributes method of the trace.Span interface
func (m *MockSpan) SetAttributes(attributes ...attribute.KeyValue) {
	m.Called(attributes)
}

// TracerProvider mocks the TracerProvider method of the trace.Span interface
func (m *MockSpan) TracerProvider() trace.TracerProvider {
	args := m.Called()
	return args.Get(0).(trace.TracerProvider)
}

// NewMockTelemetry creates a new instance of MockTelemetry
func NewMockTelemetry() *MockTelemetry {
	mockSpan := NewMockSpan()
	mockTelemetry := &MockTelemetry{}
	mockTelemetry.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(context.Background(), mockSpan)
	mockSpan.On("End", mock.Anything).Return()
	mockSpan.On("SetAttributes", mock.Anything).Return()
	return mockTelemetry
}

// NewMockSpan creates a new instance of MockSpan
func NewMockSpan() *MockSpan {
	return &MockSpan{}
}
