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
    case "es-ES":
        return es_ES_Messages{}, true
    }
    return nil, false
}

func MessagesForMust(tag string) Messages {
    switch strings.ReplaceAll(tag, "_", "-") {
    case "en-EN":
        return en_EN_Messages{}
    case "es-ES":
        return es_ES_Messages{}
    }
    panic(fmt.Errorf("unknwon language tag: " + tag))
}

func MessagesForOrDefault(tag string) Messages {
    switch strings.ReplaceAll(tag, "_", "-") {
    case "en-EN":
        return en_EN_Messages{}
    case "es-ES":
        return es_ES_Messages{}
    }
    return en_EN_Messages{}
}

type Messages interface{
    Cmds() Cmds
    First() string
    SeccondMessage() string
    MessageWithArgs(str string, num int, b bool, u any, f float64) string
    In() In
    ConditionalWithElse(messages int) string
    ConditionalWithoutElse(user string, messages int) string
}
type Cmds interface{
    SeccondLevel() string
    Multiline(arg int, arg2 string) string
}
type In interface{
    Depeer() string
    EvenDeeper() InEvenDeeper
}
type InEvenDeeper interface{
    Msg() string
}

type en_EN_Messages struct{}
func (en_EN_Messages) Cmds() Cmds {
    return en_EN_Cmds{}
}
type en_EN_Cmds struct{}
func (en_EN_Cmds) SeccondLevel() string {
    return "this message is Cmds.SeccondLevel"
}
func (en_EN_Cmds) Multiline(arg int, arg2 string) string {
    return "multiline" + "\n" +
        "string" + "\n" +
        fmt.Sprintf("even with %d", arg) + "\n" +
        "and much more!"
}
func (en_EN_Messages) First() string {
    return "first"
}
func (en_EN_Messages) SeccondMessage() string {
    return "seccond message"
}
func (en_EN_Messages) MessageWithArgs(str string, num int, b bool, u any, f float64) string {
    return fmt.Sprintf("this message embeds a string '%s', a number %d, a boolean %t, an unknwon type %v and a formatted float %.2g", str, num, b, u, f)
}
func (en_EN_Messages) In() In {
    return en_EN_In{}
}
type en_EN_In struct{}
func (en_EN_In) Depeer() string {
    return "this message is deeper but not because of dirs"
}
func (en_EN_In) EvenDeeper() InEvenDeeper {
    return en_EN_InEvenDeeper{}
}
type en_EN_InEvenDeeper struct{}
func (en_EN_InEvenDeeper) Msg() string {
    return "r/im14andthisisdeep"
}
func (en_EN_Messages) ConditionalWithElse(messages int) string {
    if messages == 0 {
        return "No new messages"
    } else if messages == 1 {
        return "One new message"
    } else {
        return fmt.Sprintf("You have %d new messages", messages)
    }
}
func (en_EN_Messages) ConditionalWithoutElse(user string, messages int) string {
    if messages == 0 {
        return "No new messages"
    } else if messages == 1 {
        return "One new message"
    } else if messages > 1000 {
        return fmt.Sprintf("%s, you seem to be popular!", user) + "\n" +
            fmt.Sprintf("You have %d new messages :o", messages)
    } else if messages > 1 {
        return fmt.Sprintf("You have %d new messages", messages)
    } else {
        panic(fmt.Errorf("no condition was true in conditional"))
    }
}


type es_ES_Messages struct{}
func (es_ES_Messages) Cmds() Cmds {
    return es_ES_Cmds{}
}
type es_ES_Cmds struct{}
func (es_ES_Cmds) SeccondLevel() string {
    return "este mensaje es Cmds.SeccondLevel"
}
func (es_ES_Cmds) Multiline(arg int, arg2 string) string {
    return fmt.Sprintf("multiline %d", arg) + "\n" +
        "string" + "\n" +
        fmt.Sprintf("even with %s", arg2) + "\n" +
        "and much more!"
}
func (es_ES_Messages) First() string {
    return "primero"
}
func (es_ES_Messages) SeccondMessage() string {
    return "segundo mensaje"
}
func (es_ES_Messages) MessageWithArgs(str string, num int, b bool, u any, f float64) string {
    return fmt.Sprintf("este mensaje tiene  un número %v, un booleano %v y una cadena de texto '%v' pero en otro orden, hasta se repite el número %v", num, b, str, num)
}
func (es_ES_Messages) In() In {
    return es_ES_In{}
}
type es_ES_In struct{}
func (es_ES_In) Depeer() string {
    return "este mensaje está más a dentro pero no por las carpetas"
}
func (es_ES_In) EvenDeeper() InEvenDeeper {
    return es_ES_InEvenDeeper{}
}
type es_ES_InEvenDeeper struct{}
func (es_ES_InEvenDeeper) Msg() string {
    return "r/im14andthisisdeep"
}
func (es_ES_Messages) ConditionalWithElse(messages int) string {
    if messages == 0 {
        return "No new messages"
    } else if messages == 1 {
        return "One new message"
    } else {
        return fmt.Sprintf("You have %d new messages", messages)
    }
}
func (es_ES_Messages) ConditionalWithoutElse(user string, messages int) string {
    if messages == 0 {
        return "No new messages"
    } else if messages == 1 {
        return "One new message"
    } else if messages > 1000 {
        return fmt.Sprintf("%s, you seem to be popular!", user) + "\n" +
            fmt.Sprintf("You have %d new messages :o", messages)
    } else if messages > 1 {
        return fmt.Sprintf("You have %d new messages", messages)
    } else {
        panic(fmt.Errorf("no condition was true in conditional"))
    }
}


