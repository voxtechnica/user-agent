package user_agent

import "strings"

// patterns are used to identify appropriate fields in the UserAgent struct.
var patterns = []match{
	{find: "Macintosh", deviceType: "Desktop", operatingSystem: "macOS"},
	{find: "iPad", deviceType: "Tablet", operatingSystem: "iPadOS"},
	{find: "iPhone", deviceType: "Mobile", operatingSystem: "iOS"},
	{find: "Mobile", deviceType: "Mobile"},
	{find: "Android", deviceType: "Tablet", operatingSystem: "Android"},  // "Mobile" catches Android Mobile first
	{find: "Windows", deviceType: "Desktop", operatingSystem: "Windows"}, // "Mobile" catches Windows Mobile first
	{find: "CrOS", deviceType: "Desktop", operatingSystem: "ChromeOS"},
	{find: "Tizen", operatingSystem: "Tizen"},
	{find: "Linux", operatingSystem: "Linux"},
	{find: "pa11y", clientType: "Bot", clientName: "Pa11y"},
	{find: "AhrefsBot", clientType: "Bot", clientName: "AhrefsBot"},
	{find: "Applebot", clientType: "Bot", clientName: "Applebot"},
	{find: "Baiduspider", clientType: "Bot", clientName: "Baiduspider"},
	{find: "adidxbot", clientType: "Bot", clientName: "AdIdxBot"},
	{find: "bingbot", clientType: "Bot", clientName: "Bingbot"},
	{find: "BingPreview", clientType: "Bot", clientName: "BingPreview"},
	{find: "Cincraw", clientType: "Bot", clientName: "Cincraw"},
	{find: "facebookexternalhit", clientType: "Bot", clientName: "FacebookBot"},
	{find: "Googlebot", clientType: "Bot", clientName: "Googlebot"},
	{find: "AdsBot-Google", clientType: "Bot", clientName: "Google-AdsBot"},
	{find: "Google-Adwords", clientType: "Bot", clientName: "Google-AdWords"},
	{find: "Google-Read-Aloud", clientType: "Bot", clientName: "Google-Read-Aloud"},
	{find: "Google-Structured-Data-Testing-Tool", clientType: "Bot", clientName: "Google-Testing"},
	{find: "HeadlessChrome", clientType: "Bot", clientName: "HeadlessChrome"},
	{find: "HubSpot", clientType: "Bot", clientName: "HubSpot"},
	{find: "Linespider", clientType: "Bot", clientName: "Linespider"},
	{find: "PagePeeker", clientType: "Bot", clientName: "PagePeeker"},
	{find: "Pinterestbot", clientType: "Bot", clientName: "Pinterestbot"},
	{find: "Seekport", clientType: "Bot", clientName: "Seekport"},
	{find: "SeoSiteCheckup", clientType: "Bot", clientName: "SeoSiteCheckup"},
	{find: "Sitebulb", clientType: "Bot", clientName: "Sitebulb"},
	{find: "SiteScoreBot", clientType: "Bot", clientName: "SiteScoreBot"},
	{find: "SMTBot", clientType: "Bot", clientName: "SMTBot"},
	{find: "Yeti", clientType: "Bot", clientName: "Yeti"},
	{find: "YisouSpider", clientType: "Bot", clientName: "YisouSpider"},
	{find: "FBSV", clientType: "App", clientName: "Facebook"}, // iOS
	{find: "FBAV", clientType: "App", clientName: "Facebook"}, // Android
	{find: "GSA/", clientType: "App", clientName: "GoogleSearch"},
	{find: "Instagram", clientType: "App", clientName: "Instagram"},
	{find: "LinkedInApp", clientType: "App", clientName: "LinkedIn"},
	{find: "Pinterest", clientType: "App", clientName: "Pinterest"},
	{find: "Snapchat", clientType: "App", clientName: "Snapchat"},
	{find: "MicroMessenger", clientType: "App", clientName: "WeChat"},
	{find: "ADG/", clientType: "Browser", clientName: "AOLDesktop"},
	{find: "Silk", clientType: "Browser", clientName: "Silk"},
	{find: "FxiOS", clientType: "Browser", clientName: "Firefox"},
	{find: "Klarna", clientType: "Browser", clientName: "Firefox"},
	{find: "Firefox", clientType: "Browser", clientName: "Firefox"},
	{find: "EdgA/", clientType: "Browser", clientName: "Edge"},
	{find: "EdgiOS/", clientType: "Browser", clientName: "Edge"},
	{find: "EdgW/", clientType: "Browser", clientName: "Edge"},
	{find: "Edg/", clientType: "Browser", clientName: "Edge"},
	{find: "Edge/", clientType: "Browser", clientName: "Edge"},
	{find: "MSIE", clientType: "Browser", clientName: "InternetExplorer"},
	{find: "Trident", clientType: "Browser", clientName: "InternetExplorer"},
	{find: "OPR/", clientType: "Browser", clientName: "Opera"},
	{find: "OPT/", clientType: "Browser", clientName: "Opera"},
	{find: "DuckDuckGo", clientType: "Browser", clientName: "DuckDuckGo"},
	{find: "SamsungBrowser", clientType: "Browser", clientName: "SamsungBrowser"},
	{find: "CriOS", clientType: "Browser", clientName: "Chrome"},
	{find: "Chrome", clientType: "Browser", clientName: "Chrome"},
	{find: "Safari", clientType: "Browser", clientName: "Safari"},
}

// match indicates the appropriate field(s) for the supplied find text.
type match struct {
	find            string
	deviceType      string
	operatingSystem string
	clientType      string
	clientName      string
}

// UserAgent provides basic information about the user, extracted from an HTTP User-Agent request header.
// For more on User-Agent strings, see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
type UserAgent struct {
	// Header is the provided User-Agent request header
	Header string `json:"header,omitempty"`

	// Fields contains parsed/cleaned segments of the User-Agent request header, used for analysis
	Fields []string `json:"fields,omitempty"`

	// ClientType indicates the application category (App, Bot, Browser, or Other)
	ClientType string `json:"clientType,omitempty"`

	// ClientName indicates the application name (Chrome, Googlebot, Edge, etc.)
	ClientName string `json:"clientName,omitempty"`

	// ClientVersion indicates the version of the application, if provided
	ClientVersion string `json:"clientVersion,omitempty"`

	// DeviceType indicates the general device category (Desktop, Mobile, Tablet, Other)
	DeviceType string `json:"deviceType,omitempty"`

	// OSName indicates the operating system running on the device (Android, Linux, iOS, macOS, Windows, etc.)
	OSName string `json:"osName,omitempty"`

	// OSVersion indicates the operating system version, if available
	OSVersion string `json:"osVersion,omitempty"`

	// URL indicates the URL provided, typically for information about a bot/crawler.
	URL string `json:"url,omitempty"`
}

// String supports the Stringer interface, providing an abbreviated user agent string.
func (ua UserAgent) String() string {
	s := ua.ClientType
	if ua.ClientName != "" {
		s = s + " " + ua.ClientName
	}
	if ua.ClientVersion != "" {
		s = s + " " + ua.ClientVersion
	}
	if ua.DeviceType != "" {
		s = s + " " + ua.DeviceType
	}
	if ua.OSName != "" {
		s = s + " " + ua.OSName
	}
	if ua.OSVersion != "" {
		s = s + " " + ua.OSVersion
	}
	if ua.URL != "" {
		s = s + " " + ua.URL
	}
	return s
}

// Parse extracts client, device, and operating system information from the User-Agent request header provided,
// returning a UserAgent. Note that the URL and versions will be empty if not provided. Other fields, however,
// will be set to "Other" if the relevant information is not provided, or if the determination is inconclusive.
func Parse(userAgent string) UserAgent {
	ua := UserAgent{Header: unquote(userAgent)}
	ua.Fields = parseFields(ua.Header)
	cleaned := strings.Join(ua.Fields, " ")
	// A URL always indicates a Bot
	if strings.Contains(ua.Header, "://") {
		ua.ClientType = "Bot"
		ua.URL = botURL(ua.Fields)
	}
	// Pattern matchers must be processed in order, and first match wins for the provided field(s)
	for _, p := range patterns {
		if strings.Contains(cleaned, p.find) {
			if p.deviceType != "" && ua.DeviceType == "" {
				ua.DeviceType = p.deviceType
			}
			if p.operatingSystem != "" && ua.OSName == "" {
				ua.OSName = p.operatingSystem
			}
			if p.clientType != "" && ua.ClientType == "" {
				ua.ClientType = p.clientType
			}
			if p.clientName != "" && ua.ClientName == "" {
				ua.ClientName = p.clientName
				ua.ClientVersion = clientVersion(ua.Fields, p.find)
			}
		}
		// Skip any remaining patterns if the UserAgent is complete
		if ua.DeviceType != "" && ua.OSName != "" && ua.ClientType != "" && ua.ClientName != "" {
			break
		}
	}
	// Post-processing: supply default values and update version numbers as appropriate
	if ua.OSName == "" {
		ua.OSName = "Other"
	} else {
		ua.OSVersion = osVersion(ua.Fields, ua.OSName)
	}
	if ua.DeviceType == "" {
		ua.DeviceType = "Desktop"
	}
	if ua.ClientName == "" {
		if ua.OSName == "iOS" || ua.OSName == "iPadOS" || ua.OSName == "macOS" {
			if ua.ClientType == "" {
				ua.ClientType = "Browser"
			}
			ua.ClientName = "Safari"
		} else if ua.OSName == "Android" {
			if ua.ClientType == "" {
				ua.ClientType = "Browser"
			}
			ua.ClientName = "Chrome"
		} else {
			ua.ClientName = "Other"
		}
	}
	if ua.ClientName == "Safari" {
		ver := version(ua.Fields) // uses Version/99.9.9 for clientVersion
		if ver != "" {
			ua.ClientVersion = ver
		}
	} else if ua.ClientName == "InternetExplorer" {
		ver := releaseVersion(ua.Fields)
		if ver != "" {
			ua.ClientVersion = ver
		}
	}
	if ua.ClientType == "" {
		ua.ClientType = "Other"
	}
	return ua
}

// unquote strips single and double quotes from the provided User-Agent string.
// Sometimes, the User-Agent string arrives unnecessarily quoted, as one would indicate a literal string in code.
func unquote(ua string) string {
	quotes := []rune("'\"")
	rs := make([]rune, 0, len(ua))
	for _, r := range ua {
		if r != quotes[0] && r != quotes[1] {
			rs = append(rs, r)
		}
	}
	return string(rs)
}

// parseFields removes common, meaningless text and parses the User-Agent string into individual segments for analysis.
func parseFields(ua string) []string {
	// Remove meaningless text
	s := strings.ReplaceAll(ua, "Mozilla/5.0", "")
	s = strings.ReplaceAll(s, "Safari/537.36", "")
	s = strings.ReplaceAll(s, "KHTML", "")
	s = strings.ReplaceAll(s, "like Gecko", "")
	s = strings.ReplaceAll(s, "compatible", "")
	s = strings.ReplaceAll(s, " like Mac OS X", "")
	s = strings.ReplaceAll(s, "CPU ", "")
	s = strings.ReplaceAll(s, "Intel ", "")
	s = strings.ReplaceAll(s, "Mac OS X", "OS")
	s = strings.ReplaceAll(s, "Windows NT", "Windows")
	s = strings.ReplaceAll(s, "WOW64", "")
	s = strings.ReplaceAll(s, "Win64", "")
	s = strings.ReplaceAll(s, "x86_64", "")
	s = strings.ReplaceAll(s, "x64", "")
	s = strings.ReplaceAll(s, "aarch64", "")
	s = strings.ReplaceAll(s, "(", " ")
	s = strings.ReplaceAll(s, ")", " ")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, ";", " ")
	fields := strings.Fields(s)
	ss := make([]string, 0, len(fields))
	for _, f := range fields {
		if f != "," && !strings.HasPrefix(f, "AppleWebKit") {
			ss = append(ss, f)
		}
	}
	return ss
}

// botURL returns a URL, if present in the User-Agent string
func botURL(fields []string) string {
	for _, f := range fields {
		if strings.Contains(f, "://") {
			i := strings.Index(f, "http") // the URL may start with a '+'
			if i == -1 {
				i = 0
			}
			return f[i:]
		}
	}
	return ""
}

// version returns the major.minor version number, indicated by "Version" in the User-Agent string
// The Safari browser uses Version to indicate its version number.
func version(fields []string) string {
	for _, f := range fields {
		if strings.HasPrefix(f, "Version/") && len(f) > 8 {
			i := strings.Index(f, "/") + 1
			// capture major.minor version, ignoring patch details
			segments := strings.Split(f[i:], ".")
			if len(segments) == 1 {
				return segments[0]
			} else if len(segments) >= 2 {
				return strings.Join(segments[:2], ".")
			}
		}
	}
	return ""
}

// releaseVersion returns the version number, indicated by "rv" in the User-Agent string.
// Release Version is used by Firefox and Internet Explorer (Trident)
func releaseVersion(fields []string) string {
	for _, f := range fields {
		if strings.HasPrefix(f, "rv:") && len(f) > 3 {
			i := strings.Index(f, ":") + 1
			return f[i:]
		}
	}
	return ""
}

// clientVersion splits a clientName/clientVersion field on the slash, returning a major.minor version,
// ignoring patch details. Most clients use this syntax (e.g. Chrome/100.0.4896.75).
func clientVersion(fields []string, clientName string) string {
	for _, f := range fields {
		if strings.Contains(f, clientName) {
			_, ver, found := strings.Cut(f, "/")
			if found {
				// capture major.minor version, ignoring patch details
				segments := strings.Split(ver, ".")
				if len(segments) == 1 {
					return segments[0]
				} else if len(segments) >= 2 {
					return strings.Join(segments[:2], ".")
				}
			}
		}
	}
	return ""
}

// osVersion returns an operating system version, if available. It's usually space-separated after the operating
// system name in the User-Agent string. Therefore, it's often the next field after the operating system name.
func osVersion(fields []string, osName string) string {
	if osName == "" {
		return ""
	}
	if osName == "iOS" || osName == "iPadOS" || osName == "macOS" {
		osName = "OS"
	}
	for i, f := range fields {
		if f == osName && i+1 < len(fields) {
			return majorMinorVersion(fields[i+1])
		}
	}
	return ""
}

// majorMinorVersion returns a cleaned and trimmed major.minor version number from the provided text.
// If the text contains any non-numeric characters, then an empty string is returned.
// Note that Apple operating system version numbers use an underscore instead of a period to separate segments.
func majorMinorVersion(ver string) string {
	if !isDigits(ver) {
		return ""
	}
	var segments []string
	if strings.Contains(ver, "_") {
		segments = strings.Split(ver, "_")
	} else {
		segments = strings.Split(ver, ".")
	}
	if len(segments) >= 2 {
		return strings.Join(segments[:2], ".")
	}
	if len(segments) == 1 {
		return segments[0] // trimmed trailing separator
	}
	return ver // original string
}

// isDigits returns true if the supplied text contains only valid numeric digits in a major.minor.patch version
func isDigits(ver string) bool {
	if ver == "" {
		return false
	}
	valid := []rune("0123456789._")
	for _, r := range ver {
		if !validRune(valid, r) {
			return false
		}
	}
	return true
}

// validRune checks to see if the supplied rune appears in a set of valid runes.
func validRune(runes []rune, r rune) bool {
	for _, i := range runes {
		if r == i {
			return true
		}
	}
	return false
}
