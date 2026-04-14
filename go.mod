module github.com/theopenlane/riverboat

go 1.25.9

replace github.com/theopenlane/riverboat/trustcenter => ./trustcenter

replace github.com/theopenlane/corejobs => ../corejobs/

replace github.com/theopenlane/go-client => ../go-client/

require (
	github.com/99designs/gqlgen v0.17.89
	github.com/gertd/go-pluralize v0.2.1
	github.com/go-chi/chi/v5 v5.2.5
	github.com/gocarina/gocsv v0.0.0-20240520201108-78e41c74b4b1
	github.com/gqlgo/gqlgenc v0.35.1
	github.com/invopop/jsonschema v0.13.0
	github.com/invopop/yaml v0.3.1
	github.com/jackc/pgx/v5 v5.9.1
	github.com/knadh/koanf/parsers/yaml v1.1.0
	github.com/knadh/koanf/providers/env/v2 v2.0.0
	github.com/knadh/koanf/providers/file v1.2.1
	github.com/knadh/koanf/providers/posflag v1.0.1
	github.com/knadh/koanf/v2 v2.3.4
	github.com/mcuadros/go-defaults v1.2.0
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/prometheus/client_golang v1.23.2
	github.com/riverqueue/river v0.33.0
	github.com/riverqueue/river/riverdriver/riverpgxv5 v0.33.0
	github.com/riverqueue/river/rivertype v0.33.0
	github.com/riverqueue/rivercontrib/otelriver v0.7.0
	github.com/rs/zerolog v1.35.0
	github.com/samber/lo v1.53.0
	github.com/slack-go/slack v0.21.1
	github.com/spf13/cobra v1.10.2
	github.com/stoewer/go-strcase v1.3.1
	github.com/stretchr/testify v1.11.1
	github.com/theopenlane/core/common v1.0.20
	github.com/theopenlane/dbx v0.1.3
	github.com/theopenlane/emailtemplates v0.3.7
	github.com/theopenlane/go-client v0.9.4
	github.com/theopenlane/httpsling v0.3.0
	github.com/theopenlane/iam v0.27.5
	github.com/theopenlane/newman v0.3.0
	github.com/theopenlane/riverboat/trustcenter v0.0.0-00010101000000-000000000000
	github.com/theopenlane/utils v0.7.0
	go.opentelemetry.io/otel v1.43.0
	go.opentelemetry.io/otel/exporters/prometheus v0.61.0
	go.opentelemetry.io/otel/sdk/metric v1.43.0
)

require (
	github.com/Yamashou/gqlgenc v0.33.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/cloudflare/cloudflare-go/v6 v6.9.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-pdf/fpdf v0.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/go-webauthn/webauthn v0.16.4 // indirect
	github.com/go-webauthn/x v0.2.3 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/google/go-tpm v0.9.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/hhrutter/lzw v1.0.0 // indirect
	github.com/hhrutter/pkcs7 v0.2.2 // indirect
	github.com/hhrutter/tiff v1.0.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/mailru/easyjson v0.9.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.21 // indirect
	github.com/mattn/go-runewidth v0.0.23 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/pdfcpu/pdfcpu v0.11.1 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/otlptranslator v1.0.0 // indirect
	github.com/prometheus/procfs v0.20.0 // indirect
	github.com/redis/go-redis/v9 v9.18.0 // indirect
	github.com/resend/resend-go/v3 v3.3.0 // indirect
	github.com/riverqueue/river/riverdriver v0.33.0 // indirect
	github.com/riverqueue/river/rivershared v0.33.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sosodev/duration v1.4.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/theopenlane/core v1.16.8 // indirect
	github.com/theopenlane/corejobs v0.1.21 // indirect
	github.com/theopenlane/echox v0.3.0 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/tinylib/msgp v1.6.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vektah/gqlparser/v2 v2.5.32 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/goleak v1.3.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/image v0.38.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hhrutter/pkcs7 => github.com/hhrutter/pkcs7 v0.2.0
