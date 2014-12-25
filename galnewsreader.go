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
	"html"
	)

const VERSION = "2.0"

/*
	routines for Szokorad galnetarchive
	http://galnetarchive.blogspot.com/feeds/posts/default
*/
func zok_getArticles(htm string) []map[string]string {
	weekdays := map[string]string{"Mon":"Monday","Tue":"Tuesday","Wed":"Wednesday","Thu":"Thursday","Fri":"Friday","Sat":"Saturday","Sun":"Sunday"}
	months := map[string]string{"Jan":"January", "Feb":"February", "Mar":"March", "Apr":"April", "May":"May", "Jun":"June", "Jul":"July", "Aug":"August", "Sep":"September", "Oct":"October", "Nov":"November", "Dec":"December"}
	re := regexp.MustCompile("<item>(.*?)</item>")
	entries := re.FindAllStringSubmatch(htm, -1)
	var articles = make([]map[string]string,0)
	if len(entries) <= 0 {
		articles = append(articles, map[string]string{"headline":"no news to show","body":"The newsfeed seems to be currently unavailable"})
		return articles
	}
	for i:=0; i<len(entries); i++ {	
		re = regexp.MustCompile("<pubDate>(.*?) \\d{2}:\\d{2}:\\d{2} [+-]\\d{4}</pubDate>")
		date := re.FindStringSubmatch(entries[i][1])[1]
		re = regexp.MustCompile("([A-Za-z]{3}), \\d{2} ([A-Za-z]{3}) \\d{4}")
		shorts := re.FindAllStringSubmatch(date, -1)
		date = strings.Replace(date, shorts[0][1], weekdays[shorts[0][1]], 1);
		date = strings.Replace(date, shorts[0][2], months[shorts[0][2]], 1);		
		re = regexp.MustCompile("<title>(.*?)</title>")
		title := re.FindStringSubmatch(entries[i][1])[1]
		headline := "stardate: "+date+".\n\n\n\nHeadline:\n"+removeTags(title)+"\n"
		re = regexp.MustCompile("<description>(.*?)</description>")
		raw := html.UnescapeString(re.FindStringSubmatch(entries[i][1])[1])
		re = regexp.MustCompile("<blockquote(.*?)>")
		raw = re.ReplaceAllLiteralString(raw, "\nQuote:\n ")
		raw = strings.Replace(raw,"</blockquote>", "\n\nEnd Quote:\n ",-1)
		body := removeTags(raw)
		article := map[string]string{"headline":headline,"body":body}
		articles = append(articles, article)
	}
	return articles
}

/*
	routines for the official Frontier Galnet News page
*/


func getLinkDate(htm string) string {
	re := regexp.MustCompile("galnet/([\\d]{4}-[\\d]{2}-[\\d]{2})\">")
	linkDate := re.FindStringSubmatch(htm)
	return linkDate[1]
}

func getHeadlines(htm string) []string {
	htm = strings.Replace(htm,"View full transmission &raquo;","",-1)
	re := regexp.MustCompile("<h3>(.*?)</h3>")
	headlines := re.FindAllString(htm, -1)
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
	article := articlesRaw[nr-1]
	re := regexp.MustCompile("news/galnet/\\d{4}-\\d{2}-\\d{2}")
	link := "http://www.elitedangerous.com/"+re.FindString(article)
	details := retrieveURL(link)
	detailsBit := strings.Split(details, "<div class=\"widget-content alt2 galnet\">")
	body := strings.Split(detailsBit[1], "<p><a href=\"/news/galnet/\">&laquo; GalNet Alert Service</a></p>")[0]
	body = strings.Replace(body,"&laquo; GalNet Alert Service","",-1)
	body = strings.Replace(body,"<blockquote>","\n\nQuote:\n\n",-1)
	body = strings.Replace(body,"</blockquote>","\n\nEnd Quote\n\n",-1)
	body = strings.Replace(body,"<h3>","\n",-1)
	body = strings.Replace(body,"transmisSion","transmission",-1)
	body = strings.Replace(body,"::","",-1)
	re = regexp.MustCompile("<figure>.*?</figure>")
	body = re.ReplaceAllLiteralString(body, "")
	
	body = removeTags( body)
	return body	
}

/*
	Utilities
 */

func outputToFile( s string) {
	b := []byte(s)
	ioutil.WriteFile("./galnews", b, 0777)
}

func retrieveURL( url string) string {
	resp,err1 := http.Get(url)
	if err1 != nil {
		outputToFile("this page could not be retrieved.")
		fmt.Println(err1)
		os.Exit(0)
	}
	bodyio := resp.Body
	buf, err2 := ioutil.ReadAll(bodyio)
	if err2 != nil {
		outputToFile("this page could not be retrieved.")
		fmt.Println(err2)
		os.Exit(0)
	}
	return string(buf)
}

func itemRequested() (int,int) {
	item := flag.Int("item",0,"headline number for summary, negative headline number for body, 0 for a list of headlines")
	zok := flag.Int("zok",1,"retrieve the articles from Szokorad galnetarchive instead of the official site. note that there are only headlines or full articles. no summaries! A negative number will return 100 full articles.")
	version := flag.String("help","no","Show the version string")
	flag.Parse()
	if *version != "no" {
		fmt.Println("galnewsreader copyright 2014 Sammy Fischer\nVersion:"+VERSION+"\n\nhttp://github.com/sammyf/galnewsreader\n\nUseage : galnewsreader -item=n -zok=n")
		os.Exit(0)
	}
	return *item, *zok
}

func removeTags(htm string) string  {
	re := regexp.MustCompile("<.*?>")
	processed := re.ReplaceAllLiteralString(htm,"")
	return processed
}

/*
	MAIN
 */


func main() {
	request, zok := itemRequested()
	details := 0
	rs := ""
	if zok==0 {
		rs = "GALNET News.\n"

		if request < 0 {
			request = -request
			details = 1
		}

		htm := retrieveURL("http://www.elitedangerous.com/news/galnet/")

		if request != 0 {
			if details == 0 {
				rs = rs+"item requested : "+strconv.Itoa(request)+"\n"
				rs = rs+getArticle(request, htm)
			} else {
				rs = rs+"full content of item : "+strconv.Itoa(request)+"\n"
				rs = rs+getDetails(request, htm)
			}
		} else {
			rs = rs+"Current Headlines : \n\n"
			headlines := getHeadlines(htm)
			for i := 0; i < len(headlines); i++ {
				galDate := getLinkDate(headlines[i])
				rs = rs+"Headline "+strconv.Itoa(i+1)+".\nStardate "+galDate+".\n"+removeTags(headlines[i])+"\n"
			}
		}
	} else {
		rs = "Zsokorad Galnet News - The fastest news in the Galaxy\n\n"
		maxresults := "10"
		if request < 0 {
			maxresults = "100"
		}
		htm := retrieveURL("http://galnetarchive.blogspot.com/feeds/posts/default?max-results="+maxresults+"&alt=rss")
		articles := zok_getArticles(htm)
		if request == 0 {
			rs = rs+"Current Headlines : \n\n"
			for i := 0; i < len(articles); i++ {
				rs = rs+"item "+strconv.Itoa(i+1)+":\n"+articles[i]["headline"]+"\n\n\n"
			}
		} else if request > 0 {
			if request > len(articles) {
				rs = rs+"This item does not exist"
			} else {
				rs = rs+articles[request-1]["headline"]+"\n\n"+articles[request-1]["body"]+"\n\n\n"
			}
		} else {
			fmt.Println(articles)
			for i := 0; i < len(articles); i++ {
rs=rs+"\n\n\nitem"+strconv.Itoa(i+1)+":\n"+articles[i]["headline"]+"\n\n"+articles[i]["body"]+"\n\n\n"
			}
		}
	}
	rs = rs + "\nEnd of Stream.\n"
	outputToFile(rs)
	fmt.Println(rs)
}
	
