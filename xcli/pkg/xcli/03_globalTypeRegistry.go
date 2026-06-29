package xcli

import (
	"time"

	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xcliconstt"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xclifn"
	"github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xcliintfc"
)

// GlobalTypeRegistry acts as the immutable central directory mapping every supported
// FlagType constant to its respective type-safe conversion and validation engine instance.
var GlobalTypeRegistry = map[xcliconstt.FlagType]xcliintfc.ValueParser{

	//
	// --- Primitive Types ---

	xcliconstt.TypeString: TypedDescriptor[string]{
		Parser:         xclifn.ParseString,
		LimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeBool: TypedDescriptor[bool]{
		Parser: xclifn.ParseBool,
	},
	xcliconstt.TypeInt: TypedDescriptor[int]{
		Parser:         xclifn.ParseInt,
		LimitValidator: xclifn.ValidateNumberLimits[int],
	},
	xcliconstt.TypeFloat: TypedDescriptor[float64]{
		Parser:         xclifn.ParseFloat,
		LimitValidator: xclifn.ValidateNumberLimits[float64],
	},

	//
	// --- Structured Data ---

	xcliconstt.TypeJSON: TypedDescriptor[string]{
		Parser:         xclifn.ParseJSON,
		LimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeDuration: TypedDescriptor[time.Duration]{
		Parser:         xclifn.ParseDuration,
		LimitValidator: xclifn.ValidateNumberLimits[time.Duration],
	},
	xcliconstt.TypeDate: TypedDescriptor[time.Time]{
		Parser:         xclifn.ParseDate,
		LimitValidator: xclifn.ValidateNumberLimits[time.Time],
	},
	xcliconstt.TypeTime: TypedDescriptor[time.Time]{
		Parser:         xclifn.ParseTime,
		LimitValidator: xclifn.ValidateNumberLimits[time.Time],
	},
	xcliconstt.TypeDateTime: TypedDescriptor[time.Time]{
		Parser:         xclifn.ParseDateTime,
		LimitValidator: xclifn.ValidateNumberLimits[time.Time],
	},
	xcliconstt.TypeEmail: TypedDescriptor[string]{
		Parser:         xclifn.ParseEmail,
		LimitValidator: xclifn.ValidateStringLimits,
	},

	//
	// --- System and Network Validations ---

	xcliconstt.TypePath: TypedDescriptor[string]{
		Parser:         xclifn.ParsePath,
		LimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeFilename: TypedDescriptor[string]{
		Parser:         xclifn.ParseFilename,
		LimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeDirname: TypedDescriptor[string]{
		Parser:         xclifn.ParseDirname,
		LimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeURL: TypedDescriptor[string]{
		Parser:         xclifn.ParseURLStandard,
		LimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeFullURL: TypedDescriptor[string]{
		Parser:         xclifn.ParseFullURL,
		LimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeRelativeURL: TypedDescriptor[string]{
		Parser:         xclifn.ParseRelativeURL,
		LimitValidator: xclifn.ValidateStringLimits,
	},

	//
	// --- Physical File System Inputs ---

	xcliconstt.TypeFilepath: TypedDescriptor[string]{
		Parser:        xclifn.ParseFilepath,
		DiskValidator: xclifn.ValidateDiskResources,
	},
	xcliconstt.TypeDirpath: TypedDescriptor[string]{
		Parser:        xclifn.ParseDirpath,
		DiskValidator: xclifn.ValidateDiskResources,
	},

	//
	// --- Array / Slice Formats ---

	xcliconstt.TypeStringArray: ArrayDescriptor[string]{
		ElementParser:      xclifn.ParseString,
		ItemLimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeIntArray: ArrayDescriptor[int]{
		ElementParser:      xclifn.ParseInt,
		ItemLimitValidator: xclifn.ValidateNumberLimits[int],
	},
	xcliconstt.TypeFloatArray: ArrayDescriptor[float64]{
		ElementParser:      xclifn.ParseFloat,
		ItemLimitValidator: xclifn.ValidateNumberLimits[float64],
	},
	xcliconstt.TypeBoolArray: ArrayDescriptor[bool]{
		ElementParser: xclifn.ParseBool,
	},
	xcliconstt.TypeDurationArray: ArrayDescriptor[time.Duration]{
		ElementParser:      xclifn.ParseDuration,
		ItemLimitValidator: xclifn.ValidateNumberLimits[time.Duration],
	},
	xcliconstt.TypeDateArray: ArrayDescriptor[time.Time]{
		ElementParser:      xclifn.ParseDate,
		ItemLimitValidator: xclifn.ValidateNumberLimits[time.Time],
	},
	xcliconstt.TypeTimeArray: ArrayDescriptor[time.Time]{
		ElementParser:      xclifn.ParseTime,
		ItemLimitValidator: xclifn.ValidateNumberLimits[time.Time],
	},
	xcliconstt.TypeDateTimeArray: ArrayDescriptor[time.Time]{
		ElementParser:      xclifn.ParseDateTime,
		ItemLimitValidator: xclifn.ValidateNumberLimits[time.Time],
	},
	xcliconstt.TypeEmailArray: ArrayDescriptor[string]{
		ElementParser:      xclifn.ParseEmail,
		ItemLimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypePathArray: ArrayDescriptor[string]{
		ElementParser:      xclifn.ParsePath,
		ItemLimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeFilenameArray: ArrayDescriptor[string]{
		ElementParser:      xclifn.ParseFilename,
		ItemLimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeFilepathArray: ArrayDescriptor[string]{
		ElementParser:     xclifn.ParseFilepath,
		ItemDiskValidator: xclifn.ValidateDiskResources,
	},
	xcliconstt.TypeDirnameArray: ArrayDescriptor[string]{
		ElementParser:      xclifn.ParseDirname,
		ItemLimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeDirpathArray: ArrayDescriptor[string]{
		ElementParser:     xclifn.ParseDirpath,
		ItemDiskValidator: xclifn.ValidateDiskResources,
	},
	xcliconstt.TypeURLArray: ArrayDescriptor[string]{
		ElementParser:      xclifn.ParseURLStandard,
		ItemLimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeFullURLArray: ArrayDescriptor[string]{
		ElementParser:      xclifn.ParseFullURL,
		ItemLimitValidator: xclifn.ValidateStringLimits,
	},
	xcliconstt.TypeRelativeURLArray: ArrayDescriptor[string]{
		ElementParser:      xclifn.ParseRelativeURL,
		ItemLimitValidator: xclifn.ValidateStringLimits,
	},
}
