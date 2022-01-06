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
	TAG_END_OF_CONTENT   = 0  // 0x00
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
	TAG_NUMBERIC_STRING  = 18 // 0x12
	TAG_PRINTABLE_STRING = 19 // 0x13
	TAG_T61_STRING       = 20 // 0x14
	TAG_VIDEOTEX_STRING  = 21 // 0x15
	TAG_IA5_STRING       = 22 // 0x16
	TAG_UTC_TIME         = 23 // 0x17
	TAG_GENERALIZED_TIME = 24 // 0x18
	TAG_BMP_STRING       = 30 // 0x1E
)

func tagName(tag int) string {
	switch tag {
	case TAG_END_OF_CONTENT:
		return "tagEndOfContent"
	case TAG_BOOLEAN:
		return "tagBoolean"
	case TAG_INTEGER:
		return "tagInteger"
	case TAG_BIT_STRING:
		return "tagBitString"
	case TAG_OCTET_STRING:
		return "tagOctetString"
	case TAG_NULL:
		return "tagNull"
	case TAG_OID:
		return "tagOid"
	case TAG_REAL:
		return "tagReal"
	case TAG_ENUMERATED:
		return "tagEnumerated"
	case TAG_UTF8_STRING:
		return "tagUtf8String"
	case TAG_TIME:
		return "tagTime"
	case TAG_SEQUENCE:
		return "tagSequence"
	case TAG_SET:
		return "tagSet"
	case TAG_NUMBERIC_STRING:
		return "tagNumbericString"
	case TAG_PRINTABLE_STRING:
		return "tarPrintableString"
	case TAG_T61_STRING:
		return "tagT61String"
	case TAG_VIDEOTEX_STRING:
		return "tagVideotexString"
	case TAG_IA5_STRING:
		return "tagIa5String"
	case TAG_UTC_TIME:
		return "tagUtcTime"
	case TAG_GENERALIZED_TIME:
		return "tagGeneralizedTime"
	case TAG_BMP_STRING:
		return "tagBmpString"
	}
	return "tagUnknown"
}
