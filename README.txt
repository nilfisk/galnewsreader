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


Version: 1.1

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
 galnewsreader [-item=n]
 if -item is specified, only the corresponding summary 'n' will be printed out. If it is not
 specified or n equals 0 only the headlines will be printed out.
 if n is negative the full body of the transmission will be printed instead of the summary.

 Using Voice Attack and Luca D's excellent plugin (https://groups.google.com/forum/#!topic/voiceattack/VotudmB84hE) 
 is highly recommended. Use the readConsole example from the plugin's sample profile or the included galnewsreader-sample.vap as template 
