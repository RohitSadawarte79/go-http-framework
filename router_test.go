package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func TestConflictDetection_SameParamStructure(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for conflicting routes, but got none")
		} else {
			t.Logf("✅ Correctly panicked: %v", r)
		}
	}()

	router := NewRouter()
	router.HandleFunc("GET", "/user/:id", dummyHandler)
	router.HandleFunc("GET", "/user/:name", dummyHandler) // should panic
}

func TestConflictDetection_DifferentMethod_NoConflict(t *testing.T) {
	router := NewRouter()
	router.HandleFunc("GET", "/user/:id", dummyHandler)
	router.HandleFunc("POST", "/user/:id", dummyHandler) // different method, OK
	t.Log("✅ Different methods don't conflict")
}

func TestConflictDetection_DifferentStatic_NoConflict(t *testing.T) {
	router := NewRouter()
	router.HandleFunc("GET", "/user/:id", dummyHandler)
	router.HandleFunc("GET", "/product/:id", dummyHandler) // different static, OK
	t.Log("✅ Different static segments don't conflict")
}

func TestConflictDetection_ParamVsStatic_NoConflict(t *testing.T) {
	router := NewRouter()
	router.HandleFunc("GET", "/user/:id", dummyHandler)
	router.HandleFunc("GET", "/user/profile", dummyHandler) // param vs static, OK
	t.Log("✅ Param vs static don't conflict")
}

func TestConflictDetection_DifferentLength_NoConflict(t *testing.T) {
	router := NewRouter()
	router.HandleFunc("GET", "/user", dummyHandler)
	router.HandleFunc("GET", "/user/:id", dummyHandler) // different length, OK
	t.Log("✅ Different segment counts don't conflict")
}

func TestConflictDetection_EmptyMethod_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty method")
		} else {
			t.Logf("✅ Correctly panicked: %v", r)
		}
	}()

	router := NewRouter()
	router.HandleFunc("", "/user", dummyHandler)
}

func TestConflictDetection_NilHandler_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil handler")
		} else {
			t.Logf("✅ Correctly panicked: %v", r)
		}
	}()

	router := NewRouter()
	router.HandleFunc("GET", "/user", nil)
}

func TestStaticBeatsParam(t *testing.T) {
	router := NewRouter()

	// Register param route FIRST — without sorting, this would always win
	router.HandleFunc("GET", "/user/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("param-handler"))
	})
	// Register static route SECOND
	router.HandleFunc("GET", "/user/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("static-handler"))
	})

	// Request /user/profile — should hit static handler, not param handler
	req := httptest.NewRequest("GET", "/user/profile", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	body := rec.Body.String()
	if body != "static-handler" {
		t.Errorf("Expected 'static-handler' but got '%s' — static route should beat param route", body)
	} else {
		t.Log("✅ Static route correctly beats param route regardless of registration order")
	}

	// Request /user/42 — should still hit param handler
	req2 := httptest.NewRequest("GET", "/user/42", nil)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	body2 := rec2.Body.String()
	if body2 != "param-handler" {
		t.Errorf("Expected 'param-handler' but got '%s'", body2)
	} else {
		t.Log("✅ Param route still works for non-static values")
	}
}

func TestParamExtractCorrectness(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("GET", "/user/:id", func(w http.ResponseWriter, r *http.Request) {
		id := URLParam(r, "id")
		w.Write([]byte(id))
	})

	req := httptest.NewRequest("GET", "/user/42", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	body := rec.Body.String()

	if body != "42" {
		t.Errorf("Expected '42' but got %s", body)
	} else if body == "42" {
		t.Log("Working got the '42' value when passed the route /user/42")
	}
}

func TestTrailingSlashNormalization(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("GET", "/user", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user-list"))
	})

	req := httptest.NewRequest("GET", "/user/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	body := rec.Body.String()

	if body != "user-list" {
		t.Errorf("was expecting 'user-list' but getting:- '%s'", body)
	} else {
		t.Logf("Test passed successfully getting:- '%s' as expected.", body)
	}
}

//multi -param extraction

func TestMultiParamExtraction(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("GET", "/user/:id/post/:postId", func(w http.ResponseWriter, r *http.Request) {
		id := URLParam(r, "id")
		postId := URLParam(r, "postId")
		body := id + ":" + postId
		w.Write([]byte(body))
	})

	req := httptest.NewRequest("GET", "/user/42/post/7", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	body := rec.Body.String()

	if body != "42:7" {
		t.Errorf("Expected '42:7' but got '%s'.", body)
	} else {
		t.Log("Test passed got '42:7' as expected.")
	}
}

func TestRootPathReturnsNotFound(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("GET", "/user", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("getting user...."))
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected 404 status not found but got: %v.", rec.Code)
	} else {
		t.Logf("Test passed got %d.", http.StatusNotFound)
	}
}

func TestPanicAtEmptyParam(t *testing.T) {
	router := NewRouter()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for empty params.")
		} else {
			t.Logf("Test passed expected panic, v: %v.", r)
		}
	}()

	router.HandleFunc("GET", "/user/:", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("getting user..."))
	})
}

func TestMethodNotAllowed(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("GET", "/user/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("getting user... got something"))
	})

	req := httptest.NewRequest("DELETE", "/user/42", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != 405 {
		t.Errorf("expecting 405 got %d.", rec.Code)
	} else {
		t.Logf("Test passed successfully corectly got 405.")
	}
}

func TestChainExecutionOrder(t *testing.T) {
	order := []string{}

	mwA := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "A")
			next.ServeHTTP(w, r)
		})
	}

	mwB := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "B")
			next.ServeHTTP(w, r)
		})
	}

	mwC := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "C")
			next.ServeHTTP(w, r)
		})
	}

	router := NewRouter()

	router.HandleFunc("GET", "/user", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.Write([]byte("Calling the handler"))
	})

	stack := Chain(mwA, mwB, mwC)(router)

	req := httptest.NewRequest("GET", "/user", nil)
	rec := httptest.NewRecorder()

	stack.ServeHTTP(rec, req)

	res := strings.Join(order, ",")

	if res != "A,B,C,handler" {
		t.Errorf("Expecting 'A,B,C,handler' but got %s.", res)
	} else {
		t.Logf("Test case passed successfully got %s.", res)
	}

}

func TestRecoveryMiddleware(t *testing.T) {
	router := NewRouter()

	router.HandleFunc("GET", "/user/:id", func(w http.ResponseWriter, r *http.Request) {
		panic("Something broke!")
	}) // this should panic

	stack := Chain(Recovery, Logger)(router)

	req := httptest.NewRequest("GET", "/user/42", nil)
	rec := httptest.NewRecorder()

	stack.ServeHTTP(rec, req)

	if rec.Code != 500 {
		t.Errorf("Expecting status code 500 but got %d.", rec.Code)
	} else {
		t.Logf("Test successfully passed! got 500 status.")
	}

}

func TestFullMiddlewareStackIntegration(t *testing.T) {
	// 1. Setup Phase
	router := NewRouter()

	// A dummy handler that writes a specific body so we know it ran
	router.HandleFunc("GET", "/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("integration-success"))
	})

	// Setup our strict CORS config
	corsConfig := CORSConfig{
		AllowedOrigins: map[string]bool{"http://localhost:3000": true},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}

	// Build the exact stack we use in main.go
	stack := Chain(NewCORS(corsConfig), Recovery, Logger, RequestId)(router)

	// 2. Execution Phase: Simulate a valid cross-origin GET request
	req := httptest.NewRequest("GET", "/api/data", nil)
	req.Header.Set("Origin", "http://localhost:3000") // Pretend we are the frontend

	rec := httptest.NewRecorder()

	// Pass the request through the ENTIRE stack
	stack.ServeHTTP(rec, req)

	// ==========================================
	// 3. YOUR TURN: ASSERTION PHASE
	// Write the checks for the following:
	// ==========================================

	// Check 1: Did the handler run? (Check that rec.Body.String() == "integration-success")
	if rec.Body.String() == "integration-success" {
		t.Logf("Passed! handler runs sucessfully go 'integration-success' j")
	} else {
		t.Errorf("Failed! expected 'integration-success' and got %s.", rec.Body.String())
	}

	// Check 2: Did CORS run and allow it? (Check that the "Access-Control-Allow-Origin" header equals "http://localhost:3000")
	if rec.Header().Values("Access-Control-Allow-Origin")[0] == "http://localhost:3000" {
		t.Logf("Passed! CORS run successfully and allowed it got 'Access-Control-Allow-Origin' as 'http://localhost:3000'")
	} else {
		t.Errorf("Failed! expected 'http://localhost:3000' but got %s.", rec.Header().Values("Access-Control-Allow-Origin")[0])
	}

	// Check 3: Did RequestId middleware run? (Check that "X-Request-ID" header is NOT empty)
	if rec.Header().Values("X-Request-ID") != nil {
		t.Logf("Passed! the RequestId middleware is running expected X-Request-ID got ID")
	} else {
		t.Errorf("Failed! the RequestId middleware is not running expected X-Request-ID but it's empty")
	}
}
