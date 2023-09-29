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

    // admin struct {
    //     password string
    // }
}


func setupViper(cfg *Config) {
    viper.SetConfigName("config")
    viper.SetConfigType("toml")

    viper.AddConfigPath("/etc/cutlink/")
    viper.AddConfigPath("$HOME/.config/cutlink/")
    viper.AddConfigPath(".")

    err := viper.ReadInConfig()
    if err != nil {
        panic(err.Error())
    }

    err = viper.Unmarshal(cfg)
    if err != nil {
        panic(err.Error())
    }
}
