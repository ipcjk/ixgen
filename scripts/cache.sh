#!/usr/bin/env bash
mkdir -p cache
cd cache
curl https://www.peeringdb.com/api/fac > fac
curl https://www.peeringdb.com/api/ix > ix
curl https://www.peeringdb.com/api/ixfac > ixfac
curl https://www.peeringdb.com/api/ixlan > ixlan
curl https://www.peeringdb.com/api/ixpfx > ixpfx
curl https://www.peeringdb.com/api/net > net
curl https://www.peeringdb.com/api/netfac > netfac
curl https://www.peeringdb.com/api/netixlan > netixlan
curl https://www.peeringdb.com/api/org > org
curl https://www.peeringdb.com/api/poc > poc