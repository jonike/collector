if ! getent passwd pganalyze > /dev/null; then
  adduser --system --quiet --home /var/lib/pganalyze-collector --no-create-home --shell /bin/bash pganalyze
fi

if ! getent group pganalyze > /dev/null; then
  addgroup --system --quiet pganalyze
  usermod -g pganalyze pganalyze
fi

mkdir -p /var/lib/pganalyze-collector
su -s /bin/sh pganalyze -c "test -O /var/lib/pganalyze-collector" || chown pganalyze /var/lib/pganalyze-collector

chown root:pganalyze /usr/bin/pganalyze-collector-helper
chmod 4750 /usr/bin/pganalyze-collector-helper

status pganalyze-collector | grep -q running
if [ $? -eq 0 ]; then
  restart -q pganalyze-collector
else
  start -q pganalyze-collector
fi
