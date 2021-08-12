OLD=$(cat ./version)
NEW=$(expr $OLD + 1)

echo "$NEW" > version
