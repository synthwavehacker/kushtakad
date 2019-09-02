package main

import (
	"net/url"
	"testing"
)

var cases = []struct {
	parent string
	uri    string
	want   string
}{
	{
		"http://www.example.com",
		"http://example.com/UserFiles/Servers/Server_543/Templates/2017/img/map-bg.png",
		"http://example.com/UserFiles/Servers/Server_543/Templates/2017/img/map-bg.png",
	},
	{
		"http://www.example.com/a/b/c/d/e/f/g/",
		"../../../../updates/concrete5.7.5.6/concrete/../concrete/images/devices/galaxy/s5-landscape.png",
		"http://www.example.com/a/b/c/updates/concrete5.7.5.6/concrete/images/devices/galaxy/s5-landscape.png",
	},
	{
		"http://www.example.com/a/b/c/d/e/f/g/",
		"/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.eot?#iefix",
		"http://www.example.com/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.eot?#iefix",
	},
	{
		"http://www.example.com/a/b/c/d/e/f/g/file.css",
		"./application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.eot?#iefix",
		"http://www.example.com/a/b/c/d/e/f/g/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.eot?#iefix",
	},
	{
		"http://www.example.com/a/b/c/d/e/f/g/file.css",
		"https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css",
		"https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css",
	},
	{
		"https://www.example.z12.zz.us/application/files/cache/css/172_style.css",
		"/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.ttf",
		"https://www.example.z12.zz.us/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.ttf",
	},
	{
		"http://www.example.com/a/b/c/d/e/f/g/file.css",
		"https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css",
		"https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css",
	},
	{
		"https://www.example.org/common/pages/Portalstatictyle.ashx?templateId=2739",
		"//example.org/UserFiles/Servers/Server_543/Templates/2017/img/map-bg.png",
		"https://example.org/UserFiles/Servers/Server_543/Templates/2017/img/map-bg.png",
	},
	{
		"https://www.example.net/Static//site/assets/styles/dashboard.css",
		"../../../GlobalAssets/webfonts/NotoSerif-Regular.ttf",
		"https://www.example.net/Static//GlobalAssets/webfonts/NotoSerif-Regular.ttf",
	},
}

func TestFunc(t *testing.T) {
	for _, c := range cases {
		puri, err := url.Parse(c.parent)
		if err != nil {
			t.Fatal(err)
		}

		got, err := absUrl(c.uri, puri)
		if err != nil {
			t.Fatal(err)
		}

		if got.String() != c.want {
			t.Errorf("absUrl(%q) == %q, want %q", c.uri, got, c.want)
		}
	}
}

/*
"/application/static/blps/img/winter-spot.jpg",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/iphone/iphone5.png",
"../../../../updates/concrete5.7.5.6/concrete//images/icons/search.png",
"../../../../updates/concrete5.7.5.6/concrete//images/icons/wrench.png",
"../../../../updates/concrete5.7.5.6/concrete//images/icons/arrow_down.png",
"/application/static/blps/img/touchbase-bg.png",
"fonts/fontawesome-webfont.woff?v=4.2.0",
"https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css",
"../../../../updates/concrete5.7.5.6/concrete//images/newsflow_logo_welcome_back.png",
"../../../../updates/concrete5.7.5.6/concrete//images/bg_icon_item_grid_overlay.png",
"fonts/icomoon.svg",
"/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.ttf",
"../../../../updates/concrete5.7.5.6/concrete//images/bg_toolbar.png",
"/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.eot",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/iphone/iphone6.png",
"../../../../updates/concrete5.7.5.6/concrete//images/icons/stack.png",
"fonts/bpa-infographic3.woff2",
"/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.eot?#iefix",
"//seattleschools.org/UserFiles/Servers/Server_543/Templates/2017/img/icon/icon-home.png",
"images/logo.svg",
"fonts/icomoon.eot",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/galaxy/s5.png",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/iphone/iphone4.png",
"fonts/fontawesome-webfont.eot?v=4.2.0",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/iphone/iphone5-landscape.png",
"https://p8cdn2static.sharpschool.com/Common/resources/DesignPortfolio/Sitestatic/CommonLib/social-media-2014/rss/defualt-rss-icon-2014.png",
"/application/static/blps/img/spotlight-pattern.png",
"fonts/fontawesome-webfont.ttf?v=4.2.0",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/iphone/iphone4-landscape.png",
"/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.svg#glyphicons_halflingsregular",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/iphone/iphone6plus-landscape.png",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/ipad/ipad-landscape.png",
"fonts/icomoon.woff2",
"fonts/icomoon.woff",
"images/logo.png",
"/application/static/blps/img/grip.png",
"/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.woff",
"fonts/fontawesome-webfont.svg?v=4.2.0#fontawesomeregular",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/ipad/ipad.png",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/iphone/iphone6plus.png",
"fonts/bpa-infographic3.svg",
"../../../../packages/blnews/blocks/blnews_list/templates/blnews_slider/images/arrows.gif",
"../../../../updates/concrete5.7.5.6/concrete//../concrete/images/devices/iphone/iphone6-landscape.png",
"fonts/bpa-infographic3.eot",
"fonts/fontawesome-webfont.eot?#iefix&v=4.2.0",
"/application/static/blps/bootstrap/app/fonts/glyphicons-halflings-regular.woff2",
"fonts/icomoon.ttf",
"fonts/bpa-infographic3.ttf",
"fonts/bpa-infographic3.woff",
"//seattleschools.org/UserFiles/Servers/Server_543/Templates/2017/img/books-bg.jpg"]
*/
