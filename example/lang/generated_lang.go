package lang

import "fmt"

func MessagesFor(tag string) (Messages, error) {
    switch tag {
    case "en-EN":
        return en_EN_Messages{}, nil
    case "es-ES":
        return es_ES_Messages{}, nil
    }
    return nil, fmt.Errorf("unknown language tag: %s", tag)
}

func MessagesForMust(tag string) Messages {
    switch tag {
    case "en-EN":
        return en_EN_Messages{}
    case "es-ES":
        return es_ES_Messages{}
    }
    panic(fmt.Errorf("unknown language tag: %s", tag))
}

func MessagesForOrDefault(tag string) Messages {
    switch tag {
    case "en-EN":
        return en_EN_Messages{}
    case "es-ES":
        return es_ES_Messages{}
    }
    return en_EN_Messages{}
}

type Messages interface {
    Cmds() Cmds
    First() string
    In() In
    MessageWithArgs(str string, num int, b any) string
    SeccondMessage() string
}

type Cmds interface {
    SeccondLevel() string
}

type In interface {
    Depeer() string
    EvenDeeper() InEvenDeeper
}

type InEvenDeeper interface {
    Msg() string
}

type en_EN_Messages struct{}
func (en_EN_Messages) Cmds() Cmds { return en_EN_Cmds{} }
func (en_EN_Messages) First() string { return "first" }
func (en_EN_Messages) In() In { return en_EN_In{} }
func (en_EN_Messages) MessageWithArgs(str string, num int, b any) string {
    return fmt.Sprintf("this message embedsa a string '%s', a number %d and a boolean %v", str, num, b)
}
func (en_EN_Messages) SeccondMessage() string { return "seccond message" }
type en_EN_Cmds struct{}
func (en_EN_Cmds) SeccondLevel() string { return "this message is Cmds.SeccondLevel" }
type en_EN_In struct{}
func (en_EN_In) Depeer() string { return "this message is deeper but not because of dirs" }
func (en_EN_In) EvenDeeper() InEvenDeeper { return en_EN_InEvenDeeper{} }
type en_EN_InEvenDeeper struct{}
func (en_EN_InEvenDeeper) Msg() string { return "r/im14andthisisdeep" }

type es_ES_Messages struct{}
func (es_ES_Messages) Cmds() Cmds { return es_ES_Cmds{} }
func (es_ES_Messages) First() string { return "primero" }
func (es_ES_Messages) In() In { return es_ES_In{} }
func (es_ES_Messages) MessageWithArgs(str string, num int, b any) string {
    return fmt.Sprintf("este mensaje tiene  un número %d, un booleano %v y una cadena de texto '%s' pero en otro orden, hasta se repite el número %d", num, b, str, num)
}
func (es_ES_Messages) SeccondMessage() string { return "segundo mensaje" }
type es_ES_Cmds struct{}
func (es_ES_Cmds) SeccondLevel() string { return "este mensaje es Cmds.SeccondLevel" }
type es_ES_In struct{}
func (es_ES_In) Depeer() string { return "este mensaje está más a dentro pero no por las carpetas" }
func (es_ES_In) EvenDeeper() InEvenDeeper { return es_ES_InEvenDeeper{} }
type es_ES_InEvenDeeper struct{}
func (es_ES_InEvenDeeper) Msg() string { return "r/im14andthisisdeep" }

