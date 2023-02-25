package ull

type LevelHooks map[Level][]Hook
type Level = int8

const (
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

type Hook interface {
	Levels() (lvs []Level)
	Fire(e Entry) (err error)
}

// Add hook plugin
func (hooks LevelHooks) Add(hook Hook) {
	for _, level := range hook.Levels() {
		hooks[level] = append(hooks[level], hook)
	}
}

// Fire execute hook plugin, if get error return err else return nil
func (hooks LevelHooks) Fire(level Level, entry Entry) (err error) {
	for _, hook := range hooks[level] {
		if err = hook.Fire(entry); err != nil {
			return err
		}
	}
	return nil
}
