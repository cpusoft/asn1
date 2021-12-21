package der

import (
	"fmt"
)

const (
	CLASS_UNIVERSAL        = 0
	CLASS_APPLICATION      = 1
	CLASS_CONTEXT_SPECIFIC = 2
	CLASS_PRIVATE          = 3
)

func classToString(class int) string {
	switch class {
	case CLASS_UNIVERSAL:
		return "Universal"
	case CLASS_APPLICATION:
		return "Application"
	case CLASS_CONTEXT_SPECIFIC:
		return "Context-Specific"
	case CLASS_PRIVATE:
		return "Private"
	default:
		return fmt.Sprintf("Class(%d)", class)
	}
}

func classShortName(class int) string {
	switch class {
	case CLASS_UNIVERSAL:
		return "UN"
	case CLASS_APPLICATION:
		return "AP"
	case CLASS_CONTEXT_SPECIFIC:
		return "CS"
	case CLASS_PRIVATE:
		return "PR"
	default:
		return fmt.Sprintf("CLASS(%d)", class)
	}
}

// Universal types (tags)
// https://en.wikipedia.org/wiki/X.690#Types
// Permitted Construction is Primitive or Both
const (
	TAG_BOOLEAN          = 1  // 0x01
	TAG_INTEGER          = 2  // 0x02
	TAG_BIT_STRING       = 3  // 0x03
	TAG_OCTET_STRING     = 4  // 0x04
	TAG_NULL             = 5  // 0x05
	TAG_OID              = 6  // 0x06
	TAG_REAL             = 9  //0x09
	TAG_ENUMERATED       = 10 // 0x0A
	TAG_UTF8_STRING      = 12 // 0x0C
	TAG_TIME             = 14 // 0x0E
	TAG_SEQUENCE         = 16 // 0x10
	TAG_SET              = 17 // 0x11
	TAG_IA5_STRING       = 22 // 0x16
	TAG_UTC_TIME         = 23 // 0x17
	TAG_GENERALIZED_TIME = 24 // 0x18
	TAG_BMP_STRING       = 30 // 0x1E
)
