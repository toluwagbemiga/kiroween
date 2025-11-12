# Analytics Service - Customization Guide

## Table of Contents
1. [Modifying Batch Configuration](#modifying-batch-configuration)
2. [Adding Custom Event Properties](#adding-custom-event-properties)
3. [Creating Custom Providers](#creating-custom-providers)
4. [Event Filtering & Transformation](#event-filtering--transformation)
5. [Advanced Retry Strategies](#advanced-retry-strategies)
6. [Performance Tuning](#performance-tuning)

---

## Modifying Batch Configuration

### Change Batch Size

```bash
# In .env
BATCH_SIZE=100  # Increase from 50 to 100 events
```

**When to increase:**
- High event volume (>1000 events/sec)
- Reduce API call frequency
- Lower network overhead

**When to decrease:**
- Low event volume
- Need faster delivery
- Memory constraints

### Change Flush Interval

```bash
# In .env
FLUSH_INTERVAL_SECONDS=5  # Decrease from 10 to 5 seconds
```

**When to decrease:**
- Need real-time analytics
- Low event volume
- Faster feedback loops

**When to increase:**
- High event volume
- Reduce API calls
- Lower costs

### Dynamic Configuration

```go
// In internal/config/config.go

// Add dynamic configuration support
type AnalyticsConfig struct {
    BatchSize         int
    FlushIntervalSec  int
    MinBatchSize      int  // NEW: Minimum batch size before flush
    MaxBatchAge       int  // NEW: Maximum age of oldest event in seconds
}

// In internal/batch_worker.go

func (w *BatchWorker) run() {
    for {
        select {
        case <-w.queue.FlushChannel():
            // Batch size reached
            w.flush()
            
        case <-w.flushTimer.C:
            // Timer expired - check if we have minimum batch
            if w.queue.Size() >= w.minBatchSize {
                w.flush()
            }
            w.resetTimer()
            
        case <-w.maxAgeTicker.C:
            // Check oldest event age
            if w.queue.OldestEventAge() > w.maxBatchAge {
                w.flush()
            }
        }
    }
}
```

---

## Adding Custom Event Properties

### Automatic Property Enrichment

```go
// In internal/grpc_handlers.go

func (s *AnalyticsServer) TrackEvent(ctx context.Context, req *pb.TrackEventRequest) (*pb.TrackEventResponse, error) {
    // Parse user properties
    properties := make(map[string]interface{})
    for key, propValue := range req.Properties {
        properties[key] = convertPropertyValue(propValue)
    }
    
    // Add automatic properties
    properties = s.enrichProperties(ctx, properties, req)
    
    event := Event{
        ID:         eventID,
        EventName:  req.EventName,
        UserID:     req.UserId,
        Properties: properties,
        Timestamp:  time.Now(),
    }
    
    s.queue.Add(event)
    return &pb.TrackEventResponse{Success: true, EventId: eventID}, nil
}

// enrichProperties adds automatic properties to events
func (s *AnalyticsServer) enrichProperties(ctx context.Context, properties map[string]interface{}, req *pb.TrackEventRequest) map[string]interface{} {
    // Add timestamp
    properties["$timestamp"] = time.Now().Unix()
    
    // Add service version
    properties["$service_version"] = "1.0.0"
    
    // Add environment
    properties["$environment"] = os.Getenv("ENVIRONMENT")
    
    // Extract metadata from gRPC context
    if md, ok := metadata.FromIncomingContext(ctx); ok {
        // Add IP address
        if ips := md.Get("x-forwarded-for"); len(ips) > 0 {
            properties["$ip"] = ips[0]
        }
        
        // Add user agent
        if agents := md.Get("user-agent"); len(agents) > 0 {
            properties["$user_agent"] = agents[0]
        }
        
        // Add request ID for tracing
        if reqIDs := md.Get("x-request-id"); len(reqIDs) > 0 {
            properties["$request_id"] = reqIDs[0]
        }
    }
    
    // Add session information (if available)
    if sessionID := getSessionIDFromContext(ctx); sessionID != "" {
        properties["$session_id"] = sessionID
    }
    
    // Add device information
    properties["$platform"] = detectPlatform(ctx)
    
    return properties
}

func detectPlatform(ctx context.Context) string {
    if md, ok := metadata.FromIncomingContext(ctx); ok {
        if agents := md.Get("user-agent"); len(agents) > 0 {
            ua := agents[0]
            if strings.Contains(ua, "Mobile") {
                return "mobile"
            } else if strings.Contains(ua, "Tablet") {
                return "tablet"
            }
            return "desktop"
        }
    }
    return "unknown"
}
```

### Custom Property Validators

```go
// In internal/validators.go

type PropertyValidator struct {
    maxStringLength int
    maxArrayLength  int
    allowedKeys     map[string]bool
}

func NewPropertyValidator() *PropertyValidator {
    return &PropertyValidator{
        maxStringLength: 1000,
        maxArrayLength:  100,
        allowedKeys:     nil, // nil = allow all
    }
}

func (v *PropertyValidator) Validate(properties map[string]interface{}) error {
    for key, value := range properties {
        // Check if key is allowed
        if v.allowedKeys != nil && !v.allowedKeys[key] {
            return fmt.Errorf("property key not allowed: %s", key)
        }
        
        // Validate value
        if err := v.validateValue(value); err != nil {
            return fmt.Errorf("invalid property %s: %w", key, err)
        }
    }
    return nil
}

func (v *PropertyValidator) validateValue(value interface{}) error {
    switch v := value.(type) {
    case string:
        if len(v) > v.maxStringLength {
            return fmt.Errorf("string too long (max %d)", v.maxStringLength)
        }
    case []interface{}:
        if len(v) > v.maxArrayLength {
            return fmt.Errorf("array too long (max %d)", v.maxArrayLength)
        }
    }
    return nil
}
```

---

## Creating Custom Providers

### Custom Provider Template

```go
// In internal/custom_provider.go

type CustomProvider struct {
    apiKey     string
    apiURL     string
    httpClient *http.Client
    logger     *zap.Logger
    testMode   bool
}

func NewCustomProvider(apiKey string, testMode bool, logger *zap.Logger) *CustomProvider {
    return &CustomProvider{
        apiKey: apiKey,
        apiURL: "https://api.custom-analytics.com/v1/events",
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        logger:   logger,
        testMode: testMode,
    }
}

func (p *CustomProvider) SendBatch(ctx context.Context, events []Event) error {
    if p.testMode {
        p.logger.Info("TEST MODE: would send batch to Custom Provider",
            zap.Int("event_count", len(events)))
        time.Sleep(100 * time.Millisecond)
        return nil
    }
    
    // Convert events to custom format
    customEvents := p.convertEvents(events)
    
    // Prepare request
    payload := map[string]interface{}{
        "api_key": p.apiKey,
        "events":  customEvents,
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal events: %w", err)
    }
    
    // Send HTTP request
    req, err := http.NewRequestWithContext(ctx, "POST", p.apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+p.apiKey)
    
    resp, err := p.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
    }
    
    p.logger.Debug("batch sent to Custom Provider",
        zap.Int("event_count", len(events)),
        zap.Int("status_code", resp.StatusCode))
    
    return nil
}

func (p *CustomProvider) convertEvents(events []Event) []map[string]interface{} {
    customEvents := make([]map[string]interface{}, len(events))
    
    for i, event := range events {
        customEvents[i] = map[string]interface{}{
            "id":         event.ID,
            "name":       event.EventName,
            "user_id":    event.UserID,
            "properties": event.Properties,
            "timestamp":  event.Timestamp.Unix(),
        }
    }
    
    return customEvents
}

func (p *CustomProvider) GetName() string {
    return "custom"
}
```

### Database Provider (PostgreSQL)

```go
// In internal/postgres_provider.go

type PostgresProvider struct {
    db     *sql.DB
    logger *zap.Logger
}

func NewPostgresProvider(connString string, logger *zap.Logger) (*PostgresProvider, error) {
    db, err := sql.Open("postgres", connString)
    if err != nil {
        return nil, err
    }
    
    return &PostgresProvider{
        db:     db,
        logger: logger,
    }, nil
}

func (p *PostgresProvider) SendBatch(ctx context.Context, events []Event) error {
    tx, err := p.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO analytics_events (id, event_name, user_id, properties, timestamp)
        VALUES ($1, $2, $3, $4, $5)
    `)
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()
    
    for _, event := range events {
        propertiesJSON, err := json.Marshal(event.Properties)
        if err != nil {
            return fmt.Errorf("failed to marshal properties: %w", err)
        }
        
        _, err = stmt.ExecContext(ctx,
            event.ID,
            event.EventName,
            event.UserID,
            propertiesJSON,
            event.Timestamp,
        )
        if err != nil {
            return fmt.Errorf("failed to insert event: %w", err)
        }
    }
    
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    p.logger.Debug("batch saved to PostgreSQL",
        zap.Int("event_count", len(events)))
    
    return nil
}

func (p *PostgresProvider) GetName() string {
    return "postgres"
}
```

---

## Event Filtering & Transformation

### Event Filters

```go
// In internal/filters.go

type EventFilter interface {
    ShouldProcess(event Event) bool
}

// Filter by event name
type EventNameFilter struct {
    allowedEvents map[string]bool
}

func NewEventNameFilter(allowedEvents []string) *EventNameFilter {
    allowed := make(map[string]bool)
    for _, name := range allowedEvents {
        allowed[name] = true
    }
    return &EventNameFilter{allowedEvents: allowed}
}

func (f *EventNameFilter) ShouldProcess(event Event) bool {
    return f.allowedEvents[event.EventName]
}

// Filter by user ID
type UserIDFilter struct {
    excludedUsers map[string]bool
}

func NewUserIDFilter(excludedUsers []string) *UserIDFilter {
    excluded := make(map[string]bool)
    for _, userID := range excludedUsers {
        excluded[userID] = true
    }
    return &UserIDFilter{excludedUsers: excluded}
}

func (f *UserIDFilter) ShouldProcess(event Event) bool {
    return !f.excludedUsers[event.UserID]
}

// Apply filters in batch worker
func (w *BatchWorker) flush() {
    batch := w.queue.GetBatch()
    if len(batch) == 0 {
        return
    }
    
    // Apply filters
    filteredBatch := w.applyFilters(batch)
    
    if len(filteredBatch) == 0 {
        w.logger.Debug("all events filtered out")
        return
    }
    
    // Send filtered batch
    w.sendBatchWithRetry(context.Background(), filteredBatch)
}

func (w *BatchWorker) applyFilters(events []Event) []Event {
    if len(w.filters) == 0 {
        return events
    }
    
    filtered := make([]Event, 0, len(events))
    for _, event := range events {
        shouldProcess := true
        for _, filter := range w.filters {
            if !filter.ShouldProcess(event) {
                shouldProcess = false
                break
            }
        }
        if shouldProcess {
            filtered = append(filtered, event)
        }
    }
    
    return filtered
}
```

### Event Transformers

```go
// In internal/transformers.go

type EventTransformer interface {
    Transform(event Event) Event
}

// Anonymize PII
type PIIAnonymizer struct {
    sensitiveKeys []string
}

func NewPIIAnonymizer() *PIIAnonymizer {
    return &PIIAnonymizer{
        sensitiveKeys: []string{"email", "phone", "ssn", "credit_card"},
    }
}

func (t *PIIAnonymizer) Transform(event Event) Event {
    if event.Properties == nil {
        return event
    }
    
    for _, key := range t.sensitiveKeys {
        if _, exists := event.Properties[key]; exists {
            event.Properties[key] = "[REDACTED]"
        }
    }
    
    return event
}

// Add computed properties
type PropertyEnricher struct{}

func (t *PropertyEnricher) Transform(event Event) Event {
    if event.Properties == nil {
        event.Properties = make(map[string]interface{})
    }
    
    // Add day of week
    event.Properties["$day_of_week"] = event.Timestamp.Weekday().String()
    
    // Add hour of day
    event.Properties["$hour_of_day"] = event.Timestamp.Hour()
    
    // Add is_weekend
    dow := event.Timestamp.Weekday()
    event.Properties["$is_weekend"] = dow == time.Saturday || dow == time.Sunday
    
    return event
}
```

---

## Advanced Retry Strategies

### Exponential Backoff with Jitter

```go
// In internal/batch_worker.go

func (w *BatchWorker) sendBatchWithRetry(ctx context.Context, batch []Event) error {
    var lastErr error
    delay := w.retryConfig.InitialDelay
    
    for attempt := 1; attempt <= w.retryConfig.MaxAttempts; attempt++ {
        err := w.provider.SendBatch(ctx, batch)
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // Check if error is retryable
        if !w.isRetryable(err) {
            return err
        }
        
        if attempt == w.retryConfig.MaxAttempts {
            break
        }
        
        // Add jitter to prevent thundering herd
        jitter := time.Duration(rand.Float64() * float64(delay) * 0.1)
        actualDelay := delay + jitter
        
        w.logger.Warn("batch send failed, retrying",
            zap.Int("attempt", attempt),
            zap.Duration("delay", actualDelay),
            zap.Error(err))
        
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(actualDelay):
        }
        
        // Exponential backoff
        delay = time.Duration(float64(delay) * w.retryConfig.BackoffFactor)
        if delay > w.retryConfig.MaxDelay {
            delay = w.retryConfig.MaxDelay
        }
    }
    
    return fmt.Errorf("max retry attempts exceeded: %w", lastErr)
}

func (w *BatchWorker) isRetryable(err error) bool {
    // Don't retry on client errors (4xx)
    if strings.Contains(err.Error(), "400") ||
       strings.Contains(err.Error(), "401") ||
       strings.Contains(err.Error(), "403") ||
       strings.Contains(err.Error(), "404") {
        return false
    }
    
    // Retry on server errors (5xx) and network errors
    return true
}
```

### Circuit Breaker Pattern

```go
// In internal/circuit_breaker.go

type CircuitBreaker struct {
    maxFailures  int
    resetTimeout time.Duration
    failures     int
    lastFailTime time.Time
    state        string // "closed", "open", "half-open"
    mu           sync.Mutex
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures:  maxFailures,
        resetTimeout: resetTimeout,
        state:        "closed",
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    
    // Check if circuit should reset
    if cb.state == "open" && time.Since(cb.lastFailTime) > cb.resetTimeout {
        cb.state = "half-open"
        cb.failures = 0
    }
    
    // Reject if circuit is open
    if cb.state == "open" {
        cb.mu.Unlock()
        return fmt.Errorf("circuit breaker is open")
    }
    
    cb.mu.Unlock()
    
    // Execute function
    err := fn()
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        
        return err
    }
    
    // Success - reset circuit
    if cb.state == "half-open" {
        cb.state = "closed"
    }
    cb.failures = 0
    
    return nil
}
```

---

## Performance Tuning

### Memory Optimization

```go
// In internal/types.go

// Use sync.Pool for event objects
var eventPool = sync.Pool{
    New: func() interface{} {
        return &Event{
            Properties: make(map[string]interface{}),
        }
    },
}

func GetEvent() *Event {
    return eventPool.Get().(*Event)
}

func PutEvent(e *Event) {
    // Clear properties
    for k := range e.Properties {
        delete(e.Properties, k)
    }
    e.ID = ""
    e.EventName = ""
    e.UserID = ""
    eventPool.Put(e)
}
```

### Batch Compression

```go
// In internal/provider.go

func (p *MixpanelProvider) SendBatch(ctx context.Context, events []Event) error {
    // Convert events
    mixpanelEvents := p.convertEvents(events)
    
    // Marshal to JSON
    jsonData, err := json.Marshal(mixpanelEvents)
    if err != nil {
        return err
    }
    
    // Compress with gzip
    var buf bytes.Buffer
    gzipWriter := gzip.NewWriter(&buf)
    if _, err := gzipWriter.Write(jsonData); err != nil {
        return err
    }
    if err := gzipWriter.Close(); err != nil {
        return err
    }
    
    // Send compressed data
    req, err := http.NewRequestWithContext(ctx, "POST", p.apiURL, &buf)
    req.Header.Set("Content-Encoding", "gzip")
    req.Header.Set("Content-Type", "application/json")
    
    // ... send request
}
```

### Monitoring & Metrics

```go
// In internal/metrics.go

type Metrics struct {
    eventsReceived    prometheus.Counter
    eventsProcessed   prometheus.Counter
    batchesFlushed    prometheus.Counter
    flushDuration     prometheus.Histogram
    queueSize         prometheus.Gauge
    providerErrors    *prometheus.CounterVec
}

func NewMetrics() *Metrics {
    return &Metrics{
        eventsReceived: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "analytics_events_received_total",
            Help: "Total number of events received",
        }),
        eventsProcessed: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "analytics_events_processed_total",
            Help: "Total number of events processed",
        }),
        batchesFlushed: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "analytics_batches_flushed_total",
            Help: "Total number of batches flushed",
        }),
        flushDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name:    "analytics_flush_duration_seconds",
            Help:    "Duration of batch flush operations",
            Buckets: prometheus.DefBuckets,
        }),
        queueSize: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "analytics_queue_size",
            Help: "Current size of event queue",
        }),
        providerErrors: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "analytics_provider_errors_total",
                Help: "Total number of provider errors",
            },
            []string{"provider"},
        ),
    }
}
```

---

## Summary

The analytics service is now fully customizable:

✅ **Batch Configuration**: Adjustable size and timing  
✅ **Custom Properties**: Automatic enrichment  
✅ **Custom Providers**: Easy to add new destinations  
✅ **Event Filtering**: Process only relevant events  
✅ **Advanced Retry**: Exponential backoff with jitter  
✅ **Performance**: Memory optimization and compression  

All customizations maintain the non-blocking architecture! 🎃
