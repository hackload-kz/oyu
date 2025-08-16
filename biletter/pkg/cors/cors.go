package cors

import (
	"biletter/internal/config"
	"github.com/rs/cors"
)

func GetCorsSettings(cfg *config.Config) *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:     cfg.Cors.AllowedOrigins,
		AllowedMethods:     cfg.Cors.AllowedMethods,
		AllowedHeaders:     cfg.Cors.AllowedHeaders,
		ExposedHeaders:     cfg.Cors.ExposedHeaders,
		AllowCredentials:   cfg.Cors.AllowCredentials,
		OptionsPassthrough: cfg.Cors.OptionsPassthrough,
		Debug:              cfg.Cors.Debug,
	})
	return c
}
