package config

import (
	"errors"
	"os"

	"example.com/kode-notes/repository"
	"example.com/kode-notes/server"
)

type Config struct {
	Spellcheck bool
}

func ParseArgs(sv_cfg *server.Config,
	db_cfg *repository.Config, app_cfg *Config) error {

	args := os.Args[1:]
	if len(args) < 1 {
		return nil
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-s", "--server-host":
			if i+1 > len(args)-1 {
				return errors.New(
					"no host specified")
			}

			sv_cfg.Host = args[i+1]
			i++

		case "-p", "--server-port":
			if i+1 > len(args)-1 {
				return errors.New(
					"no port specified")
			}

			sv_cfg.Port = args[i+1]
			i++

		case "-D", "--db-host":
			if i+1 > len(args)-1 {
				return errors.New(
					"No database password" +
						" specified")
			}

			db_cfg.Host = args[i+1]
			i++

		case "-P", "--db-passwd":
			if i+1 > len(args)-1 {
				return errors.New(
					"No database password" +
						" specified")
			}

			db_cfg.Password = args[i+1]
			i++

		case "-C", "--no-spellcheck":
			app_cfg.Spellcheck = false

		default:
			return errors.New("Unknown argument - " +
				args[i])
		}
	}

	return nil
}
