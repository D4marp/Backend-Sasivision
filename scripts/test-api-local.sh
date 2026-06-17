#!/usr/bin/env bash
# Full local API smoke test — jalankan setelah server di http://localhost:8080
set -euo pipefail

BASE="${API_BASE_URL:-http://127.0.0.1:8080}"
PASS=0
FAIL=0

check() {
  local name="$1"
  local expect="$2"
  local method="$3"
  local path="$4"
  local body="${5:-}"
  local auth="${6:-}"

  local args=(-sS -w "\n%{http_code}" -X "$method" "$BASE$path" -H "Content-Type: application/json")
  if [[ -n "$auth" ]]; then
    args+=(-H "Authorization: Bearer $auth")
  fi
  if [[ -n "$body" ]]; then
    args+=(-d "$body")
  fi

  local out
  out="$(curl "${args[@]}")"
  local code="${out##*$'\n'}"
  local resp="${out%$'\n'*}"

  if [[ "$code" == "$expect" ]]; then
    echo "✅ $name — HTTP $code"
    PASS=$((PASS + 1))
  else
    echo "❌ $name — expected HTTP $expect, got $code"
    echo "   $resp"
    FAIL=$((FAIL + 1))
  fi
}

echo "=== SasiVision API local test ==="
echo "Base: $BASE"
echo ""

check "Health" 200 GET "/health"
check "Quiz categories" 200 GET "/api/quiz/categories"
check "Quiz questions Post-Test" 200 GET "/api/quiz/questions/Post-Test"
check "Content markers" 200 GET "/api/content/markers"
check "Content videos" 200 GET "/api/content/videos"
check "Feature switches" 200 GET "/api/features/switches"
check "Sign-in invalid" 401 POST "/api/auth/sign-in" '{"email":"x@y.com","password":"wrongpass1"}'

echo ""
echo "--- Auth: demo user ---"
SIGNIN=$(curl -sS -X POST "$BASE/api/auth/sign-in" \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@sasivision.com","password":"Sasivision123"}')
echo "$SIGNIN" | grep -q '"status":"success"' && echo "✅ Demo sign-in" && PASS=$((PASS+1)) || { echo "❌ Demo sign-in"; echo "$SIGNIN"; FAIL=$((FAIL+1)); }
DEMO_TOKEN=$(echo "$SIGNIN" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('token',''))" 2>/dev/null || true)

echo ""
echo "--- Auth: editor ---"
EDITOR=$(curl -sS -X POST "$BASE/api/auth/sign-in" \
  -H "Content-Type: application/json" \
  -d '{"email":"editor@sasivision.com","password":"Sasivision123"}')
echo "$EDITOR" | grep -q '"role":"editor"' && echo "✅ Editor sign-in" && PASS=$((PASS+1)) || { echo "❌ Editor sign-in"; echo "$EDITOR"; FAIL=$((FAIL+1)); }
EDITOR_TOKEN=$(echo "$EDITOR" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('token',''))" 2>/dev/null || true)

echo ""
echo "--- Auth: admin ---"
ADMIN=$(curl -sS -X POST "$BASE/api/auth/sign-in" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@sasivision.com","password":"Sasivision123"}')
echo "$ADMIN" | grep -q '"role":"admin"' && echo "✅ Admin sign-in" && PASS=$((PASS+1)) || { echo "❌ Admin sign-in"; echo "$ADMIN"; FAIL=$((FAIL+1)); }
ADMIN_TOKEN=$(echo "$ADMIN" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',{}).get('token',''))" 2>/dev/null || true)

if [[ -n "$DEMO_TOKEN" ]]; then
  check "Quiz history (auth)" 200 GET "/api/quiz/history/demo@sasivision.com" "" "$DEMO_TOKEN"
fi

if [[ -n "$EDITOR_TOKEN" ]]; then
  check "Admin quiz categories (editor)" 200 GET "/api/admin/quiz/categories" "" "$EDITOR_TOKEN"
  check "Admin quiz questions (editor)" 200 GET "/api/admin/quiz/questions?category_id=1" "" "$EDITOR_TOKEN"
fi

if [[ -n "$ADMIN_TOKEN" ]]; then
  check "Admin users (admin)" 200 GET "/api/admin/users" "" "$ADMIN_TOKEN"
  check "Admin analytics (admin)" 200 GET "/api/admin/analytics" "" "$ADMIN_TOKEN"
  check "Admin stats (admin)" 200 GET "/api/admin/stats" "" "$ADMIN_TOKEN"
fi

check "Protected without token" 401 GET "/api/admin/users"
check "404 route" 404 GET "/api/not-exists"

echo ""
echo "=== Quiz submit flow ==="
if [[ -n "$DEMO_TOKEN" ]]; then
  QDATA=$(curl -sS "$BASE/api/quiz/questions/Post-Test")
  QCOUNT=$(echo "$QDATA" | python3 -c "import sys,json; d=json.load(sys.stdin); print(len(d.get('data',[])))" 2>/dev/null || echo 0)
  if [[ "$QCOUNT" -ge 5 ]]; then
    echo "✅ Post-Test has $QCOUNT questions" && PASS=$((PASS+1))
  else
    echo "❌ Post-Test questions = $QCOUNT (expected >= 5)" && FAIL=$((FAIL+1))
  fi
fi

echo ""
echo "=============================="
echo "PASS: $PASS  FAIL: $FAIL"
if [[ "$FAIL" -gt 0 ]]; then
  exit 1
fi
echo "All local API checks passed."
