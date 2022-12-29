#!/bin/bash

for i in {1..10}
do
  home=mage-$i

  rm -rf $home/data
done

