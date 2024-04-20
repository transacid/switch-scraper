# NETGEAR switch metrics scrapper

## About
this is a foolish attempt to get metrics out of a managed netgear switch. I have tested it with my `GS108Ev3`.
> [!CAUTION]
> there is **no guarantee** that this will work with your switch.

Since netgear does not provide an API it basically is webscraping from the port-monitor page. The scraper logs itself in and uses the returned cookie to continue scraping.
## Known bugs
this whole construct is super fragile, sometimes the switch doesn't play nice and you get "banned" for around 5 minutes. it often continues.

## Usage
the scraper expects two ENV variables:
| Key       | Value                                               |
| --------- | --------------------------------------------------- |
| REMOTE_IP | this is the IP of the switch                        |
| PASSWORD  | the password used to authenticate against the webUI |

the prometheus endpoint is at `:8080/metrics`.
## Finally
No guarantees, no warranty, no nothing. You are on your own.
If you have suggestions I'm happy for pull requests. It make take a super long while till I respond. Please don't file issues, I do this in my spare time. Feel free to fork :)

[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech/)