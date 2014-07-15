i18n
====

Middleware i18n is for [Macaron](https://github.com/Unknwon/macaron) app Internationalization and Localization.

[API Reference](https://gowalker.org/github.com/macaron-contrib/i18n)

## Usage

```go
import (
	"github.com/Unknwon/macaron"
	"github.com/macaron-contrib/i18n"
)

func main() {
  	m := macaron.Classic()
  	m.Use(i18n.I18n(i18n.LocaleOptions{
		Langs:    []string{"en-US", "zh-CN"},
		Names:    []string{"English", "简体中文"},
	}))
	
	m.Get("/", func(locale i18n.Locale) string {
		return "current language is" + locale.Lang
	})

	m.Run()
}
```

## Options

`i18n.I18n` comes with a variety of configuration options:

```go
// ...
m.Use(i18n.I18n(i18n.LocaleOptions{
	Directory:	"conf/locale", // Directory to load locale files.
	Langs:		[]string{"en-US", "zh-CN"}, // Langauges that will be supported, order is meaningful.
	Names:		[]string{"English", "简体中文"}, // Human friendly names corresponding to Langs list.
	Format:		"locale_%s.ini", // Locale file naming style.
	Parameter:	"lang", // Name of language parameter name in URL.
	Redirect:	false, // Redirect when user uses get parameter to specify language.
	TmplName:	"i18n", // Name that maps into template variable.
}))
// ...
```

## Loading Locale Files

By default, locale files should be put in `conf/locale`:

```
conf/
  |
  |__ locale/
  		 |
  		 |__ locale_en-US.ini
  		 |
   		 |__ locale_zh-CN.ini
```

## Specification

See [Unknwon/i18n](https://github.com/Unknwon/i18n) for specification of translation.

By using this middleware, you can use `{{.i18n.Tr "tranlate me"}}` to do template translation. See [gogs.io](https://github.com/gogits/gogsweb) as an example.

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.