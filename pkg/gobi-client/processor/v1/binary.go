package processor_v1

import "fmt"

func (p *Processor) ProcessClientBinaryMessage(message []byte) error {
	return fmt.Errorf("binary messages are not supported for version 1")
}
