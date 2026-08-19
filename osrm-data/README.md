# OSRM map data

The `osrm` service in [docker-compose.yml](../docker-compose.yml) serves routes from a
preprocessed `.osrm` dataset in this directory. Preprocessing is a one-time step done with
the same `osrm-backend` image, not something docker-compose runs for you.

1. Download an OSM extract for the region you need, e.g. Thailand from
   [Geofabrik](https://download.geofabrik.de/asia/thailand.html):

   ```sh
   curl -L -o osrm-data/map.osm.pbf https://download.geofabrik.de/asia/thailand-latest.osm.pbf
   ```

2. Extract, partition, and customize it (MLD algorithm, matching the `osrm-routed` command
   in docker-compose.yml):

   ```sh
   docker run --rm -v "$PWD/osrm-data:/data" ghcr.io/project-osrm/osrm-backend osrm-extract -p /opt/car.lua /data/map.osm.pbf
   docker run --rm -v "$PWD/osrm-data:/data" ghcr.io/project-osrm/osrm-backend osrm-partition /data/map.osrm
   docker run --rm -v "$PWD/osrm-data:/data" ghcr.io/project-osrm/osrm-backend osrm-customize /data/map.osrm
   ```

   Use `/opt/bicycle.lua` or `/opt/foot.lua` instead of `/opt/car.lua` to match
   `OSRM_PROFILE=cycling`/`walking`.

3. Start the service and point the backend at it:

   ```sh
   docker compose up -d osrm
   ```

   ```
   OSRM_BASE_URL="http://localhost:5001"
   ```

Files in this directory (`*.osm.pbf`, `*.osrm*`) are gitignored - each environment builds
its own dataset.
