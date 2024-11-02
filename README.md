# go-n-i18n

A code generation tool writen in go and inspired by [ParaglideJS](https://inlang.com/m/gerre34r/library-inlang-paraglideJs)
for type safe and feature rich internationalization.

## How it works

Messages are defined in a JSON file,
go-n-i18n will extract the messages from these files and generate code with it.
Here is an example:

```JSON
{
  "where-am-i": "Assume this json is in the file \"en-EN.json\"",
  "nested-messages": {
    "simple": "This is just a simple message nested into \"nested-messages\"",
    "parametrized": "This message has an amount parameter of type int: {amount:int}"
  },
  "multiline-message": [
    "Hello {user:str}!",
    "Messages can be multiline",
    "And each one can have parameters",
    "This one has a float formated with 2 decimals! {amount:float64:.2f}"
  ],
  "?conditional-messages": {
    "amount == 0": "If amount is 0, this message is used",
    "amount == 1": "This message is returned if the amount is 1",
    "": [
      "This is the \"else\" branch",
      "This multiline message is used",
      "And shows the amount: {amount:int}"
    ]
  }
}
```

When running go-n-i18n, you'll get code that looks like this:

```go
// Utility methods
func MessagesFor(tag string) (Messages, bool) { ... }

func MessagesForMust(tag string) Messages { ... }

func MessagesForOrDefault(tag string) Messages { ... }

type Messages interface{
    WhereAmI() string
    NestedMessages() nestedMessages
    MultilineMessage(user string, amount float64) string
    ConditionalMessages(amount int) string
}
type nestedMessages interface{
    Simple() string
    Parametrized(amount int) string
}

// Struct that implements Messages returning the messages defined in the language file
type en_EN_Messages struct{}
// More code... See examples/lang/generated_lang.go for all of it
```

Now you can get an instance of your messages and use them!

```go
func main() {
  bundle := lang.MessagesForMust("en-EN")

  fmt.Println(bundle.WhereAmI())
  // Assume this json is in the file "en-EN.json"

  fmt.Println(bundle.NestedMessages().Parametrized(4))
  // This message has an amount parameter of type int: 4

  fmt.Println(bundle.ConditionalMessages(100))
  /*
    This is the "else" branch
    This multiline message is used
    And shows the amount: 100
  */

    fmt.Println(bundle.MultilineMessage("MrNemo64", 13.1267))
  /*
    Hello MrNemo64!
    Messages can be multiline
    And each one can have parameters
    This one has a float formated with 2 decimals! 13.13
  */
}
```

## Installing and using

Install by cloning the repository and running `make install` or by running `go install github.com/MrNemo64/go-n-i18n/cmd/i18n@v0.0.3`.

To use it you must invoke the generator. This can be done by using a specifig file in your language folder as such:

```go
package lang

// use en-EN as default language and start looking for language files in the current directory
//go:generate i18n -default-language en-EN -messages .
```

Or by manually running the command.

## More information

See the [wiki](https://github.com/MrNemo64/go-n-i18n/wiki) or the [docs](https://github.com/MrNemo64/go-n-i18n/tree/main/docs) folder for more details on how to use the tool.
