package krc

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var miarry = []byte{
	'@', 'G', 'a', 'w', '^', '2', 't', 'G', 'Q', '6', '1', '-', 'Î', 'Ò', 'n', 'i',
}

var ErrWrongFmt = fmt.Errorf("not krc format")

func GetLrcFromKrc(filename string) (string, error) {
	f, _ := os.Open(filename)

	defer f.Close()

	fileInfo, _ := f.Stat()
	zipByte := make([]byte, fileInfo.Size())

	top := make([]byte, 4)

	_, err := f.Read(top)
	if err != nil {
		return "", err
	}
	// fmt.Println(string(top))

	if string(top) != "krc1" {
		return "", ErrWrongFmt
	}

	_, err = f.Read(zipByte)
	if err != nil {
		return "", err
	}

	length := fileInfo.Size()
	var b bytes.Buffer
	for i := 0; i < int(length); i++ {
		l := i % 16
		zipByte[i] = zipByte[i] ^ miarry[l]
		b.WriteByte(zipByte[i])
	}

	r, err := zlib.NewReader(&b)
	if err != nil {
		return "", err
	}

	var outBuffer bytes.Buffer
	outBuffer.ReadFrom(r)

	var ret string

	ret = string(outBuffer.Bytes())
	return ret, nil
}

// 将毫秒转换成mm:ss.xx格式
func formatTime(msTime int) string {

	timeString := fmt.Sprintf("%vms", msTime)

	timeDuation, _ := time.ParseDuration(timeString)
	durationString := timeDuation.String()

	if msTime < 1000 {
		return fmt.Sprintf("00:%v", timeDuation.Seconds())
	}

	var ret string
	if !strings.Contains(durationString, "m") {
		ret = "00:" + durationString[0:len(durationString)-1]
	} else {
		index := strings.Index(durationString, "m")
		ret = durationString[0:index] + ":" + durationString[index+1:len(durationString)-1]
	}

	return ret
}

//eg: [513,803]<0,150,0>艾<150,150,0>丽<300,201,0>雅 <501,151,0>- <652,151,0>秋
func krcTolrc(lyrics string) string {
	// 计算此行开始的时间
	lineBeginIndex := strings.Index(lyrics, "[")
	lineEndIndex := strings.Index(lyrics, "]")
	lineBeginSlice := strings.Split(lyrics[lineBeginIndex+1:lineEndIndex], ",")
	lineBeginTime, _ := strconv.Atoi(lineBeginSlice[0])

	replaceTime := formatTime(lineBeginTime)

	lyrics = strings.Replace(lyrics, lyrics[lineBeginIndex+1:lineEndIndex], replaceTime, 1)
	// fmt.Println(ret)
	for {
		end := len(lyrics)
		wordBegin := strings.Index(lyrics[0:end], "<")
		wordEnd := strings.Index(lyrics, ">")

		if wordBegin == -1 || wordEnd == -1 {
			break
		}

		wordTimeSlice := strings.Split(lyrics[wordBegin+1:wordEnd], ",")

		wordBeginTime, _ := strconv.Atoi(wordTimeSlice[0])
		replaceTime := formatTime(wordBeginTime + lineBeginTime)
		replace := fmt.Sprintf("$%v@", replaceTime)

		lyrics = strings.Replace(lyrics, lyrics[wordBegin:wordEnd+1], replace, 1)

	}

	lyrics = strings.Replace(lyrics, "$", "<", -1)
	lyrics = strings.Replace(lyrics, "@", ">", -1)

	return lyrics
}
