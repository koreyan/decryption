package main

import (
	"fmt"
	"time"

	"github.com/alexmullins/zip"
)

var token string = "abcdefghijklnmopqrstuvwxyz0123456789"

func main() {

	startTime := time.Now()
	var zipFileName string = "emergency_storage_key.zip"

	// zip 파일 열기
	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer r.Close()

	// 압축파일에는 파일이 여러개 있음
	// 그 파일들이 r.File에 슬라이스로 저장됨
	// 근데 현 상황에선 압축파일 안에 파일이 하나 밖에 없으니까
	// 그냥 아래 코드로 함
	f := r.File[0]
	ch := make(chan string)

	for i := 0; i < len(token); i++ {
		go makePassword(f, i, ch)
	}

	fmt.Println("암호찾기 실행중")

	ticker := time.Tick(time.Minute)
	var password string

	count := 1
	for {
		select {
		case <-ticker:
			fmt.Printf("%d분 지남\n", count)
			count++

		case password = <-ch:
			fmt.Println(password, "올바른 password")
			close(ch) // 채널 닫아주기
			endTime := time.Now()
			elapsedTime := endTime.Sub(startTime)
			fmt.Println("암호 해제에 걸린 시간:", elapsedTime)
			var eight int
			fmt.Scan(&eight)
			return
		}
	}

}

// password 만들어서 fileDecrypt를 호출하여 password가 맞는지 확인하는 함수
func makePassword(f *zip.File, start int, ch chan string) {

	var tryPassword []byte = make([]byte, 6)

	tryPassword[0] = token[start] // 6번째 자리

	var password string
	for p5 := 0; p5 < len(token); p5++ {
		tryPassword[1] = token[p5]
		for p4 := 0; p4 < len(token); p4++ {
			tryPassword[2] = token[p4]
			for p3 := 0; p3 < len(token); p3++ {
				tryPassword[3] = token[p3]
				for p2 := 0; p2 < len(token); p2++ {
					tryPassword[4] = token[p2]
					for p1 := 0; p1 < len(token); p1++ {
						tryPassword[5] = token[p1]
						password = fileDecrypt(f, string(tryPassword))
						if password != "" {
							ch <- password
							return
						}
					}
				}
			}
		}
	}

}

func fileDecrypt(f *zip.File, password string) string {

	// 암호를 사용하여 zip 파일을 엽니다.
	if f.IsEncrypted() { // f.IsEncrypted() => f가 암호화 되었는지 아닌지 확인함 암호화 되면 참
		f.SetPassword(password) // password 적용
	}

	// 파일 열기 => 암호가 올바르지 않으면 여기서 에러 뜸
	rc, err := f.Open()
	if err != nil {
		// fmt.Println(password, err)
		return ""
	}

	// // 새로운 파일 만들기
	// outFile, err := os.Create(f.Name)
	// if err != nil {
	// 	fmt.Println(err)
	// 	rc.Close()
	// 	return
	// }

	// // zip 파일 내의 내용을 새로운 파일로 복사한다.
	// _, err = io.Copy(outFile, rc)
	// if err != nil {
	// 	fmt.Println(err)
	// 	rc.Close()
	// 	outFile.Close()
	// 	return
	// }

	// // 파일 닫기
	// outFile.Close()
	rc.Close()
	return password
}
