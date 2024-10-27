package logger

func Fatal(v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Fatal(v...)
	if fileLogger != nil {
		fileLogger.Fatal(v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Fatalf(format, v...)
	if fileLogger != nil {
		fileLogger.Fatalf(format, v...)
	}
}

func Error(v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Error(v...)
	if fileLogger != nil {
		fileLogger.Error(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Errorf(format, v...)
	if fileLogger != nil {
		fileLogger.Errorf(format, v...)
	}
}

func Warn(v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Warn(v...)
	if fileLogger != nil {
		fileLogger.Warn(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Warnf(format, v...)
	if fileLogger != nil {
		fileLogger.Warnf(format, v...)
	}
}

func Info(v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Info(v...)
	if fileLogger != nil {
		fileLogger.Info(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Infof(format, v...)
	if fileLogger != nil {
		fileLogger.Infof(format, v...)
	}
}

func Debug(v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Debug(v...)
	if fileLogger != nil {
		fileLogger.Debug(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Debugf(format, v...)
	if fileLogger != nil {
		fileLogger.Debugf(format, v...)
	}
}

func Trace(v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Trace(v...)
	if fileLogger != nil {
		fileLogger.Trace(v...)
	}
}

func Tracef(format string, v ...interface{}) {
	if stdLogger == nil {
		Init(nil, "")
	}
	stdLogger.Tracef(format, v...)
	if fileLogger != nil {
		fileLogger.Tracef(format, v...)
	}
}
