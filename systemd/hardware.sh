#!/bin/sh

# Disable red led
echo 0 > /sys/devices/platform/leds/leds/PWR/brightness

# Set display brightness to low
# With help from https://gitlab.com/anthonydigirolamo/rpi-hardware-pwm
/home/bob/rpi-hardware-pwm/pwm 19 1000000 160000
