# CurveApi.cf

CurveApi was built to provide simple transparent API for [CurveFever](http://curvefever.com).
And it's used on [CurveApi.cf](http://curveapi.cf)

----

## Requirements
1. MongoDB up and running
2. GO package

----
## Installation
1. Get github.com/drone/routes

>go get github.com/julienschmidt/httprouter
	
2. Get mgo

>go get gopkg.in/mgo.v2

3. Build

>go build
	
----
## Usage

### Fetch player profile by id

>http://curveapi.cf/user/793301

**Result:**

    {
 	  "uid": "793301",
      "name": "maciekmm_tk",
      "premium": true,
      "champion": false,
      "picture": "http://curvefever.com/sites/default/files/pictures/picture-793301-1429279990.png",
      "ranks": {
        "1v1_asia": {
          "rank": 700,
          "bonus": 500,
    ...

----

### Fetch player profile by name

>http://curveapi.cf/username/maciekmm_tk

**Result:**

    {
 	  "uid": "793301",
      "name": "maciekmm_tk",
      "premium": true,
      "champion": false,
      "picture": "http://curvefever.com/sites/default/files/pictures/picture-793301-1429279990.png",
      "ranks": {
        "1v1_asia": {
          "rank": 700,
          "bonus": 500,
    ...
	