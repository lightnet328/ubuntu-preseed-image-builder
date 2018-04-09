package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type Bool struct {
	IsNil   bool
	Boolean bool
}

func NewBool(s string) *Bool {
	return &Bool{
		IsNil:   s == "" || !(s == "true" || s == "false"),
		Boolean: s == "true",
	}
}

func (b *Bool) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var aux interface{}
	if err = unmarshal(&aux); err != nil {
		return err
	}
	switch raw := aux.(type) {
	case nil:
		*b = *NewBool("")
	case bool:
		*b = *NewBool(strconv.FormatBool(raw))
	case string:
		*b = *NewBool(raw)
	}
	return nil
}

func (b Bool) MarshalYAML() ([]byte, error) {
	return json.Marshal(b)
}

func (b Bool) String() string {
	if b.IsNil {
		return ""
	}
	if b.Boolean {
		return "true"
	}
	return "false"
}

type LocalizationType struct {
	Locale           string
	SupportedLocales string `yaml:"supported_locales"`
}

type KeyboardType struct {
	Layout string
	Model  string
}

type NetType struct {
	UseAutoconfig     *Bool `yaml:"use_autoconfig"`
	DisableDHCP       *Bool `yaml:"disable_dhcp"`
	Interface         string
	DisableAutoconfig *Bool  `yaml:"disable_autoconfig"`
	IPAddress         string `yaml:"ip_address"`
	Netmask           string
	Gateway           string
	NameServers       string `yaml:"name_servers"`
	Hostname          string
}

type HTTPType struct {
	Hostname string
}

type MirrorType struct {
	HTTP HTTPType `yaml:"http"`
}

type RootType struct {
	Password string
}

type UserType struct {
	Fullname          string
	Name              string
	Password          string
	AllowPasswordWeak bool `yaml:"allow_password_weak"`
}

type TimeType struct {
	Zone string
}

type PackageType struct {
	Additional      []string
	Upgrade         string
	LanguagePacks   string `yaml:"language_packs"`
	LanguageSupport *Bool  `yaml:"language_support"`
	Update          string
}

type Env struct {
	Localization LocalizationType
	Keyboard     KeyboardType
	Net          NetType
	Mirror       MirrorType
	Root         RootType
	User         UserType
	Time         TimeType
	Package      PackageType
}

func (env Env) ReadFile(path string) (Env, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return env, err
	}
	err = yaml.Unmarshal(file, &env)
	return env, err
}

func (dst Env) Merge(src ...Env) (Env, error) {
	var err error
	for _, s := range src {
		err = mergo.Merge(&dst, s, mergo.WithOverride)
		if err != nil {
			return dst, err
		}
	}
	return dst, err
}

func buildPreseedConfig(config, secret string) {
	var err error

	envDefault := Env{
		Localization: LocalizationType{
			Locale:           "en_US",
			SupportedLocales: "",
		},
		Keyboard: KeyboardType{
			Layout: "",
			Model:  "",
		},
		Net: NetType{
			UseAutoconfig:     NewBool(""),
			DisableDHCP:       NewBool(""),
			Interface:         "auto",
			DisableAutoconfig: NewBool(""),
			IPAddress:         "192.168.1.42",
			Netmask:           "255.255.255.0",
			Gateway:           "192.168.1.1",
			NameServers:       "192.168.1.1",
			Hostname:          "somehost",
		},
		Mirror: MirrorType{
			HTTP: HTTPType{
				Hostname: "archive.ubuntu.com",
			},
		},
		Root: RootType{
			Password: "r00tme",
		},
		User: UserType{
			Fullname:          "Ubuntu",
			Name:              "ubuntu",
			Password:          "insecure",
			AllowPasswordWeak: true,
		},
		Time: TimeType{
			Zone: "US/Eastern",
		},
		Package: PackageType{
			Additional:      nil,
			Upgrade:         "none",
			LanguagePacks:   "",
			LanguageSupport: NewBool(""),
			Update:          "none",
		},
	}

	var envData Env
	if envData, err = envData.ReadFile(config); err != nil {
		panic(err)
	}

	var envSecretData Env
	if envSecretData, err = envSecretData.ReadFile(secret); err != nil {
		panic(err)
	}

	env := Env{}
	if env, err = env.Merge(envDefault, envData, envSecretData); err != nil {
		panic(err)
	}

	out, err := os.Create("preseed.cfg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	t := template.New("preseed.cfg.tmpl").Funcs(template.FuncMap{
		"join": func(a []string) string { return strings.Join(a, " ") },
		"last": func(i int, s interface{}) bool { return i == reflect.ValueOf(s).Len()-1 },
	})
	t = template.Must(t.ParseFiles("preseed.cfg.tmpl"))
	if err = t.Execute(out, env); err != nil {
		panic(err)
	}
}
