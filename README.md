# Ideal format to implement

The JSON

```json
{
  "key-1": "message", // simple key->message
  "key-2": "message with {arg}", // key->message but with an `any` argument called `arg`,
  "key-3": "message with {arg:int} used twice {arg}", // key->message but with an `int` argument called `arg`, only need to specify the type ones
  "key-4": "message with {arg1:int} and {arg2:float64:.2f}", // key->message but with an `int` argument called `arg1` and a `float64` argument called `arg2` using the format `.2f`
  "key-5": {
    "key-1": "nested message" // nested message with key `key-5.key-1`
  },
  "?key-6": {
    // conditional key->message pair. Here, depending on the value of the argument `messages` one of the 3 options is selected
    "messages==0": "You don't have any messages.",
    "messages==1": "You have one new message.",
    "": "You have {messages:int} new messages." // else branch
  },
  "?key-7": {
    // conditional key->message pair. Here, depending on the value of the argument `messages` one of the 3 options is selected. Since we don't use the conditional variable in any message, it can't be infered from these. For these cases, an extra `_args` entry is added
    "_args": ["messages:int"],
    "messages==0": "You don't have any messages.",
    "messages==1": "You have one new message.",
    "": "You have several new messages." // else branch
  }
}
```

Should produce the output

```go
func MessagesFor(tag string) (Messages, error) {
    if tag == "en-EN" {
        return en_EN_Messages{}, nil
    }
    return nil, fmt.Errorf("unknown language tag: %s", tag)
}

func MessagesForMust(tag string) Messages {
    if tag == "en-EN" {
        return en_EN_Messages{}, nil
    }
    panic(fmt.Errorf("unknown language tag: %s", tag))
}

func MessagesForOrDefault(tag string) Messages {
    if tag == "en-EN" {
        return en_EN_Messages{}
    }
    return en_EN_Messages{} // when calling the generator, en-EN was specified as default
}

type Messages interface {
    Key1() string
    Key2(arg any) string
    Key3(arg int) string
    Key4(arg1 int, arg2 float64) string
    Key5() MessagesKey5
    Key6(messages int) string
    Key7(messages int) string
}

type MessagesKey5 interface {
    Key1() string
}

type en_EN_Messages struct {}
type en_EN_Messages_Key5 struct {}
func (en_EN_Messages) Key1() string { return "message" }
func (en_EN_Messages) Key2(arg any) string { return fmt.Sprintf("message with %v", arg) }
func (en_EN_Messages) Key3(arg int) string { return fmt.Sprintf("message with %d used twice %d", arg, arg)  }
func (en_EN_Messages) Key4(arg1 int, arg2 float64) string { return fmt.Sprintf("message with %d and %.2f", arg1, arg2)  }
func (en_EN_Messages) Key5() MessagesKey5 { return en_EN_Messages_Key5{}   }
func (en_EN_Messages_Key5) Key1() string { return "nested message" }
func (en_EN_Messages) Key6(messages int) string {
    if messages == 0 {
        return "You don't have any messages."
    } else if messages == 1 {
        return "You have one new message."
    } else {
        return fmt.Sprintf("You have %d new messages.", messages)
    }
}
func (en_EN_Messages) Key7(messages int) string {
        if messages == 0 {
        return "You don't have any messages."
    } else if messages == 1 {
        return "You have one new message."
    } else {
        return "You have several new messages."
    }
}
```
