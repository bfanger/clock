
setup:
	go install github.com/bokwoon95/wgo@latest
	brew install pkg-config sdl2 sdl2_image sdl2_ttf

install:
	go install ./cmd/clock
	go install ./cmd/garbage-truck
	go install ./cmd/school-schedule
	go install ./cmd/reminders
	go install ./cmd/weather

dev:
	wgo run ./cmd/clock

test:
	go test ./pkg/tween ./internal/app ./internal/schedule

