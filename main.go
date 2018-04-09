package main

import "flag"

func main() {
	var (
		config string
		secret string
		image  string
		suffix string
	)
	flag.StringVar(&config, "config", "env.yml", "Config file name of preseed")
	flag.StringVar(&config, "c", "env.yml", "Config file name of preseed")
	flag.StringVar(&secret, "secret", "env.secret.yml", "Secret config file name of preseed")
	flag.StringVar(&secret, "s", "env.secret.yml", "Secret config file name of preseed")
	flag.StringVar(&image, "image", "ubuntu-16.04.3-server-amd64", "Ubuntu image to customize")
	flag.StringVar(&image, "i", "ubuntu-16.04.3-server-amd64", "Ubuntu image to customize")
	flag.StringVar(&suffix, "suffix", "", "Suffix added to the output file")
	flag.Parse()

	buildPreseedConfig(config, secret)
	regenerateISO(image, suffix)
}
