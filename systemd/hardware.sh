#!/bin/sh

# Disable red led
echo 0 > /sys/class/leds/led1/brightness

# Set diplay brighness to low
gpio -g pwm 19 30
gpio -g mode 19 pwm