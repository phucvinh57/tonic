package docs

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TickLabVN/tonic/core/errors"
)

type ValidateFlag struct {
	// Fields
	EqCsField     bool `json:"eqcsfield,omitempty"`
	EqField       bool `json:"eqfield,omitempty"`
	FieldContains bool `json:"fieldcontains,omitempty"`
	GtCsField     bool `json:"gtcsfield,omitempty"`
	GtField       bool `json:"gtfield,omitempty"`
	GteCsField    bool `json:"gtecsfield,omitempty"`
	GteField      bool `json:"gtefield,omitempty"`
	LtCsField     bool `json:"ltcsfield,omitempty"`
	LtField       bool `json:"ltfield,omitempty"`
	LteCsField    bool `json:"ltecsfield,omitempty"`
	LteField      bool `json:"ltefield,omitempty"`
	NeCsField     bool `json:"necsfield,omitempty"`
	NeField       bool `json:"nefield,omitempty"`

	// Network
	Cidr            bool `json:"cidr"`
	CidrV4          bool `json:"cidrv4"`
	CidrV6          bool `json:"cidrv6"`
	DataURI         bool `json:"datauri"`
	Fqdn            bool `json:"fqdn"`
	Hostname        bool `json:"hostname"`
	HostnamePort    bool `json:"hostname_port"`
	HostnameRfc1123 bool `json:"hostname_rfc1123"`
	IP              bool `json:"ip"`
	IP4Addr         bool `json:"ip4_addr"`
	IP6Addr         bool `json:"ip6_addr"`
	IPAddr          bool `json:"ip_addr"`
	IPv4            bool `json:"ipv4"`
	IPv6            bool `json:"ipv6"`
	MAC             bool `json:"mac"`
	TCP4Addr        bool `json:"tcp4_addr"`
	TCP6Addr        bool `json:"tcp6_addr"`
	TCPAddr         bool `json:"tcp_addr"`
	UDP4Addr        bool `json:"udp4_addr"`
	UDP6Addr        bool `json:"udp6_addr"`
	UDPAddr         bool `json:"udp_addr"`
	UnixAddr        bool `json:"unix_addr"`
	URI             bool `json:"uri"`
	URL             bool `json:"url"`
	HTTPURL         bool `json:"http_url"`
	URLEncoded      bool `json:"url_encoded"`
	URNRfc2141      bool `json:"urn_rfc2141"`

	// String
	Alpha           bool     `json:"alpha"`
	Alphanum        bool     `json:"alphanum"`
	AlphanumUnicode bool     `json:"alphanumunicode"`
	AlphaUnicode    bool     `json:"alphaunicode"`
	ASCII           bool     `json:"ascii"`
	Boolean         bool     `json:"boolean"`
	Contains        []string `json:"contains"`
	ContainsAny     []string `json:"containsany"`
	ContainsRune    []string `json:"containsrune"`
	EndsNotWith     []string `json:"endsnotwith"`
	EndsWith        []string `json:"endswith"`
	Excludes        []string `json:"excludes"`
	ExcludesAll     []string `json:"excludesall"`
	ExcludesRune    []string `json:"excludesrune"`
	Lowercase       bool     `json:"lowercase"`
	Multibyte       bool     `json:"multibyte"`
	Number          bool     `json:"number"`
	Numeric         bool     `json:"numeric"`
	PrintASCII      bool     `json:"printascii"`
	StartsNotWith   []string `json:"startsnotwith"`
	StartsWith      []string `json:"startswith"`
	Uppercase       bool     `json:"uppercase"`

	// Formats
	Base64                     bool `json:"base64"`
	Base64URL                  bool `json:"base64url"`
	Base64RawURL               bool `json:"base64rawurl"`
	BIC                        bool `json:"bic"`
	BCP47Language              bool `json:"bcp47_language_tag"`
	BTCAddress                 bool `json:"btc_addr"`
	BTCAddressBech32           bool `json:"btc_addr_bech32"`
	CreditCard                 bool `json:"credit_card"`
	MongoDB                    bool `json:"mongodb"`
	MongoDBConnectionString    bool `json:"mongodb_connection_string"`
	Cron                       bool `json:"cron"`
	SpiceDBObject              bool `json:"spicedb"`
	Datetime                   bool `json:"datetime"`
	E164                       bool `json:"e164"`
	EIN                        bool `json:"ein"`
	Email                      bool `json:"email"`
	EthereumAddress            bool `json:"eth_addr"`
	Hexadecimal                bool `json:"hexadecimal"`
	Hexcolor                   bool `json:"hexcolor"`
	HSL                        bool `json:"hsl"`
	HSLA                       bool `json:"hsla"`
	HTML                       bool `json:"html"`
	HTMLEncoded                bool `json:"html_encoded"`
	ISBN                       bool `json:"isbn"`
	ISBN10                     bool `json:"isbn10"`
	ISBN13                     bool `json:"isbn13"`
	ISSN                       bool `json:"issn"`
	ISO31661Alpha2             bool `json:"iso3166_1_alpha2"`
	ISO31661Alpha3             bool `json:"iso3166_1_alpha3"`
	ISO31661AlphaNumeric       bool `json:"iso3166_1_alpha_numeric"`
	ISO31662                   bool `json:"iso3166_2"`
	ISO4217                    bool `json:"iso4217"`
	JSON                       bool `json:"json"`
	JWT                        bool `json:"jwt"`
	Latitude                   bool `json:"latitude"`
	Longitude                  bool `json:"longitude"`
	LuhnChecksum               bool `json:"luhn_checksum"`
	PostcodeISO3166Alpha2      bool `json:"postcode_iso3166_alpha2"`
	PostcodeISO3166Alpha2Field bool `json:"postcode_iso3166_alpha2_field"`
	RGB                        bool `json:"rgb"`
	RGBA                       bool `json:"rgba"`
	SSN                        bool `json:"ssn"`
	Timezone                   bool `json:"timezone"`
	UUID                       bool `json:"uuid"`
	UUID3                      bool `json:"uuid3"`
	UUID3RFC4122               bool `json:"uuid3_rfc4122"`
	UUID4                      bool `json:"uuid4"`
	UUID4RFC4122               bool `json:"uuid4_rfc4122"`
	UUID5                      bool `json:"uuid5"`
	UUID5RFC4122               bool `json:"uuid5_rfc4122"`
	UUIDRFC4122                bool `json:"uuid_rfc4122"`
	MD4                        bool `json:"md4"`
	MD5                        bool `json:"md5"`
	SHA256                     bool `json:"sha256"`
	SHA384                     bool `json:"sha384"`
	SHA512                     bool `json:"sha512"`
	RIPEMD128                  bool `json:"ripemd128"`
	RIPEMD160                  bool `json:"ripemd160"`
	TIGER128                   bool `json:"tiger128"`
	TIGER160                   bool `json:"tiger160"`
	TIGER192                   bool `json:"tiger192"`
	Semver                     bool `json:"semver"`
	ULID                       bool `json:"ulid"`
	CVE                        bool `json:"cve"`

	// Comparison
	Eq           []string `json:"eq,omitempty"`
	EqIgnoreCase []string `json:"eq_ignore_case,omitempty"`
	Ne           []string `json:"ne,omitempty"`
	NeIgnoreCase []string `json:"ne_ignore_case,omitempty"`
	Gt           string   `json:"gt,omitempty"`
	Gte          string   `json:"gte,omitempty"`
	Lt           string   `json:"lt,omitempty"`
	Lte          string   `json:"lte,omitempty"`

	// Others
	Dir                bool            `json:"dir"`
	DirPath            bool            `json:"dirpath"`
	File               bool            `json:"file"`
	FilePath           bool            `json:"filepath"`
	Image              bool            `json:"image"`
	IsDefault          bool            `json:"isdefault"`
	Len                string          `json:"len"`
	Max                string          `json:"max"`
	Min                string          `json:"min"`
	OneOf              []string        `json:"oneof,omitempty"`
	Required           bool            `json:"required"`
	RequiredIf         *ComparisonFlag `json:"required_if,omitempty"`
	RequiredUnless     *ComparisonFlag `json:"required_unless,omitempty"`
	RequiredWith       []string        `json:"required_with,omitempty"`
	RequiredWithAll    []string        `json:"required_with_all,omitempty"`
	RequiredWithout    []string        `json:"required_without,omitempty"`
	RequiredWithoutAll []string        `json:"required_without_all,omitempty"`
	ExcludedIf         *ComparisonFlag `json:"excluded_if,omitempty"`
	ExcludedUnless     *ComparisonFlag `json:"excluded_unless,omitempty"`
	ExcludedWith       []string        `json:"excluded_with,omitempty"`
	ExcludedWithAll    []string        `json:"excluded_with_all,omitempty"`
	ExcludedWithout    []string        `json:"excluded_without,omitempty"`
	ExcludedWithoutAll []string        `json:"excluded_without_all,omitempty"`
	Unique             bool            `json:"unique"`
}

type ComparisonFlag struct {
	FieldName string   `json:"fieldname"`
	Values    []string `json:"values"`
}

func (v *ValidateFlag) GetFormat() string {
	if v.Alpha {
		return "alpha"
	}
	if v.Alphanum {
		return "alphanum"
	}
	if v.AlphanumUnicode {
		return "alphanum_unicode"
	}
	if v.AlphaUnicode {
		return "alpha_unicode"
	}
	if v.ASCII {
		return "ascii"
	}
	if v.Boolean {
		return "boolean"
	}
	if v.Number {
		return "number"
	}
	if v.Numeric {
		return "numeric"
	}
	if v.PrintASCII {
		return "print_ascii"
	}
	if v.Uppercase {
		return "uppercase"
	}
	if v.Lowercase {
		return "lowercase"
	}
	if v.Multibyte {
		return "multibyte"
	}

	if v.Base64 {
		return "base64"
	}
	if v.Base64URL {
		return "base64url"
	}
	if v.Base64RawURL {
		return "base64rawurl"
	}
	if v.BIC {
		return "bic"
	}
	if v.BCP47Language {
		return "bcp47_language_tag"
	}
	if v.BTCAddress {
		return "btc_addr"
	}
	if v.BTCAddressBech32 {
		return "btc_addr_bech32"
	}
	if v.CreditCard {
		return "credit_card"
	}
	if v.MongoDB {
		return "mongodb"
	}
	if v.MongoDBConnectionString {
		return "mongodb_connection_string"
	}
	if v.Cron {
		return "cron"
	}
	if v.SpiceDBObject {
		return "spicedb"
	}
	if v.Datetime {
		return "datetime"
	}
	if v.E164 {
		return "e164"
	}
	if v.EIN {
		return "ein"
	}
	if v.Email {
		return "email"
	}
	if v.EthereumAddress {
		return "eth_addr"
	}
	if v.Hexadecimal {
		return "hexadecimal"
	}
	if v.Hexcolor {
		return "hexcolor"
	}
	if v.HSL {
		return "hsl"
	}
	if v.HSLA {
		return "hsla"
	}
	if v.HTML {
		return "html"
	}
	if v.HTMLEncoded {
		return "html_encoded"
	}
	if v.ISBN {
		return "isbn"
	}
	if v.ISBN10 {
		return "isbn10"
	}
	if v.ISBN13 {
		return "isbn13"
	}
	if v.ISSN {
		return "issn"
	}
	if v.ISO31661Alpha2 {
		return "iso3166_1_alpha2"
	}
	if v.ISO31661Alpha3 {
		return "iso3166_1_alpha3"
	}
	if v.ISO31661AlphaNumeric {
		return "iso3166_1_alpha_numeric"
	}
	if v.ISO31662 {
		return "iso3166_2"
	}
	if v.ISO4217 {
		return "iso4217"
	}
	if v.JSON {
		return "json"
	}
	if v.JWT {
		return "jwt"
	}
	if v.Latitude {
		return "latitude"
	}
	if v.Longitude {
		return "longitude"
	}
	if v.LuhnChecksum {
		return "luhn_checksum"
	}
	if v.PostcodeISO3166Alpha2 {
		return "postcode_iso3166_alpha2"
	}
	if v.PostcodeISO3166Alpha2Field {
		return "postcode_iso3166_alpha2_field"
	}
	if v.RGB {
		return "rgb"
	}
	if v.RGBA {
		return "rgba"
	}
	if v.SSN {
		return "ssn"
	}
	if v.Timezone {
		return "timezone"
	}
	if v.UUID {
		return "uuid"
	}
	if v.UUID3 {
		return "uuid3"
	}
	if v.UUID3RFC4122 {
		return "uuid3_rfc4122"
	}
	if v.UUID4 {
		return "uuid4"
	}
	if v.UUID4RFC4122 {
		return "uuid4_rfc4122"
	}
	if v.UUID5 {
		return "uuid5"
	}
	if v.UUID5RFC4122 {
		return "uuid5_rfc4122"
	}
	if v.UUIDRFC4122 {
		return "uuid_rfc4122"
	}
	if v.MD4 {
		return "md4"
	}
	if v.MD5 {
		return "md5"
	}
	if v.SHA256 {
		return "sha256"
	}
	if v.SHA384 {
		return "sha384"
	}
	if v.SHA512 {
		return "sha512"
	}
	if v.RIPEMD128 {
		return "ripemd128"
	}
	if v.RIPEMD160 {
		return "ripemd160"
	}
	if v.TIGER128 {
		return "tiger128"
	}
	if v.TIGER160 {
		return "tiger160"
	}
	if v.TIGER192 {
		return "tiger192"
	}
	if v.Semver {
		return "semver"
	}
	if v.ULID {
		return "ulid"
	}
	if v.CVE {
		return "cve"
	}

	// Network
	if v.Cidr {
		return "cidr"
	}
	if v.CidrV4 {
		return "cidrv4"
	}
	if v.CidrV6 {
		return "cidrv6"
	}
	if v.DataURI {
		return "data-uri"
	}
	if v.Fqdn {
		return "fqdn"
	}
	if v.Hostname {
		return "hostname"
	}
	if v.HostnamePort {
		return "hostname-port"
	}
	if v.HostnameRfc1123 {
		return "hostname-rfc1123"
	}
	if v.IP {
		return "ip"
	}
	if v.IP4Addr {
		return "ipv4"
	}
	if v.IP6Addr {
		return "ipv6"
	}
	if v.IPAddr {
		return "ip"
	}
	if v.IPv4 {
		return "ipv4"
	}
	if v.IPv6 {
		return "ipv6"
	}
	if v.MAC {
		return "mac"
	}
	if v.TCP4Addr {
		return "tcp4-addr"
	}
	if v.TCP6Addr {
		return "tcp6-addr"
	}
	if v.TCPAddr {
		return "tcp-addr"
	}
	if v.UDP4Addr {
		return "udp4-addr"
	}
	if v.UDP6Addr {
		return "udp6-addr"
	}
	if v.UDPAddr {
		return "udp-addr"
	}
	if v.UnixAddr {
		return "unix-addr"
	}
	if v.URI {
		return "uri"
	}
	if v.URL {
		return "url"
	}
	if v.HTTPURL {
		return "http-url"
	}
	if v.URLEncoded {
		return "url-encoded"
	}
	if v.URNRfc2141 {
		return "urn-rfc2141"
	}

	return ""
}

func (v *ValidateFlag) GetPattern() string {
	if v.Alpha {
		return "^[a-zA-Z]*$"
	}
	if v.Alphanum {
		return "^[a-zA-Z0-9]*$"
	}
	if v.AlphanumUnicode {
		return "^[\\p{L}0-9]*$"
	}
	if v.AlphaUnicode {
		return "^[\\p{L}]*$"
	}
	if v.ASCII {
		return "^[\\x00-\\x7F]*$"
	}
	if len(v.Contains) > 0 {
		return ".*" + strings.Join(v.Contains, " ") + ".*"
	}
	if len(v.ContainsAny) > 0 {
		return ".*[" + strings.Join(v.ContainsAny, " ") + "].*"
	}
	if len(v.ContainsRune) > 0 {
		return ".*" + strings.Join(v.ContainsRune, "") + ".*"
	}
	if len(v.EndsWith) > 0 {
		return ".*" + strings.Join(v.EndsWith, "|") + "$"
	}
	if v.Lowercase {
		return "^[a-z]*$"
	}
	if v.Multibyte {
		return "^[\\p{L}0-9]*$"
	}
	if v.Number {
		return "^[0-9]*$"
	}
	if v.Numeric {
		return "^[0-9]*$"
	}
	if v.PrintASCII {
		return "^[\\x20-\\x7E]*$"
	}
	if len(v.StartsWith) > 0 {
		return "^" + strings.Join(v.StartsWith, "|") + ".*"
	}
	if v.Uppercase {
		return "^[A-Z]*$"
	}
	return ""
}

func (v *ValidateFlag) GetDescription() string {
	if len(v.Contains) > 0 {
		return "Contain '" + strings.Join(v.Contains, " ") + "'"
	}
	if len(v.ContainsAny) > 0 {
		return "Contain any chars in '" + strings.Join(v.ContainsAny, " ") + "'"
	}
	if len(v.ContainsRune) > 0 {
		return "Contain '" + strings.Join(v.ContainsRune, "") + "'"
	}
	if len(v.EndsNotWith) > 0 {
		return "Not end with '" + strings.Join(v.EndsNotWith, " ") + "'"
	}
	if len(v.EndsWith) > 0 {
		return "End with '" + strings.Join(v.EndsWith, " ") + "'"
	}
	if len(v.Excludes) > 0 {
		return "Not contain '" + strings.Join(v.Excludes, " ") + "'"
	}
	if len(v.ExcludesAll) > 0 {
		return "Not contain any chars in '" + strings.Join(v.ExcludesAll, " ") + "'"
	}
	if len(v.ExcludesRune) > 0 {
		return "Not contain '" + strings.Join(v.ExcludesRune, "") + "'"
	}
	if len(v.StartsNotWith) > 0 {
		return "Not start with '" + strings.Join(v.StartsNotWith, " ") + "'"
	}
	if len(v.StartsWith) > 0 {
		return "Start with '" + strings.Join(v.StartsWith, " ") + "'"
	}

	formatDescriptions := map[string]string{
		"base64":                        "Base64 String",
		"base64url":                     "Base64URL String",
		"base64rawurl":                  "Base64RawURL String",
		"bic":                           "Business Identifier Code (ISO 9362)",
		"bcp47_language_tag":            "Language tag (BCP 47)",
		"btc_addr":                      "Bitcoin Address",
		"btc_addr_bech32":               "Bitcoin Bech32 Address (segwit)",
		"credit_card":                   "Credit Card Number",
		"mongodb":                       "MongoDB ObjectID",
		"mongodb_connection_string":     "MongoDB Connection String",
		"cron":                          "Cron",
		"spicedb":                       "SpiceDb ObjectID/Permission/Type",
		"datetime":                      "Datetime",
		"e164":                          "e164 formatted phone number",
		"email":                         "E-mail String",
		"eth_addr":                      "Ethereum Address",
		"hexadecimal":                   "Hexadecimal String",
		"hexcolor":                      "Hexcolor String",
		"hsl":                           "HSL String",
		"hsla":                          "HSLA String",
		"html":                          "HTML Tags",
		"html_encoded":                  "HTML Encoded",
		"isbn":                          "International Standard Book Number",
		"isbn10":                        "International Standard Book Number 10",
		"isbn13":                        "International Standard Book Number 13",
		"issn":                          "International Standard Serial Number",
		"iso3166_1_alpha2":              "Two-letter country code (ISO 3166-1 alpha-2)",
		"iso3166_1_alpha3":              "Three-letter country code (ISO 3166-1 alpha-3)",
		"iso3166_1_alpha_numeric":       "Numeric country code (ISO 3166-1 numeric)",
		"iso3166_2":                     "Country subdivision code (ISO 3166-2)",
		"iso4217":                       "Currency code (ISO 4217)",
		"json":                          "JSON",
		"jwt":                           "JSON Web Token (JWT)",
		"latitude":                      "Latitude",
		"longitude":                     "Longitude",
		"luhn_checksum":                 "Luhn Algorithm Checksum (for strings and (u)int)",
		"postcode_iso3166_alpha2":       "Postcode",
		"postcode_iso3166_alpha2_field": "Postcode",
		"rgb":                           "RGB String",
		"rgba":                          "RGBA String",
		"ssn":                           "Social Security Number SSN",
		"timezone":                      "Timezone",
		"uuid":                          "Universally Unique Identifier UUID",
		"uuid3":                         "Universally Unique Identifier UUID v3",
		"uuid3_rfc4122":                 "Universally Unique Identifier UUID v3 RFC4122",
		"uuid4":                         "Universally Unique Identifier UUID v4",
		"uuid4_rfc4122":                 "Universally Unique Identifier UUID v4 RFC4122",
		"uuid5":                         "Universally Unique Identifier UUID v5",
		"uuid5_rfc4122":                 "Universally Unique Identifier UUID v5 RFC4122",
		"uuid_rfc4122":                  "Universally Unique Identifier UUID RFC4122",
		"md4":                           "MD4 hash",
		"md5":                           "MD5 hash",
		"sha256":                        "SHA256 hash",
		"sha384":                        "SHA384 hash",
		"sha512":                        "SHA512 hash",
		"ripemd128":                     "RIPEMD-128 hash",
		"ripemd160":                     "RIPEMD-160 hash",
		"tiger128":                      "TIGER128 hash",
		"tiger160":                      "TIGER160 hash",
		"tiger192":                      "TIGER192 hash",
		"semver":                        "Semantic Versioning 2.0.0",
		"ulid":                          "Universally Unique Lexicographically Sortable Identifier ULID",
		"cve":                           "Common Vulnerabilities and Exposures Identifier (CVE id)",
		"cidr":                          "Classless Inter-Domain Routing CIDR",
		"cidrv4":                        "Classless Inter-Domain Routing CIDRv4",
		"cidrv6":                        "Classless Inter-Domain Routing CIDRv6",
		"data-uri":                      "Data URI",
		"fqdn":                          "Fully Qualified Domain Name (FQDN)",
		"hostname":                      "Hostname RFC 952",
		"hostname-port":                 "HostPort",
		"hostname-rfc1123":              "Hostname RFC 1123",
		"ip":                            "Internet Protocol Address IP",
		"ipv4":                          "Internet Protocol Address IPv4",
		"ipv6":                          "Internet Protocol Address IPv6",
		"mac":                           "Media Access Control Address MAC",
		"tcp4-addr":                     "Transmission Control Protocol Address TCPv4",
		"tcp6-addr":                     "Transmission Control Protocol Address TCPv6",
		"tcp-addr":                      "Transmission Control Protocol Address TCP",
		"udp4-addr":                     "User Datagram Protocol Address UDPv4",
		"udp6-addr":                     "User Datagram Protocol Address UDPv6",
		"udp-addr":                      "User Datagram Protocol Address UDP",
		"unix-addr":                     "Unix domain socket end point Address",
		"uri":                           "URI String",
		"url":                           "URL String",
		"http-url":                      "HTTP URL String",
		"url-encoded":                   "URL Encoded",
		"urn-rfc2141":                   "Urn RFC 2141 String",
	}

	return formatDescriptions[v.GetFormat()]
}

func ParseValidateTag(tag string) (*ValidateFlag, error) {
	tag = strings.TrimSpace(tag)
	if len(tag) == 0 {
		return nil, nil
	}
	parts := strings.Split(tag, ",")
	rawTag := make(map[string]any)

	aliases := map[string]string{
		"alphanum_unicode": "alphanumunicode",
		"alpha_unicode":    "alphaunicode",
		"contains_any":     "containsany",
		"contains_rune":    "containsrune",
		"ends_not_with":    "endsnotwith",
		"ends_with":        "endswith",
		"excludes_all":     "excludesall",
		"excludes_rune":    "excludesrune",
		"print_ascii":      "printascii",
		"starts_not_with":  "startsnotwith",
		"starts_with":      "startswith",
	}

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
		kv := strings.Split(parts[i], "=")
		key := kv[0]
		if alias, ok := aliases[key]; ok {
			key = alias
		}
		if len(kv) > 1 {
			values := strings.Split(kv[1], " ")
			if (key == "gt" || key == "gte" || key == "lt" || key == "lte" || key == "len" || key == "min" || key == "max") && len(values) == 1 {
				rawTag[key] = values[0]
			} else {
				rawTag[key] = values
			}
		} else {
			rawTag[key] = true
		}
	}

	validateTag := &ValidateFlag{}
	b, err := json.Marshal(rawTag)
	if err != nil {
		return nil, fmt.Errorf("%w: marshalling error %v", errors.ErrParseValidateTag, err)
	}
	err = json.Unmarshal(b, validateTag)
	if err != nil {
		return nil, fmt.Errorf("%w: unmarshalling error %v", errors.ErrParseValidateTag, err)
	}
	return validateTag, nil
}
