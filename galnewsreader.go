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

const VERSION = "3.01"

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
func off_getArticles(htm string) []map[string]string {
	months := map[string]string{"JAN":"0 1", "FEB":"0 2", "MAR":"0 3", "APR":"0 4", "MAY":"0 5", "JUN":"0 6", "JUL":"0 7", "AUG":"0 8", "SEP":"0 9", "OCT":"10", "NOV":"11", "DEC":"12"}
	htm = strings.Replace(htm,"\n"," ",-1)
	htm = strings.Replace(htm,"\r"," ",-1)
	re := regexp.MustCompile("div class=\"article\">(.+?)<div class=\"widget\">")
	entries := re.FindAllStringSubmatch(htm, -1)
	var articles = make([]map[string]string,0)
	if len(entries) <= 0 {
		articles = append(articles, map[string]string{"headline":"no news to show","body":"The newsfeed seems to be currently unavailable"})
		return articles
	}
	for i:=0; i<len(entries); i++ {
		re = regexp.MustCompile("<div class=\"i_right\" style=\"margin: 5px\"><p class=\"small\" style=\"color:#888;\">(\\d{2}) ([A-Z]{3}) (\\d{4})</p></div>")
		dateSubs := re.FindStringSubmatch(entries[i][1])
		shorts := dateSubs[2]
		date := dateSubs[3]+" "+months[shorts]+" "+dateSubs[1]		
		re = regexp.MustCompile("read=([a-zA-Z0-9]+?)\">(.*?)</a>")
		title := re.FindStringSubmatch(entries[i][1])[2]
		headline := "stardate: "+date+".\n\n\n\nHeadline:\n"+title+"\n\n"
		re := regexp.MustCompile("\\W")
		tmp := re.ReplaceAllLiteralString(title,"\\W")
		re = regexp.MustCompile(tmp+"</a></h3><p>(.*?)</div>")
		raw := html.UnescapeString(re.FindStringSubmatch(entries[i][1])[1])
		raw = strings.Replace(raw,"</p>", "\n\n",-1)
		body := removeTags(raw)
		article := map[string]string{"headline":headline,"body":body}
		articles = append(articles, article)
	}
	return articles
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
	item := flag.Int("item",0,"headline number for summary,a  negative headline number will return all full articles, 0 for a list of headlines")
	zok := flag.Int("zok",0,"retrieve the articles from Szokorad galnetarchive instead of the official site. note that there are only headlines or full articles. no summaries! A negative number will return 100 full articles.")
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
	rs := ""
	var articles []map[string]string
	if zok==0 {
		rs = "GALNET News.\n"
		htm := retrieveURL("http://www.elitedangerous.com/news/galnet/")
		articles = off_getArticles(htm)
	} else {
		rs = "Zsokorad Galnet Archive\n\n"
		maxresults := "10"
		if request < 0 {
			maxresults = "100"
		}
		htm := retrieveURL("http://galnetarchive.blogspot.com/feeds/posts/default?max-results="+maxresults+"&alt=rss")
		articles = zok_getArticles(htm)
	}
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
	rs = rs + "\nEnd of Stream.\n"
	outputToFile(rs)
	fmt.Println(rs)
}
	
