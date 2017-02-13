package config

import (
    "log"
    "os"
    "encoding/json"
)

type mainConfig struct {
    Server ServerConfig `json:"server"`
}

type ServerConfig struct {
    ServerHost string `json:"host"`
    ServerPort int `json:"port"`
}

var MainConfigPath = "config.json"
var MainConfig = new(mainConfig)

func init() {
    log.Println("[CRAWLER] Init config")
    mc := new(mainConfig)
    if err := mc.read(MainConfigPath); err != nil {
        panic(err)
    }
    MainConfig = mc
}

//Read config file from input path
func (mc *mainConfig) read(path string) error {
    log.Println("[Config] Reading config: ", path)
    c, err := os.Open(path)
    if err != nil {
        log.Printf("[Config] Error: Could not read a config file %s.\n%s\n", path, err)
        return err
    }

    decoder := json.NewDecoder(c)
    if err := decoder.Decode(mc); err != nil {
        log.Printf("[Config] Error: Could not decode a config file %s.\n%s\n", path, err)
        return  err
    }

    return nil
}
