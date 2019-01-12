package main

import (
	"os"
	"os/exec"
	"bytes"
	"flag"
	"log"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"mime/multipart"
	"regexp"
	"net/textproto"
	"io"
)

func pull(bucket string, desc string) {
	var result bytes.Buffer
	cmd := exec.Command("./qrsctl", "listprefix", bucket, "")
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	totalDownloads := 0
	// extrace `marker:` signature
	images := strings.Split(result.String()[10:], "\n")
	for _, image := range images {
		if strings.Contains(image, "/") {
			dir := strings.Split(image, "/")[0]
			path := desc + "/" + dir
			if _, err := os.Stat(path); os.IsNotExist(err) {
				os.Mkdir(path, 0766)
			}
		}
		var dResult bytes.Buffer
		dCmd := exec.Command("./qrsctl", "get", bucket, image, desc + "/" + image)
		dCmd.Stdout = &dResult
		err := dCmd.Run()
		if err != nil {
			log.Println(err, ", downloading " + image + " failed, download continue")
		} else {
			totalDownloads++
			fmt.Println("image: ", image, " download success")
		}
	}
	fmt.Println("total downloads: ", totalDownloads)
}

func upload(name string, path string) {
	var form bytes.Buffer
	writer := multipart.NewWriter(&form)
	filePart := make(textproto.MIMEHeader)
	filePart.Set("Content-Disposition", fmt.Sprintf(`form-data; name="photo"; filename="%s"`, name))
	filePart.Set("Content-Type", "image/png")
	formImage, _ := writer.CreatePart(filePart)
	image, _ := os.Open(path)
	defer image.Close()
	io.Copy(formImage, image)
	writer.WriteField("name", name)
	writer.WriteField("albumId", "{{YOUR_ALBUM_ID}}")
	writer.Close()

	req, _ := http.NewRequest("POST", "http://x.yupoo.com/api/photos", &form)
	auth := http.Cookie{Name: "impress", Value: "{{ YOUR_YUPOO_TOKEN }}"}
	req.AddCookie(&auth)
	req.Header.Set("Authorization", "{{ YOUR_OAUTH_TOKEN }}")
	req.Header.Set("content-type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body), " , uploading ", path, " failed, continue uploading")
		return
	}

	fmt.Println(name, " upload success")
}

func push(desc string) {
	files, err := ioutil.ReadDir(desc)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)
	for _, file := range files {
		isImg, _ := regexp.MatchString("\\.jpe?g|png|gif|webp|svg", file.Name())
		if !isImg {
			continue
		}
		if file.IsDir() {
			push(desc + "/" + file.Name())
			continue
		}
		go upload(file.Name(), desc + "/" + file.Name())
	}
	<-done
}

func main() {
	method := flag.String("m", "login", "action you want to take(login/pull/push)")
	qiniu := flag.String("q", "images", "target bucket to pull")
	desc := flag.String("d", "dist", "location where qiniu images going to save")
	flag.Parse()
	switch *method {
	case "pull":
		if _, err := os.Stat("./" + *desc); os.IsNotExist(err) {
			os.Mkdir("./" + *desc, 0766)
		}
		pull(*qiniu, *desc)
	case "push":
		push(*desc)
	}
}
