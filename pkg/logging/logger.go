package logging

import (
	"fmt"
	"log"
)

func Info(componentName string, format string, args ...interface{})  {
	msg := fmt.Sprintf(format, args...)
	logmsg := fmt.Sprintf(" [INFO] '%s' %s", componentName, msg)
	log.Println(logmsg)
}

func Error(componentName string, format string, args ...interface{})  {
	msg := fmt.Sprintf(format, args...)
	logmsg := fmt.Sprintf("[ERROR] '%s' %s", componentName, msg)
	fmt.Println(logmsg)
}

func Fatal(componentName string, format string, args ...interface{})  {
	msg := fmt.Sprintf(format, args...)
	logmsg := fmt.Sprintf("[FATAL] '%s' %s", componentName, msg)
	log.Fatalln(logmsg)
}