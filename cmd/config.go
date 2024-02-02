package main


import (
    "github.com/spf13/viper"
)


type Config struct {
    Server struct {
        Port    int
        Addr    string
    }

    Database struct {
        MainDB      string
        SessionsDB  string
    }

    Management struct {
        NoSignup        bool
        RateLimitMax    int
    }

    Admin struct {
        Route   string
    }

    Tls struct {
        Cert    string
        Key     string
    }
}


func setupViper(cfg *Config, configPath string) {
    viper.SetConfigName("config")
    viper.SetConfigType("toml")

    viper.AddConfigPath("/etc/cutlink/")
    viper.AddConfigPath(".")

    if configPath != "" {
        viper.SetConfigFile(configPath)
    }

    err := viper.ReadInConfig()
    if err != nil {
        panic(err.Error())
    }

    err = viper.Unmarshal(cfg)
    if err != nil {
        panic(err.Error())
    }
}
