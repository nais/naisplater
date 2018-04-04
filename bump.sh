OLD=$(cat ./version | cut -d'.' -f1)
NEW=$(expr $OLD + 1).0.0

echo "$NEW" > version
