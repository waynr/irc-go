IRC
===

Go-IRC is tiny IRC library, inspired by the goty_. For an example application,
see the examples_ directory


How to use it?
--------------

* connect: `c, _ := irc.Dial("irc.freenode.net:6667")`
* write: `c.ToSend <- "JOIN #mychan"`
* read: `msg := <-c.Received`


.. _goty: https://github.com/RecursiveForest/goty.git
.. _examples: https://github.com/husio/go-irc/tree/master/examples
