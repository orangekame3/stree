#!/bin/bash

start_date="2023-12-01"
end_date="2024-01-16"

current_date="$start_date"

while [ "$current_date" != "$(gdate -I -d "$end_date + 1 day")" ]; do
  # Convert date to the required format: YYYY/MM/DD
  formatted_date=$(gdate -d "$current_date" +"%Y/%m/%d")
  
  # Create the directory structure
  mkdir -p "$formatted_date"
  
  # Create the dummy.txt file in the directory
  touch "$formatted_date/dummy.txt"
  
  # Move to the next date
  current_date=$(gdate -I -d "$current_date + 1 day")
done
