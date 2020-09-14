package app

type Config struct {
	Log     Log     `yaml:"log"`
	URLRepo URLRepo `yaml:"urlrepo"`
	API     API     `yaml:"api"`
}

type Log struct {
	File string `yaml:"file"`
}

type API struct {
	HTTPPort string `yaml:"httpport"`
}

type URLRepo struct {
	Path string `yaml:"path"`
}
