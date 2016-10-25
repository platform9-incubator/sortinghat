Whistle-Log

Whistle-Log is an application that is useful for analyzing OpenStack (and generally any log).

It contains a web application and corresponding analysis web-service.

The Web-Service ingest ERROR logs from any source, finds correlation between various ERROR logs
and lets user annotate what is important and not important.

This is a very powerful tool that can be customized for any log analysis. Currently it works with 
PaperTrail, but can be easily ported to consume data from any log aggreation service.

It uses Mongo database as the backend to store and analyze log statements.

