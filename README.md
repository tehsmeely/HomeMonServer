#HomeMon Server
A simple server built with Golang. Its primary function is as a relay/router between [my website's](http://jontyheiser.co.uk/) host and the [gpioMonitor](https://github.com/tehsmeely/GPIOMonitor) running on the same host as this server. Digest auth provides a level of protection in acquiring the data from the home mon daemon.

Using inbuild golang net and net/http packages for both networking interactions.

Built with:
+ [go-http-auth](http://github.com/abbot/go-http-auth) - For Digest auth on http connections]
+ [Viper](https://github.com/spf13/viper) - For easy config files
