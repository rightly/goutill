package goutill

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"time"
)

var File = file{}

type file struct{}

func (file) Open(filename string) (*os.File, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		// 파일이 없어서 import 못할시엔 nil
		if os.IsNotExist(err) {
			return nil, nil
		}
		ferr := f.Close()
		if ferr != nil {
			err = ferr
		}
		return nil, err
	}

	return f, nil
}

func (file) Import(filename string) (b []byte, t time.Time, err error) {
	f, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		// 파일이 없어서 import 못할시엔 nil
		if os.IsNotExist(err) {
			return nil, time.Time{}, nil
		}
		ferr := f.Close()
		if ferr != nil {
			err = ferr
		}
		return nil, time.Time{}, err
	}

	st, err := f.Stat()
	if err != nil {
		// 파일이 없어서 import 못할시엔 nil
		if os.IsNotExist(err) {
			return nil, time.Time{}, nil
		}
		ferr := f.Close()
		if ferr != nil {
			err = ferr
		}
		return nil, time.Time{}, err
	}

	b = make([]byte, st.Size())
	_, err = f.Read(b)
	modTime := File.ModTimeWithFile(f)

	if err != nil {
		ferr := f.Close()
		if ferr != nil {
			err = ferr
		}
		return nil, modTime, err
	}

	ferr := f.Close()
	if ferr != nil {
		err = ferr
	}
	return b, modTime, err
}

func (file) Export(filename string, v interface{}) (t time.Time, err error) {
	var writeBytes []byte

	f, err := os.Create(filename)

	defer func() {
		ferr := f.Close()
		if ferr != nil {
			err = ferr
		}
	}()

	if err != nil {
		return time.Time{}, err
	}

	if _, ok := v.([]byte); !ok {
		writeBytes, err = json.Marshal(v)
		if err != nil {
			return time.Time{}, err
		}
	} else {
		writeBytes = v.([]byte)
	}

	_, err = f.Write(writeBytes)
	if err != nil {
		return time.Time{}, err
	}

	return File.ModTimeWithFile(f), f.Sync()
}

func (file) ExportJsonPretty(filename string, v interface{}) (t time.Time, err error) {
	var writeBytes []byte

	// write 할 객체 byte 로 변경
	if _, ok := v.([]byte); !ok {
		writeBytes, err = json.MarshalIndent(v, "", "\t")
		if err != nil {
			return time.Time{}, err
		}
	} else {
		return time.Time{}, errors.New("is already bytes")
	}

	return File.Export(filename, writeBytes)
}

func (file) Append(filename string, v interface{}) (t time.Time, err error) {
	//now := time.Now().Format("2006-01-02")
	var writeBytes []byte

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)

	defer func() {
		ferr := f.Close()
		if ferr != nil {
			err = ferr
		}
	}()

	// 파일 오픈 실패
	if err != nil {
		return time.Time{}, err
	}

	// []byte 형일경우 Marshal
	if _, ok := v.([]byte); !ok {
		writeBytes, err = json.Marshal(v)
		if err != nil {
			return time.Time{}, err
		}
		// reflect.Value 형 일경우 값을 참조해 Marshal
	} else if _, ok := v.(reflect.Value); ok {
		writeBytes, err = json.Marshal(v.(reflect.Value).Interface())
	} else {
		writeBytes = v.([]byte)
	}

	_, err = f.Write(writeBytes)
	if err != nil {
		return time.Time{}, err
	}
	return File.ModTimeWithFile(f), f.Sync()
}

func (file) AppendJson(filename string, v interface{}) (t time.Time, err error) {
	var writeBytes []byte

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)

	defer func() {
		ferr := f.Close()
		if ferr != nil {
			if ferr.Error() == "invalid argument" {
				err = nil
				File.Export(filename, v)
			} else {
				err = ferr
			}
		}
	}()

	if err != nil {
		return time.Time{}, err
	}

	offset, _ := f.Stat()
	// 만약 append 를 시도했으나 file 크기가 0 또는 2({} 만 있는 경우)인 경우 새로 파일을 쓴다.
	if offset.Size() == 0 || offset.Size() == 2 {
		return File.Export(filename, v)
	}

	// 기존 파일 마지막 } 대신 , 를 찍고 ( windows 에선 정상 동작하지 않는듯함 )
	comma := []byte{',', ' '}
	_, err = f.WriteAt(comma, offset.Size()-2)
	if err != nil {
		return time.Time{}, err
	}

	// write 할 객체 byte 로 변경
	if _, ok := v.([]byte); !ok {
		writeBytes, err = json.Marshal(v)
		if err != nil {
			return time.Time{}, err
		}
	} else {
		writeBytes = v.([]byte)
	}

	// 맨 앞 { 를 제거하고 write
	writeString := strings.Replace(string(writeBytes), "{", "", 1)
	_, err = f.WriteString(writeString)
	if err != nil {
		return time.Time{}, err
	}

	return File.ModTimeWithFile(f), f.Sync()
}

func (file) AppendJsonPretty(filename string, v interface{}) (t time.Time, err error) {
	var writeBytes []byte

	// write 할 객체 byte 로 변경
	if _, ok := v.([]byte); !ok {
		writeBytes, err = json.MarshalIndent(v, "", "\t")
		if err != nil {
			return time.Time{}, err
		}
	} else {
		return time.Time{}, errors.New("is already bytes")
	}

	return File.AppendJson(filename, writeBytes)
}

func (file) ModTimeWithString(filename string) time.Time {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	// 파일 오픈 실패
	if err != nil {
		return time.Time{}
	}

	st, err := f.Stat()
	if err != nil {
		return time.Time{}
	}
	defer f.Close()

	return st.ModTime()

}

func (file) ModTimeWithFile(f *os.File) time.Time {
	st, err := f.Stat()
	if err != nil {
		return time.Time{}
	}
	lc, _ := time.LoadLocation("GMT")

	return st.ModTime().In(lc)
}
