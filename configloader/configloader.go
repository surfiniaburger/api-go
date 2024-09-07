// configloader/configloader.go
package configloader

import (
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    PublicHost             string
    Port                   string
    DBURL                  string
    JWTSecret              string
    JWTExpirationInSeconds int64
}

// Make this function public by renaming it to InitConfig
func InitConfig() Config {
    godotenv.Load()

    // Determine which DB URL to use: external for local, internal for Render
    dbURL := getEnv("DB_URL_EXTERNAL", "")
    if os.Getenv("RENDER") == "true" {
        dbURL = getEnv("DB_URL_INTERNAL", dbURL)
    }

    return Config{
        PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
        Port:                   getEnv("PORT", "8080"),
        DBURL:                  dbURL,
        JWTSecret:              getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
        JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
    }
}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
    if value, ok := os.LookupEnv(key); ok {
        i, err := strconv.ParseInt(value, 10, 64)
        if err != nil {
            return fallback
        }
        return i
    }
    return fallback
}
