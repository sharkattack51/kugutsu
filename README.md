kugutsu
===

kugutsu is input automation server.

## Options

|Params|Description|Default|
|:--|:--|:--|
|/m|listen mode `UDP` or `HTTP`|UDP|
|/p|listen port number|5000|

## Send Message Examples

- UDP: send message string
    - ex) -> `CLICK:100:200`
- HTTP: send request with a query string:`msg`
    - ex) http://127.0.0.1:5000?msg=CLICK:100:200

|Message|Automation|
|:--|:--|
|CLICK|Mouse click at screen center|
|CLICK:100:200|Mouse click at (100, 200)|
|WCLICK:100:200|Mouse double click at (100, 200)|
|KEY:a|Press key[ `a` ]|
|KEY:a,LSHIFT,RALT,CTRL|Press key[ `a` + `left-shift` + `right-alt` + `ctrl` ]|
|KEY:a:3|Press key[ `a` ] hold 3 sec|
|KEY:a,SHIFT:1|Press key[ `a` + `shift` ] hold 1 sec|
|KEY:abc|Type key[ `a` `b` `c` ]|
