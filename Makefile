
# Detect Apple Silicon and use static build.
# Fixes "This header is only meant to be used on x86 and x64 architecture"" error. 
# Related https://github.com/veandco/go-sdl2/issues/479
goarg := $(shell go env|grep GOARCH=)
ifeq ($(goarg),GOARCH='arm64')
  static := -tags static 
  export CGO_LDFLAGS="-Wl,-no_warn_duplicate_libraries"
endif



setup:
	go install ${static} github.com/bokwoon95/wgo@latest
	brew install pkg-config sdl2 sdl2_image sdl2_ttf

install:
	go install ${static} ./cmd/clock
	go install ./cmd/garbage-truck
	go install ./cmd/school-schedule
	go install ./cmd/reminders
	go install ./cmd/weather

dev:
	wgo run ${static} ./cmd/clock

test:
	go test ${static} ./pkg/tween ./internal/app ./internal/schedule

