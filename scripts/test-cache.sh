#!/bin/bash

echo "=== Testing Maps Service Cache ==="

# 1. ลบ cache เก่า
echo "1. Clearing old cache..."
podman exec -it podman-compose_redis_1 redis-cli FLUSHALL

# 2. ตรวจสอบว่าไม่มี cache
echo "2. Checking cache is empty..."
podman exec -it podman-compose_redis_1 redis-cli KEYS "*"

# 3. ทำ API call แรก (จะเรียก HERE API)
echo "3. Making first API call (should call HERE API)..."
curl -X GET "http://localhost:8082/api/v1/maps/search/query?q=วัดพระแก้ว+กรุงเทพมหานคร&lat=13.7563&lng=100.5018&limit=5&lang=th" \
  -H "Content-Type: application/json" | jq '.meta.fromCache'

# 4. ตรวจสอบ cache ที่สร้างขึ้น
echo "4. Checking cache created..."
podman exec -it podman-compose_redis_1 redis-cli KEYS "maps:*"

# 5. ทำ API call ที่สอง (ควรใช้ cache)
echo "5. Making second API call (should use cache)..."
curl -X GET "http://localhost:8082/api/v1/maps/search/query?q=วัดพระแก้ว+กรุงเทพมหานคร&lat=13.7563&lng=100.5018&limit=5&lang=th" \
  -H "Content-Type: application/json" | jq '.meta.fromCache'

# 6. ดู TTL ของ cache
echo "6. Checking cache TTL..."
podman exec -it podman-compose_redis_1 redis-cli TTL "maps:search:วัดพระแก้ว กรุงเทพมหานคร:13.7563:100.5018:5:th"

# 7. ดูข้อมูลใน cache
echo "7. Cache data sample:"
podman exec -it podman-compose_redis_1 redis-cli GET "maps:search:วัดพระแก้ว กรุงเทพมหานคร:13.7563:100.5018:5:th" | jq '.cachedAt, .expiresAt, .results | length'

echo "=== Cache Test Complete ==="
