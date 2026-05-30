package xlog

import (
	"os"
	"testing"
)

var (
	FnMockable_IsCLISupportsColors = &fnMockable_IsCLISupportsColors
	FnMockable_IsDir               = &fnMockable_IsDir
	FnMockable_CreateDirPath       = &fnMockable_CreateDirPath
	FnMockable_OpenFileWrite       = &fnMockable_OpenFileWrite
)

// MockFunction substitui temporariamente uma função por um mock de teste
func MockFunction[T any](t *testing.T, original *T, mock T) {
	t.Helper()
	oldValue := *original
	*original = mock
	t.Cleanup(func() {
		*original = oldValue
	})
}

// ExportIsCLISupportsColors expõe a lógica privada de detecção de cores para testes externos.
func ExportIsCLISupportsColors() bool {
	return isCLISupportsColors()
}

// ExportFixANSIEscape expõe a lógica privada de correção de códigos ANSI para testes externos.
func ExportFixANSIEscape(str string) string {
	return fixANSIEscape(str)
}

// ExportGenerateLogMessage expõe a lógica privada de geração de mensagens para testes externos.
func (o *LogHandler) ExportGenerateLogMessage(timeStr string, levelStr string, msg string) (string, string) {
	return o.generateLogMessage(timeStr, levelStr, msg)
}

// ClosePrivateRegistryFile fecha o descritor de arquivo privado para simular erros de escrita em testes.
func (o *LogHandler) ClosePrivateRegistryFile() error {
	if o.logRegistryFile != nil {
		return o.logRegistryFile.Close()
	}
	return nil
}

// SetLogRegistryFile permite que a suíte de testes injete um Pipe em memória
// sem expor a propriedade privada para os usuários da biblioteca em produção.
func (o *LogHandler) SetLogRegistryFile(f *os.File) {
	o.logRegistryFile = f
}
