# subly
A lightweight subscription layer with conventions for [NATS](http://nats.io/)

Package subly helps with subscribing methods on a struct type as callbacks for NATS, with some naming conventions. A sample usage would look like:

```go
s := NewSubscriber(ctx, econn)
s.Subscribe(&timeService{econn})
```

Assuming we have a service like:

```go
type someService struct{}

func (*someService) SubActionMessage(p *person) {
    /**/
}

func (*someService) RepActionMessageQueue(subject, reply string, p *person) {
    /**/
}
```

then SubActionMessage would get subscribed to subject `someservice.subaction` and RepActionMessageQueue would get subscribed to subject `someservice.repaction`. Subject naming convension is <struct type name>.<method name> all lower case,
with words message and queue removed from the end.

If a method name ends in Message, it will subscribe to subject as a normall
subscriber (just receiving). If a method name ends in MessageQueue, it will subscribe
to subject as a member of a queue and the queue name will be <struct type name>_<method name>.

Message methods are expected to have one of four signatures.

```go
type person struct {
	Name string `json:"name,omitempty"`
	Age  uint   `json:"age,omitempty"`
}

handler := func(m *Msg)
handler := func(p *person)
handler := func(subject string, o *obj)
handler := func(subject, reply string, o *obj)
```

Which are NATS's conventions for callbacks.

And the callback methods will unsubscribe from subject when context got canceled.

## status
alpha - thinkering, in it's early stages of real world usage/feedback