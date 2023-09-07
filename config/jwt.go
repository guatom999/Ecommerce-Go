package config

type IJwtConfig interface {
	SecretKey() []byte
	AdminKey() []byte
	ApiKey() []byte
	AccessExpiresAt() int
	RefreshExpireAt() int
	SetAccessExpires(t int)
	SetRefreshExpire(t int)
}

type jwt struct {
	adminKey         string
	secretKey        string
	apiKey           string
	accessExpiresAt  int //second
	refreshExpiresAt int //millisecond
}

func (c *config) Jwt() IJwtConfig {
	return c.jwt
}

func (j *jwt) SecretKey() []byte      { return []byte(j.secretKey) }
func (j *jwt) AdminKey() []byte       { return []byte(j.adminKey) }
func (j *jwt) ApiKey() []byte         { return []byte(j.apiKey) }
func (j *jwt) AccessExpiresAt() int   { return j.accessExpiresAt }
func (j *jwt) RefreshExpireAt() int   { return j.refreshExpiresAt }
func (j *jwt) SetAccessExpires(t int) { j.accessExpiresAt = t }
func (j *jwt) SetRefreshExpire(t int) { j.refreshExpiresAt = t }
