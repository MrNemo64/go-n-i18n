# Messages

## Message types

### Literal messages

These messages are just a literal string

```json
{
  "key": "message"
}
```

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
    "ammount > 0": "The amount is positive ({amount:int})",
    "amount < 0": "The amount is negative ({amount})",
    "": "The amount is 0"
  }
}
```

## Message nesting / Grouping messages
