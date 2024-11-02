# Messages

## Message types

### Literal messages

These messages are just a literal string

```json
{
  "key": "message"
}
```

<details>
  <summary>Generated code</summary>

```go
/** Code generated using https://github.com/MrNemo64/go-n-i18n
* Any changes to this file will be lost on the next tool run */

package lang

import (
	"fmt"
	"strings"
)

func MessagesFor(tag string) (Messages, bool) {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}, true
	}
	return nil, false
}

func MessagesForMust(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	panic(fmt.Errorf("unknwon language tag: " + tag))
}

func MessagesForOrDefault(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	return en_EN_Messages{}
}

type Messages interface {
	Key() string
}

type en_EN_Messages struct{}

func (en_EN_Messages) Key() string {
	return "message"
}
```

</details>

### Parametrized messages

These messages hold one or more parameters. Parameters are specified by following the format `{name:type:format}` where the type and format are optional.

- `{name}`: Parameter of unknown type named `name`
- `{name:str}`: Parameter of type string named `name`
- `{amount:float64:.2f}`: Parameter of type float with 64 bits with a format rounded to 2 decimals

The same parameter can be used several times on the same language, using diferent formats but always the same type. The type only needs to be specified ones in one language and all languages will use the same type. It is recomended to specify in the default language all the types and just reference the parameters by name in the rest of languages.

```json
{
  "key": "message with a parameter of type float with 64 bits and rounded to 2 decimals {value:float64:.2f}"
}
```

<details>
  <summary>Generated code</summary>

```go
/** Code generated using https://github.com/MrNemo64/go-n-i18n
 * Any changes to this file will be lost on the next tool run */

package lang

import (
    "fmt"
    "strings"
)

func MessagesFor(tag string) (Messages, bool) {
    switch strings.ReplaceAll(tag, "_", "-") {
    case "en-EN":
        return en_EN_Messages{}, true
    }
    return nil, false
}

func MessagesForMust(tag string) Messages {
    switch strings.ReplaceAll(tag, "_", "-") {
    case "en-EN":
        return en_EN_Messages{}
    }
    panic(fmt.Errorf("unknwon language tag: " + tag))
}

func MessagesForOrDefault(tag string) Messages {
    switch strings.ReplaceAll(tag, "_", "-") {
    case "en-EN":
        return en_EN_Messages{}
    }
    return en_EN_Messages{}
}

type Messages interface{
    Key(value float64) string
}

type en_EN_Messages struct{}
func (en_EN_Messages) Key(value float64) string {
    return fmt.Sprintf("message with a parameter of type float with 64 bits and rounded to 2 decimals %.2f", value)
}
```

</details>

#### Allowed arguments

| Name    | Type    | Aliases        | Default format |
| ------- | ------- | -------------- | -------------- |
| any     | any     | unknown        | v              |
| string  | string  | str            | s              |
| boolean | bool    | boolean        | t              |
| integer | int     | int            | d              |
| float   | float64 | f64, f, double | g              |

More arguments will be aded with time

### Multiline messages

These messages span multiple lines. Each line may be a [literal message](#literal-messages) or a [parametrized message](#parametrized-messages). All lines share the same parameters

```json
{
  "key": [
    "first line of the message",
    "the seccond line has a parameter {arg:int} of type int",
    "and the thir line reuses that parameter {arg}"
  ]
}
```

<details>
  <summary>Generated code</summary>

```go
/** Code generated using https://github.com/MrNemo64/go-n-i18n
 * Any changes to this file will be lost on the next tool run */

package lang

import (
	"fmt"
	"strings"
)

func MessagesFor(tag string) (Messages, bool) {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}, true
	}
	return nil, false
}

func MessagesForMust(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	panic(fmt.Errorf("unknwon language tag: " + tag))
}

func MessagesForOrDefault(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	return en_EN_Messages{}
}

type Messages interface {
	Key(arg int) string
}

type en_EN_Messages struct{}

func (en_EN_Messages) Key(arg int) string {
	return "first line of the message" + "\n" +
		fmt.Sprintf("the seccond line has a parameter %d of type int", arg) + "\n" +
		fmt.Sprintf("and the thir line reuses that parameter %d", arg)
}
```

</details>

### Conditional messages

These messages allow to change the message itself based on a condition and have its key prefixed by a `?`. Useful, for example, for quantitnes. Each condition value may be a [literal message](#literal-messages), a [parametrized message](#parametrized-messages) or a [multiline message](#multiline-messages). All condition values share the same parameters.

Conditions and their respective associated message are specified as key-value pairs in an object. An empty key can be specified to indicate the "else" message, the message to be used if none of the conditions evaluate to true. If no else message is specified, an else branch is added with a call to panic. Conditions are writen in the code as they're found in the json, in the same order and copying each one into the if statement.

```json
{
  "?key": {
    "messages > 100": "You have a lot of new messages ({messages:int})!",
    "messages > 10": "You have {messages} new messages",
    "messages == 1": "You have one new message",
    "messages == 0": "No new messages"
  },
  "?key-with-else-branch": {
    "amount > 0": "The amount is positive ({amount:int})",
    "amount < 0": "The amount is negative ({amount})",
    "": "The amount is 0"
  }
}
```

<details>
  <summary>Generated code</summary>

```go
/** Code generated using https://github.com/MrNemo64/go-n-i18n
 * Any changes to this file will be lost on the next tool run */

package lang

import (
	"fmt"
	"strings"
)

func MessagesFor(tag string) (Messages, bool) {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}, true
	}
	return nil, false
}

func MessagesForMust(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	panic(fmt.Errorf("unknwon language tag: " + tag))
}

func MessagesForOrDefault(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	return en_EN_Messages{}
}

type Messages interface {
	Key(messages int) string
	KeyWithElseBranch(amount int) string
}

type en_EN_Messages struct{}

func (en_EN_Messages) Key(messages int) string {
	if messages > 100 {
		return fmt.Sprintf("You have a lot of new messages (%d)!", messages)
	} else if messages > 10 {
		return fmt.Sprintf("You have %d new messages", messages)
	} else if messages == 1 {
		return "You have one new message"
	} else if messages == 0 {
		return "No new messages"
	} else {
		panic(fmt.Errorf("no condition was true in conditional"))
	}
}
func (en_EN_Messages) KeyWithElseBranch(amount int) string {
	if amount > 0 {
		return fmt.Sprintf("The amount is positive (%d)", amount)
	} else if amount < 0 {
		return fmt.Sprintf("The amount is negative (%d)", amount)
	} else {
		return "The amount is 0"
	}
}
```

</details>

## Message nesting / Grouping messages

Messages can be grouped or nested by nesting json objets.
By nesting messages, a separation is done and each group of messages is placed into their own interface and structs.
This way autocompletion of messages is not polluted with hunderds of messages and its easyer to navigate them.
It also means that each part of a program can receive only the interface with the messages it needs.

```json
{
  "key-level-1": {
    "key-level-2": {
      "key-level-3": "Assume this message is in the file `en-EN.json`"
    }
  }
}
```

To get the message we need to call `messages.KeyLevel1().KeyLevel2().KeyLevel3()`.

Another way of nesting messages is using folders to nest files.

```json
{
  "key-level-3": "Assume this message is in the file `key-level-1/key-level-2/en-EN.json`"
}
```

To get this message we also need to call `messages.KeyLevel1().KeyLevel2().KeyLevel3()`.

Nested levels can be defined by using nested json objects, nesting files in folders or both.

<details>
  <summary>Generated code</summary>

```go
/** Code generated using https://github.com/MrNemo64/go-n-i18n
 * Any changes to this file will be lost on the next tool run */

package lang

import (
	"fmt"
	"strings"
)

func MessagesFor(tag string) (Messages, bool) {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}, true
	}
	return nil, false
}

func MessagesForMust(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	panic(fmt.Errorf("unknwon language tag: " + tag))
}

func MessagesForOrDefault(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	return en_EN_Messages{}
}

type Messages interface {
	KeyLevel1() keyLevel1
}
type keyLevel1 interface {
	KeyLevel2() keyLevel1keyLevel2
}
type keyLevel1keyLevel2 interface {
	KeyLevel3() string
}

type en_EN_Messages struct{}

func (en_EN_Messages) KeyLevel1() keyLevel1 {
	return en_EN_keyLevel1{}
}

type en_EN_keyLevel1 struct{}

func (en_EN_keyLevel1) KeyLevel2() keyLevel1keyLevel2 {
	return en_EN_keyLevel1keyLevel2{}
}

type en_EN_keyLevel1keyLevel2 struct{}

func (en_EN_keyLevel1keyLevel2) KeyLevel3() string {
	return "Assume this message is in the file `en-EN.json`"
}
```

</details>

### Interface renaming

By default the name used to create the interface of nested groups of messages is the full path of the group of messages.
In the example above, 3 interfaces would have been generated: `Messages`, `keyLevel1` and `keyLevel1keyLevel2` (if `public-non-named-interfaces` is specified when generating the code, the names would been `Messages`, `KeyLevel1` and `KeyLevel1KeyLevel2` to make all of them public).

When nesting too much, these interface names can get long.
Since we may want to use some of the generated interfaces in our code, we can provide a name for them in the json by putting `:name` after the key of the group of messages.

```json
{
  "key-level-1:l1": {
    "key-level-2:l2": {
      "key-level-3": "Assume this message is in the file `en-EN.json`"
    }
  }
}
```

In this case since we renamed the keys to `l1` and `l2` the generated interfaces will be named `Messages`, `L1` and `L2`.
If we want to rename a group specified by folders, since `:` is not a valid character for folder names, we can rename the key in the parent group of messages with an empty json object.

```json
{
  "key-level-3": "Assume this message is in the file `key-level-1/key-level-2/en-EN.json`"
}
```

```json
{
  "key-level-1:l1": {
    "key-level-2:l2": {}
  }
}
```

Here we renamed both groups of messages even though these groups are defined by folders and not by nesting json objects.

<details>
  <summary>Generated code</summary>

```go
/** Code generated using https://github.com/MrNemo64/go-n-i18n
 * Any changes to this file will be lost on the next tool run */

package lang

import (
	"fmt"
	"strings"
)

func MessagesFor(tag string) (Messages, bool) {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}, true
	}
	return nil, false
}

func MessagesForMust(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	panic(fmt.Errorf("unknwon language tag: " + tag))
}

func MessagesForOrDefault(tag string) Messages {
	switch strings.ReplaceAll(tag, "_", "-") {
	case "en-EN":
		return en_EN_Messages{}
	}
	return en_EN_Messages{}
}

type Messages interface {
	KeyLevel1() L1
}
type L1 interface {
	KeyLevel2() L2
}
type L2 interface {
	KeyLevel3() string
}

type en_EN_Messages struct{}

func (en_EN_Messages) KeyLevel1() L1 {
	return en_EN_L1{}
}

type en_EN_L1 struct{}

func (en_EN_L1) KeyLevel2() L2 {
	return en_EN_L2{}
}

type en_EN_L2 struct{}

func (en_EN_L2) KeyLevel3() string {
	return "Assume this message is in the file `en-EN.json`"
}
```

</details>
