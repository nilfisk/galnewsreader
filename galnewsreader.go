/*
    GalNewsReader - a pseudo news reader for Frontier's Elite:Dangerous Galnews website
    Copyright (C) 2014  Sammy Fischer

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"regexp"
	"flag"
	"os"
	)

const VERSION = "1.1"

func itemRequested() int {
	item := flag.Int("item",0,"headline number for summary, negative headline number for body, 0 for a list of headlines")
	version := flag.String("help","no","Show the version string")	
	flag.Parse()
	if *version != "no" {
		fmt.Println("galnewsreader copyright 2014 Sammy Fischer\nVersion:"+VERSION+"\n\nhttp://github.com/sammyf/galnewsreader\n\nUseage : galnewsreader -item=n")
		os.Exit(0)
	}
	return *item
}

func removeTags(htm string) string  {
	re := regexp.MustCompile("<.*?>")
	processed := re.ReplaceAllLiteralString(htm,"")
	return processed
}

func getLinkDate(htm string) string {
	re := regexp.MustCompile("galnet/([\\d]{4}-[\\d]{2}-[\\d]{2})\">")
	linkDate := re.FindStringSubmatch(htm)
	return linkDate[1]
}

func getHeadlines(htm string) []string {
	htm = strings.Replace(htm,"View full transmission &raquo;","",-1)
	re := regexp.MustCompile("<h3>(.*?)</h3>")
	headlines := re.FindAllString(htm, 10)
	return headlines
}

func getArticle(nr int, htm string) string {	
	articlesRaw := strings.Split(htm,"<article>")
	if nr > len(articlesRaw) {
		return "This article does not exist."
	}
	article := articlesRaw[nr]
	article = strings.Replace(article,"View full transmission &raquo;","",-1)
	article = strings.Replace(article,"<h3>", "Headline: ",-1)
	article = strings.Replace(article,"<p class=\"small hiLite\">", "Star Date: ",-1)
	
	article = removeTags(article)	
	return article
}

func getDetails(nr int, htm string) string {
	articlesRaw := getHeadlines(htm)
	if nr > len(articlesRaw) {
		return "This article does not exist."
	}
	if nr < 1 {
		return "This article does not exist."
	}
	article := articlesRaw[nr]
	re := regexp.MustCompile("news/galnet/\\d{4}-\\d{2}-\\d{2}")
	link := "http://www.elitedangerous.com/"+re.FindString(article)
	details := retrieveURL(link)
	detailsBit := strings.Split(details, "<div class=\"widget-content alt2 galnet\">")
	body := strings.Split(detailsBit[1], "<p><a href=\"/news/galnet/\">&laquo; GalNet Alert Service</a></p>")[0]
	body = strings.Replace(body,"&laquo; GalNet Alert Service","",-1)
	body = strings.Replace(body,"<blockquote>","\n\nQuote:\n\n",-1)
	body = strings.Replace(body,"</blockquote>","\n\nEnd Quote\n\n",-1)
	body = strings.Replace(body,"<h3>","\n\n\n",-1)
	re = regexp.MustCompile("<figure>.*?</figure>")
	body = re.ReplaceAllLiteralString(body, "")
	
	body = removeTags( body)
	return body	
}

func retrieveURL( url string) string {
	resp,err1 := http.Get(url)
	if err1 != nil {
		fmt.Println("this page could not be retrieved.")
		fmt.Println(err1)
		os.Exit(0)
	}
	bodyio := resp.Body
	buf, err2 := ioutil.ReadAll(bodyio)
	if err2 != nil {
		fmt.Println("this page could not be retrieved.")
		fmt.Println(err2)
		os.Exit(0)
	}
	return string(buf)
}

func main() {
	request := itemRequested()
	details := 0
	
	fmt.Println("GALNET News")

	if request < 0 {
		request = -request
		details = 1
	}

	htm := retrieveURL("http://www.elitedangerous.com/news/galnet/")

	if request != 0 {
		if details == 0 {
			fmt.Println("item requested : "+strconv.Itoa(request))	
			fmt.Println(getArticle(request, htm))
		} else {
			fmt.Println("full content of item : "+strconv.Itoa(request))	
			fmt.Println(getDetails(request, htm))
		}
	} else {
		fmt.Println("Current Headlines : \n\n")	
		headlines := getHeadlines(htm)
		for i:=1; i<len(headlines); i++ {
			galDate := getLinkDate(headlines[i])
			fmt.Println("Headline "+strconv.Itoa(i)+".\nStardate "+galDate+".\n"+removeTags(headlines[i])+"\n")
		}
	}
}
	
