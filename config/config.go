package config

import (
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func LoadConfig(path string) IAppConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load dotenv failed : %v", err)
	}

	_ = envMap

	return &config{
		app: &app{
			host: envMap["APP_HOST"],
			port: func() int {

				p, err := strconv.Atoi(envMap["APP_PORT"])

				if err != nil {
					log.Fatalf("load dotenv failed : %v", err)
				}
				return p
			}(),
			name:    envMap["APP_NAME"],
			version: envMap["APP_VERSION"],
			readTimeout: func() time.Duration {

				r, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])

				if err != nil {
					log.Fatalf("load read timeout failed : %v", err)
				}
				return time.Duration(int64(r) * int64(math.Pow10(9)))
			}(),
			writeTimeout: func() time.Duration {

				w, err := strconv.Atoi(envMap["APP_WRTIE_TIMEOUT"])

				if err != nil {
					log.Fatalf("load write timeout failed : %v", err)
				}
				return time.Duration(int64(w) * int64(math.Pow10(9)))
			}(),
			bodyLimit: func() int {
				b, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
				if err != nil {
					log.Fatalf("load bodylimit failed : %v", err)
				}
				return b
			}(),
			fileLimit: func() int {
				f, err := strconv.Atoi(envMap["APP_FILE_LIMIT"])
				if err != nil {
					log.Fatalf("load file limit failed : %v", err)
				}
				return f
			}(),
			gcpbucket: envMap["APP_GCP_BUCKET"],
		},
		db: &db{
			host: envMap["DB_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["DB_PORT"])
				if err != nil {
					log.Fatalf("load port failed : %v", err)
				}
				return p
			}(),
			protocal: envMap["DB_PROTOCOL"],
			username: envMap["DB_USERNAME"],
			password: envMap["DB_PASSWORD"],
			database: envMap["DB_DATABASE"],
			sslMode:  envMap["DB_SSL_MODE"],
			maxConnection: func() int {
				maxConnec, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
				if err != nil {
					log.Fatalf("load maxConnec failed : %v", err)
				}
				return maxConnec
			}(),
		},
		jwt: &jwt{
			adminKey:  envMap["JWT_ADMIN_KEY"],
			secretKey: envMap["JWT_SECRET_KEY"],
			apiKey:    envMap["JWT_API_KEY"],
			accessExpiresAt: func() int {
				t, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("load access expires at failed : %v", err)
				}
				return t
			}(),
			refreshExpireAt: func() int {
				t, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("load access expires at failed : %v", err)
				}
				return t
			}(),
		},
	}
}

type IConfig interface {
	App() IAppConfig
	Db() IDbConfig
	Jwt() IJwtConfig
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

type IAppConfig interface {
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int
	fileLimit    int
	gcpbucket    string
}

func (c *config) App() IAppConfig {
	return nil
}

type IDbConfig interface {
}

type db struct {
	host          string
	port          int
	protocal      string
	username      string
	password      string
	database      string
	sslMode       string
	maxConnection int
}

func (c *config) Db() IDbConfig {
	return nil
}

type IJwtConfig interface {
}

type jwt struct {
	adminKey        string
	secretKey       string
	apiKey          string
	accessExpiresAt int //second
	refreshExpireAt int //millisecond
}

func (c *config) Jwt() IJwtConfig {
	return nil
}
