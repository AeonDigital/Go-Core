package xclifn

import "github.com/AeonDigital/Go-Core/xcli/pkg/xcli/xcliintfc"

// ExportIsLess provides a test bridge to evaluate the internal private isLess helper behavior.
func ExportIsLess[T xcliintfc.Quantifiable](a, b T) bool {
	return isLess(a, b)
}

// ExportIsGreater provides a test bridge to evaluate the internal private isGreater helper behavior.
func ExportIsGreater[T xcliintfc.Quantifiable](a, b T) bool {
	return isGreater(a, b)
}

// ExportFormatQuantifiable provides a test bridge to evaluate the internal private formatQuantifiable helper behavior.
func ExportFormatQuantifiable[T xcliintfc.Quantifiable](val, limit T, timeLayout string) (string, string) {
	return formatQuantifiable(val, limit, timeLayout)
}
