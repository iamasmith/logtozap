# logtozap - golang log package routing to Zap loggers

We love [Zap logger](https://github.com/uber-go/zap) for producing our structured logs.

Many of the packages we use have logging implemented using the golang log package and as such any events they log mix in with our structured logs when using Zap.

This package routes standard log messages through to a Zap logger that you provide ensuring..

* Messages are trimmed of whitespace including newlines.
* Dates and Times are removed from the input message so that they don't appear in the message field within the structured log event and only as fields of that event.
* Call stack is unwound to the original caller, not this wrapper

e.g.
```
....
18: // foo.go
19: log.Print("Hello from unstructured land")
....
```
```
{"level":"warn","ts":1683968275.461857,"caller":"demo/foo.go:19","msg":"Hello from unstructured land","logger":"zap.SugaredLogger"}
```

The approach works by routing standard log output via an io.Writer to Zap. It will lack control of using Zap directly as it will not allow you to set 'level' on an individual event basis and it is solely entended to catch log messages from projects that you might have included in your own and retain the source of the log call within these messages to assist with debugging.

For clarity, you should use Zap directly in your project for best control but for packages you need to include that use the golang log package then you can use this approach for consistency.

# Use

Firstly set up your logger, as normal with any extra fields that you want to be common and in keeping with events you are going to be logging from your app using the regular Zap logger.

Here we use a sugared logger but you can also use a desugared one.

```
c := zap.NewProductionConfig()
logger := zap.Must(c.Build()).Sugar().With("myfield", "hassomething")
```
logtozap.ToSugared allows you to specify your Sugrared logger and a level field for all messages going to that log.
```
logtozap.ToSugared(logger, zapcore.InfoLevel)
```
Some packages that you include might let you specify your own logger but may be dependent on the log package interface. One good example is [go-redis](https://github.com/redis/go-redis) where it is possible to call redis.SetLogger prior to setting up a client and have all messages route to that.
If you add extra arguments to the ToSugared call then you can route messages from additional logs created with log.New().
In this form the call won't route the default log as well so if you wanted to route a logger for redis and the default log for other packages that you included then you would do something like this.
```
c := zap.NewProductionConfig()
logger := zap.Must(c.Build()).Sugar()
redisLogger := log.New(log.Writer(), "REDIS:", 2)
logtozap.ToSugared(logger, zapcore.WarnLevel, redisLogger)
redis.SetLogger(redisLogger)
logtozap.ToSugared(logger, zapcore.DebugLevel)
```
This would also cover reusing code that had separate logs for stdout and stderr allowing you to route them with appropriate levels into Zap.


If you are using the 'desugared' logging for your app then you should use ToLogger() instead of ToSugared - the parameters are the same.

# WithSkip variants

If you encounter another interface requirement in your logging you can add additional skip levels to the oned discovered by the trainer routines by using ToSugaredWithSkip or ToLoggerWithSkip that take a skip factor following the level parameter.

### redis-go/v9 example

I recently encoutered in redis-go v9 the requirement to pass a logger that includes a context, something that the standard log package does not do but can easily be achieved with a small wrapper.

Here's how that can be overcome.
```
// Create a new standard logger
lp := log.New(log.Writer(), "", 0)
// Wrap it using WithSkip value of 1
logtozap.ToSugaredWithSkip(logger, zapcore.WarnLevel, 1, lp)
// Use a wrapped logger calling the original logger
redis.SetLogger(&redisLogger{l: lp})
```
```
type redisLogger struct {
	l *log.Logger
}

func (r redisLogger) Printf(ctx context.Context, s string, args ...interface{}) {
	r.l.Printf(s, args...)
}
```
Adding the wrapper, however, means that Redis is calling the previously wrapped log function with extra depth and without using the WithSkip feature to increase the stack depth one will get the line number of the Printf statement in the wrapper and not the calling function within the redis library in the message.