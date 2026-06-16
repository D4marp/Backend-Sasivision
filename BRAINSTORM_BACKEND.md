# SasiVision Backend - Brainstorming & Future Enhancements

**Version:** 1.0  
**Status:** Planning Phase  
**Date Created:** 2026-06-10  

---

## 📝 Current Backend Status

**Implemented:**
- ✅ Project structure with Golang + Gin
- ✅ Docker setup with MySQL
- ✅ Base models (User, Quiz, Video, Marker, etc.)
- ✅ Handler stubs for all API endpoints
- ✅ Middleware (CORS, Rate Limiting, Auth)
- ✅ SQL migrations (all 13 tables)
- ✅ Configuration management (.env)

**In Development:**
- 🔄 Authentication service (bcrypt + JWT)
- 🔄 Database repositories
- 🔄 Business logic services

**Not Started:**
- ❌ Comprehensive testing
- ❌ Advanced caching (Redis)
- ❌ Real-time features (WebSockets)
- ❌ Analytics & monitoring
- ❌ Multi-language support

---

## 🧠 Brainstorming: Auto-Login Strategy

### Recommended: Option B (Token-Based)

**Flow:**
1. User logs in with email + password
2. Backend generates cryptographically secure token
3. Token stored in `user_sessions` table with expiry (30 days default)
4. Token returned in response + stored in app SharedPreferences
5. On app restart, app sends token to `/api/verify-token`
6. Backend validates token against DB + checks expiry
7. If valid → user auto-logged in to HomeScreen
8. If invalid → clear session, show SignInScreen

**Advantages:**
- ✅ Secure (no password stored locally)
- ✅ Server-side control (can revoke anytime)
- ✅ Scalable
- ✅ Supports multiple sessions
- ✅ Standard JWT approach

**Implementation Priority:** P0 (Week 1)

### Enhancement Ideas (Future)

#### 1. Biometric Multi-Factor Auth
```golang
// Add to users table
fingerprint_enrolled bool
face_recognition_enrolled bool

// New endpoints
POST /api/auth/biometric/enroll
POST /api/auth/biometric/verify
```

#### 2. Session Ping / Keepalive
```golang
// Background job
func PingActiveSessions() {
    // Every hour: check sessions still valid
    // Update last_activity timestamp
    // If stale > 30 days: soft delete
}

// Client-side: Send ping when user active
POST /api/session/ping
```

#### 3. Device-Specific Sessions
```golang
// Add to user_sessions table
device_id VARCHAR(255)
device_name VARCHAR(255)
device_os VARCHAR(100)
last_ip VARCHAR(45)

// Only allow login from registered devices
// OR require additional verification on new device
```

#### 4. Smart Timeout
```golang
// Configuration
SHORT_TIMEOUT = 15 minutes  // For sensitive operations
MEDIUM_TIMEOUT = 1 hour     // General operations
LONG_TIMEOUT = 7 days       // Remember me

// Endpoints
POST /api/auth/token/refresh (if token expired but within refresh window)
POST /api/auth/re-authenticate (force user to enter password)
```

#### 5. Activity-Based Session
```golang
// Reset expiry on each activity
PUT /api/session/touch

// Benefits:
// - Active users can stay logged in indefinitely
// - Inactive sessions auto-expire for security
```

---

## 🚀 Phase 2: Advanced Features

### A. Analytics & Monitoring

```golang
// Track user behavior
POST /api/analytics/event
{
  "event_type": "quiz_started|quiz_completed|video_watched",
  "metadata": {...}
}

// Admin dashboard
GET /api/admin/analytics
{
  "total_users": 5000,
  "active_users_today": 234,
  "avg_session_duration": "15m",
  "quiz_completion_rate": 0.72,
  "most_popular_quiz": "Basics",
  "most_watched_video": "History of Sasirangan"
}
```

### B. Caching Layer (Redis)

**Why:**
- Feature switches checked on every API call
- Quiz categories rarely change
- Video list static (mostly)

**Implementation:**
```golang
// Cache feature switches for 1 hour
// Cache quiz categories for 24 hours
// Cache video list for 6 hours

// Using Redis
const (
  CACHE_FEATURE_SWITCHES = "cache:feature_switches"
  CACHE_QUIZ_CATEGORIES = "cache:quiz_categories"
  CACHE_VIDEOS = "cache:videos"
)

func GetFeatureSwitches(ctx context.Context) []FeatureSwitch {
  // 1. Try Redis cache
  // 2. If miss, query DB
  // 3. Store in Redis with TTL
  // 4. Return
}
```

### C. Search & Filtering

```golang
// Full-text search on quiz questions
GET /api/quiz/search?q=batik+pattern
Response: Matching questions across all categories

// Filter videos by source
GET /api/videos?source=Educational%20Team

// Filter markers by language
GET /api/markers?language=en
```

### D. Leaderboard System

```golang
// Track top scorers
POST /api/leaderboard/submit

// Get leaderboard
GET /api/leaderboard?category=Basics&limit=10
Response:
[
  {
    "rank": 1,
    "email": "user1@email.com",
    "score": 95,
    "attempts": 3
  },
  ...
]
```

### E. Achievements/Badges

```golang
// Track achievements
POST /api/achievements/unlock
{
  "achievement_type": "perfect_score|speedrun|quiz_master",
  "quiz_category": "Basics"
}

// Get user achievements
GET /api/user/achievements
```

---

## 🔐 Phase 3: Security Enhancements

### A. Rate Limiting Per User

```golang
// Current: IP-based (100 req/min per IP)
// Future: User-based

const (
  RATE_LIMIT_LOGIN_ATTEMPT = 5 // per 15 minutes
  RATE_LIMIT_API_GENERAL = 500  // per hour
  RATE_LIMIT_QUIZ_SUBMIT = 20   // per day
)
```

### B. HTTPS Enforcement

```golang
// All API calls HTTPS only
// Strict-Transport-Security header
// HSTS preload list
```

### C. SQL Injection Prevention

```golang
// Already using parameterized queries
// But audit code for raw SQL usage
// Implement query builder pattern
```

### D. XSS Prevention

```golang
// Sanitize user input
// Content-Security-Policy headers
// HTML entity encoding
```

### E. CSRF Protection

```golang
// Add CSRF tokens for state-changing operations
POST /api/auth/logout
  X-CSRF-Token: [token]
```

---

## 📊 Phase 4: Performance Optimization

### A. Database Optimization

```sql
-- Add indexes for common queries
CREATE INDEX idx_user_quizzes 
  ON quiz_attempts(user_id, finish_date DESC);

-- Partition large tables by date
ALTER TABLE quiz_attempts 
  PARTITION BY RANGE(YEAR(finish_date)) (
    PARTITION p2025 VALUES LESS THAN (2025),
    PARTITION p2026 VALUES LESS THAN (2026),
    PARTITION pmax VALUES LESS THAN MAXVALUE
  );
```

### B. Query Optimization

```golang
// N+1 problem: Load user WITH their quiz attempts
// Before: O(n) queries
// After: 1 query with JOIN

SELECT qa.*, qc.name 
FROM quiz_attempts qa
LEFT JOIN quiz_categories qc ON qa.category_id = qc.id
WHERE qa.user_id = ?
ORDER BY qa.finish_date DESC;
```

### C. Connection Pooling

```golang
// Current: 25 max, 5 idle
// Monitor connection usage
// Adjust based on load testing
db.SetMaxOpenConns(50)
db.SetMaxIdleConns(10)
```

### D. API Response Caching

```golang
// Cache GET responses at HTTP level
Cache-Control: public, max-age=3600

// Vary by query params
Vary: Accept-Encoding, Accept-Language
```

---

## 🧪 Phase 5: Testing Strategy

### Unit Tests
```golang
// Test each service method
✓ TestAuthService_SignIn_ValidCredentials
✓ TestAuthService_SignIn_InvalidCredentials
✓ TestQuizService_CalculateScore
✓ TestSessionService_VerifyToken_Valid
✓ TestSessionService_VerifyToken_Expired
```

### Integration Tests
```golang
// Test API endpoints with real DB
✓ Test POST /api/sign-in → success flow
✓ Test POST /api/user-quiz-attempts → submission + validation
✓ Test GET /api/user-quiz-attempts/{email} → retrieves all attempts
```

### Load Testing
```
Tool: Apache JMeter / k6

Scenario 1: 100 concurrent users
- Each makes 10 API calls
- Measure response time, error rate

Scenario 2: Quiz submission spike
- 50 quiz submissions simultaneously
- Verify data consistency

Target:
- P95 response time < 500ms
- Error rate < 0.1%
```

---

## 🌐 Phase 6: Scalability

### Horizontal Scaling

```golang
// Load balancer (Nginx)
// Multiple API instances behind LB
// Stateless API servers
// Shared database
// Redis for session cache
```

### Database Replication

```
Primary MySQL (write)
↓
Replica 1 (read)
Replica 2 (read)
Replica 3 (read backup)
```

### Microservices (Future)

```
Current: Monolithic Gin server

Future:
├── Auth Service
├── Quiz Service  
├── Content Service
├── Analytics Service
└── Admin Service

Communication: REST API / gRPC / Message Queue (RabbitMQ)
```

---

## 📱 Phase 7: Mobile App Integration

### Push Notifications

```golang
// When new content available
POST /api/notifications/send
{
  "target_users": "all|by_category|by_segment",
  "title": "New quiz available!",
  "body": "Try our new Advanced Patterns quiz"
}

// Client-side
firebase_cloud_messaging.listen()
```

### Offline Support

```golang
// API provides sync delta
GET /api/sync/delta?last_sync=2026-06-09T10:00:00Z
Response: Only changed resources since last sync

// Benefits:
// - Smaller payloads
// - Faster sync
// - Better offline UX
```

---

## 📋 Implementation Roadmap

### Week 1-2: Core Auth
- [ ] Implement bcrypt password hashing
- [ ] Generate & verify JWT tokens
- [ ] Create session records in DB
- [ ] Test auto-login flow

### Week 3-4: Quiz & Content APIs
- [ ] Implement quiz endpoints
- [ ] Quiz scoring algorithm
- [ ] Implement video/marker endpoints
- [ ] Feature switch logic

### Week 5-6: Testing & Optimization
- [ ] Unit tests (>70% coverage)
- [ ] Integration tests
- [ ] Load testing
- [ ] Performance tuning

### Week 7-8: Deployment
- [ ] Docker/K8s setup
- [ ] CI/CD pipeline
- [ ] Staging environment
- [ ] Production deployment

### Week 9+: Phase 2 Features
- [ ] Analytics setup
- [ ] Redis caching
- [ ] Search/filtering
- [ ] Leaderboard

---

## 🤝 Team Coordination Points

### Flutter ↔ Backend Sync

1. **API Contract:**
   - Flutter sends: `{"email": "...", "password": "..."}`
   - Backend responds: `{"status": "success", "email": "...", "token": "..."}`
   - Test both sides on same JSON structure

2. **Error Handling:**
   - Flutter expects error format: `{"status": "error", "message": "...", "code": "ERR_CODE"}`
   - Backend must consistently return this format

3. **Token Management:**
   - Flutter stores token in SharedPreferences
   - Backend validates token signature + expiry on each protected endpoint
   - Token refresh endpoint for extending sessions

4. **Data Validation:**
   - Backend validates all inputs (email format, password strength, etc.)
   - Flutter also validates before sending (fail fast)

---

## 📊 Success Metrics

**Phase 1 (Week 1-2):**
- ✅ 100% auth endpoints working
- ✅ Auto-login flow tested end-to-end
- ✅ 0 data loss on session table

**Phase 2 (Week 3-4):**
- ✅ All quiz endpoints working
- ✅ Score calculation correct (verified vs Android app)
- ✅ < 100ms response time for quiz questions

**Overall (Week 8):**
- ✅ P95 response time < 500ms
- ✅ Uptime > 99.5%
- ✅ Test coverage > 70%
- ✅ 0 security vulnerabilities (OWASP Top 10)

---

## 🔗 Related Documents

- [FLUTTER_PRD.md](../FLUTTER_PRD.md) - Flutter features & requirements
- [FLUTTER_ERD.md](../FLUTTER_ERD.md) - Database schema detailed
- [FLUTTER_MIGRATION_ANALYSIS.md](../FLUTTER_MIGRATION_ANALYSIS.md) - Android app analysis
- [.agent.md](../.agent.md) - AI agent context for development
- [SKILL.md](../SKILL.md) - Development skills/tools

---

**Backend Brainstorm Status:** 📝 In Planning  
**Ready to Start Development:** ✅ Yes  
**Next Step:** Implement Week 1 auth features

