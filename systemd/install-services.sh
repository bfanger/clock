#!/bin/sh

# Start the graphical clock interface
go install -v ~/clock/cmd/clock/clock.go
sudo ln -s ~/clock/systemd/clock.service /etc/systemd/system/clock.service
sudo systemctl enable clock
# Dim the hyperpixel screen (needs apt get wiringpi)
sudo ln -s ~/clock/systemd/hardware.service /etc/systemd/system/hardware.service
sudo systemctl enable hardware
# reminders (bedtime notifications)
go install -v ~/clock/cmd/reminders/reminders.go
sudo ln -s ~/clock/systemd/reminders.service /etc/systemd/system/reminders.service
sudo systemctl enable reminders
# school starting notifications from magister calendar
go install -v ~/clock/cmd/school-schedule/school-schedule.go
sudo ln -s ~/clock/systemd/school-schedule.service /etc/systemd/system/school-schedule.service
sudo systemctl enable school-schedule
# garbage-truck (garbage & recycle pickup notification)
go install -v ~/clock/cmd/garbage-truck/garbage-truck.go
sudo ln -s ~/clock/systemd/garbage-truck.service /etc/systemd/system/garbage-truck.service
sudo systemctl enable garbage-truck
# weather (freezing icon, car window might need defrosting, needs OPENWEATHERMAP_APPID in .env)
go install -v ~/clock/cmd/weather/weather.go
sudo ln -s ~/clock/systemd/weather.service /etc/systemd/system/weather.service
sudo systemctl enable weather
# sonos-volume (Show the volume of the sonos soundbar)
go install -v ~/clock/cmd/sonos-volume/sonos-volume.go
sudo ln -s ~/clock/systemd/sonos-volume.service /etc/systemd/system/sonos-volume.service
sudo systemctl enable sonos-volume
