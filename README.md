Zabbix Passenger monitoring utility
===================================

This small utility parses the output of passenger-status and produces values
that can be used with the accompanying zabbix template (zabbix-template.xml).

Copy the binary to /usr/local/bin and put the userparameter_passenger.conf file
to the correct location to be included in your agent configuration.

The helper binary looks for passenger-status in the path or a RVM wrapper for
the default Ruby.

Disclaimer
----------

I am not a Go developer, so this code is pretty crude and probably not idiomatic
Go.
