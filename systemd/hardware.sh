#!/bin/sh

# Disable red led
echo 0 > /sys/devices/platform/leds/leds/PWR/brightness

# Set display brightness to low
gpio -g pwm 19 30
gpio -g mode 19 pwm