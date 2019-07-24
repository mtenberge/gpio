// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/mtenberge/gpio"
	"github.com/mtenberge/gpio/spi/adc0832"
	"github.com/warthog618/config"
	"github.com/warthog618/config/blob"
	"github.com/warthog618/config/blob/decoder/json"
	"github.com/warthog618/config/blob/loader/file"
	"github.com/warthog618/config/dict"
	"github.com/warthog618/config/env"
	"github.com/warthog618/config/pflag"
)

// This example reads both channels from an ADC0832 connected to the RPI by four
// data lines - CSZ, CLK, DI, and DO. The default pin assignments are defined in
// loadConfig, but can be altered via configuration (env, flag or config file).
// The DI and DO may be tied to reduce the pin count by one, though I prefer to
// keep the two separate to remove the chance of accidental conflict.
// All pins other than DO are outputs so do not run this example on a board
// where those pins serve other purposes.
func main() {
	cfg := loadConfig()
	err := gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()
	tclk := cfg.MustGet("tclk").Duration()
	tset := cfg.MustGet("tset").Duration()
	if tset < tclk {
		tset = tclk
	}
	a := adc0832.New(
		tclk,
		tset,
		uint8(cfg.MustGet("clk").Uint()),
		uint8(cfg.MustGet("csz").Uint()),
		uint8(cfg.MustGet("di").Uint()),
		uint8(cfg.MustGet("do").Uint()))
	defer a.Close()
	ch0 := a.Read(0)
	ch1 := a.Read(1)
	fmt.Printf("ch0=0x%02x, ch1=0x%02x\n", ch0, ch1)
}

func loadConfig() *config.Config {
	defaultConfig := map[string]interface{}{
		"tclk": "2500ns",
		"tset": "2500ns", // should be at least tclk - enforced in main
		"clk":  gpio.GPIO6,
		"csz":  gpio.GPIO5,
		"di":   gpio.GPIO19,
		"do":   gpio.GPIO13,
	}
	def := dict.New(dict.WithMap(defaultConfig))
	shortFlags := map[byte]string{
		'c': "config-file",
	}
	fget, err := pflag.New(pflag.WithShortFlags(shortFlags))
	if err != nil {
		panic(err)
	}
	// environment next
	eget, err := env.New(env.WithEnvPrefix("ADC0832_"))
	if err != nil {
		panic(err)
	}
	// highest priority sources first - flags override environment
	sources := config.NewStack(fget, eget)
	cfg := config.NewConfig(config.Decorate(sources, config.WithDefault(def)))

	// config file may be specified via flag or env, so check for it
	// and if present add it with lower priority than flag and env.
	configFile, err := cfg.Get("config.file")
	jsondec := json.NewDecoder()
	if err == nil {
		// explicitly specified config file - must be there
		f, err := file.New(configFile.String())
		if err != nil {
			panic(err)
		}
		jget, err := blob.New(f, jsondec)
		if err != nil {
			panic(err)
		}
		sources.Append(jget)
	} else {
		// implicit and optional default config file
		f, _ := file.New("adc0832.json")
		jget, err := blob.New(f, jsondec)
		if err == nil {
			sources.Append(jget)
		} else {
			if _, ok := err.(*os.PathError); !ok {
				panic(err)
			}
		}
	}
	cfg = cfg.GetConfig("", config.WithMust())
	return cfg
}
