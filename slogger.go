package slogger

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type Logger struct {
	Prefix string
}

// NewLogger Инициализация логгера
func NewLogger(prefix string) Logger {
	return Logger{
		Prefix: prefix + ":",
	}
}

// timeForLog Возвращает текущее время в формате "2006-01-02T15:04:05.000000Z07:00"
func timeForLog() (t string) {
	t = time.Now().Format("2006-01-02T15:04:05.000000Z07:00")
	return
}

// log Логирует данные из args, используя: level - уровень логирование, prefix - префикс, установленный в логере
func log(level string, prefix string, args ...interface{}) {
	//fmt.Printf("%s %s >> %s %+v\n", level, timeForLog(), prefix, args)
	str := fmt.Sprintf("[%s] [%s] >> %s", level, timeForLog(), prefix)
	for _, arg := range args {
		if arg == nil {
			str += " <nil>"
			continue
		}
		switch t := arg.(type) {
		case error:
			str += fmt.Sprintf(" %v", t)
			continue
		}
		str += " " + stringFromReflectValue(reflect.ValueOf(arg))
	}
	fmt.Println(str)
}

// Info Логирование данных на уровне [INFO]
func (l *Logger) Info(args ...interface{}) {
	level := "INFO"
	log(level, l.Prefix, args...)
}

// Err Логирование данных на уровне [ERROR]
func (l *Logger) Err(args ...interface{}) {
	level := "ERROR"
	log(level, l.Prefix, args...)
}

// Fatal Логирование данных на уровне [FATAL], c завершением работы приложения
func (l *Logger) Fatal(args ...interface{}) {
	level := "FATAL"
	log(level, l.Prefix, args...)
	os.Exit(1)
}

// Line Выводит линию в консоли
func (l *Logger) Line() {
	fmt.Println("-------------------------------------------------------------------")
}

// AddToPrefix Добавляет новый префикс к существующему
func (l *Logger) AddToPrefix(prefix string) {
	l.Prefix = l.Prefix[:len(l.Prefix)-1] + " " + prefix + ":"
}

// stringFromReflectValue Используется методом log для правильного отображения разичных типов данных и скрытия некоторых данных, таких как пароли, в структурах
func stringFromReflectValue(v reflect.Value) (str string) {
	val := reflect.Indirect(v)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		str += val.Type().String()
		temp := "["
		for i := 0; i < val.Len(); i++ {
			temp += stringFromReflectValue(val.Index(i))
			if i != val.Len()-1 {
				temp += ", "
			}
		}
		temp += "]"
		str += temp
	case reflect.Struct:
		str += val.Type().String()
		str += "{"
		num := val.NumField()
		for i := 0; i < num; i++ {
			name := val.Type().Field(i).Name
			str += name + ":"
			if strings.Contains(strings.Join([]string{"password", "pass", "authorization"}, ";"), strings.ToLower(name)) {
				str += "******"
			} else {
				str += fmt.Sprintf("%v", val.Field(i))
				//str += fmt.Sprintf("%+v", reflect.Indirect(val.Field(i)))
			}
			if i != num-1 {
				str += " "
			}
		}
		str += "}"
	default:
		str += fmt.Sprintf("%+v", val)
	}
	return
}
