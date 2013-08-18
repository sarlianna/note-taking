note-taking
===========

a command-line note taking utility.  You pass it a message and a tag to store a note, and can have it recall the latest note(s),
latest note(s) specific to a tag, etc.

Mostly meant to be a centralized collection point for important information, both automated and otherwise.


Build/Prereqs
-------------

Currently setup to use postgresql and pq (github.com/lib/pq).  To install pq, run 'go get github.com/lib/pq'.

Table schema:
(id serial,  
tag varchar,  
message varchar,  
age timestamp without timezone NOT NULL DEFAULT now() )  


TODO
--------------

-pull database config out into another file
-write tests
-add features detailed in notes.go's TODO list.
