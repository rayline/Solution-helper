package main

import "fmt"
import "net/http"
import "net/url"
import "net/http/cookiejar"
import "bytes"
import "regexp"
import "os"

//import "io/ioutil"

import "strings"

//import "github.com/PuerkitoBio/goquery"

func main() {
	var contestNum, username, password string
	fmt.Print("ID of the contest:")
	fmt.Scanln(&contestNum)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		CheckRedirect: nil,
		Jar:           jar,
	}
	//login to access
	fmt.Print("username:")
	fmt.Scanln(&username)
	fmt.Print("password:")
	fmt.Scanln(&password)
	resp, err := client.PostForm("http://biancheng.love/user/login",
		url.Values{"username": {username}, "password": {password}, "returnUrl": {"http://biancheng.love/contest/" + contestNum + "/problem"}})
	if err != nil {
		fmt.Println("Error Occured\n")
	}
	buf := new(bytes.Buffer)
	resp.Write(buf)
	doc := buf.String()
	addrExp, _ := regexp.Compile("/problem/./")
	contentExp, _ := regexp.Compile("(?s)<div class=.markdown-body.*?</div></div>")
	submissionExp, _ := regexp.Compile(`>[0-9]+<`)
	codeExp, err := regexp.Compile(`(?s)<code class="c\+\+">.*</code>`)
	link := addrExp.FindAllString(doc, 26)
	if len(link) > 0 {
		fmt.Printf("%d problems found\n", len(link))
	} else {
		fmt.Println("No problems found!\n" +
			"It's possible that:\n" +
			"Access Denied\n" +
			"Network Problem\n")
	}
	f, _ := os.Create("output.html")
	r := strings.NewReplacer("&lt;", "<", "&gt;", ">")
	//writing HTML head here
	f.WriteString(`<!DOCTYPE html>
		<html lang="zh-ch">
		<head>
		<meta charset="utf-8">
		<link rel="stylesheet" href="http://biancheng.love/stylesheets/problem/detail.css">
		</head>
		<body>`)
	for i := 0; i < len(link); i++ {
		resp, err = client.Get("http://biancheng.love/contest/" + contestNum + link[i] + "index")
		if err != nil {
			fmt.Printf("Error Getting %s\n", "http://biancheng.love/contest"+contestNum+link[i]+"index")
		}
		buf.Reset()
		resp.Write(buf)
		doc = buf.String()
		//f.WriteString(r.Replace(contentExp.FindString(doc)))
		f.WriteString(contentExp.FindString(doc))

		//fetch last Accepted code
		resp, err = client.PostForm("http://biancheng.love/contest/"+contestNum+link[i]+"submission", url.Values{"nickname": {""}, "result": {"AC"}, "language": {""}})
		buf.Reset()
		//resp.Write(f)
		resp.Write(buf)
		doc = buf.String()
		submission := submissionExp.FindString(doc)
		submission = strings.Trim(submission, "<>")
		var code string
		if submission != "1" {
			resp, err = client.Get("http://biancheng.love/submission/" + submission)
			buf.Reset()
			resp.Write(buf)
			code = buf.String()
			code = r.Replace(codeExp.FindString(code))
		}
		if submission == "1" {
			//var codePath string
			//L1:
			fmt.Println("No Accepted Submission Found Online for problem " + string(link[i][9]) + "! The code section will be left clear and you'll have to do add it on yourself")
			code = "在此处添加代码"
			/*fmt.Scanf("%s", &codePath)
			if len(codePath) > 0 {
				codeFile, err := os.Open(codePath)
				if err != nil {
					fmt.Println(err)
					goto L1
				}
				bbuf, err := ioutil.ReadAll(codeFile)
				f.Write(bbuf)
			}*/
		}
		f.WriteString(`<div class="markdown-body containing">
				<!--TODO: Write your solution here-->
				<h2>解题思路</h2>
				<p>在这里完成解题思路的编写</p>
				<h2>通过代码</h2><pre>`)
		f.WriteString(code)
		f.WriteString(`</pre></div>`)
		fmt.Fprint(f, "\n<br/><br/>\n")
		//fmt.Println("http://biancheng.love/contest/" + contestNum + link[i] + "submission/getSubmissionApi")
	}
	//print HTML tail
	f.WriteString(`</body>
		</html>`)
}
