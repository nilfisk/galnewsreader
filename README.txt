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

Version: 2.0
The source file should always be in the 7z archive! Please contact me if it is missing

Changelog :
	2.0
	  new:
	  * added Szokorad's Galnet Archive as alternative feed
	  * set the galnet archive as default feed as Frontier took
		down the official page
	1.2
	  new:
	  * the output is writen to a file ./galnews
	  fixes:
	  * the numbering between headlines and news was off
	  * reduced the number of new lines.
	  * fix for the weird spelling of "transmission"
	  * fix for the weird :: in some headlines

For old users : 
	As of Dec.24th Frontier seems to have taken down the Galnet News Page, so a new (default!) news feed has been added, namely Cmdr Szokorad galnetarchive which can be found at 
		http://galnetarchive.blogspot.com/feeds/posts/default
	Note that the parameters have been modified too!
 	  
Description:
 Small GO Program allowing to read the headlines from Frontier's Elite:Dangerous GalNews
 page or single article summaries. The idea is to have Windows TextToSpeech Engine read
 the news aloud while flying.
 A small powershell script to "say" a string is included but there doesn't seem to be any
 way to call it without windows taking focus from E:D and dropping to the desktop.
 

License:
 why GPL3? why *NOT*????

How to compile for Windows:
in order to disable the console window which pops up in Windows when the reader is started, you need to compile it with:

         go build -ldflags -H=windowsgui galnewsreader.go

Usage:
 galnewsreader [-item=n] [-zok={0|1}]
 
 if zok=1 (Default bevahiour) then the news from Szokorad's Galnet Archive will be used instead
 of Frontier's official page.
 
 if -item is specified, only the corresponding summary 'n' will be printed out. If it is not
 specified or n equals 0 only the headlines will be printed out.
 For the Frontier Galnet News page :
	if n is negative the full body of the transmission will be printed instead of the summary.
 For Szokorad's Galnet Archive :
	if n is negative the full text of the 100 first entries will be printed.

	
	
 Using Voice Attack and Luca D's excellent plugin (https://groups.google.com/forum/#!topic/voiceattack/VotudmB84hE) 
 is highly recommended. Use the readConsole example from the plugin's sample profile or the included galnewsreader-sample.vap as template 

 
 Enjoy,
 Sammy Fischer