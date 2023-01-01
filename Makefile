install:
	go install -ldflags "-X github.com/mage-labs/mgtool/config/generate.ConfigTemplatesDir=$(CURDIR)/config/templates"