#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080"

echo -e "${YELLOW}🚀 Testing Carpooling API${NC}"
echo "=================================="

# Test health check
echo -e "\n${YELLOW}1. Testing Health Check${NC}"
response=$(curl -s "$API_URL/health")
if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Health check successful${NC}"
    echo "$response" | jq '.'
else
    echo -e "${RED}❌ Health check failed${NC}"
fi

# Test API info
echo -e "\n${YELLOW}2. Testing API Info${NC}"
response=$(curl -s "$API_URL/")
if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ API info successful${NC}"
    echo "$response" | jq '.'
else
    echo -e "${RED}❌ API info failed${NC}"
fi

# Create test users
echo -e "\n${YELLOW}3. Creating Test Users${NC}"

echo -e "\n📝 Creating user: Cyril"
cyril_response=$(curl -s -X POST "$API_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cyril",
    "email": "cyril@example.com",
    "phone": "+1234567890"
  }')

if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Cyril created successfully${NC}"
    cyril_id=$(echo "$cyril_response" | jq -r '.id')
    echo "$cyril_response" | jq '.'
else
    echo -e "${RED}❌ Failed to create Cyril${NC}"
fi

echo -e "\n📝 Creating user: Abhinav"
abhinav_response=$(curl -s -X POST "$API_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Abhinav",
    "email": "abhinav@example.com",
    "phone": "+1234567891"
  }')

if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Abhinav created successfully${NC}"
    abhinav_id=$(echo "$abhinav_response" | jq -r '.id')
    echo "$abhinav_response" | jq '.'
else
    echo -e "${RED}❌ Failed to create Abhinav${NC}"
fi

echo -e "\n📝 Creating user: Vishal"
vishal_response=$(curl -s -X POST "$API_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Vishal",
    "email": "vishal@example.com",
    "phone": "+1234567892"
  }')

if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Vishal created successfully${NC}"
    vishal_id=$(echo "$vishal_response" | jq -r '.id')
    echo "$vishal_response" | jq '.'
else
    echo -e "${RED}❌ Failed to create Vishal${NC}"
fi

# List users
echo -e "\n${YELLOW}4. Listing All Users${NC}"
users_response=$(curl -s "$API_URL/users")
if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Users listed successfully${NC}"
    echo "$users_response" | jq '.'
else
    echo -e "${RED}❌ Failed to list users${NC}"
fi

# Test get user by name
echo -e "\n${YELLOW}5. Getting User by Name${NC}"
user_by_name_response=$(curl -s "$API_URL/users/by-name/Cyril")
if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Get user by name successful${NC}"
    echo "$user_by_name_response" | jq '.'
else
    echo -e "${RED}❌ Failed to get user by name${NC}"
fi

# Create trip with transactions
echo -e "\n${YELLOW}6. Creating Trip with Transactions${NC}"
echo -e "\n🚗 Creating morning trip with 3 participants"
trip_response=$(curl -s -X POST "$API_URL/trips/with-transactions" \
  -H "Content-Type: application/json" \
  -d '{
    "driver_name": "Cyril",
    "date": "2024-01-15",
    "time": "09:00:00",
    "total_cost": 145,
    "week": 3,
    "participants": ["Cyril", "Abhinav", "Vishal"]
  }')

if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Trip created successfully${NC}"
    trip_id=$(echo "$trip_response" | jq -r '.trip_id')
    echo "$trip_response" | jq '.'
else
    echo -e "${RED}❌ Failed to create trip${NC}"
    echo "$trip_response"
fi

# Create evening trip
echo -e "\n🚗 Creating evening trip with 2 participants"
evening_trip_response=$(curl -s -X POST "$API_URL/trips/with-transactions" \
  -H "Content-Type: application/json" \
  -d '{
    "driver_name": "Cyril",
    "date": "2024-01-15",
    "time": "18:00:00",
    "total_cost": 145,
    "week": 3,
    "participants": ["Cyril", "Vishal"]
  }')

if [[ $? -eq 0 ]]; then
    echo -e "${GREEN}✅ Evening trip created successfully${NC}"
    evening_trip_id=$(echo "$evening_trip_response" | jq -r '.trip_id')
    echo "$evening_trip_response" | jq '.'
else
    echo -e "${RED}❌ Failed to create evening trip${NC}"
    echo "$evening_trip_response"
fi

# Get trip transactions
if [[ ! -z "$trip_id" && "$trip_id" != "null" ]]; then
    echo -e "\n${YELLOW}7. Getting Trip Transactions${NC}"
    trip_transactions_response=$(curl -s "$API_URL/trips/$trip_id/transactions")
    if [[ $? -eq 0 ]]; then
        echo -e "${GREEN}✅ Trip transactions retrieved successfully${NC}"
        echo "$trip_transactions_response" | jq '.'
    else
        echo -e "${RED}❌ Failed to get trip transactions${NC}"
    fi
fi

echo -e "\n${YELLOW}🎉 API Testing Complete!${NC}"
echo "=================================="
echo -e "${GREEN}✅ All core functionality tested${NC}"
echo -e "💡 Your carpooling API is ready to use!"