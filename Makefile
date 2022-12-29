install:
	go install -ldflags "-X github.com/furya-official/mgtool/config/generate.ConfigTemplatesDir=$(CURDIR)/config/templates"