package integration_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/workforce-pro/workforce-payroll/internal/httpapi"
	"github.com/workforce-pro/workforce-payroll/internal/identity"
	"github.com/workforce-pro/workforce-payroll/internal/persistence"
)

const (
	integrationIssuer   = "https://identity.example.test/realms/workforce"
	integrationAudience = "workforce-api"
	integrationTenant   = "018f6f4a-7b2c-7d11-8a33-123456789abc"
	otherTenant         = "018f6f4a-7b2c-7d11-9b44-abcdef123456"
)

type integrationFixture struct {
	pool       *pgxpool.Pool
	server     *httptest.Server
	privateKey jwk.Key
}

func TestTask21PrincipalAndTenantBoundaryIntegration(t *testing.T) {
	fixture := newIntegrationFixture(t)

	t.Run("human Principal receives only server-derived Tenant", func(t *testing.T) {
		response := fixture.request(t, "human-subject", identity.PrincipalKindHuman, "/v1/context", "")
		assertResponse(t, response, http.StatusOK, `"tenant_id":"`+integrationTenant+`"`)
	})

	t.Run("Integration Principal receives only server-derived Tenant", func(t *testing.T) {
		response := fixture.request(t, "integration-subject", identity.PrincipalKindIntegration, "/v1/context", "")
		assertResponse(t, response, http.StatusOK, `"principal_kind":"INTEGRATION"`)
		// The helper deliberately adds a conflicting tenant_id token claim;
		// persistence, not that claim, remains authoritative.
		response = fixture.request(t, "integration-subject", identity.PrincipalKindIntegration, "/v1/context", "")
		assertResponse(t, response, http.StatusOK, `"tenant_id":"`+integrationTenant+`"`)
	})

	t.Run("matching caller-supplied Tenant is compared then accepted", func(t *testing.T) {
		response := fixture.request(t, "human-subject", identity.PrincipalKindHuman, "/v1/tenants/"+integrationTenant+"/context", integrationTenant)
		assertResponse(t, response, http.StatusOK, `"tenant_id":"`+integrationTenant+`"`)
	})

	t.Run("cross-Tenant path and header substitutions have identical non-disclosing denial", func(t *testing.T) {
		before := fixture.principalSnapshot(t)
		pathResponse := fixture.request(t, "human-subject", identity.PrincipalKindHuman, "/v1/tenants/"+otherTenant+"/context", "")
		headerResponse := fixture.request(t, "human-subject", identity.PrincipalKindHuman, "/v1/context", otherTenant)
		pathBody := assertResponse(t, pathResponse, http.StatusNotFound, `"code":"RESOURCE_NOT_FOUND"`)
		headerBody := assertResponse(t, headerResponse, http.StatusNotFound, `"code":"RESOURCE_NOT_FOUND"`)
		if pathBody != headerBody {
			t.Fatalf("cross-Tenant outcomes differ: path=%q header=%q", pathBody, headerBody)
		}
		if after := fixture.principalSnapshot(t); after != before {
			t.Fatalf("boundary denial changed authoritative state: before=%q after=%q", before, after)
		}
	})

	t.Run("disabled and unknown Principals share safe unauthenticated response", func(t *testing.T) {
		disabled := fixture.request(t, "disabled-subject", identity.PrincipalKindHuman, "/v1/context", "")
		unknown := fixture.request(t, "unknown-subject", identity.PrincipalKindHuman, "/v1/context", "")
		disabledBody := assertResponse(t, disabled, http.StatusUnauthorized, `"code":"UNAUTHENTICATED"`)
		unknownBody := assertResponse(t, unknown, http.StatusUnauthorized, `"code":"UNAUTHENTICATED"`)
		if disabledBody != unknownBody {
			t.Fatalf("authentication denials disclose registry state: disabled=%q unknown=%q", disabledBody, unknownBody)
		}
	})

	t.Run("request without bearer token is rejected", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, fixture.server.URL+"/v1/context", nil)
		if err != nil {
			t.Fatal(err)
		}
		response, err := fixture.server.Client().Do(request)
		if err != nil {
			t.Fatal(err)
		}
		assertResponse(t, response, http.StatusUnauthorized, `"code":"UNAUTHENTICATED"`)
	})
	t.Run("Principal-store outage fails closed without misreporting credentials", func(t *testing.T) {
		fixture.pool.Close()
		response := fixture.request(t, "human-subject", identity.PrincipalKindHuman, "/v1/context", "")
		assertResponse(t, response, http.StatusServiceUnavailable, `"code":"SERVICE_UNAVAILABLE"`)
	})
}

func newIntegrationFixture(t *testing.T) integrationFixture {
	t.Helper()
	databaseURL, cleanup := integrationDatabase(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	migrator, err := persistence.NewMigrator(pool)
	if err != nil {
		t.Fatal(err)
	}
	if err := migrator.Up(ctx); err != nil {
		t.Fatal(err)
	}
	if err := migrator.Up(ctx); err != nil {
		t.Fatalf("idempotent migration retry failed: %v", err)
	}

	seedPrincipal := func(id, tenant, subject, kind, status string) {
		t.Helper()
		_, err := pool.Exec(ctx, `
			INSERT INTO principal_bindings (principal_id, tenant_id, issuer, subject, principal_kind, status)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, id, tenant, integrationIssuer, subject, kind, status)
		if err != nil {
			t.Fatal(err)
		}
	}
	seedPrincipal("018f6f4a-7b2c-7d11-aa55-001122334455", integrationTenant, "human-subject", "HUMAN", "ACTIVE")
	seedPrincipal("018f6f4a-7b2c-7d11-bb66-112233445566", integrationTenant, "integration-subject", "INTEGRATION", "ACTIVE")
	seedPrincipal("018f6f4a-7b2c-7d11-cc77-223344556677", otherTenant, "disabled-subject", "HUMAN", "DISABLED")

	rawKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := jwk.Import(rawKey)
	if err != nil {
		t.Fatal(err)
	}
	if err := privateKey.Set("kid", "integration-key"); err != nil {
		t.Fatal(err)
	}
	if err := privateKey.Set("alg", jwa.RS256()); err != nil {
		t.Fatal(err)
	}
	publicKey, err := privateKey.PublicKey()
	if err != nil {
		t.Fatal(err)
	}
	keys := jwk.NewSet()
	if err := keys.AddKey(publicKey); err != nil {
		t.Fatal(err)
	}
	verifier, err := identity.NewJWTVerifier(integrationIssuer, integrationAudience, keys, 0)
	if err != nil {
		t.Fatal(err)
	}
	repository, err := identity.NewPostgresPrincipalRepository(pool)
	if err != nil {
		t.Fatal(err)
	}
	authenticator, err := identity.NewAuthenticator(verifier, repository)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := httpapi.NewServer(authenticator)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return integrationFixture{pool: pool, server: server, privateKey: privateKey}
}

func (fixture integrationFixture) request(t *testing.T, subject string, kind identity.PrincipalKind, path, suppliedTenant string) *http.Response {
	t.Helper()
	now := time.Now().UTC()
	token, err := jwt.NewBuilder().
		Issuer(integrationIssuer).
		Subject(subject).
		Audience([]string{integrationAudience}).
		IssuedAt(now.Add(-time.Second)).
		Expiration(now.Add(time.Minute)).
		Claim("principal_kind", string(kind)).
		Claim("tenant_id", otherTenant).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	serialized, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), fixture.privateKey))
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest(http.MethodGet, fixture.server.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+string(serialized))
	if suppliedTenant != "" {
		request.Header.Set("X-Tenant-ID", suppliedTenant)
	}
	response, err := fixture.server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func (fixture integrationFixture) principalSnapshot(t *testing.T) string {
	t.Helper()
	var snapshot string
	if err := fixture.pool.QueryRow(context.Background(), `
		SELECT string_agg(
			concat_ws('|', principal_id::text, tenant_id::text, issuer, subject, principal_kind, status, created_at::text, updated_at::text),
			',' ORDER BY principal_id
		)
		FROM principal_bindings
	`).Scan(&snapshot); err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func assertResponse(t *testing.T, response *http.Response, wantStatus int, wantBodyFragment string) string {
	t.Helper()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != wantStatus {
		t.Fatalf("status = %d, want %d, body=%s", response.StatusCode, wantStatus, body)
	}
	if !strings.Contains(string(body), wantBodyFragment) {
		t.Fatalf("body = %q, want fragment %q", body, wantBodyFragment)
	}
	if response.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", response.Header.Get("Cache-Control"))
	}
	return string(body)
}

func integrationDatabase(t *testing.T) (string, func()) {
	t.Helper()
	if databaseURL := os.Getenv("TEST_DATABASE_URL"); databaseURL != "" {
		return databaseURL, func() {}
	}
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatal("TEST_DATABASE_URL is unset and docker-compatible container CLI is unavailable")
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()
	name := fmt.Sprintf("workforce-task21-%d", time.Now().UnixNano())
	command := exec.Command("docker", "run", "--detach", "--rm", "--name", name,
		"--env", "POSTGRES_PASSWORD=integration-secret",
		"--env", "POSTGRES_DB=workforce_task21",
		"--publish", fmt.Sprintf("127.0.0.1:%d:5432", port),
		"postgres:17.6-bookworm")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("start PostgreSQL container: %v: %s", err, output)
	}
	cleanup := func() {
		_ = exec.Command("docker", "stop", "--time", "2", name).Run()
	}
	databaseURL := fmt.Sprintf("postgres://postgres:integration-secret@127.0.0.1:%d/workforce_task21?sslmode=disable", port)
	deadline := time.Now().Add(45 * time.Second)
	for time.Now().Before(deadline) {
		pool, err := pgxpool.New(context.Background(), databaseURL)
		if err == nil {
			err = pool.Ping(context.Background())
			pool.Close()
		}
		if err == nil {
			return databaseURL, cleanup
		}
		time.Sleep(250 * time.Millisecond)
	}
	cleanup()
	t.Fatal("PostgreSQL container did not become ready")
	return "", func() {}
}
