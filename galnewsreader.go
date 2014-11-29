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
	)

func itemRequested() int {
	item := flag.Int("item",0,"headline number or -1 for a list of headlines")
	flag.Parse()
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
	headlines := re.FindAllString(htm, -1)
	return headlines
}

func getArticle(nr int, htm string) string {
	htm = strings.Replace(htm,"View full transmission &raquo;","",-1)
	htm = strings.Replace(htm,"<h3>", "News Item, Stardate ",-1)
	htm = strings.Replace(htm,"<h3>", "Content : ",-1)

	articlesRaw := strings.Split(htm,"<article>")

	if nr > len(articlesRaw) {
		return "This article does not exist."
	}
	
	article := removeTags(articlesRaw[nr])	
	return article
}

func main() {
	request := itemRequested()

	fmt.Println("GALNET News")

	if request < 0 {
		fmt.Println("This article does not exist.")
		return
	}
	resp,err1 := http.Get("http://www.elitedangerous.com/news/galnet/")
	if err1 != nil {
		fmt.Println("ERROR 1")
		fmt.Println(err1)		
	}
	bodyio := resp.Body
	buf, err2 := ioutil.ReadAll(bodyio)
	if err2 != nil {
		fmt.Println("ERROR 1")
		fmt.Println(err2)		
	}

	htm := string(buf)

	if request != 0 {
		fmt.Println("item requested : "+strconv.Itoa(request))	
		fmt.Println(getArticle(request, htm))		
	} else {
		fmt.Println("Current Headlines : \n\n")	
		headlines := getHeadlines(htm)
		for i:=1; i<len(headlines); i++ {
			galDate := getLinkDate(headlines[i])
			fmt.Println("Headline "+strconv.Itoa(i)+".\nStardate "+galDate+".\n"+removeTags(headlines[i])+"\n")
		}
	}
}
	
