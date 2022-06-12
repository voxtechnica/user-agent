package ua

import (
	"testing"
)

// blackhole is a variable used to ensure that the compiler doesn't optimize Parse out of the benchmark
var blackhole UserAgent

func parseCompare(userAgents, expected []string, t *testing.T) {
	for i, userAgent := range userAgents {
		ua := Parse(userAgent)
		s := ua.String()
		if s != expected[i] {
			t.Errorf("expected/received #%d:\n%s\n%s", i, expected[i], s)
		}
	}
}

func Test_unquote(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single quotes",
			input:    "'Mozilla/5.0 (compatible; Seekport Crawler; http://seekport.com/'",
			expected: "Mozilla/5.0 (compatible; Seekport Crawler; http://seekport.com/",
		},
		{
			name:     "double quotes",
			input:    "\"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4889.0 Safari/537.36\"",
			expected: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4889.0 Safari/537.36",
		},
		{
			name:     "no quotes",
			input:    "SiteScoreBot v20210315 - https://sitescore.ai",
			expected: "SiteScoreBot v20210315 - https://sitescore.ai",
		},
		{
			name:     "unbalanced quotes",
			input:    "'Mozilla/5.0 (compatible; Seekport Crawler; http://seekport.com/",
			expected: "Mozilla/5.0 (compatible; Seekport Crawler; http://seekport.com/",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := unquote(c.input)
			if result != c.expected {
				t.Errorf("expected/received #%s:\n%s\n%s", c.name, c.expected, result)
			}
		})
	}
}

func Test_parseFields(t *testing.T) {
	cases := []struct {
		name  string
		input string
		count int
	}{
		{
			name:  "empty string",
			input: "",
			count: 0,
		},
		{
			name:  "empty parentheses",
			input: "()",
			count: 0,
		},
		{
			name:  "ends with open parenthesis",
			input: "Mozilla/5.0 (",
			count: 0,
		},
		{
			name:  "unbalanced parentheses",
			input: "Mozilla/5.0 (compatible; Seekport Crawler; http://seekport.com/",
			count: 3,
		},
		{
			name:  "ends with parenthesis",
			input: "Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)",
			count: 2,
		},
		{
			name:  "multiple parentheses",
			input: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15 (Applebot/0.1; +http://www.apple.com/go/applebot)",
			count: 7,
		},
		{
			name:  "nested parentheses",
			input: "Mozilla/5.0 (Linux; Android 10; moto e (XT2052DL)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.101 Mobile Safari/537.36",
			count: 8,
		},
		{
			name:  "no parentheses",
			input: "pa11y/6.1.1 ",
			count: 1,
		},
		{
			name:  "square brackets",
			input: "Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 LightSpeed [FBAN/MessengerLiteForiOS;FBAV/361.1.0.40.107;FBBV/370907200;FBDV/iPhone11,2;FBMD/iPhone;FBSN/iOS;FBSV/15.4.1;FBSS/3;FBCR/;FBID/phone;FBLC/en;FBOP/0]",
			count: 18,
		},
		{
			name:  "double spaces",
			input: "Mozilla/5.0 (Linux; Android 10; V2027; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/87.0.4280.141 Mobile Safari/537.36 VivoBrowser/8.3.1.0",
			count: 9,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fields := parseFields(c.input)
			if len(fields) != c.count {
				t.Errorf("expected %d fields, received %d", c.count, len(fields))
			}
		})
	}
}

func Test_majorMinorVersion(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"macOS", "10_11_6", "10.11"},
		{"iOS", "15_1", "15.1"},
		{"Linux", "x86_64", ""},
		{"Windows10", "10", "10"},
		{"WindowsXP", "5.1", "5.1"},
		{"iPhone", "iPhone", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := majorMinorVersion(c.input)
			if result != c.expected {
				t.Errorf("expected/received #%s:\n%s\n%s", c.name, c.expected, result)
			}
		})
	}
}

// TestApplebot tests Apple's official bot User-Agent strings.
// Reference: https://support.apple.com/en-us/HT204683
func TestApplebot(t *testing.T) {
	uas := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.1 Safari/605.1.15 (Applebot/0.1)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15 (Applebot/0.1; +http://www.apple.com/go/applebot)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_4_1 like Mac OS X) AppleWebKit/605.1.15Z (KHTML, like Gecko) Version/13.1 Mobile/15E148 Safari/604.1 (Applebot/0.1)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B410 Safari/600.1.4 (Applebot/0.1; +http://www.apple.com/go/applebot)",
	}
	expected := []string{
		"Bot Applebot 0.1 Desktop macOS 10.14",
		"Bot Applebot 0.1 Desktop macOS 10.15 http://www.apple.com/go/applebot",
		"Bot Applebot 0.1 Mobile iOS 13.4",
		"Bot Applebot 0.1 Mobile iOS 8.1 http://www.apple.com/go/applebot",
	}
	parseCompare(uas, expected, t)
}

// TestBingbot tests Microsoft's official bot User-Agent strings.
// Reference: https://www.bing.com/webmasters/help/which-crawlers-does-bing-use-8c184ec0
func TestBingbot(t *testing.T) {
	uas := []string{
		// Bingbot is our standard crawler and handles most of our crawling needs each day.
		// "W.X.Y.Z" would be substituted with the latest Microsoft Edge version we're using, such as “100.0.4896.127"
		"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)", // historical; discontinue by Fall 2022
		"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm) Chrome/W.X.Y.Z Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/W.X.Y.Z Mobile Safari/537.36 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",

		// AdIdxBot is the crawler used by Bing Ads. AdIdxBot is responsible for crawling ads and following through to
		// websites from those ads for quality control purposes.
		"Mozilla/5.0 (compatible; adidxbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A465 Safari/9537.53 (compatible; adidxbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Mozilla/5.0 (Windows Phone 8.1; ARM; Trident/7.0; Touch; rv:11.0; IEMobile/11.0; NOKIA; Lumia 530) like Gecko (compatible; adidxbot/2.0; +http://www.bing.com/bingbot.htm)",

		// BingPreview is used to generate page snapshots.
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534+ (KHTML, like Gecko) BingPreview/1.0b",
		"Mozilla/5.0 (Windows Phone 8.1; ARM; Trident/7.0; Touch; rv:11.0; IEMobile/11.0; NOKIA; Lumia 530) like Gecko BingPreview/1.0b",
	}
	expected := []string{
		"Bot Bingbot 2.0 Desktop Other http://www.bing.com/bingbot.htm",
		"Bot Bingbot 2.0 Desktop Other http://www.bing.com/bingbot.htm",
		"Bot Bingbot 2.0 Mobile Android 6.0 http://www.bing.com/bingbot.htm",
		"Bot AdIdxBot 2.0 Desktop Other http://www.bing.com/bingbot.htm",
		"Bot AdIdxBot 2.0 Mobile iOS 7.0 http://www.bing.com/bingbot.htm",
		"Bot AdIdxBot 2.0 Mobile Windows http://www.bing.com/bingbot.htm",
		"Bot BingPreview 1.0b Desktop Windows 6.1",
		"Bot BingPreview 1.0b Mobile Windows",
	}
	parseCompare(uas, expected, t)
}

// TestFacebookbot tests Facebook's Crawler User-Agent strings.
// Reference: https://developers.facebook.com/docs/sharing/webmasters/crawler/
func TestFacebookbot(t *testing.T) {
	uas := []string{
		"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
		"facebookexternalhit/1.1",
	}
	expected := []string{
		"Bot FacebookBot 1.1 Desktop Other http://www.facebook.com/externalhit_uatext.php",
		"Bot FacebookBot 1.1 Desktop Other",
	}
	parseCompare(uas, expected, t)
}

// TestGooglebot tests Google's crawler User-Agent strings
// Reference: https://developers.google.com/search/docs/advanced/crawling/overview-google-crawlers
func TestGooglebot(t *testing.T) {
	uas := []string{
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4851.0 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/101.0.4951.64 Safari/537.36",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"AdsBot-Google (+http://www.google.com/adsbot.html)",
		"Mozilla/5.0 (Linux; Android 5.0; SM-G920A) AppleWebKit (KHTML, like Gecko) Chrome Mobile Safari (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P; AdsBot-Google; +http://www.google.com/bot.html) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Mobile/15E148 Safari/604.1 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		"Mozilla/5.0 (Linux; Android 7.0; SM-G930V Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36 (compatible; Google-Read-Aloud; +https://support.google.com/webmasters/answer/1061943)",
	}
	expected := []string{
		"Bot Googlebot 2.1 Mobile Android 6.0 http://www.google.com/bot.html",
		"Bot Googlebot 2.1 Desktop Other http://www.google.com/bot.html",
		"Bot Googlebot 2.1 Desktop Other http://www.google.com/bot.html",
		"Bot Google-AdsBot Desktop Other http://www.google.com/adsbot.html",
		"Bot Google-AdsBot Mobile Android 5.0 http://www.google.com/mobile/adsbot.html",
		"Bot Google-AdsBot Mobile Android 6.0 http://www.google.com/bot.html",
		"Bot Google-AdsBot Mobile iOS 14.7 http://www.google.com/mobile/adsbot.html",
		"Bot Google-Read-Aloud Mobile Android 7.0 https://support.google.com/webmasters/answer/1061943",
	}
	parseCompare(uas, expected, t)
}

// TestBots tests a variety of small crawler User-Agent strings
func TestBots(t *testing.T) {
	uas := []string{
		"Mozilla/5.0 (compatible; Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 (compatible; Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)",
		"Mozilla/5.0 (compatible; Cincraw/1.0; +http://cincrawdata.net/bot/)",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/101.0.4950.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/78.0.3882.0 Safari/537.36",
		"Mozilla/5.0 (compatible; HubSpot Crawler; +https://www.hubspot.com)",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.0 Safari/537.36 (compatible; Linespider/1.1; +https://lin.ee/4dwXkTH)",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36 (compatible; PagePeeker/3.0; +https://pagepeeker.com/robots/)",
		"Mozilla/5.0 (compatible; Pinterestbot/1.0; +http://www.pinterest.com/bot.html)",
		"'Mozilla/5.0 (compatible; Seekport Crawler; http://seekport.com/'",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36 SeoSiteCheckup (https://seositecheckup.com)",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3835.0 Safari/537.36 (compatible; Sitebulb/1.1; +https://sitebulb.com)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 +https://sitebulb.com",
		"SiteScoreBot v20210315 - https://sitescore.ai",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.75 Safari/537.36 (compatible; SMTBot/1.0; http://www.similartech.com/smtbot)",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.0 Safari/537.36 (compatible; Yeti/1.1; +http://naver.me/spd)",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 YisouSpider/5.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/56.0.2924.75 Mobile/14E5239e YisouSpider/5.0 Safari/602.1",
	}
	expected := []string{
		"Bot Baiduspider 2.0 Desktop Other http://www.baidu.com/search/spider.html",
		"Bot Baiduspider 2.0 Mobile iOS 9.1 http://www.baidu.com/search/spider.html",
		"Bot Cincraw 1.0 Desktop Other http://cincrawdata.net/bot/",
		"Bot HeadlessChrome 101.0 Desktop Linux",
		"Bot HeadlessChrome 78.0 Desktop Windows 10.0",
		"Bot HubSpot Desktop Other https://www.hubspot.com",
		"Bot Linespider 1.1 Desktop Windows 6.1 https://lin.ee/4dwXkTH",
		"Bot PagePeeker 3.0 Desktop Windows 6.3 https://pagepeeker.com/robots/",
		"Bot Pinterestbot 1.0 Desktop Other http://www.pinterest.com/bot.html",
		"Bot Seekport Desktop Other http://seekport.com/",
		"Bot SeoSiteCheckup Desktop Linux https://seositecheckup.com",
		"Bot Sitebulb 1.1 Desktop Windows 10.0 https://sitebulb.com",
		"Bot Chrome 41.0 Mobile Android 6.0 https://sitebulb.com",
		"Bot SiteScoreBot Desktop Other https://sitescore.ai",
		"Bot SMTBot 1.0 Desktop Windows 10.0 http://www.similartech.com/smtbot",
		"Bot Yeti 1.1 Desktop Windows 6.1 http://naver.me/spd",
		"Bot YisouSpider 5.0 Desktop Windows 6.1",
		"Bot YisouSpider 5.0 Mobile iOS 10.3",
	}
	parseCompare(uas, expected, t)
}

// TestApplications tests various application User-Agent strings.
func TestApplications(t *testing.T) {
	uas := []string{
		"Mozilla/5.0 (iPad; CPU OS 15_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/19D52 [FBAN/FBIOS;FBAV/358.0.0.29.112;FBBV/357001037;FBDV/iPad13,1;FBMD/iPad;FBSN/iPadOS;FBSV/15.3.1;FBSS/2;FBID/tablet;FBLC/en_US;FBOP/5;FBRV/358280696]",
		"Mozilla/5.0 (iPad; CPU OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/19E258 [FBAN/FBIOS;FBDV/iPad8,1;FBMD/iPad;FBSN/iPadOS;FBSV/15.4.1;FBSS/2;FBID/tablet;FBLC/en_US;FBOP/5]",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 LightSpeed [FBAN/MessengerLiteForiOS;FBAV/360.1.0.28.109;FBBV/369327353;FBDV/iPhone13,3;FBMD/iPhone;FBSN/iOS;FBSV/15.4.1;FBSS/3;FBCR/;FBID/phone;FBLC/en;FBOP/0]",
		"Mozilla/5.0 (Linux; Android 9; KFMAWI Build/PS7322; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/98.0.4758.101 Safari/537.36 [FB_IAB/Orca-Android;FBAV/357.0.0.13.112;]",
		"Mozilla/5.0 (Linux; Android 12; SM-S908U1 Build/SP1A.210812.016; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/101.0.4951.61 Mobile Safari/537.36 [FB_IAB/Orca-Android;FBAV/360.0.0.10.113;]",
		"Mozilla/5.0 (Linux; Android 11; Pixel 4 Build/RQ3A.211001.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.127 Mobile Safari/537.36 [FB_IAB/FB4A;FBAV/364.1.0.25.132;]",
		"Mozilla/5.0 (iPad; CPU OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) GSA/213.0.449417121 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPod; CPU iPhone OS 12_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) GSA/160.0.373863126 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) GSA/213.0.449417121 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 12; SM-N981U1 Build/SP1A.210812.016; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.127 Mobile Safari/537.36 Instagram 233.0.0.13.112 Android (31/12; 420dpi; 1080x2182; samsung; SM-N981U1; c1q; qcom; en_US; 367202479)",
		"Mozilla/5.0 (iPad; CPU OS 14_8 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Instagram 228.0.0.15.111 (iPad6,12; iOS 14_8; en_US; en-US; scale=2.00; 750x1334; 359294435)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Instagram 236.0.0.19.111 (iPhone11,8; iOS 15_4_1; en_US; en-US; scale=2.00; 828x1792; 371179233) NW/1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 [LinkedInApp]/9.22.7334",
		"Mozilla/5.0 (Linux; Android 12; SM-N986U Build/SP1A.210812.016; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.127 Mobile Safari/537.36 [Pinterest/Android]",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Safari/604.1 [Pinterest/iOS]",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Snapchat/11.78.0.29 (like Safari/8613.1.17.0.8, panda)",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x6305002e)",
		"Mozilla/5.0 (Linux; Android 11; M2102J2SC Build/RKQ1.200826.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/3195 MMWEBSDK/20220204 Mobile Safari/537.36 MMWEBID/6026 MicroMessenger/8.0.20.2100(0x2800143D) Process/toolsmp WeChat/arm64 Weixin NetType/WIFI Language/zh_CN ABI/arm64",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x18001239) NetType/4G Language/zh_CN",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Flipboard/4.2.143(1704)",
	}
	expected := []string{
		"App Facebook 15.3 Tablet iPadOS 15.3",
		"App Facebook 15.4 Tablet iPadOS 15.4",
		"App Facebook 15.4 Mobile iOS 15.4",
		"App Facebook 357.0 Tablet Android 9",
		"App Facebook 360.0 Mobile Android 12",
		"App Facebook 364.1 Mobile Android 11",
		"App GoogleSearch 213.0 Tablet iPadOS 15.5",
		"App GoogleSearch 160.0 Mobile iOS 12.5",
		"App GoogleSearch 213.0 Mobile iOS 15.5",
		"App Instagram Mobile Android 12",
		"App Instagram Tablet iPadOS 14.8",
		"App Instagram Mobile iOS 15.4",
		"App LinkedIn 9.22 Mobile iOS 15.4",
		"App Pinterest Android Mobile Android 12",
		"App Pinterest iOS Mobile iOS 15.4",
		"App Snapchat 11.78 Mobile iOS 15.4",
		"App WeChat 7.0 Desktop Windows 6.1",
		"App WeChat 8.0 Mobile Android 11",
		"App WeChat 8.0 Mobile iOS 15.3",
		"Browser Safari Desktop macOS 10.15", // Unrecognized application (Flipboard)
	}
	parseCompare(uas, expected, t)
}

// TestBrowsers tests various less-common browser User-Agent strings
func TestBrowsers(t *testing.T) {
	uas := []string{
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 ADG/11.0.3684 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 9; KFTRWI) AppleWebKit/537.36 (KHTML, like Gecko) Silk/98.7.2 like Chrome/98.0.4758.136 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 7.1.2; AFTMM) AppleWebKit/537.36 (KHTML, like Gecko) Silk/100.1.153 like Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/4.0 (compatible; GoogleToolbar 4.0.1019.5266-big; Windows XP 5.1; MSIE 6.0.2900.2180)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 10.0; WOW64; Trident/7.0; .NET4.0C; .NET4.0E; wbx 1.0.0; wbxapp 1.0.0; Zoom 3.6.0)",
		"Mozilla/5.0 (Windows NT 10.0; Trident/7.0; MALC; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; EIE10;ENUSMCM; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; Touch; rv:11.0) like Gecko",
		"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Trident/6.0; Mozilla/5.0 (Linux; U; Android 4.0.3; de-ch; HTC Sensation Build/IML74K) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30)",
		"Mozilla/5.0 (Linux; Andr0id 9; BRAVIA 4K UR3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36 OPR/46.0.2207.0 OMI/4.21.0.273.DIA6.199 Model/Sony-BRAVIA-4K-UR3",
		"Mozilla/5.0 (Linux; U; Android 12; SM-N975U Build/SP1A.210812.016; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/96.0.4664.104 Mobile Safari/537.36 OPR/61.0.2254.59937",
		"Mozilla/5.0 (Linux; Android 12; Pixel 4a (5G)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 Mobile Safari/537.36 OPR/68.2.3557.64219",
		"Mozilla/5.0 (Linux; Android 12; SM-G991B Build/SP1A.210812.016) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Mobile Safari/537.36 OPT/2.9",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Safari/604.1 OPT/3.2.11",
		"Mozilla/5.0 (iPad; CPU OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 OPT/3.2.13",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36 OPR/85.0.4341.18",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36 OPR/85.0.4341.75",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36 OPR/85.0.4341.75 (Edition avira-2)",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36 OPR/38.0.2220.41",
		"Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/101.0.4951.61 DuckDuckGo/5 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 DuckDuckGo/7 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 DuckDuckGo/7",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Mobile/15E148 DuckDuckGo/7 Safari/605.1.15",
		"Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/101.0.4951.61 Mobile DuckDuckGo/5 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.127 DuckDuckGo/5 Safari/537.36",
	}
	expected := []string{
		"Browser AOLDesktop 11.0 Desktop Windows 6.2",
		"Browser Silk 98.7 Tablet Android 9",
		"Browser Silk 100.1 Tablet Android 7.1",
		"Browser InternetExplorer Desktop Windows",
		"Browser InternetExplorer Desktop Windows 10.0",
		"Browser InternetExplorer 11.0 Desktop Windows 10.0",
		"Browser InternetExplorer 11.0 Desktop Windows 6.1",
		"Browser InternetExplorer 11.0 Desktop Windows 10.0",
		"Browser InternetExplorer Mobile Android 4.0",
		"Browser Opera 46.0 Desktop Linux",
		"Browser Opera 61.0 Mobile Android 12",
		"Browser Opera 68.2 Mobile Android 12",
		"Browser Opera 2.9 Mobile Android 12",
		"Browser Opera 3.2 Mobile iOS 15.4",
		"Browser Opera 3.2 Tablet iPadOS 15.4",
		"Browser Opera 85.0 Desktop Linux",
		"Browser Opera 85.0 Desktop macOS 10.15",
		"Browser Opera 85.0 Desktop Windows 10.0",
		"Browser Opera 38.0 Desktop Windows 6.1",
		"Browser DuckDuckGo 5 Desktop Linux",
		"Browser DuckDuckGo 7 Desktop macOS 10.15",
		"Browser DuckDuckGo 7 Mobile iOS 14.6",
		"Browser DuckDuckGo 7 Mobile iOS 15.5",
		"Browser DuckDuckGo 5 Mobile Android 12",
		"Browser DuckDuckGo 5 Tablet Android 12",
	}
	parseCompare(uas, expected, t)
}

// TestFirefox tests a variety of Firefox (and derivative) User-Agent strings
func TestFirefox(t *testing.T) {
	uas := []string{
		"Mozilla/5.0 (Linux; Android 12; SM-S908B Build/SP1A.210812.016; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.127 Mobile Safari/537.36 Klarna/22.16.177",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Klarna/22.18.281",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/24.1 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/100.1 Mobile/15E148 Safari/605.1.15",
		"Mozilla/5.0 (iPad; CPU OS 15_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/14.0b12646 Mobile/15E148 Safari/605.1.15",
		"Mozilla/5.0 (Android 12; Mobile; rv:102.0) Gecko/102.0 Firefox/102.0",
		"Mozilla/5.0 (Mobile; ALCATEL 4056W; rv:84.0) Gecko/84.0 Firefox/84.0 KAIOS/3.0",
		"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:99.0) Gecko/20100101 Firefox/99.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:99.0) Gecko/20100101 Firefox/99.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:99.0) Gecko/20100101 Firefox/99.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/102.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:98.0) Gecko/20100101 Firefox/98.0 OpenWave/93.4.3993.94",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:100.0) Gecko/20100101 Firefox/100.0",
		"Mozilla/5.0 (Windows NT 5.1; rv:52.0) Gecko/20100101 Firefox/52.0",
		"Mozilla/5.0 (Windows; rv:81.0) Gecko/20100101 Firefox/81.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:65.0) Gecko/20100101 Firefox/65.0 IceDragon/65.0.2",
		"Mozilla/5.0 (X11; Linux x86_64; rv:38.0) Gecko/20100101 Firefox/38.0 Iceweasel/38.6.1",
		"Mozilla/5.0 (Windows NT 5.1; rv:68.0) Gecko/20100101 Goanna/4.8 Firefox/68.0 Mypal/29.3.0",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64; rv:68.0) Gecko/20100101 Goanna/4.8 Firefox/68.0 PaleMoon/29.4.4",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:91.0) Gecko/20100101 Firefox/91.0 Waterfox/91.10.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0 Waterfox/91.10.0",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:56.0; Waterfox) Gecko/20100101 Firefox/56.2.7",
		"Mozilla / 5.0(X11;Ubuntu;Linuxi686;rv:81.0) Gecko / 20100101Firefox / 81.0",
	}
	expected := []string{
		"Browser Firefox 22.16 Mobile Android 12",
		"Browser Firefox 22.18 Mobile iOS 15.4",
		"Browser Firefox 24.1 Desktop macOS 10.15",
		"Browser Firefox 100.1 Mobile iOS 15.5",
		"Browser Firefox 14.0b12646 Tablet iPadOS 15.0",
		"Browser Firefox 102.0 Mobile Android 12",
		"Browser Firefox 84.0 Mobile Other",
		"Browser Firefox 99.0 Desktop Linux",
		"Browser Firefox 99.0 Desktop Linux",
		"Browser Firefox 100.0 Desktop Linux",
		"Browser Firefox 99.0 Desktop macOS 10.15",
		"Browser Firefox 102.0 Desktop Windows 10.0",
		"Browser Firefox 98.0 Desktop Windows 10.0",
		"Browser Firefox 100.0 Desktop Windows 6.1",
		"Browser Firefox 52.0 Desktop Windows 5.1",
		"Browser Firefox 81.0 Desktop Windows",
		"Browser Firefox 65.0 Desktop Windows 10.0",
		"Browser Firefox 38.0 Desktop Linux",
		"Browser Firefox 68.0 Desktop Windows 5.1",
		"Browser Firefox 68.0 Desktop Windows 6.3",
		"Browser Firefox 91.0 Desktop macOS 10.15",
		"Browser Firefox 91.0 Desktop Windows 10.0",
		"Browser Firefox 56.2 Desktop Windows 6.1",
		"Browser Firefox Desktop Linux",
	}
	parseCompare(uas, expected, t)
}

// TestEdge tests a variety of Microsoft Edge User-Agent strings
func TestEdge(t *testing.T) {
	uas := []string{
		"Mozilla/5.0 (Linux; Android 10; Redmi Note 9 Pro Build/QKQ1.191215.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.127 Mobile Safari/537.36 EdgW/1.0",
		"Mozilla/5.0 (Linux; Android 12; SM-N986U Build/SP1A.210812.016; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/101.0.4951.61 Mobile Safari/537.36 EdgW/1.0",
		"Mozilla/5.0 (Linux; Android 12; Pixel 6 Pro Build/SP2A.220405.004; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.127 Mobile Safari/537.36 EdgW/1.0",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 7 Build/MOB30X; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/99.0.4844.88 Safari/537.36 EdgW/1.0",
		"Mozilla/5.0 (Linux; Android 10; SM-T830) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 EdgA/100.0.1185.50",
		"Mozilla/5.0 (iPad; CPU OS 12_5_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 EdgiOS/46.3.30 Mobile/15E148 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36 Edg/101.0.1210.53 PTST/220520.174724",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5032.0 Safari/537.36 Edg/103.0.1255.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.41 Safari/537.36 Edg/101.0.1210.32 Herring/100.1.3060.61",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; WebView/3.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.19043",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.19041",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36 Edg/100.0.1185.29",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36 Edg/99.0.1150.55/8mqLkJuL-86",
		"Mozilla/5.0 (Linux; Android 12; SM-N986U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Mobile Safari/537.36 EdgA/101.0.1210.47",
		"Mozilla/5.0 (Linux; Android 12; SM-S901U1 Build/SP1A.210812.016; ) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/95.0.4638.69 Mobile Safari/537.36 EdgA/95.0.1020.48 BingSapphire/22.6.400511306",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_8 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) EdgiOS/101.0.1210.47 Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 15_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) EdgiOS/101.0.1210.47 Version/15.0 Mobile/15E148 Safari/604.1",
	}
	expected := []string{
		"Browser Edge 1.0 Mobile Android 10",
		"Browser Edge 1.0 Mobile Android 12",
		"Browser Edge 1.0 Mobile Android 12",
		"Browser Edge 1.0 Tablet Android 6.0",
		"Browser Edge 100.0 Tablet Android 10",
		"Browser Edge 46.3 Tablet iPadOS 12.5",
		"Browser Edge 101.0 Desktop Linux",
		"Browser Edge 103.0 Desktop macOS 10.15",
		"Browser Edge 101.0 Desktop Windows 10.0",
		"Browser Edge 18.19043 Desktop Windows 10.0",
		"Browser Edge 18.19041 Desktop Windows 10.0",
		"Browser Edge 100.0 Desktop Windows 6.1",
		"Browser Edge 99.0 Desktop Windows 6.1",
		"Browser Edge 101.0 Mobile Android 12",
		"Browser Edge 95.0 Mobile Android 12",
		"Browser Edge 101.0 Mobile iOS 14.8",
		"Browser Edge 101.0 Tablet iPadOS 15.3",
	}
	parseCompare(uas, expected, t)
}

// TestChrome tests a variety of Google Chrome User-Agent strings
func TestChrome(t *testing.T) {
	uas := []string{
		"Mozilla/5.0 (X11; CrOS aarch64 14526.89.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.133 Safari/537.36",
		"Mozilla/5.0 (X11; CrOS armv7l 13597.84.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.106 Safari/537.36",
		"Mozilla/5.0 (X11; CrOS x86_64 14695.25.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64)  AppleWebKit/537.36 (KHTML, like Gecko; Google Web Preview)  Chrome/99.0.4844.84 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3694.0 Safari/537.36 Chrome-Lighthouse",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12_3_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36",
		"Mozilla/5.0 (DomainUser;Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36/sgrRRPuj-29",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36 LikeWise/100.6.1315.16",
		"Mozilla/5.0 (Linux; Android 12; SM-F916U1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.61 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; KingPad_K10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64; HTC/Sensation/3.32.163.12; fr-fr) AppleWebKit/534.24 (KHTML, like Gecko) Chrome/11.0.696.34 Safari/534.24",
		"Mozilla/5.0 (Linux; NetCast; U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36 SmartTV/10.0 Colt/2.0",
	}
	expected := []string{
		"Browser Chrome 100.0 Desktop ChromeOS",
		"Browser Chrome 98.0 Desktop ChromeOS",
		"Browser Chrome 102.0 Desktop ChromeOS",
		"Browser Chrome 101.0 Desktop Linux",
		"Browser Chrome 99.0 Desktop Linux",
		"Browser Chrome 74.0 Desktop macOS 10.13",
		"Browser Chrome 89.0 Desktop macOS 12.3",
		"Browser Chrome 67.0 Desktop Windows 10.0",
		"Browser Chrome 101.0 Desktop Windows 10.0",
		"Browser Chrome 101.0 Desktop Windows 10.0",
		"Browser Chrome 101.0 Desktop Windows 10.0",
		"Browser Chrome 101.0 Tablet Android 12",
		"Browser Chrome 100.0 Tablet Android 10",
		"Browser Chrome 11.0 Desktop Linux",
		"Browser Chrome 79.0 Desktop Linux",
	}
	parseCompare(uas, expected, t)
}

// TestSafari tests a variety of Apple Safari User-Agent Strings
func TestSafari(t *testing.T) {
	uas := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPod touch; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15G77",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 15_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
		"Mozilla/5.0 (iPad; CPU OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4.1 Mobile/15E148 Safari/605.1.15 BingSapphire/1.0.400324001",
	}
	expected := []string{
		"Browser Safari Desktop macOS 10.15",
		"Browser Safari 15.4 Desktop macOS 10.15",
		"Browser Safari 15.5 Desktop macOS 10.15", // mobile spoofing desktop
		"Browser Safari 15.4 Mobile iOS 15.4",
		"Browser Safari Mobile iOS 11.4",
		"Browser Safari 15.5 Mobile iOS 15.6",
		"Browser Safari Tablet iPadOS 15.3",
		"Browser Safari 15.4 Tablet iPadOS 15.4",
	}
	parseCompare(uas, expected, t)
}

// BenchmarkParse checks performance on parsing different User-Agent strings.
// Note that some are detected earlier in the cascade (e.g. bots and applications).
func BenchmarkParse(b *testing.B) {
	agents := []struct {
		name string
		ua   string
	}{
		{"Googlebot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"},
		{"Chrome", "Mozilla/5.0 (X11; CrOS x86_64 14695.25.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"},
		{"Firefox", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/102.0"},
		{"Safari", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.5 Mobile/15E148 Safari/604.1"},
	}
	for _, a := range agents {
		b.Run("Parse-"+a.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ua := Parse(a.ua)
				// avoid optimizing Parse out of the loop
				blackhole = ua
			}
		})
	}
}
