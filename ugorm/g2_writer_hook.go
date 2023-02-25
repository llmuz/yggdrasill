package ugorm

import (
	"fmt"
	"os"

	"gorm.io/gorm/logger"

	"github.com/llmuz/yggdrasill/ull"
)

type Hooks map[logger.LogLevel][]Hook
type Hook interface {
	Fire(e ull.Entry) (err error)
	Levels() (levels []logger.LogLevel)
}

func (c Hooks) Fire(level logger.LogLevel, e ull.Entry) {
	hooks := c[level]
	for _, h := range hooks {
		if err := h.Fire(e); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to fire %s", err)
			continue
		}
	}
}
