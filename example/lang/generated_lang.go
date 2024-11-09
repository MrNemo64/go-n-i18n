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
    WhereAmI() string
    NestedMessages() nestedMessages
    MultiLineMessage(user string, amount float64) string
    ConditionalMessages(amount int) string
    ConditionalMessagesWithConditionArg(amount int, notUsed any) string
}
type nestedMessages interface{
    Simple() string
    Parametrized(amount int) string
    ParametrizedWithArgs(notUsed int) string
}

type en_EN_Messages struct{}
func (en_EN_Messages) WhereAmI() string {
    return "Assume this json is in the file \"en-EN.json\""
}
func (en_EN_Messages) NestedMessages() nestedMessages {
    return en_EN_nestedMessages{}
}
type en_EN_nestedMessages struct{}
func (en_EN_nestedMessages) Simple() string {
    return "This is just a simple message nested into \"nested-messages\""
}
func (en_EN_nestedMessages) Parametrized(amount int) string {
    return fmt.Sprintf("This message has an amount parameter of type int: %d", amount)
}
func (en_EN_nestedMessages) ParametrizedWithArgs(notUsed int) string {
    return fmt.Sprintf("This message is parametrized by `%d` even if the variable is not used", notUsed)
}
func (en_EN_Messages) MultiLineMessage(user string, amount float64) string {
    return fmt.Sprintf("Hello %s!", user) + "\n" +
        "Messages can be multi-line" + "\n" +
        "And each one can have parameters" + "\n" +
        fmt.Sprintf("This one has a float formatted with 2 decimals! %.2f", amount)
}
func (en_EN_Messages) ConditionalMessages(amount int) string {
    if amount == 0 {
        return "If amount is 0, this message is used"
    } else if amount == 1 {
        return "This message is returned if the amount is 1"
    } else {
        return "This is the \"else\" branch" + "\n" +
            "This multi-line message is used" + "\n" +
            fmt.Sprintf("And shows the amount: %d", amount)
    }
}
func (en_EN_Messages) ConditionalMessagesWithConditionArg(amount int, notUsed any) string {
    if amount == 0 {
        return "If amount is 0, this message is used"
    } else if amount == 1 {
        return "This message is returned if the amount is 1"
    } else {
        return "This is the \"else\" branch" + "\n" +
            "This multi-line message is used" + "\n" +
            "But the ammount is not displayed"
    }
}


